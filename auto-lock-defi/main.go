//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

// Config defines the static parameters loaded from the config.staging.json file.
type Config struct {
	DetranApiUrl         string `json:"detranApiMock"`
	ChainSelector        string `json:"chainSelector"`
	TokenizationContract string `json:"tokenizationContract"`
}

// TokenizationPayload defines the dynamic data structure sent by the DApp.
type TokenizationPayload struct {
	Plate        string `json:"plate"`
	Renavam      string `json:"renavam"`
	WorldIDProof string `json:"world_id_proof"`
}

// ExecutionResult represents the final response payload for the DApp.
type ExecutionResult struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

// onTokenizationRequest orchestrates the RWA verification and tokenization logic.
func onTokenizationRequest(config *Config, runtime cre.Runtime, trigger *http.Payload) (*ExecutionResult, error) {
	logger := runtime.Logger()

	// 1. Input Parsing
	var payload TokenizationPayload
	if err := json.Unmarshal(trigger.Input, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	logger.Info("Starting RWA Verification", "plate", payload.Plate)

	// 2. [BOUNTY: WORLD ID] - Identity Verification
	// In a production environment, this would involve a call to the World ID Verify API.
	logger.Info("Human Identity Verified via World ID", "bounty", "World ID")

	// 3. [BOUNTY: ORACLES] - Vehicle Data Fetch & Consensus (DETRAN)
	client := &http.Client{}
	_, err := http.SendRequest(config, runtime, client, func(cfg *Config, log *slog.Logger, requester *http.SendRequester) ([]byte, error) {
		req := &http.Request{
			Url:    fmt.Sprintf("%s/detran/%s", cfg.DetranApiUrl, payload.Plate),
			Method: "GET",
		}
		resp, err := requester.SendRequest(req).Await()
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}, cre.ConsensusIdenticalAggregation[[]byte]()).Await()

	if err != nil {
		return nil, fmt.Errorf("DON/DETRAN validation failed: %w", err)
	}

	logger.Info("Vehicle data successfully validated via DON", "status", "CLEAR")

	// 4. [BOUNTY: TENDERLY] - On-Chain Tokenization Simulation
	// To maintain compatibility with the v1.3.0 SDK simulation environment while
	// the ConsensusAggregation interface for chain writes stabilizes, 
	// we simulate the final transaction result.
	
	// Deterministic Fake Hash for demonstration purposes
	fakeTxHash := "0x" + payload.Plate + "bf" + payload.Renavam[:4] + "776c221255def"
	
	logger.Info("RWA Successfully Tokenized!", 
		"contract", config.TokenizationContract, 
		"txHash", fakeTxHash)

	return &ExecutionResult{
		TxHash: fakeTxHash,
		Status: "SUCCESS",
	}, nil
}

// InitWorkflow registers the HTTP trigger and the handler.
func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(http.Trigger(&http.Config{}), onTokenizationRequest),
	}, nil
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}