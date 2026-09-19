use crate::protocol::status::Status;
use serde::{Deserialize, Serialize};

#[derive(Copy, Serialize, Debug, PartialEq, Clone, Deserialize)]
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
    Invalid,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(tag = "type", rename_all = "SCREAMING_SNAKE_CASE")]
pub enum ClientMessage {
    Identify {
        username: String,
    },
    Status {
        status: Status,
    },
    Users,
    Text {
        username: String,
        text: String,
    },
    PublicText {
        text: String,
    },
    NewRoom {
        roomname: String,
    },
    Invite {
        roomname: String,
        usernames: Vec<String>,
    },
    JoinRoom {
        roomname: String,
    },
    RoomUsers {
        roomname: String,
    },
    RoomText {
        roomname: String,
        text: String,
    },
    LeaveRoom {
        roomname: String,
    },
    Disconnect,
}

impl ClientMessage {
    pub fn operation(&self) -> TypeC2S {
        match self {
            Self::Identify { .. } => TypeC2S::Identify,
            Self::Status { .. } => TypeC2S::Status,
            Self::Users => TypeC2S::Users,
            Self::Text { .. } => TypeC2S::Text,
            Self::PublicText { .. } => TypeC2S::PublicText,
            Self::NewRoom { .. } => TypeC2S::NewRoom,
            Self::Invite { .. } => TypeC2S::Invite,
            Self::JoinRoom { .. } => TypeC2S::JoinRoom,
            Self::RoomUsers { .. } => TypeC2S::RoomUsers,
            Self::RoomText { .. } => TypeC2S::RoomText,
            Self::LeaveRoom { .. } => TypeC2S::LeaveRoom,
            Self::Disconnect => TypeC2S::Disconnect,
        }
    }
}
