use crate::protocol::{incoming::ClientMessage, outcoming::ServerMessage};
use serde_json;

pub fn serialize(msg: &ServerMessage) -> Result<String, serde_json::Error> {
    serde_json::to_string(msg)
}

pub fn deserealize(data: &str) -> Result<ClientMessage, serde_json::Error> {
    serde_json::from_str::<ClientMessage>(data)
}
