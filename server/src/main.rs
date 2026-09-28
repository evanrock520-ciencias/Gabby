mod client;
mod connection;
mod hub;
mod protocol;
mod room;

use std::sync::{Arc, Mutex};

use tokio::net::TcpListener;
use tokio::signal;
use tokio::sync::broadcast;
use tokio::task::JoinSet;

use crate::hub::Hub;

#[tokio::main]
async fn main() {
    let addr = std::env::args()
        .nth(1)
        .unwrap_or_else(|| "127.0.0.1:9090".to_string());

    let listener = TcpListener::bind(&addr).await.unwrap();
    let hub = Arc::new(Mutex::new(Hub::new()));
    let (shutdown_tx, _) = broadcast::channel(1);
    let mut tasks = JoinSet::new();

    println!("Waiting a connection on {}", listener.local_addr().unwrap());

    loop {
        tokio::select! {
            res = listener.accept() => {
                match res {
                    Ok((socket, addr)) => {
                        println!("Received a connection petition from {}", addr);
                        let client_hub = Arc::clone(&hub);
                        let shutdown_rx = shutdown_tx.subscribe();

                        tasks.spawn(async move {
                            if let Err(e) = connection::handle(socket, client_hub, shutdown_rx).await {
                                eprintln!("Failed to handle the connection: {}", e);
                            }
                        });
                    }
                    Err(e) => {
                        eprintln!("Error accepting connection: {}", e);
                    }
                }
            }

            _ = signal::ctrl_c() => {
                println!("The server will shutdown.");
                break;
            }
        }
    }

    let _ = shutdown_tx.send(());
    while let Some(_) = tasks.join_next().await {}

    println!("The server shut down.")
}
