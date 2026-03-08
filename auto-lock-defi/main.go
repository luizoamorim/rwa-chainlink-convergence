//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	// generated binding
	"auto-lock-defi/contracts/evm/src/generated/vehicle_token_consumer"
)

//////////////////////////////////////////////////////////////
// DEMO NULLIFIER STORAGE
//////////////////////////////////////////////////////////////

// var usedNullifiers = make(map[string]bool)

//////////////////////////////////////////////////////////////
// CONFIG
//////////////////////////////////////////////////////////////

type Config struct {
	DetranApiUrl         string `json:"detranApiMock"`
	ChainSelector        string `json:"chainSelector"`
	TokenizationContract string `json:"tokenizationContract"`
	WorldIDRpID          string `json:"worldIdRpId"`
	ActionName           string `json:"actionName"`
}

//////////////////////////////////////////////////////////////
// WORLD ID STRUCTS
//////////////////////////////////////////////////////////////

type WorldIDResponse struct {
	Identifier string `json:"identifier"`
	Nullifier  string `json:"nullifier"`
	MerkleRoot string `json:"merkle_root"`
	Proof      string `json:"proof"`
	SignalHash string `json:"signal_hash"`
}

type WorldIDProof struct {
	ProtocolVersion string            `json:"protocol_version"`
	Nonce           string            `json:"nonce"`
	Action          string            `json:"action"`
	Environment     string            `json:"environment"`
	Responses       []WorldIDResponse `json:"responses"`
}

//////////////////////////////////////////////////////////////
// REQUEST PAYLOAD
//////////////////////////////////////////////////////////////

type TokenizationPayload struct {
	Plate   string       `json:"plate"`
	Renavam string       `json:"renavam"`
	Wallet  string       `json:"wallet"`
	Proof   WorldIDProof `json:"proof"`
	Action  string       `json:"action"`
}

//////////////////////////////////////////////////////////////
// RESPONSE
//////////////////////////////////////////////////////////////

type ExecutionResult struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

//////////////////////////////////////////////////////////////
// WORLD ID VERIFICATION
//////////////////////////////////////////////////////////////

func verifyWorldID(cfg *Config, runtime cre.Runtime, wallet string, proof WorldIDProof) error {

	if len(proof.Responses) == 0 {
		return fmt.Errorf("no world ID responses found")
	}

	if proof.Action != cfg.ActionName {
		return fmt.Errorf("invalid world id action")
	}

	fmt.Println("🔎 Verifying World Id proof")

	body, _ := json.Marshal(proof)

	client := &http.Client{}

	_, err := http.SendRequest(
		cfg,
		runtime,
		client,
		func(cfg *Config, logger *slog.Logger, requester *http.SendRequester) ([]byte, error) {

			req := &http.Request{
				Url: fmt.Sprintf(
					"https://developer.world.org/api/v4/verify/%s",
					cfg.WorldIDRpID,
				),
				Method: "POST",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: body,
			}

			resp, err := requester.SendRequest(req).Await()
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("world id verification failed")
			}

			return resp.Body, nil
		},
		cre.ConsensusIdenticalAggregation[[]byte](),
	).Await()

	if err != nil {
		return err
	}

	fmt.Println("✅ World ID verified")

	return nil
}

//////////////////////////////////////////////////////////////
// DETRAN FETCH
//////////////////////////////////////////////////////////////

type DetranResponse struct {
	Plate     string  `json:"plate"`
	Status    string  `json:"status"`
	Fines     float64 `json:"fines"`
	ModelCode string  `json:"model_code"`
	Price     float64 `json:"price"`
}

func fetchDetran(
	cfg *Config,
	runtime cre.Runtime,
	plate string,
) (*DetranResponse, error) {

	client := &http.Client{}

	respBytes, err := http.SendRequest(
		cfg,
		runtime,
		client,
		func(cfg *Config, logger *slog.Logger, requester *http.SendRequester) ([]byte, error) {

			req := &http.Request{
				Url:    fmt.Sprintf("%s/detran/%s", cfg.DetranApiUrl, plate),
				Method: "GET",
			}

			resp, err := requester.SendRequest(req).Await()
			if err != nil {
				return nil, err
			}

			return resp.Body, nil
		},
		cre.ConsensusIdenticalAggregation[[]byte](),
	).Await()

	if err != nil {
		return nil, err
	}

	var detran DetranResponse

	err = json.Unmarshal(respBytes, &detran)
	if err != nil {
		return nil, err
	}

	if detran.Status != "clear" {
		return nil, fmt.Errorf("vehicle blocked")
	}

	return &detran, nil
}

//////////////////////////////////////////////////////////////
// ON-CHAIN WRITE
//////////////////////////////////////////////////////////////

func writeReportOnChain(
	config *Config,
	runtime cre.Runtime,
	payload TokenizationPayload,
	detran *DetranResponse,
) (string, error) {

	selector, err := evm.ChainSelectorFromName(config.ChainSelector)
	if err != nil {
		return "", err
	}

	evmClient := &evm.Client{
		ChainSelector: selector,
	}

	consumerAddress := common.HexToAddress(config.TokenizationContract)

	consumerContract, err := vehicle_token_consumer.NewVehicleTokenConsumer(
		evmClient,
		consumerAddress,
		nil,
	)
	if err != nil {
		return "", err
	}

	gasConfig := &evm.GasConfig{
		GasLimit: 500000,
	}

	writePromise := consumerContract.WriteReportFromVehicleReport(
		runtime,
		vehicle_token_consumer.VehicleReport{
			Owner:   common.HexToAddress(payload.Wallet),
			Plate:   payload.Plate,
			Renavam: payload.Renavam,
			Value:   big.NewInt(int64(detran.Price)),
			Uri:     "https://metadata.example/vehicle.json",
		},
		gasConfig,
	)

	resp, err := writePromise.Await()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("0x%x", resp.TxHash), nil
}

//////////////////////////////////////////////////////////////
// HANDLER
//////////////////////////////////////////////////////////////

func onTokenizationRequest(
	config *Config,
	runtime cre.Runtime,
	trigger *http.Payload,
) (*ExecutionResult, error) {

	var payload TokenizationPayload

	if err := json.Unmarshal(trigger.Input, &payload); err != nil {
		return nil, err
	}

	fmt.Println("🚀 Received tokenization request PAYLOAD:", payload)

	// 1️⃣ Verify WorldID
	if err := verifyWorldID(config, runtime, payload.Wallet, payload.Proof); err != nil {
		return nil, err
	}

	// 2️⃣ Fetch DETRAN
	detran, err := fetchDetran(config, runtime, payload.Plate)
	if err != nil {
		return nil, err
	}

	fmt.Println("DETRAN OK:", detran.Plate, "Price:", detran.Price)

	// 3️⃣ Write on-chain
	txHash, err := writeReportOnChain(config, runtime, payload, detran)
	if err != nil {
		return nil, err
	}

	return &ExecutionResult{
		TxHash: txHash,
		Status: "SUCCESS",
	}, nil
}

//////////////////////////////////////////////////////////////
// WORKFLOW
//////////////////////////////////////////////////////////////

func InitWorkflow(
	config *Config,
	logger *slog.Logger,
	secretsProvider cre.SecretsProvider,
) (cre.Workflow[*Config], error) {

	return cre.Workflow[*Config]{
		cre.Handler(
			http.Trigger(&http.Config{}),
			onTokenizationRequest,
		),
	}, nil
}

//////////////////////////////////////////////////////////////
// ENTRYPOINT
//////////////////////////////////////////////////////////////

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
