use std::sync::{Arc, Mutex};

use tokio::net::tcp::{OwnedReadHalf, OwnedWriteHalf};

use tokio::{
    io::{AsyncBufReadExt, AsyncWriteExt, BufReader},
    net::TcpStream,
    sync::mpsc,
};
use uuid::Uuid;

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
    let (tx, rx) = mpsc::unbounded_channel::<TypeS2C>();

    let client_id = match identify(&mut reader, &mut writer, tx, &hub).await {
        Some(id) => id,
        None => return Ok(()),
    };

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
/// El id del cliente conectado.
///
async fn identify(
    reader: &mut BufReader<OwnedReadHalf>,
    writer: &mut OwnedWriteHalf,
    tx: mpsc::UnboundedSender<TypeS2C>,
    hub: &Arc<Mutex<Hub>>,
) -> Option<Uuid> {
    let mut line: String = String::new();
    reader.read_line(&mut line).await.ok()?;

    let msg: ClientMessage = serializer::deserialize(line.trim()).ok()?;

    let username: String = msg.username?;
    let client: Client = Client::new(username, tx);
    let client_id: &Uuid = client.id();

    send_msg(
        writer,
        new_response(TypeC2S::Identify, MessageResult::Success, None),
    )
    .await;

    // TODO: Mandar la lista de usuarios
    // TODO: Avisar a los otros usuarios de la nueva conexión
    Some(*client_id)
}

async fn send_msg(writer: &mut OwnedWriteHalf, msg: TypeS2C) {
    if let Ok(json) = serializer::serialize(&msg) {
        let line = format!("{json}\n");
        writer.write_all(line.as_bytes()).await.ok();
    }
}
