use axum::{
    extract::ws::{Message, WebSocket, WebSocketUpgrade},
    response::IntoResponse,
    routing::get,
    Router,
};
use dotenvy::dotenv;
use ethers::prelude::*;
use futures_util::StreamExt;
use std::{env, net::SocketAddr, sync::Arc, time::Duration};
use tokio::sync::broadcast;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    dotenv().ok();

    let ws_rpc = env::var("TENDERLY_WS_URL")?;
    let contract_address: Address = env::var("VEHICLE_NFT_ADDRESS")?.parse()?;

    // Broadcast channel para WebSocket clients
    let (tx, _) = broadcast::channel::<String>(100);

    // Spawn listener
    let tx_clone = tx.clone();
    tokio::spawn(async move {
        run_listener(ws_rpc, contract_address, tx_clone).await;
    });

    // WebSocket server
    let app = Router::new().route("/ws", get(move |ws| ws_handler(ws, tx.subscribe())));

    let addr = SocketAddr::from(([0, 0, 0, 0], 4000));
    println!("🚀 WebSocket server running on ws://localhost:4000/ws");

    axum::serve(tokio::net::TcpListener::bind(addr).await?, app).await?;

    Ok(())
}

async fn ws_handler(
    ws: WebSocketUpgrade,
    mut rx: broadcast::Receiver<String>,
) -> impl IntoResponse {
    ws.on_upgrade(move |mut socket: WebSocket| async move {
        while let Ok(msg) = rx.recv().await {
            let _ = socket.send(Message::Text(msg)).await;
        }
    })
}

//////////////////////////////////////////////////////////////
// MAIN LISTENER LOOP
//////////////////////////////////////////////////////////////

async fn run_listener(ws_rpc: String, contract_address: Address, tx: broadcast::Sender<String>) {
    loop {
        println!("🔌 Connecting to blockchain...");

        match connect_and_listen(&ws_rpc, contract_address, tx.clone()).await {
            Ok(_) => {}
            Err(err) => {
                println!("⚠ Listener error: {:?}", err);
            }
        }

        println!("🔁 Reconnecting in 3 seconds...");
        tokio::time::sleep(Duration::from_secs(3)).await;
    }
}

//////////////////////////////////////////////////////////////
// CONNECT AND STREAM EVENTS
//////////////////////////////////////////////////////////////

async fn connect_and_listen(
    ws_rpc: &str,
    contract_address: Address,
    tx: broadcast::Sender<String>,
) -> anyhow::Result<()> {
    let ws = Ws::connect(ws_rpc).await?;
    let provider = Provider::new(ws);
    let provider = Arc::new(provider);

    let transfer_topic = H256::from_slice(&ethers::utils::keccak256(
        "Transfer(address,address,uint256)",
    ));

    let zero_topic = H256::from_slice(&[0u8; 32]);

    let filter = Filter::new()
        .address(contract_address)
        .topic0(transfer_topic)
        .topic1(zero_topic);

    println!("📜 Checking past mint events...");

    let past_logs = provider.get_logs(&filter).await?;

    for log in past_logs {
        handle_log(&log, &tx);
    }

    println!("👂 Listening for new NFT mint events...");

    let mut stream = provider.subscribe_logs(&filter).await?;

    while let Some(log) = stream.next().await {
        handle_log(&log, &tx);
    }

    Ok(())
}

//////////////////////////////////////////////////////////////
// LOG HANDLER
//////////////////////////////////////////////////////////////

fn handle_log(log: &Log, tx: &broadcast::Sender<String>) {
    let Some(topic) = log.topics.get(3) else {
        return;
    };

    let token_id = U256::from(topic.0);

    println!("🔥 NFT MINT DETECTED");
    println!("Token ID: {}", token_id);
    println!("Tx: {:?}", log.transaction_hash);

    let payload = serde_json::json!({
        "type": "VEHICLE_TOKENIZED",
        "token_id": token_id,
        "tx_hash": log.transaction_hash
    });

    let _ = tx.send(payload.to_string());
}
