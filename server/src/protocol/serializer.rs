use crate::protocol::{incoming::ClientMessage, outcoming::ServerMessage, status::Status};
use serde_json;

pub fn serialize(msg: &ServerMessage) -> Result<String, serde_json::Error> {
    serde_json::to_string(msg)
}

pub fn deserealize(data: &str) -> Result<ClientMessage, serde_json::Error> {
    serde_json::from_str::<ClientMessage>(data)
}

#[cfg(test)]
mod tests {
    use crate::protocol::{
        incoming::TypeC2S::{self, Status},
        status::Status::Away,
    };

    use super::*;

    // Para comprobar que las piezas encajen, las cadenas seran serializadas a partir de los ClientMessages del cliente.
    #[test]
    fn test_deserealize_identify_message() {
        let data = r#"{"type":"IDENTIFY","username":"Evan"}"#;
        let msg = deserealize(data).unwrap();
        assert_eq!(TypeC2S::Identify, msg.type_c2s);
        assert_eq!(Some("Evan".to_string()), msg.username);
    }

    #[test]
    fn test_deserealize_status_message() {
        let data = r#"{"type":"STATUS","status":"AWAY"}"#;
        let msg = deserealize(data).unwrap();
        assert_eq!(TypeC2S::Status, msg.type_c2s);
        assert_eq!(Some(Away), msg.status);
    }

    #[test]
    fn test_deserealize_users_message() {
        let data = r#"{"type":"USERS"}"#;
        let msg = deserealize(data).unwrap();
        assert_eq!(TypeC2S::Users, msg.type_c2s);
    }

    #[test]
    fn test_deserealize_text_message() {
        let data = r#"{"type":"TEXT","username":"Evan","text":"Pasa las pruebas"}"#;
        let msg = deserealize(data).unwrap();
        assert_eq!(TypeC2S::Text, msg.type_c2s);
        assert_eq!(Some("Evan".to_string()), msg.username);
        assert_eq!(Some("Pasa las pruebas".to_string()), msg.text);
    }

    #[test]
    fn test_deserealize_public_text_message() {
        let data = r#"{"type":"PUBLIC_TEXT","text":"Pasa las pruebas"}"#;
        let msg = deserealize(data).unwrap();
        assert_eq!(TypeC2S::PublicText, msg.type_c2s);
        assert_eq!(Some("Pasa las pruebas".to_string()), msg.text);
    }

    #[test]
    fn test_deserealize_new_room_message() {
        let data = r#"{"type":"NEW_ROOM","roomname":"Sala 1"}"#;
        let msg = deserealize(data).unwrap();
        assert_eq!(TypeC2S::NewRoom, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_invite_message() {
        let data = r#"{"type":"INVITE","roomname":"Sala 1","usernames":["Derek","Yahel","Erik"]}"#;
        let msg = deserealize(data).unwrap();

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
        let msg = deserealize(data).unwrap();

        assert_eq!(TypeC2S::JoinRoom, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_room_users_message() {
        let data = r#"{"type":"ROOM_USERS","roomname":"Sala 1"}"#;
        let msg = deserealize(data).unwrap();

        assert_eq!(TypeC2S::RoomUsers, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_room_text_message() {
        let data = r#"{"type":"ROOM_TEXT","roomname":"Sala 1","text":"Pasa las pruebas"}"#;
        let msg = deserealize(data).unwrap();

        assert_eq!(TypeC2S::RoomText, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
        assert_eq!(Some("Pasa las pruebas".to_string()), msg.text);
    }

    #[test]
    fn test_deserealize_leave_room_message() {
        let data = r#"{"type":"LEAVE_ROOM","roomname":"Sala 1"}"#;
        let msg = deserealize(data).unwrap();

        assert_eq!(TypeC2S::LeaveRoom, msg.type_c2s);
        assert_eq!(Some("Sala 1".to_string()), msg.room);
    }

    #[test]
    fn test_deserealize_disconnect_message() {
        let data = r#"{"type":"DISCONNECT"}"#;
        let msg = deserealize(data).unwrap();

        assert_eq!(TypeC2S::Disconnect, msg.type_c2s);
    }
}
