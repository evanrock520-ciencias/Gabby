use crate::protocol::status::Status;
use serde::Serialize;

#[derive(Serialize, Debug, Clone, PartialEq)]
pub enum TypeS2C {
    #[serde(rename = "RESPONSE")]
    Response,
    #[serde(rename = "NEW_USER")]
    NewUser,
    #[serde(rename = "NEW_STATUS")]
    NewStatus,
    #[serde(rename = "USER_LIST")]
    UserList,
    #[serde(rename = "TEXT_FROM")]
    TextFrom,
    #[serde(rename = "PUBLIC_TEXT_FROM")]
    PublicTextFrom,
    #[serde(rename = "JOINED_ROOM")]
    JoinedRoom,
    #[serde(rename = "ROOM_USER_LIST")]
    RoomUserList,
    #[serde(rename = "ROOM_TEXT_FROM")]
    RoomTextFrom,
    #[serde(rename = "LEFT_ROOM")]
    LeftRoom,
    #[serde(rename = "DISCONNECTED")]
    Disconnected,
    #[serde(rename = "INVITATION")]
    Invitation,
}

#[derive(Serialize, Debug)]
pub struct ServerMessage {
    #[serde(rename = "type")]
    pub tipo: TypeS2C,

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
