# Chainlink Runtime Environment (CRE)

This project uses the **Chainlink Runtime Environment (CRE)** to orchestrate secure interactions between Web2 services and on-chain smart contracts.

The CRE acts as the automation layer responsible for:

- HTTP API interactions
- Identity verification
- Oracle data validation
- Smart contract execution

This allows our project to build a **deterministic workflow that bridges off-chain verification with on-chain settlement**.

---

# Why CRE?

In a Real World Asset tokenization system we must combine multiple trust layers:

• Identity verification  
• External data validation  
• Blockchain settlement  

Normally this requires several backend services.

The **Chainlink Runtime Environment** simplifies this by providing a **single secure workflow engine** capable of executing hybrid Web2 + Web3 pipelines.

Benefits:

- Deterministic execution
- Secure off-chain computation
- Native blockchain interaction
- Capability-based architecture

---

# CRE Capabilities Used

### HTTP Capability

Used to call external APIs.

Examples:

- World ID verification
- DETRAN vehicle registry API

---

### EVM Capability

Used to interact with smart contracts.

Responsibilities:

- Send reports to contracts
- Trigger NFT minting

---

### WASM Execution

The workflow logic runs inside a **WASM module written in Go**.

Advantages:

- deterministic execution
- sandboxed runtime
- reproducible workflows

---

# Mock Forwarder

In production Chainlink architecture the flow is:

Chainlink DON → Forwarder → Consumer Contract

The forwarder validates oracle reports before sending them to the target contract.

Since this project runs on a **Tenderly Virtual Testnet**, we use a **Mock Forwarder** for development.

This keeps the architecture aligned with production systems while allowing local testing.

---

# CRE Workflow Architecture

```mermaid
flowchart LR

A[Worker Payload] --> B[CRE Trigger]

B --> C[World ID Verification]

C --> D[DETRAN API]

D --> E[Generate Vehicle Report]

E --> F[Forwarder Contract]

F --> G[Consumer Contract]

G --> H[Vehicle NFT Mint]


---

# 🔄 Full Tokenization Flow (Mermaid)

Esse mostra **todo o fluxo desde frontend até o mint**.

```markdown
# Full Tokenization Flow

```mermaid
sequenceDiagram

participant Frontend
participant Worker
participant CRE
participant WorldID
participant DETRAN
participant Forwarder
participant Consumer
participant NFT

Frontend->>Worker: Tokenization Request

Worker->>CRE: HTTP Trigger Payload

CRE->>WorldID: Verify Proof
WorldID-->>CRE: Proof Valid

CRE->>DETRAN: Fetch Vehicle Data
DETRAN-->>CRE: Vehicle Value + Status

CRE->>Forwarder: Send Verified Report

Forwarder->>Consumer: processReport()

Consumer->>NFT: mintVehicle()

NFT-->>Frontend: NFT Minted

---

# Vehicle Report

The workflow generates a structured report sent to the contract.

Example:

```
struct VehicleReport {
    address owner;
    string plate;
    string renavam;
    uint256 value;
    string uri;
}
```

This report contains verified vehicle data retrieved from the DETRAN API.

---

# Why CRE Matters for RWA

Tokenizing real-world assets requires:

• trusted data  
• identity verification  
• deterministic execution  

The CRE provides the infrastructure necessary to bridge **real-world validation with blockchain settlement**.

---

# Summary

The Chainlink Runtime Environment is the **core automation engine of this project**.

It securely combines:

- World ID identity verification
- DETRAN vehicle validation
- On-chain NFT minting

All executed in a **single verifiable workflow**.

---

# Important References

These are the main resources used during the development of this project.

### Project Configuration

https://docs.chain.link/cre/reference/project-configuration-go#checking-your-versions

---

### Onchain Write Capability

https://docs.chain.link/cre/guides/workflow/using-evm-client/onchain-write/overview-go

---

### Forwarder Directory

https://docs.chain.link/cre/guides/workflow/using-evm-client/forwarder-directory-go#simulation-mainnets

---

### Consumer Contract (ReceiverTemplate)

https://docs.chain.link/cre/guides/workflow/using-evm-client/onchain-write/building-consumer-contracts#3-using-receivertemplate