# Load environment variables
include .env
export

# Terminal Colors
GREEN  := \033[0;32m
RED    := \033[0;31m
NC     := \033[0m # No Color

.PHONY: build-contracts deploy simulate-rwa test-env

# Main command
test-env: deploy simulate-rwa

build-contracts:
	@echo "🛠️  Compiling Smart Contracts..."
	@forge build

deploy: build-contracts
	@echo "🚀 Deploying to Tenderly Virtual TestNet..."
	@$(eval CONTRACT_ADDRESS=$(shell forge script contracts/evm/script/DeployVehicleNFT.s.sol:DeployVehicleNFT --rpc-url $(TENDERLY_RPC_URL) --broadcast --non-interactive | grep "Deployed to:" | awk '{print $$3}'))
	@echo "✅ Contract deployed at: $(GREEN)$(CONTRACT_ADDRESS)$(NC)"
	@echo "📝 Updating auto-lock-defi/config.staging.json..."
	@if [ "$(shell uname)" = "Darwin" ]; then \
		sed -i '' 's/"tokenizationContract": ".*"/"tokenizationContract": "$(CONTRACT_ADDRESS)"/' auto-lock-defi/config.staging.json; \
	else \
		sed -i 's/"tokenizationContract": ".*"/"tokenizationContract": "$(CONTRACT_ADDRESS)"/' auto-lock-defi/config.staging.json; \
	fi

simulate-rwa:
	@echo "🌐 Starting Mock DETRAN Backend..."
	@go run mocks/main.go & sleep 3
	@echo "🧪 Running Chainlink CRE Simulation..."
	@if cre workflow simulate auto-lock-defi --target staging-settings --http-payload '{"plate":"ABC1234", "renavam":"123456789", "world_id_proof":"valid_proof_here"}' --trigger-index 0 --non-interactive; then \
		echo "$(GREEN)✅ SUCCESS: RWA Workflow validated on-chain!$(NC)"; \
	else \
		echo "$(RED)❌ ERROR: Simulation failed. Check logs above.$(NC)"; \
		pkill -f "go run mocks/main.go" || true; \
		exit 1; \
	fi
	@echo "🧹 Cleaning up processes..."
	@pkill -f "go run mocks/main.go" || true