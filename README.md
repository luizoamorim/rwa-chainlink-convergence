# 🚗 AutoLock DeFi - RWA Automotive Tokenization

![Chainlink](https://img.shields.io/badge/Orchestration-Chainlink_CRE-blue)
![Solidity](https://img.shields.io/badge/Smart_Contracts-Solidity-black)
![Foundry](https://img.shields.io/badge/Framework-Foundry-red)
![Golang](https://img.shields.io/badge/Language-Go-00ADD8)
![WASM](https://img.shields.io/badge/Compiled_to-WASM-purple)
![WorldID](https://img.shields.io/badge/Identity-World_ID-brightgreen)
![Tenderly](https://img.shields.io/badge/Simulation-Tenderly-orange)
![thirdweb](https://img.shields.io/badge/Frontend-thirdweb-111111)
![License](https://img.shields.io/badge/License-MIT-success)

AutoLock DeFi is a protocol for **tokenizing real-world vehicles as blockchain assets**.

The system converts vehicle ownership data into **ERC-721 NFTs**, enabling the creation of **Real World Assets (RWA)** that can be integrated with decentralized finance.

The platform integrates:

• **Chainlink Runtime Environment (CRE)** for decentralized workflow orchestration  
• **World ID** for human verification  
• **Thirdweb** for wallet and frontend integration  
• **Tenderly Virtual Testnet** for blockchain infrastructure

---

# 🏗 Architecture Overview

The protocol architecture is composed of four main layers:

```
Frontend
↓
Worker API
↓
Chainlink CRE Workflow
↓
Smart Contracts
```

Detailed architecture documentation:

📖 **System Architecture**  
➡️ [README_ARCHITECTURE.md](README_ARCHITECTURE.md)

---

# 📚 Component Documentation

Each component of the system has its own documentation.

| Component | Documentation |
|--------|--------|
| CRE Workflow | [auto-lock-defi/README.md](auto-lock-defi/README.md) |
| Frontend | [frontend/README_FRONTEND.md](frontend/README_FRONTEND.md) |
| Worker | [worker/README.md](worker/README.md) |
| API Mocks | [mocks/README.md](mocks/README.md) |
| Smart Contracts | [contracts/README.md](contracts/README.md) |

---

# ⚠️ Setup Requirements

This project depends on external platforms.

Before running the project you must configure accounts for:

• Thirdweb  
• World ID  
• Tenderly  

Setup guides:

- 📦 [Thirdweb Setup](README_THIRDWEB.md)
- 🧑‍🚀 [World ID Setup](README_WORLD_ID.md)
- ⛓️ [Tenderly Setup](README_TENDERLY.md)

---

# 📁 Project Structure

```
rwa-chainlink-convergence
│
├── frontend
│   Next.js Web3 interface
│
├── worker
│   Backend orchestrator triggering CRE workflows
│
├── mocks
│   Mock APIs for vehicle registry data
│
├── event-listener
│   Rust service monitoring blockchain events
│
├── auto-lock-defi
│   Chainlink Runtime Environment workflow
│
└── contracts
    Solidity smart contracts
```

---

# 🚀 Quick Start

Clone the repository:

```bash
git clone https://github.com/<your-repo>
cd rwa-chainlink-convergence
```

Install dependencies:

```bash
make install
```

Deploy smart contracts:

```bash
make deploy
```

Start the platform:

```bash
make up
```

---

# 🧪 Run CRE Simulation

You can manually simulate the RWA workflow:

```bash
make simulate-rwa
```

---

# 👨‍💻 Hackathon Project

This project was developed as part of the **Chainlink CRE ecosystem**.

The goal is to demonstrate how decentralized workflows can be used to tokenize **real-world assets**.

Core features:

• vehicle registry verification  
• human identity verification  
• decentralized oracle execution  
• automated NFT minting