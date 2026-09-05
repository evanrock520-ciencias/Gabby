use crate::protocol::{incoming::TypeC2S, result::Result as ProtocolResult, status::Status};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(tag = "type", rename_all = "SCREAMING_SNAKE_CASE")]
pub enum TypeS2C {
    Response {
        operation: TypeC2S,
        result: ProtocolResult,
        #[serde(skip_serializing_if = "Option::is_none")]
        extra: Option<String>,
    },
    NewUser {
        username: String,
    },
    NewStatus {
        username: String,
        status: Status,
    },
    UserList {
        usernames: Vec<String>,
    },
    TextFrom {
        username: String,
        text: String,
    },
    PublicTextFrom {
        username: String,
        text: String,
    },
    JoinedRoom {
        username: String,
        roomname: String,
    },
    RoomUserList {
        roomname: String,
        usernames: Vec<String>,
    },
    RoomTextFrom {
        roomname: String,
        username: String,
        text: String,
    },
    LeftRoom {
        username: String,
        roomname: String,
    },
    Disconnected {
        username: String,
    },
    Invitation {
        username: String,
        roomname: String,
    },
}

impl TypeS2C {
    pub fn new_user_message(username: String) -> Self {
        Self::NewUser { username }
    }

    pub fn new_status_message(username: String, status: Status) -> Self {
        Self::NewStatus { username, status }
    }

    pub fn user_list_message(usernames: Vec<String>) -> Self {
        Self::UserList { usernames }
    }

    pub fn text_from_message(username: String, text: String) -> Self {
        Self::TextFrom { username, text }
    }

    pub fn public_text_from_message(username: String, text: String) -> Self {
        Self::PublicTextFrom { username, text }
    }

    pub fn joined_room_message(username: String, roomname: String) -> Self {
        Self::JoinedRoom { username, roomname }
    }

    pub fn room_user_list_message(roomname: String, usernames: Vec<String>) -> Self {
        Self::RoomUserList {
            roomname,
            usernames,
        }
    }

    pub fn room_text_from_message(roomname: String, username: String, text: String) -> Self {
        Self::RoomTextFrom {
            roomname,
            username,
            text,
        }
    }

    pub fn left_room_message(username: String, roomname: String) -> Self {
        Self::LeftRoom { username, roomname }
    }

    pub fn disconnected_message(username: String) -> Self {
        Self::Disconnected { username }
    }

    pub fn response_message(
        operation: TypeC2S,
        result: ProtocolResult,
        extra: Option<String>,
    ) -> Self {
        Self::Response {
            operation,
            result,
            extra,
        }
    }

    pub fn invitation_message(username: String, roomname: String) -> Self {
        Self::Invitation { username, roomname }
    }
}
