use crate::protocol::status::Status;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Debug, PartialEq, Clone, Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum TypeC2S {
    Identify,
    Status,
    Users,
    Text,
    PublicText,
    NewRoom,
    Invite,
    JoinRoom,
    RoomUsers,
    RoomText,
    LeaveRoom,
    Disconnect,
}

#[derive(Serialize, Debug, Deserialize)]
pub struct ClientMessage {
    #[serde(rename = "type")]
    pub type_c2s: TypeC2S,

    #[serde(rename = "username", skip_serializing_if = "Option::is_none")]
    pub username: Option<String>,

    #[serde(rename = "roomname", skip_serializing_if = "Option::is_none")]
    pub room: Option<String>,

    #[serde(rename = "status", skip_serializing_if = "Option::is_none")]
    pub status: Option<Status>,

    #[serde(rename = "text", skip_serializing_if = "Option::is_none")]
    pub text: Option<String>,

    #[serde(rename = "usernames", skip_serializing_if = "Option::is_none")]
    pub usernames: Option<Vec<String>>,
}
