use std::sync::{Arc, Mutex};

use tokio::{
    io::{AsyncBufReadExt, AsyncRead, AsyncWrite, AsyncWriteExt, BufReader, split},
    sync::mpsc,
};

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
use crate::{hub, protocol::status::Status};

/// Maneja la conexión del cliente.
///
/// # Arguments
///
/// * `stream` - El stream bidireccional establecido con el cliente.
/// * `hub` - Referencia compartida al hub central del servidor.
///
pub async fn handle<S>(stream: S, hub: Arc<Mutex<Hub>>) -> std::io::Result<()>
where
    S: AsyncRead + AsyncWrite + Unpin + Send + 'static,
{
    let (reader, mut writer) = split(stream);
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
async fn identify<R, W>(
    reader: &mut BufReader<R>,
    writer: &mut W,
    tx: mpsc::UnboundedSender<TypeS2C>,
    _hub: &Arc<Mutex<Hub>>,
) -> Option<String>
where
    R: AsyncRead + Unpin,
    W: AsyncWrite + Unpin,
{
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

            Some(username)
        }
        Err(e) => {
            println!("The username {} is already used.", username);
            send_msg(writer, new_response(TypeC2S::Identify, e, None)).await;
            None
        }
    }
}

/// Maneja el mensaje USERS.
///
/// # Arguments
///
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
///
async fn handle_users_list<W>(hub: &Arc<Mutex<Hub>>, writer: &mut W)
where
    W: AsyncWrite + Unpin,
{
    let users = {
        let hub = hub.lock().unwrap();
        hub.usernames()
    };

    send_msg(writer, TypeS2C::UserList { usernames: users }).await
}

/// Maneja el mensaje STATUS.
///
/// # Arguments
///
/// * `username` - El nombre de usuario del cliente.
/// * `status` - El nuevo status del cliente.
/// * `hub` - Referencia compartida al hub central del servidor.
///
///
async fn handle_status(username: &str, status: Status, hub: &Arc<Mutex<Hub>>) {
    {
        let mut hub = hub.lock().unwrap();
        if hub.change_status(username, status) {
            hub.broadcast(
                &TypeS2C::NewStatus {
                    username: username.into(),
                    status,
                },
                username,
            );
        }
    }
}

/// Maneja el mensaje PUBLIC_TEXT.
///
/// # Arguments
///
/// * `username` - El nombre de usuario del cliente.
/// * `text` - El mensaje de texto enviado.
/// * `hub` - Referencia compartida al hub central del servidor.
///
///
async fn handle_public_text(username: &str, text: &str, hub: &Arc<Mutex<Hub>>) {
    {
        let hub = hub.lock().unwrap();
        hub.broadcast(
            &TypeS2C::PublicTextFrom {
                username: username.to_string(),
                text: text.to_string(),
            },
            username,
        );
    }
}

/// Maneja el mensaje CREATE_ROOM.
///
/// # Arguments
///
/// * `username` - El nombre de usuario del cliente.
/// * `roomname` - El nombre de la sala.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
async fn handle_create_room<W>(
    username: &str,
    roomname: &str,
    hub: &Arc<Mutex<Hub>>,
    writer: &mut W,
) where
    W: AsyncWrite + Unpin,
{
    let result = {
        let mut hub = hub.lock().unwrap();
        hub.register_room(roomname, username)
    };

    match result {
        Ok(_) => {
            send_msg(
                writer,
                new_response(TypeC2S::NewRoom, MessageResult::Success, None),
            )
            .await;
        }
        Err(e) => {
            send_msg(writer, new_response(TypeC2S::NewRoom, e, None)).await;
        }
    }
}

/// Maneja el mensaje TEXT.
///
/// # Arguments
///
/// * `sender` - El nombre de usuario del cliente que manda el mensaje.
/// * `roomname` - El nombre de usuario del cliente que recibe el mensaje.
/// * ``text - El mensaje de texto enviado.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
async fn handle_text<W>(
    sender: &str,
    receiver: &str,
    text: &str,
    hub: &Arc<Mutex<Hub>>,
    writer: &mut W,
) where
    W: AsyncWrite + Unpin,
{
    let msg = TypeS2C::TextFrom {
        username: sender.to_string(),
        text: text.to_string(),
    };

    let result = {
        let hub = hub.lock().unwrap();
        hub.send_to(&msg, sender, receiver)
    };

    match result {
        Ok(_) => {
            send_msg(
                writer,
                new_response(TypeC2S::Text, MessageResult::Success, None),
            )
            .await;
        }
        Err(e) => {
            // TODO: Manejar ambos tipos de errores, por el extra
            send_msg(writer, new_response(TypeC2S::Text, e, None)).await;
        }
    }
}

/// Maneja el mensaje INVITE.
///
/// # Arguments
///
/// * `username` - El nombre de usuario del cliente que manda el mensaje.
/// * `roomname` - La sala a la que se invita el cliente.
/// * `guests` - La lista de clientes invitados.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
async fn handle_invite<W>(
    username: &str,
    roomname: &str,
    guests: Vec<String>,
    hub: &Arc<Mutex<Hub>>,
    writer: &mut W,
) where
    W: AsyncWrite + Unpin,
{
    let result = {
        let mut hub = hub.lock().unwrap();
        hub.invite(roomname, guests.clone())
    };

    match result {
        Ok(_) => {
            send_msg(
                writer,
                new_response(TypeC2S::Invite, MessageResult::Success, None),
            )
            .await;

            {
                let hub = hub.lock().unwrap();
                let msg = &TypeS2C::Invitation {
                    username: username.to_string(),
                    roomname: roomname.to_string(),
                };

                hub.send_to_members(msg, username, guests).unwrap();
            }
        }

        Err(e) => {
            send_msg(writer, new_response(TypeC2S::Invite, e, None)).await;
        }
    }
}

/// Maneja el mensaje JOIN ROOM.
///
/// # Arguments
///
/// * `username` - El usuario que quiere unirse.
/// * `roomname` - La sala a unirse.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
async fn handle_join_room<W>(username: &str, roomname: &str, hub: &Arc<Mutex<Hub>>, writer: &mut W)
where
    W: AsyncWrite + Unpin,
{
    let result = {
        let mut hub = hub.lock().unwrap();
        hub.be_member_of(roomname, username)
    };

    match result {
        Ok(_) => {
            send_msg(
                writer,
                new_response(TypeC2S::JoinRoom, MessageResult::Success, None),
            )
            .await;

            let msg = &TypeS2C::JoinedRoom {
                username: username.to_string(),
                roomname: roomname.to_string(),
            };

            {
                let hub = hub.lock().unwrap();
                hub.to_room(msg, roomname, username).unwrap();
            }
        }

        Err(e) => {
            send_msg(writer, new_response(TypeC2S::JoinRoom, e, None)).await;
        }
    }
}

/// Maneja el mensaje LEAVE ROOM.
///
/// # Arguments
///
/// * `username` - El usuario que quiere abandonar la sala.
/// * `roomname` - La sala que quiere abandonar.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
async fn handle_leave_room<W>(username: &str, roomname: &str, hub: &Arc<Mutex<Hub>>, writer: &mut W)
where
    W: AsyncWrite + Unpin,
{
    let result = {
        let mut hub = hub.lock().unwrap();
        hub.leave_room(roomname, username)
    };

    match result {
        Ok(_) => {
            send_msg(
                writer,
                new_response(TypeC2S::LeaveRoom, MessageResult::Success, None),
            )
            .await;

            let msg = &TypeS2C::LeftRoom {
                username: username.to_string(),
                roomname: roomname.to_string(),
            };
            let hub = hub.lock().unwrap();
            hub.to_room(msg, roomname, username).unwrap();
        }

        Err(e) => {
            send_msg(writer, new_response(TypeC2S::LeaveRoom, e, None)).await;
        }
    }
}

/// Maneja el mensaje ROOM USERS.
///
/// # Arguments
///
/// * `username` - El usuario que pide la lista de usuarios.
/// * `roomname` - La sala sobre la que se pide la lista de usuarios.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
async fn handle_room_users<W>(username: &str, roomname: &str, hub: &Arc<Mutex<Hub>>, writer: &mut W)
where
    W: AsyncWrite + Unpin,
{
    let result = {
        let hub = hub.lock().unwrap();
        hub.room_usernames(roomname, username)
    };

    match result {
        Ok(usernames) => {
            let msg = TypeS2C::RoomUserList {
                roomname: roomname.to_string(),
                usernames,
            };

            send_msg(writer, msg).await;
        }

        Err(e) => {
            send_msg(writer, new_response(TypeC2S::RoomUsers, e, None)).await;
        }
    }
}

/// Maneja el mensaje ROOM TEXT.
///
/// # Arguments
///
/// * `username` - El usuario que manda el mensaje.
/// * `roomname` - La sala a la que manda el mensaje.
/// * `text` - El texto del mensaje que se manda.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
async fn handle_room_text<W>(
    username: &str,
    roomname: &str,
    text: &str,
    hub: &Arc<Mutex<Hub>>,
    writer: &mut W,
) where
    W: AsyncWrite + Unpin,
{
    let msg = &TypeS2C::RoomTextFrom {
        roomname: roomname.to_string(),
        username: username.to_string(),
        text: text.to_string(),
    };

    let result = {
        let hub = hub.lock().unwrap();
        if !hub.is_member_of(username, roomname) {
            Err(MessageResult::NotJoined)
        } else {
            hub.to_room(msg, roomname, username)
        }
    };

    match result {
        Ok(_) => {}
        Err(e) => {
            send_msg(writer, new_response(TypeC2S::RoomText, e, None)).await;
        }
    }
}

/// Enruta los mensajes a su handler correspondiente.
///
/// # Arguments
///
/// * `msg` - El mensaje deserializado que mando el cliente.
/// * `username` - El nombre de usuario del cliente.
/// * `hub` - Referencia compartida al hub central del servidor.
/// * `writer` - Escritor del socket del cliente.
///
///
async fn route_msg<W>(msg: ClientMessage, username: &str, hub: &Arc<Mutex<Hub>>, writer: &mut W)
where
    W: AsyncWrite + Unpin,
{
    match msg {
        ClientMessage::Users => handle_users_list(hub, writer).await,
        ClientMessage::Status { status } => handle_status(username, status, hub).await,
        ClientMessage::PublicText { text } => handle_public_text(username, &text, hub).await,
        ClientMessage::NewRoom { roomname } => {
            handle_create_room(username, &roomname, hub, writer).await
        }
        ClientMessage::Text {
            username: receiver,
            text,
        } => handle_text(username, &receiver, &text, hub, writer).await,
        ClientMessage::Invite {
            roomname,
            usernames,
        } => {
            handle_invite(username, &roomname, usernames, hub, writer).await;
        }
        ClientMessage::JoinRoom { roomname } => {
            handle_join_room(username, &roomname, hub, writer).await;
        }
        ClientMessage::LeaveRoom { roomname } => {
            handle_leave_room(username, &roomname, hub, writer).await;
        }
        ClientMessage::RoomUsers { roomname } => {
            handle_room_users(username, roomname.as_str(), hub, writer).await;
        }
        ClientMessage::RoomText { roomname, text } => {
            handle_room_text(username, &roomname, &text, hub, writer).await;
        }
        _ => return,
    }
}

async fn send_msg<W>(writer: &mut W, msg: TypeS2C)
where
    W: AsyncWrite + Unpin,
{
    if let Ok(json) = serializer::serialize(&msg) {
        let line = format!("{json}\n");
        writer.write_all(line.as_bytes()).await.ok();
    }
}

#[cfg(test)]
mod test {
    // TODO: Agregar pruebas para la conexión usando tokio::io::duplex
}
