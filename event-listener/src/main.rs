use axum::{
    extract::ws::{Message, WebSocket, WebSocketUpgrade},
    response::IntoResponse,
    routing::get,
    Router,
};
use dotenvy::dotenv;
use ethers::prelude::*;
use std::{env, net::SocketAddr, sync::Arc};
use tokio::sync::broadcast;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    dotenv().ok();

    let ws_rpc = env::var("TENDERLY_WS_URL")?;
    let contract_address: Address = env::var("VEHICLE_NFT_ADDRESS")?.parse()?;

    let ws = Ws::connect(ws_rpc).await?;
    let provider = Provider::new(ws);
    let provider = Arc::new(provider);

    // Broadcast channel para WebSocket clients
    let (tx, _) = broadcast::channel::<String>(100);

    // Spawn listener
    let provider_clone = provider.clone();
    let tx_clone = tx.clone();

    tokio::spawn(async move {
        listen_events(provider_clone, contract_address, tx_clone).await;
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

async fn listen_events(
    provider: Arc<Provider<Ws>>,
    contract_address: Address,
    tx: broadcast::Sender<String>,
) {
    let filter = Filter::new().address(contract_address);

    let mut stream = provider.subscribe_logs(&filter).await.unwrap();

    println!("👂 Listening for contract events...");

    while let Some(log) = stream.next().await {
        println!("📦 Event received: {:?}", log);

        let payload = serde_json::json!({
            "type": "VEHICLE_TOKENIZED",
            "tx_hash": log.transaction_hash,
        });

        let _ = tx.send(payload.to_string());
    }
}
