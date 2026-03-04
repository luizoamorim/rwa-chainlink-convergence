//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	protos "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	// 🔥 IMPORT DO BINDING GERADO
	"auto-lock-defi/contracts/evm/src/generated/vehicle_token_consumer"
)

type Config struct {
	DetranApiUrl         string `json:"detranApiMock"`
	ChainSelector        string `json:"chainSelector"`
	TokenizationContract string `json:"tokenizationContract"`
	WorldIDAppID         string `json:"worldIdAppId"`
}

type WorldIDProof struct {
	NullifierHash     string `json:"nullifier_hash"`
	MerkleRoot        string `json:"merkle_root"`
	Proof             string `json:"proof"`
	VerificationLevel string `json:"verification_level"`
}

type TokenizationPayload struct {
	Plate   string       `json:"plate"`
	Renavam string       `json:"renavam"`
	Wallet  string       `json:"wallet"`
	Proof   WorldIDProof `json:"proof"`
}

type ExecutionResult struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

func verifyWorldID(cfg *Config, runtime cre.Runtime, proof WorldIDProof) error {

	secretPromise := runtime.GetSecret(&protos.SecretRequest{
		Id: "WORLD_ID_API_KEY",
	})

	secret, err := secretPromise.Await()
	if err != nil {
		return err
	}

	apiKey := secret.Value

	payload := map[string]interface{}{
		"nullifier_hash":     proof.NullifierHash,
		"merkle_root":        proof.MerkleRoot,
		"proof":              proof.Proof,
		"verification_level": proof.VerificationLevel,
		"action":             "tokenizevehicle",
	}

	body, _ := json.Marshal(payload)

	client := &http.Client{}

	_, err = http.SendRequest(
		cfg,
		runtime,
		client,
		func(cfg *Config, log *slog.Logger, requester *http.SendRequester) ([]byte, error) {

			req := &http.Request{
				Url:    fmt.Sprintf("https://developer.worldcoin.org/api/v2/verify/%s", cfg.WorldIDAppID),
				Method: "POST",
				Headers: map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer " + apiKey,
				},
				Body: body,
			}

			resp, err := requester.SendRequest(req).Await()
			if err != nil {
				return nil, err
			}

			return resp.Body, nil
		},
		cre.ConsensusIdenticalAggregation[[]byte](),
	).Await()

	return err
}

func writeReportOnChain(
	config *Config,
	runtime cre.Runtime,
	payload TokenizationPayload,
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
			Value:   big.NewInt(45000),
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
func onTokenizationRequest(config *Config, runtime cre.Runtime, trigger *http.Payload) (*ExecutionResult, error) {

	var payload TokenizationPayload
	if err := json.Unmarshal(trigger.Input, &payload); err != nil {
		return nil, err
	}

	if err := verifyWorldID(config, runtime, payload.Proof); err != nil {
		return nil, err
	}

	txHash, err := writeReportOnChain(config, runtime, payload)
	if err != nil {
		return nil, err
	}

	return &ExecutionResult{
		TxHash: txHash,
		Status: "SUCCESS",
	}, nil
}

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(
			http.Trigger(&http.Config{}),
			onTokenizationRequest,
		),
	}, nil
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
