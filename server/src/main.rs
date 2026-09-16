mod client;
mod connection;
mod hub;
mod protocol;
mod room;

use std::sync::{Arc, Mutex};

use tokio::net::TcpListener;

use crate::hub::Hub;

#[tokio::main]
async fn main() {
    let addr = std::env::args()
        .nth(1)
        .unwrap_or_else(|| "127.0.0.1:9090".to_string());

    let listener = TcpListener::bind(&addr).await.unwrap();
    let hub = Arc::new(Mutex::new(Hub::new()));

    println!("Waiting a connection on {}", listener.local_addr().unwrap());

    loop {
        let (socket, addr) = listener.accept().await.unwrap();
        println!("Received a connection petition from {}", addr);
        let client_hub = Arc::clone(&hub);
        tokio::spawn(async move {
            if let Err(e) = connection::handle(socket, client_hub).await {
                eprintln!("Failed to handle the connection: {}", e);
            }
        });
    }
}
