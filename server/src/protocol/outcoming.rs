use crate::protocol::status::Status;
use serde::Serialize;

#[derive(Serialize, Debug, Clone, PartialEq)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum TypeS2C {
    Response,
    NewUser,
    NewStatus,
    UserList,
    TextFrom,
    PublicTextFrom,
    JoinedRoom,
    RoomUserList,
    RoomTextFrom,
    LeftRoom,
    Disconnected,
    Invitation,
}

#[derive(Serialize, Debug)]
pub struct ServerMessage {
    #[serde(rename = "type")]
    pub type_s2c: TypeS2C,

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

impl ServerMessage {
    pub fn new_user_message(username: String) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::NewUser,
            username: Some(username),
            room: None,
            status: None,
            text: None,
            usernames: None,
        }
    }

    pub fn new_status_message(username: String, status: Status) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::NewStatus,
            username: Some(username),
            room: None,
            status: Some(status),
            text: None,
            usernames: None,
        }
    }

    pub fn user_list_message(usernames: Option<Vec<String>>) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::UserList,
            username: None,
            room: None,
            status: None,
            text: None,
            usernames: usernames,
        }
    }

    pub fn text_from_message(username: String, text: String) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::TextFrom,
            username: Some(username),
            room: None,
            status: None,
            text: Some(text),
            usernames: None,
        }
    }

    pub fn public_text_from_message(username: String, text: String) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::PublicTextFrom,
            username: Some(username),
            room: None,
            status: None,
            text: Some(text),
            usernames: None,
        }
    }

    pub fn joined_room_message(username: String, roomname: String) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::JoinedRoom,
            username: Some(username),
            room: Some(roomname),
            status: None,
            text: None,
            usernames: None,
        }
    }

    pub fn room_user_list_message(roomname: String, usernames: Vec<String>) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::RoomUserList,
            username: None,
            room: Some(roomname),
            status: None,
            text: None,
            usernames: Some(usernames),
        }
    }

    pub fn room_text_from_message(
        roomname: String,
        username: String,
        text: String,
    ) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::RoomTextFrom,
            username: Some(username),
            room: Some(roomname),
            status: None,
            text: Some(text),
            usernames: None,
        }
    }

    pub fn left_room_message(username: String, roomname: String) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::LeftRoom,
            username: Some(username),
            room: Some(roomname),
            status: None,
            text: None,
            usernames: None,
        }
    }

    pub fn disconnected_message(username: String) -> ServerMessage {
        ServerMessage {
            type_s2c: TypeS2C::Disconnected,
            username: Some(username),
            room: None,
            status: None,
            text: None,
            usernames: None,
        }
    }
}
