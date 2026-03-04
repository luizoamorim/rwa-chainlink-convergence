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

	// Generated contract bindings
	"auto-lock-defi/contracts/evm/src/generated/vehicle_token_consumer"
)

//////////////////////////////////////////////////////////////
// CONFIGURATION STRUCT
//////////////////////////////////////////////////////////////

// Config is loaded from config.json.
// It contains runtime parameters injected by CRE.
type Config struct {
	DetranApiUrl         string `json:"detranApiMock"`
	ChainSelector        string `json:"chainSelector"`
	TokenizationContract string `json:"tokenizationContract"`
	WorldIDRpID          string `json:"worldIdRpId"`
}

//////////////////////////////////////////////////////////////
// WORLD ID PROOF STRUCTS (v4-compatible format)
//////////////////////////////////////////////////////////////

// WorldIDResponse represents a single proof response
// returned by the World ID 4.0 widget.
type WorldIDResponse struct {
	Identifier string `json:"identifier"`
	Nullifier  string `json:"nullifier"`
	MerkleRoot string `json:"merkle_root"`
	Proof      string `json:"proof"`
	SignalHash string `json:"signal_hash"`
}

// WorldIDProof represents the full proof payload
// sent from the frontend to the CRE workflow.
type WorldIDProof struct {
	ProtocolVersion string            `json:"protocol_version"`
	Action          string            `json:"action"`
	Environment     string            `json:"environment"`
	Responses       []WorldIDResponse `json:"responses"`
}

//////////////////////////////////////////////////////////////
// TOKENIZATION PAYLOAD
//////////////////////////////////////////////////////////////

// TokenizationPayload is the full input
// received by the HTTP trigger.
type TokenizationPayload struct {
	Plate   string       `json:"plate"`
	Renavam string       `json:"renavam"`
	Wallet  string       `json:"wallet"`
	Proof   WorldIDProof `json:"proof"`
}

//////////////////////////////////////////////////////////////
// WORKFLOW RESULT STRUCT
//////////////////////////////////////////////////////////////

// ExecutionResult is returned as JSON
// to the worker after successful execution.
type ExecutionResult struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

//////////////////////////////////////////////////////////////
// WORLD ID VERIFICATION (v4)
//////////////////////////////////////////////////////////////

// verifyWorldID verifies a World ID proof
// using the official v4 endpoint (RP-based verification).
func verifyWorldID(cfg *Config, runtime cre.Runtime, proof WorldIDProof) error {

	// Ensure at least one proof response exists
	if len(proof.Responses) == 0 {
		return fmt.Errorf("no World ID responses found")
	}

	response := proof.Responses[0]

	// Build verification payload
	payload := map[string]interface{}{
		"nullifier":   response.Nullifier,
		"merkle_root": response.MerkleRoot,
		"proof":       response.Proof,
	}

	body, _ := json.Marshal(payload)

	client := &http.Client{}

	// Send request through CRE HTTP capability (consensus-aware)
	_, err := http.SendRequest(
		cfg,
		runtime,
		client,
		func(cfg *Config, logger *slog.Logger, requester *http.SendRequester) ([]byte, error) {

			req := &http.Request{
				Url:    fmt.Sprintf("https://developer.worldcoin.org/api/v4/verify/%s", cfg.WorldIDRpID),
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

			return resp.Body, nil
		},
		cre.ConsensusIdenticalAggregation[[]byte](),
	).Await()

	return err
}

//////////////////////////////////////////////////////////////
// ON-CHAIN WRITE
//////////////////////////////////////////////////////////////

// writeReportOnChain writes vehicle data
// to the deployed smart contract.
func writeReportOnChain(
	config *Config,
	runtime cre.Runtime,
	payload TokenizationPayload,
) (string, error) {

	// Convert chain selector string to CRE selector
	selector, err := evm.ChainSelectorFromName(config.ChainSelector)
	if err != nil {
		return "", err
	}

	evmClient := &evm.Client{
		ChainSelector: selector,
	}

	// Convert contract address
	consumerAddress := common.HexToAddress(config.TokenizationContract)

	// Instantiate contract binding
	consumerContract, err := vehicle_token_consumer.NewVehicleTokenConsumer(
		evmClient,
		consumerAddress,
		nil,
	)
	if err != nil {
		return "", err
	}

	// Define gas configuration
	gasConfig := &evm.GasConfig{
		GasLimit: 500000,
	}

	// Execute contract write
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

//////////////////////////////////////////////////////////////
// MAIN WORKFLOW HANDLER
//////////////////////////////////////////////////////////////

// onTokenizationRequest is executed when the HTTP trigger fires.
func onTokenizationRequest(
	config *Config,
	runtime cre.Runtime,
	trigger *http.Payload,
) (*ExecutionResult, error) {

	var payload TokenizationPayload

	// Parse incoming JSON payload
	if err := json.Unmarshal(trigger.Input, &payload); err != nil {
		return nil, err
	}

	// Step 1: Verify World ID proof
	if err := verifyWorldID(config, runtime, payload.Proof); err != nil {
		return nil, err
	}

	// Step 2: Write vehicle report on-chain
	txHash, err := writeReportOnChain(config, runtime, payload)
	if err != nil {
		return nil, err
	}

	return &ExecutionResult{
		TxHash: txHash,
		Status: "SUCCESS",
	}, nil
}

//////////////////////////////////////////////////////////////
// WORKFLOW INITIALIZATION
//////////////////////////////////////////////////////////////

// InitWorkflow registers the HTTP trigger handler.
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
// ENTRY POINT
//////////////////////////////////////////////////////////////

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
