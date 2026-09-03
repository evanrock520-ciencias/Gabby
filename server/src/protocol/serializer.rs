use crate::protocol::{incoming::ClientMessage, outcoming::ServerMessage};
use serde_json;

pub fn serialize(msg: &ServerMessage) -> Result<String, serde_json::Error> {
    serde_json::to_string(msg)
}

pub fn deserialize(data: &str) -> Result<ClientMessage, serde_json::Error> {
    serde_json::from_str::<ClientMessage>(data)
}

#[cfg(test)]
mod tests {

    use crate::protocol::{
        incoming::TypeC2S::{self},
        status::Status::{Away, Busy},
    };

    use super::*;

    // Para comprobar que las piezas encajen, las cadenas seran serializadas a partir de los ClientMessages del cliente.
    #[test]
    fn test_deserealize_identify_message() {
        let data = r#"{"type":"IDENTIFY","username":"Evan"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(TypeC2S::Identify, msg.type_c2s);
        assert_eq!(Some("Evan".to_string()), msg.username);
    }

    #[test]
    fn test_deserealize_status_message() {
        let data = r#"{"type":"STATUS","status":"AWAY"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(TypeC2S::Status, msg.type_c2s);
        assert_eq!(Some(Away), msg.status);
    }

    #[test]
    fn test_deserealize_users_message() {
        let data = r#"{"type":"USERS"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(TypeC2S::Users, msg.type_c2s);
    }

    #[test]
    fn test_deserealize_text_message() {
        let data = r#"{"type":"TEXT","username":"Evan","text":"Pasa las pruebas"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(TypeC2S::Text, msg.type_c2s);
        assert_eq!(Some("Evan".to_string()), msg.username);
        assert_eq!(Some("Pasa las pruebas".to_string()), msg.text);
    }

    #[test]
    fn test_deserealize_public_text_message() {
        let data = r#"{"type":"PUBLIC_TEXT","text":"Pasa las pruebas"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(TypeC2S::PublicText, msg.type_c2s);
        assert_eq!(Some("Pasa las pruebas".to_string()), msg.text);
    }

    #[test]
    fn test_deserealize_new_room_message() {
        let data = r#"{"type":"NEW_ROOM","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();
        assert_eq!(TypeC2S::NewRoom, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_invite_message() {
        let data = r#"{"type":"INVITE","roomname":"Sala 1","usernames":["Derek","Yahel","Erik"]}"#;
        let msg = deserialize(data).unwrap();

        assert_eq!(TypeC2S::Invite, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
        assert_eq!(
            Some(vec![
                "Derek".to_string(),
                "Yahel".to_string(),
                "Erik".to_string(),
            ]),
            msg.usernames
        );
    }

    #[test]
    fn test_deserealize_join_room_message() {
        let data = r#"{"type":"JOIN_ROOM","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();

        assert_eq!(TypeC2S::JoinRoom, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_room_users_message() {
        let data = r#"{"type":"ROOM_USERS","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();

        assert_eq!(TypeC2S::RoomUsers, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_room_text_message() {
        let data = r#"{"type":"ROOM_TEXT","roomname":"Sala 1","text":"Pasa las pruebas"}"#;
        let msg = deserialize(data).unwrap();

        assert_eq!(TypeC2S::RoomText, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
        assert_eq!(Some("Pasa las pruebas".to_string()), msg.text);
    }

    #[test]
    fn test_deserealize_leave_room_message() {
        let data = r#"{"type":"LEAVE_ROOM","roomname":"Sala 1"}"#;
        let msg = deserialize(data).unwrap();

        assert_eq!(TypeC2S::LeaveRoom, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_disconnect_message() {
        let data = r#"{"type":"DISCONNECT"}"#;
        let msg = deserialize(data).unwrap();

        assert_eq!(TypeC2S::Disconnect, msg.type_c2s);
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
        let msg = ServerMessage::new_user_message("Evan".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "NEW_USER")));
        assert!(data.contains(&field("username", "Evan")));
    }

    #[test]
    fn test_serialize_new_status_message() {
        let msg = ServerMessage::new_status_message("Evan".to_string(), Busy);
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "NEW_STATUS")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("status", "BUSY")));
    }

    #[test]
    fn test_serialize_user_list_message() {
        let usernames = vec!["Derek".to_string(), "Yahel".to_string(), "Luis".to_string()];
        let msg = ServerMessage::user_list_message(Some(usernames));
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "USER_LIST")));
        assert!(data.contains(&field_array(
            "usernames",
            &["Derek".to_string(), "Yahel".to_string(), "Luis".to_string(),],
        )));
    }

    #[test]
    fn test_serialize_text_from_message() {
        let msg = ServerMessage::text_from_message("Evan".to_string(), "Hola amigos".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "TEXT_FROM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("text", "Hola amigos")));
    }

    #[test]
    fn test_serialize_public_text_from_message() {
        let msg =
            ServerMessage::public_text_from_message("Evan".to_string(), "Hola amigos".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "PUBLIC_TEXT_FROM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("text", "Hola amigos")));
    }

    #[test]
    fn test_serialize_joined_room_message() {
        let msg = ServerMessage::joined_room_message("Evan".to_string(), "Sala 1".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "JOINED_ROOM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("roomname", "Sala 1")));
    }

    #[test]
    fn test_serialize_room_user_list_message() {
        let usernames = vec!["Derek".to_string(), "Yahel".to_string(), "Luis".to_string()];
        let msg = ServerMessage::room_user_list_message("Sala 1".to_string(), usernames);
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "ROOM_USER_LIST")));
        assert!(data.contains(&field("roomname", "Sala 1")));
        assert!(data.contains(&field_array(
            "usernames",
            &["Derek".to_string(), "Yahel".to_string(), "Luis".to_string(),],
        )));
    }

    #[test]
    fn test_serialize_room_text_from_message() {
        let msg = ServerMessage::room_text_from_message(
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
        let msg = ServerMessage::left_room_message("Evan".to_string(), "Sala 1".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "LEFT_ROOM")));
        assert!(data.contains(&field("username", "Evan")));
        assert!(data.contains(&field("roomname", "Sala 1")));
    }

    #[test]
    fn test_serialize_disconnected_message() {
        let msg = ServerMessage::disconnected_message("Evan".to_string());
        let data = serialize(&msg).unwrap();
        assert!(data.contains(&field("type", "DISCONNECTED")));
        assert!(data.contains(&field("username", "Evan")));
    }
}
