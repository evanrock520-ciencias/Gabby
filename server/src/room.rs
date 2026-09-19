use std::collections::HashSet;

use crate::protocol::result::MessageResult;

/// Representa un cuarto de chat.
#[derive(Clone)]
pub struct Room {
    roomname: String,
    members: HashSet<String>,
    invited: HashSet<String>,
}

impl Room {
    pub fn new(roomname: String, creator_username: String) -> Self {
        Room {
            roomname,
            members: HashSet::from([creator_username]),
            invited: HashSet::new(),
        }
    }

    /// Añade a un miembro al cuarto.
    ///
    /// # Arguments
    ///
    /// * `username` - El identificador del usuario a agregar.
    ///
    /// # Returns
    ///
    /// Devuelve si la operación fue exitosa o no.
    ///
    pub fn add_member(&mut self, username: &str) -> Result<bool, MessageResult> {
        if !self.invited.contains(username) {
            return Err(MessageResult::NotInvited);
        }

        self.invited.remove(username);
        Ok(self.members.insert(username.to_string()))
    }

    /// Invita a un miembro al cuarto.
    ///
    /// # Arguments
    ///
    /// * `username` - El identificador del usuario a invitar.
    ///
    ///
    pub fn invitate(&mut self, username: &str) {
        self.invited.insert(username.to_string());
    }

    /// Elimina a un usuario del cuarto.
    ///
    /// # Arguments
    ///
    /// * `username` - El identificador del usuario a eliminar.
    ///
    ///
    pub fn remove_member(&mut self, username: &str) {
        self.members.remove(username);
    }

    /// Verifica que un usuario sea miembro del cuarto.
    ///
    /// # Arguments
    ///
    /// * `username` - El identificador del usuario a verificar si es miembro.
    ///
    /// # Returns
    ///
    /// Devuelve si es miembro o no.
    ///
    pub fn is_member(&self, username: &str) -> bool {
        self.members.contains(username)
    }

    /// Verifica que un usuario haya sido invitado a el cuarto.
    ///
    /// # Arguments
    ///
    /// * `username` - El identificador del usuario a verificar si ha sido invitado.
    ///
    /// # Returns
    ///
    /// Devuelve si fue invitado o no.
    ///
    pub fn is_invited(&self, username: &str) -> bool {
        self.invited.contains(username)
    }

    pub fn roomname(&self) -> &str {
        &self.roomname
    }

    pub fn members(&self) -> &HashSet<String> {
        &self.members
    }

    pub fn is_empty(&self) -> bool {
        self.members().is_empty()
    }
}

#[cfg(test)]
mod tests {
    use crate::{protocol::result::MessageResult, room::Room};

    #[test]
    fn test_create_new_room() {
        let creator = "alice".to_string();
        let room = Room::new("Sala".to_string(), creator.clone());
        assert_eq!("Sala", room.roomname);
        assert!(room.is_member(&creator));
    }

    #[test]
    fn test_cant_join_if_not_invited() {
        let creator = "alice".to_string();
        let friend = "bob".to_string();

        let mut room = Room::new("Sala".to_string(), creator);
        let result = room.add_member(&friend);

        assert_eq!(result, Err(MessageResult::NotInvited));

        room.invitate(&friend);
        assert!(room.is_invited(&friend));

        let _ = room.add_member(&friend);
        assert!(room.is_member(&friend));
    }

    #[test]
    fn test_remove_user() {
        let creator = "alice".to_string();
        let friend = "bob".to_string();

        let mut room = Room::new("Sala".to_string(), creator);
        room.invitate(&friend);

        let _ = room.add_member(&friend);
        assert!(room.is_member(&friend));

        room.remove_member(&friend);
        assert!(!room.is_member(&friend));
    }
}
