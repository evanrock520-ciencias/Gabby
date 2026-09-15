use std::any::Any;
use std::sync::{Arc, Mutex};

use tokio::net::tcp::{OwnedReadHalf, OwnedWriteHalf};

use tokio::{
    io::{AsyncBufReadExt, AsyncWriteExt, BufReader},
    net::TcpStream,
    sync::mpsc,
};

use crate::hub;
use crate::{
    client::Client,
    hub::Hub,
    protocol::{
        incoming::{ClientMessage, TypeC2S},
        outcoming::TypeS2C,
        response::new_response,
        result::MessageResult,
        serializer,
    },
};

/// Maneja la conexión del cliente.
///
/// # Arguments
///
/// * `socket` - La conexión TCP establecida con el cliente.
/// * `hub` - Referencia compartida al hub central del servidor.
///
///
pub async fn handle(socket: TcpStream, hub: Arc<Mutex<Hub>>) -> std::io::Result<()> {
    println!(
        "Received a connection petition from {}",
        socket.peer_addr()?
    );

    let (reader, mut writer) = socket.into_split();
    let mut reader = BufReader::new(reader);
    let (tx, mut rx) = mpsc::unbounded_channel::<TypeS2C>();

    let _username = match identify(&mut reader, &mut writer, tx, &hub).await {
        Some(name) => name,
        None => return Ok(()),
    };

    let mut line = String::new();
    loop {
        line.clear();

        tokio::select! {
            result = reader.read_line(&mut line) => {
                match result {
                    Ok(0) => break,
                    Ok(_) => {
                      // TODO: Un match para emparejar el tipo de mensaje con un handle
                      match serializer::deserialize(&line) {
                        Ok(msg) => {
                            println!("Valid message: {}", line.trim());
                            route_msg(msg, &_username, &hub, &mut writer).await;
                        },
                        Err(_) => {
                             eprintln!("Unexpected message format {}", line.trim());
                        }
                      }
                    }
                    Err(_) => break,
                }
            }

            Some(msg) = rx.recv() => {
                send_msg(&mut writer, msg).await;
            }

            else => break,
        }
    }

    Ok(())
}

/// Maneja la identificación del usuario y recibe el primer mensaje.
///
/// # Arguments
///
/// * `reader` - Lector del socket del cliente.
/// * `writer` - Escritor del socket del cliente.
/// * `tx` - Extremo emisor del canal por el cual el hub podrá enviarle mensajes a este cliente.
/// * `hub` - Referencia compartida al hub central del servidor.
///
/// # Returns
///
/// El nombre de usuario del cliente conectado.
///
async fn identify(
    reader: &mut BufReader<OwnedReadHalf>,
    writer: &mut OwnedWriteHalf,
    tx: mpsc::UnboundedSender<TypeS2C>,
    _hub: &Arc<Mutex<Hub>>,
) -> Option<String> {
    let mut line: String = String::new();
    reader.read_line(&mut line).await.ok()?;

    let msg: ClientMessage = serializer::deserialize(line.trim()).ok()?;

    let username: String = match msg {
        ClientMessage::Identify { username } => username,
        _ => return None,
    };
    let _client: Client = Client::new(username.clone(), tx);

    let register_result = {
        let mut hub = _hub.lock().unwrap();
        hub.register(_client)
    };

    match register_result {
        Ok(_) => {
            println!("The username {} was registered.", username);
            send_msg(
                writer,
                new_response(TypeC2S::Identify, MessageResult::Success, None),
            )
            .await;

            let usernames = {
                let hub = _hub.lock().unwrap();
                hub.usernames()
            };

            println!("Current users {:?}", usernames);
            send_msg(writer, TypeS2C::UserList { usernames }).await;

            {
                let hub = _hub.lock().unwrap();
                hub.broadcast(
                    &TypeS2C::NewUser {
                        username: username.clone(),
                    },
                    &username,
                );
            }
        }
        Err(e) => {
            println!("The username {} is already used.", username);
            send_msg(writer, new_response(TypeC2S::Identify, e, None)).await;
        }
    }

    Some(username)
}

async fn handle_users_list(username: &str, hub: &Arc<Mutex<Hub>>, writer: &mut OwnedWriteHalf) {
    println!("{} requires the user list", username);

    let users = {
        let hub = hub.lock().unwrap();
        hub.usernames()
    };

    send_msg(writer, TypeS2C::UserList { usernames: users }).await
}

async fn route_msg(
    msg: ClientMessage,
    username: &str,
    hub: &Arc<Mutex<Hub>>,
    writer: &mut OwnedWriteHalf,
) {
    match msg {
        ClientMessage::Users => handle_users_list(username, hub, writer).await,
        _ => return,
    }
}

async fn send_msg(writer: &mut OwnedWriteHalf, msg: TypeS2C) {
    if let Ok(json) = serializer::serialize(&msg) {
        let line = format!("{json}\n");
        writer.write_all(line.as_bytes()).await.ok();
    }
}
