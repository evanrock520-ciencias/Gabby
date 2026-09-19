use serde::{Deserialize, Serialize};

#[derive(Serialize, Debug, Clone, PartialEq, Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum MessageResult {
    Success,
    UserAlreadyExists,
    NoSuchUser,
    NoSuchRoom,
    RoomAlreadyExists,
    NotInvited,
    NotJoined,
    NotIdentified,
    Invalid,
}
