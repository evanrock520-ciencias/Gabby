use crate::protocol::{incoming::ClientMessage, outcoming::TypeS2C};
use serde_json;

pub fn serialize(msg: &TypeS2C) -> Result<String, serde_json::Error> {
    serde_json::to_string(msg)
}

pub fn deserialize(data: &str) -> Result<ClientMessage, serde_json::Error> {
    serde_json::from_str::<ClientMessage>(data)
}

#[cfg(test)]
mod tests {

    use std::collections::HashMap;

    use super::*;
    use crate::protocol::{
        incoming::TypeC2S,
        outcoming::TypeS2C,
        status::Status::{Active, Away, Busy},
    };

    // Para comprobar que las piezas encajen, las cadenas seran serializadas a partir de los ClientMessages del cliente.
    #[test]
    fn test_deserealize_identify_message() {
        let data = r#"{"type":"IDENTIFY","username":"Evan"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::Identify {
                username: "Evan".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::Identify, msg.operation());
    }

    #[test]
    fn test_deserealize_status_message() {
        let data = r#"{"type":"STATUS","status":"AWAY"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(ClientMessage::Status { status: Away }, msg);
        assert_eq!(TypeC2S::Status, msg.operation());
    }

    #[test]
    fn test_deserealize_users_message() {
        let data = r#"{"type":"USERS"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(ClientMessage::Users, msg);
        assert_eq!(TypeC2S::Users, msg.operation());
    }

    #[test]
    fn test_deserealize_text_message() {
        let data = r#"{"type":"TEXT","username":"Evan","text":"Pasa las pruebas"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::Text {
                username: "Evan".to_string(),
                text: "Pasa las pruebas".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::Text, msg.operation());
    }

    #[test]
    fn test_deserealize_public_text_message() {
        let data = r#"{"type":"PUBLIC_TEXT","text":"Pasa las pruebas"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::PublicText {
                text: "Pasa las pruebas".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::PublicText, msg.operation());
    }

    #[test]
    fn test_deserealize_new_room_message() {
        let data = r#"{"type":"NEW_ROOM","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::NewRoom {
                roomname: "Sala 1".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::NewRoom, msg.operation());
    }

    #[test]
    fn test_deserealize_invite_message() {
        let data = r#"{"type":"INVITE","roomname":"Sala 1","usernames":["Derek","Yahel","Erik"]}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::Invite {
                roomname: "Sala 1".to_string(),
                usernames: vec!["Derek".to_string(), "Yahel".to_string(), "Erik".to_string()]
            },
            msg
        );
        assert_eq!(TypeC2S::Invite, msg.operation());
    }

    #[test]
    fn test_deserealize_join_room_message() {
        let data = r#"{"type":"JOIN_ROOM","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::JoinRoom {
                roomname: "Sala 1".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::JoinRoom, msg.operation());
    }

    #[test]
    fn test_deserealize_room_users_message() {
        let data = r#"{"type":"ROOM_USERS","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::RoomUsers {
                roomname: "Sala 1".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::RoomUsers, msg.operation());
    }

    #[test]
    fn test_deserealize_room_text_message() {
        let data = r#"{"type":"ROOM_TEXT","roomname":"Sala 1","text":"Pasa las pruebas"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::RoomText {
                roomname: "Sala 1".to_string(),
                text: "Pasa las pruebas".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::RoomText, msg.operation());
    }

    #[test]
    fn test_deserealize_leave_room_message() {
        let data = r#"{"type":"LEAVE_ROOM","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(
            ClientMessage::LeaveRoom {
                roomname: "Sala 1".to_string()
            },
            msg
        );
        assert_eq!(TypeC2S::LeaveRoom, msg.operation());
    }

    #[test]
    fn test_deserealize_disconnect_message() {
        let data = r#"{"type":"DISCONNECT"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(ClientMessage::Disconnect, msg);
        assert_eq!(TypeC2S::Disconnect, msg.operation());
    }

    // Tests de Serialización

    fn field(key: &str, value: &str) -> String {
        format!(r#""{}":"{}""#, key, value)
    }

    fn field_array(key: &str, value: &[String]) -> String {
        let json_value = serde_json::to_string(value).unwrap();
        format!(r#""{}":{}"#, key, json_value)
    }

    #[test]
    fn test_serialize_new_user_message() {
        let msg = TypeS2C::new_user_message("Evan".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "NEW_USER")));
        assert!(data.contains(&field("username", "Evan")));
    }

    #[test]
    fn test_serialize_new_status_message() {
        let msg = TypeS2C::new_status_message("Evan".to_string(), Busy);
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "NEW_STATUS")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("status", "BUSY")));
    }

    #[test]
    fn test_serialize_user_list_message() {
        let usernames = HashMap::from([("Derek".to_string(), Active), ("Yahel".to_string(), Busy)]);
        let msg = TypeS2C::user_list_message(usernames);
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "USER_LIST")));
        assert!(data.contains(&field("Derek", "ACTIVE")));
        assert!(data.contains(&field("Yahel", "BUSY")));
    }

    #[test]
    fn test_serialize_text_from_message() {
        let msg = TypeS2C::text_from_message("Evan".to_string(), "Hola amigos".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "TEXT_FROM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("text", "Hola amigos")));
    }

    #[test]
    fn test_serialize_public_text_from_message() {
        let msg = TypeS2C::public_text_from_message("Evan".to_string(), "Hola amigos".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "PUBLIC_TEXT_FROM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("text", "Hola amigos")));
    }

    #[test]
    fn test_serialize_joined_room_message() {
        let msg = TypeS2C::joined_room_message("Evan".to_string(), "Sala 1".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "JOINED_ROOM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("roomname", "Sala 1")));
    }

    #[test]
    fn test_serialize_room_user_list_message() {
        let mut usernames = HashMap::new();
        usernames.insert("Derek".to_string(), Active);
        usernames.insert("Yahel".to_string(), Away);
        usernames.insert("Luis".to_string(), Busy);
        let msg = TypeS2C::room_user_list_message("Sala 1".to_string(), usernames);
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "ROOM_USER_LIST")));
        assert!(data.contains(&field("roomname", "Sala 1")));
    }

    #[test]
    fn test_serialize_room_text_from_message() {
        let msg = TypeS2C::room_text_from_message(
            "Sala 1".to_string(),
            "Evan".to_string(),
            "Hola amigos".to_string(),
        );
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "ROOM_TEXT_FROM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("text", "Hola amigos")));
        assert!(data.contains(&field("roomname", "Sala 1")));
    }

    #[test]
    fn test_serialize_left_room_message() {
        let msg = TypeS2C::left_room_message("Evan".to_string(), "Sala 1".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "LEFT_ROOM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("roomname", "Sala 1")));
    }

    #[test]
    fn test_serialize_disconnected_message() {
        let msg = TypeS2C::disconnected_message("Evan".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "DISCONNECTED")));
        assert!(data.contains(&field("username", "Evan")));
    }

    #[test]
    fn test_serialize_response_message() {
        let msg = crate::protocol::response::new_response(
            TypeC2S::Identify,
            crate::protocol::result::MessageResult::Success,
            None,
        );
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "RESPONSE")));
        assert!(data.contains(&field("operation", "IDENTIFY")));
        assert!(data.contains(&field("result", "SUCCESS")));
        assert!(!data.contains("extra"));
    }

    #[test]
    fn test_serialize_response_message_with_extra() {
        let msg = TypeS2C::response_message(
            TypeC2S::Identify,
            crate::protocol::result::MessageResult::UserAlreadyExists,
            Some("Username already taken".to_string()),
        );
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "RESPONSE")));
        assert!(data.contains(&field("operation", "IDENTIFY")));
        assert!(data.contains(&field("result", "USER_ALREADY_EXISTS")));
        assert!(data.contains(&field("extra", "Username already taken")));
    }

    #[test]
    fn test_serialize_invitation_message() {
        let msg = TypeS2C::invitation_message("Evan".to_string(), "Sala 1".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "INVITATION")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("roomname", "Sala 1")));
    }

    #[test]
    fn test_roundtrip_types2c_variants() {
        let original = TypeS2C::NewUser {
            username: "Evan".to_string(),
        };
        let json_str = serialize(&original).unwrap();
        let parsed: TypeS2C = serde_json::from_str(&json_str).unwrap();
        assert_eq!(original, parsed);
    }
}
