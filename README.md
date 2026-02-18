# 🚗 AutoLock DeFi

## RWA Automotive Tokenization

![Chainlink](https://img.shields.io/badge/Orchestration-Chainlink_CRE-blue)
![Solidity](https://img.shields.io/badge/Smart_Contracts-Solidity-black)
![Foundry](https://img.shields.io/badge/Framework-Foundry-red)
![Golang](https://img.shields.io/badge/Language-Go-00ADD8)
![WASM](https://img.shields.io/badge/Compiled_to-WASM-purple)
![WorldID](https://img.shields.io/badge/Identity-World_ID-brightgreen)
![Tenderly](https://img.shields.io/badge/Simulation-Tenderly-orange)
![thirdweb](https://img.shields.io/badge/Frontend-thirdweb-111111)
![License](https://img.shields.io/badge/License-MIT-success)

------------------------------------------------------------------------

AutoLock DeFi is a decentralized lending protocol that bridges Brazilian
automotive equity with global liquidity. By tokenizing vehicle titles as
Real World Assets (RWA), we allow users to access instant stablecoin
liquidity while ensuring legal compliance through automated registry and
verifiable data.

------------------------------------------------------------------------

# 🏗 Architecture Overview

Built on the **Chainlink Runtime Environment (CRE)**, leveraging
decentralized orchestration for off-chain verification and on-chain
settlement.

------------------------------------------------------------------------

# 🔄 Trigger-Action-Target Model

## 1️⃣ Trigger (The "When")

Authenticated HTTP API Trigger capturing:

-   User intent\
-   Vehicle data (Plate / Renavam)\
-   World ID proof

## 2️⃣ Action (The "What")

### Identity Verification

World ID validation to prevent Sybil attacks.

### Oracle Fetch & Consensus

Multiple DON nodes fetch vehicle data (DETRAN/FIPE mocks) and reach
Byzantine Fault Tolerant consensus.

## 3️⃣ Target (The "Where")

EVM chain write target that:

-   Mints RWA NFT\
-   Releases liquidity\
-   Executes on Tenderly Virtual TestNet

------------------------------------------------------------------------

# 📊 System Workflow

``` mermaid
sequenceDiagram
    participant U as Vehicle Owner
    participant D as DApp (thirdweb)
    participant W as World ID
    participant C as Chainlink CRE (Go WASM)
    participant M as Mock API (DETRAN/FIPE)
    participant T as Tenderly Virtual TestNet

    U->>D: Enter Plate & Renavam
    D->>W: Request Personhood Verification
    W-->>D: Proof of Personhood
    D->>C: POST /gateway (Payload + Proof)
    
    C->>C: Verify World ID Proof
    
    C->>M: GET /detran/{plate}
    M-->>C: Vehicle Data & Price
    
    C->>T: mintVehicleNFT(owner, vehicleData)
    T-->>C: txHash
    C-->>D: 200 OK (ExecutionResult)
```

------------------------------------------------------------------------

# 🛠 Configuration

## workflow.yaml

Unlocks permissions:

-   networking:http\
-   blockchain:evm

## project.yaml

-   Chain Selector: 999999\
-   Forwarder: MockKeystoneForwarder

------------------------------------------------------------------------

# 🏆 Bounty Integrations

-   World ID (Sybil Resistance)\
-   Chainlink Oracles (Verifiable Data)\
-   Tenderly Virtual TestNets (Execution Transparency)\
-   thirdweb (Web3 UX)

------------------------------------------------------------------------

# 🚀 Technical Stack

-   Chainlink CRE\
-   Golang (WASM wasip1)\
-   Solidity (Foundry)\
-   World ID\
-   Tenderly\
-   thirdweb SDK

------------------------------------------------------------------------

# 📝 How to Run

## 1️⃣ Setup

Create `.env`:

TENDERLY_RPC_URL=\
PRIVATE_KEY=

## 2️⃣ Run Full Suite

make test-env

## 3️⃣ Manual

go run mocks/main.go\
cre workflow simulate --target staging-settings --broadcast

------------------------------------------------------------------------

# 🎯 Chainlink Constellation Hackathon Submission

AutoLock DeFi --- Bridging Brazilian automotive assets with global DeFi
liquidity.
