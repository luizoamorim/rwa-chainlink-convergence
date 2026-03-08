# Event Listener

The **Event Listener** is a Rust service responsible for monitoring blockchain events emitted by the `VehicleNFT` smart contract.

Its main purpose is to detect when a **new vehicle NFT is minted** and notify connected clients in real time through **WebSocket broadcasts**.

This allows the frontend to react instantly when a vehicle tokenization transaction is confirmed on-chain.

---

# Responsibilities

The event listener performs the following tasks:

• connects to the blockchain using a WebSocket RPC  
• subscribes to NFT `Transfer` events  
• filters mint events (`from = address(0)`)  
• extracts the token ID and transaction hash  
• broadcasts the result to WebSocket clients  

---

# Technology Stack

| Technology | Purpose |
|-----------|--------|
| Rust | Event listener implementation |
| Axum | HTTP / WebSocket server |
| Tokio | Async runtime |
| Ethers-rs | Ethereum RPC client |
| WebSocket | Real-time event broadcasting |

---

# WebSocket Endpoint

The listener exposes a WebSocket endpoint:

```
ws://localhost:4000/ws
```

Clients connecting to this endpoint will receive real-time events whenever a new NFT is minted.

---

# Event Detection

The listener subscribes to the standard ERC-721 **Transfer** event.

Mint events are detected when the `from` address equals `0x0000000000000000000000000000000000000000`.

Event signature:

```
Transfer(address,address,uint256)
```

Filter configuration:

```
topic0 = Transfer event signature
topic1 = zero address (mint)
```

---

# Event Payload

When a mint event is detected, the listener broadcasts a JSON message.

Example:

```json
{
  "type": "VEHICLE_TOKENIZED",
  "token_id": 1,
  "tx_hash": "0x..."
}
```

This message can be consumed by the frontend or other services.

---

# Environment Variables

The listener requires the following variables:

```
TENDERLY_WS_URL
VEHICLE_NFT_ADDRESS
```

Example `.env`:

```
TENDERLY_WS_URL=wss://...
VEHICLE_NFT_ADDRESS=0x...
```

---

# Running the Listener

Start the service using:

```
cargo run
```

The server will start on:

```
ws://localhost:4000/ws
```

---

# Summary

The event listener acts as the **real-time bridge between the blockchain and the frontend**.

By streaming NFT mint events through WebSockets, it allows the application to provide immediate feedback when a vehicle has been successfully tokenized on-chain.