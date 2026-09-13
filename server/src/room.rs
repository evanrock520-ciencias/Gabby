use std::collections::HashSet;
use uuid::Uuid;

use crate::protocol::result::MessageResult;

/// Representa un cuarto de chat.
pub struct Room {
    roomname: String,
    members: HashSet<Uuid>,
    invited: HashSet<Uuid>,
}

impl Room {
    fn new(roomname: String, creator_id: Uuid) -> Self {
        Room {
            roomname,
            members: HashSet::from([creator_id]),
            invited: HashSet::new(),
        }
    }

    /// Añade a un miembro al cuarto.
    ///
    /// # Arguments
    ///
    /// * `id` - El identificador del usuario a agregar.
    ///
    /// # Returns
    ///
    /// Devuelve si la operación fue exitosa o no.
    ///
    fn add_member(&mut self, id: &Uuid) -> Result<bool, MessageResult> {
        if !self.invited.contains(id) {
            return Err(MessageResult::NotInvited);
        }

        self.invited.remove(id);
        Ok(self.members.insert(*id))
    }

    /// Invita a un miembro al cuarto.
    ///
    /// # Arguments
    ///
    /// * `id` - El identificador del usuario a invitar.
    ///
    ///
    fn invitate(&mut self, id: &Uuid) {
        self.invited.insert(*id);
    }

    /// Elimina a un usuario del cuarto.
    ///
    /// # Arguments
    ///
    /// * `id` - El identificador del usuario a eliminar.
    ///
    ///
    fn remove_member(&mut self, id: &Uuid) {
        self.members.remove(id);
    }

    /// Verifica que un usuario sea miembro del cuarto.
    ///
    /// # Arguments
    ///
    /// * `id` - El identificador del usuario a verificar si es miembro.
    ///
    /// # Returns
    ///
    /// Devuelve si es miembro o no.
    ///
    fn is_member(&self, id: &Uuid) -> bool {
        self.members.contains(id)
    }

    /// Verifica que un usuario haya sido invitado a el cuarto.
    ///
    /// # Arguments
    ///
    /// * `id` - El identificador del usuario a verificar si ha sido invitado.
    ///
    /// # Returns
    ///
    /// Devuelve si fue invitado o no.
    ///
    fn is_invited(&self, id: &Uuid) -> bool {
        self.invited.contains(id)
    }
}

#[cfg(test)]
mod tests {
    use uuid::Uuid;

    use crate::{protocol::result::MessageResult, room::Room};

    #[test]
    fn test_create_new_room() {
        let id = Uuid::new_v4();
        let room = Room::new("Sala".to_string(), id);
        assert_eq!("Sala".to_string(), room.roomname);
        assert!(room.is_member(&id));
    }

    #[test]
    fn test_cant_join_if_not_invited() {
        let creator_id = Uuid::new_v4();
        let friend_id = Uuid::new_v4();

        let mut room = Room::new("Sala".to_string(), creator_id);
        let result = room.add_member(&friend_id);

        assert_eq!(result, Err(MessageResult::NotInvited));

        room.invitate(&friend_id);
        assert!(room.is_invited(&friend_id));

        let _ = room.add_member(&friend_id);
        assert!(room.is_member(&friend_id));
    }

    #[test]
    fn test_remove_user() {
        let creator_id = Uuid::new_v4();
        let friend_id = Uuid::new_v4();

        let mut room = Room::new("Sala".to_string(), creator_id);
        room.invitate(&friend_id);

        let _ = room.add_member(&friend_id);
        assert!(room.is_member(&friend_id));

        room.remove_member(&friend_id);
        assert!(!room.is_member(&friend_id));
    }
}
