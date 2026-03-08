# --------------------------------------
# LOAD ENV
# --------------------------------------
-include .env
export

# --------------------------------------
# TERMINAL COLORS
# --------------------------------------
GREEN  := \033[0;32m
RED    := \033[0;31m
NC     := \033[0m

# --------------------------------------
# PHONY
# --------------------------------------
.PHONY: install install-frontend install-go install-rust \
build-contracts deploy export-abi generate-bindings simulate-rwa \
worker detran listener frontend up down test-env demo

# --------------------------------------
# INSTALL (PROJECT SETUP)
# --------------------------------------
install: install-frontend install-go install-rust build-contracts generate-bindings
	@echo "$(GREEN)✅ Project environment ready!$(NC)"

# --------------------------------------
# FRONTEND SETUP
# --------------------------------------
install-frontend:
	@echo "📦 Installing frontend dependencies..."
	cd frontend && npm install

# --------------------------------------
# GO DEPENDENCIES
# --------------------------------------
install-go:
	@echo "🐹 Installing Go dependencies..."
	@echo "→ root modules"
	go mod tidy
	@echo "→ worker modules"
	cd worker && go mod tidy
	@echo "→ detran mock modules"
	cd mocks && go mod tidy
	@echo "→ CRE workflow modules"
	cd auto-lock-defi && go mod tidy

# --------------------------------------
# RUST LISTENER BUILD
# --------------------------------------
install-rust:
	@echo "🦀 Building event listener..."
	cd event-listener && cargo build

# --------------------------------------
# BUILD CONTRACTS
# --------------------------------------
build-contracts:
	@echo "🛠️  Compiling Smart Contracts..."
	forge build

# --------------------------------------
# DEPLOY (NFT + CONSUMER + OWNERSHIP)
# --------------------------------------
deploy: build-contracts
	@echo "🚀 Deploying VehicleNFT..."
	@forge script contracts/evm/script/DeployVehicleNFT.s.sol:DeployVehicleNFT \
		--rpc-url $(TENDERLY_RPC_URL) \
		--private-key $(CRE_ETH_PRIVATE_KEY) \
		--broadcast --non-interactive
	@NFT_ADDRESS=$$(jq -r '.transactions[0].contractAddress' \
		broadcast/DeployVehicleNFT.s.sol/99911155111/run-latest.json); \
	if [ -z "$$NFT_ADDRESS" ] || [ "$$NFT_ADDRESS" = "null" ]; then \
		echo "$(RED)❌ Failed to capture NFT address$(NC)"; \
		exit 1; \
	fi; \
	echo "✅ NFT deployed at: $(GREEN)$$NFT_ADDRESS$(NC)"; \
	if grep -q "^VEHICLE_NFT_ADDRESS=" .env; then \
		sed -i.bak "s|^VEHICLE_NFT_ADDRESS=.*|VEHICLE_NFT_ADDRESS=$$NFT_ADDRESS|" .env; \
	else \
		echo "VEHICLE_NFT_ADDRESS=$$NFT_ADDRESS" >> .env; \
	fi; \
	rm -f .env.bak
	@echo "🚀 Deploying VehicleTokenConsumer..."
	@NFT_ADDR=$$(grep VEHICLE_NFT_ADDRESS .env | cut -d '=' -f2); \
	VEHICLE_NFT_ADDRESS=$$NFT_ADDR forge script contracts/evm/script/DeployVehicleConsumer.s.sol:DeployVehicleConsumer \
		--rpc-url $(TENDERLY_RPC_URL) \
		--private-key $(CRE_ETH_PRIVATE_KEY) \
		--broadcast \
		--non-interactive
	@CONSUMER_ADDRESS=$$(jq -r '.transactions[0].contractAddress' \
		broadcast/DeployVehicleConsumer.s.sol/99911155111/run-latest.json); \
	if [ -z "$$CONSUMER_ADDRESS" ] || [ "$$CONSUMER_ADDRESS" = "null" ]; then \
		echo "$(RED)❌ Failed to capture Consumer address$(NC)"; \
		exit 1; \
	fi; \
	echo "✅ Consumer deployed at: $(GREEN)$$CONSUMER_ADDRESS$(NC)"; \
	if grep -q "^CONSUMER_ADDRESS=" .env; then \
		sed -i.bak "s|^CONSUMER_ADDRESS=.*|CONSUMER_ADDRESS=$$CONSUMER_ADDRESS|" .env; \
	else \
		echo "CONSUMER_ADDRESS=$$CONSUMER_ADDRESS" >> .env; \
	fi; \
	rm -f .env.bak
	@echo "🔄 Transferring NFT ownership..."
	@NFT_ADDR=$$(grep VEHICLE_NFT_ADDRESS .env | cut -d '=' -f2); \
	CONS_ADDR=$$(grep CONSUMER_ADDRESS .env | cut -d '=' -f2); \
	VEHICLE_NFT_ADDRESS=$$NFT_ADDR CONSUMER_ADDRESS=$$CONS_ADDR \
	forge script contracts/evm/script/TransferOwnership.s.sol:TransferOwnership \
		--rpc-url $(TENDERLY_RPC_URL) \
		--private-key $(CRE_ETH_PRIVATE_KEY) \
		--broadcast --non-interactive
	@echo "📝 Updating config.staging.json..."
	@CONS_ADDR=$$(grep CONSUMER_ADDRESS .env | cut -d '=' -f2); \
	sed -i.bak "s|\"tokenizationContract\": \".*\"|\"tokenizationContract\": \"$$CONS_ADDR\"|" auto-lock-defi/config.staging.json; \
	rm -f auto-lock-defi/config.staging.json.bak
	@echo "$(GREEN)🎉 Deploy completed successfully!$(NC)"
	@$(MAKE) export-abi
	@$(MAKE) generate-bindings

# --------------------------------------
# EXPORT ABI
# --------------------------------------
export-abi:
	@echo "📦 Exporting VehicleTokenConsumer ABI..."
	@mkdir -p auto-lock-defi/contracts/evm/src/abi
	@jq '.abi' contracts/evm/out/VehicleTokenConsumer.sol/VehicleTokenConsumer.json \
		> auto-lock-defi/contracts/evm/src/abi/VehicleTokenConsumer.abi
	@echo "✅ ABI exported successfully!"

# --------------------------------------
# GENERATE CRE BINDINGS
# --------------------------------------
generate-bindings:
	@echo "⚙️ Generating CRE bindings..."
	@cd auto-lock-defi && cre generate-bindings evm
	@cd auto-lock-defi && go mod tidy
	@echo "✅ Bindings generated."

# --------------------------------------
# CRE SIMULATION
# --------------------------------------
simulate-rwa:
	@echo "🌐 Starting Mock DETRAN Backend..."
	@go run mocks/main.go & sleep 3
	@echo "🧪 Running Chainlink CRE Simulation..."
	@if cre workflow simulate auto-lock-defi \
		--target staging-settings \
		-e .env \
		--broadcast \
		--http-payload '{"plate":"ABC1234","renavam":"123456789","wallet":"0x0000000000000000000000000000000000000001","proof":{"nullifier_hash":"0xabc","merkle_root":"0xdef","proof":"0xghi","verification_level":"device"}}' \
		--trigger-index 0 \
		--non-interactive; then \
		echo "$(GREEN)✅ SUCCESS: RWA Workflow validated!$(NC)"; \
	else \
		echo "$(RED)❌ ERROR: Simulation failed. Check logs above.$(NC)"; \
		pkill -f "go run mocks/main.go" || true; \
		exit 1; \
	fi
	@echo "🧹 Cleaning up processes..."
	@pkill -f "go run mocks/main.go" || true

# --------------------------------------
# INDIVIDUAL SERVICES
# --------------------------------------
worker:
	cd worker && go run .

detran:
	cd detran-mock && go run .

listener:
	cd event-listener && cargo run

frontend:
	cd frontend && npm run dev

# --------------------------------------
# FULL STACK
# --------------------------------------
up:
	@echo "🚀 Starting full stack..."
	@make -j4 worker detran listener frontend

down:
	@echo "🛑 Stopping services..."
	@pkill -f worker || true
	@pkill -f detran-mock || true
	@pkill -f event-listener || true
	@pkill -f next || true

# --------------------------------------
# TEST ENVIRONMENT
# --------------------------------------
test-env: deploy simulate-rwa

# --------------------------------------
# DEMO (ONE COMMAND FOR HACKATHON)
# --------------------------------------
demo:
	$(MAKE) install
	$(MAKE) deploy
	$(MAKE) up