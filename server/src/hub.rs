use std::collections::HashMap;

use crate::{
    client::Client,
    protocol::{outcoming::TypeS2C, result::MessageResult, status::Status},
    room::Room,
};

#[derive(Clone)]
pub struct Hub {
    clients: HashMap<String, Client>,
    rooms: HashMap<String, Room>,
}

impl Hub {
    pub fn new() -> Self {
        Hub {
            clients: HashMap::new(),
            rooms: HashMap::new(),
        }
    }

    /// Retorna si el cliente pertenece al hub.
    ///
    /// # Arguments
    ///
    /// * `username` - El nombre de usuario del cliente.
    ///
    /// # Returns
    ///
    /// `true` si es un usuario del hub.
    /// `false` si no es usuario del hub.
    pub fn is_user(&self, username: &str) -> bool {
        return self.clients.contains_key(username);
    }

    /// Retorna la cantidad de clientes en el hub.
    ///
    /// # Returns
    ///
    /// La cantidad de clientes en el hub.
    pub fn connected_users(&self) -> usize {
        return self.clients.len();
    }

    /// Retorna si la sala pertenece al hub.
    ///
    /// # Arguments
    ///
    /// * `roomname` - El nombre de la sala.
    ///
    /// # Returns
    ///
    /// `true` si es una sala del hub.
    /// `false` si no es sala del hub.
    pub fn is_room(&self, roomname: &str) -> bool {
        return self.rooms.contains_key(roomname);
    }

    /// Retorna el usuario pertenece a una sala.
    ///
    /// # Arguments
    ///
    /// * `username` - El nombre de usuario del cliente.
    /// * `roomname` - El nombre de la sala.
    ///
    /// # Returns
    ///
    /// `true` si es una sala del hub.
    /// `false` si no es sala del hub.
    pub fn is_member_of(&self, username: &str, roomname: &str) -> bool {
        self.rooms
            .get(roomname)
            .is_some_and(|room| room.is_member(username))
    }

    /// Propaga un mensaje a todos los clientes.
    ///
    /// # Arguments
    ///
    /// * `msg` - El mensaje a propagar.
    /// * `except` - El cliente excluido de la propagación.
    ///
    pub fn broadcast(&self, msg: &TypeS2C, except: &str) {
        for (id, client) in &self.clients {
            if id != except {
                client.send(msg);
            }
        }
    }

    /// Propaga un mensaje a todos los clientes en una sala.
    ///
    /// # Arguments
    ///
    /// * `msg` - El mensaje a propagar.
    /// * `except` - El cliente excluido de la propagación.
    ///
    pub fn to_room(
        &self,
        msg: &TypeS2C,
        roomname: &str,
        except: &str,
    ) -> Result<bool, MessageResult> {
        let Some(room) = self.rooms.get(roomname) else {
            return Err(MessageResult::NoSuchRoom);
        };

        for user in room.members() {
            if let Some(client) = self.clients.get(user) {
                if client.username() == except {
                    continue;
                }
                client.send(msg);
            }
        }

        Ok(true)
    }

    /// Envía un mensaje a otro cliente.
    ///
    /// # Arguments
    ///
    /// * `msg` - El mensaje a propagar.
    /// * `sender` - El usuario que manda el mensaje.
    /// * `receiver` - El usuario que recibe el mensaje.
    ///
    pub fn send_to(
        &self,
        msg: &TypeS2C,
        sender: &str,
        receiver: &str,
    ) -> Result<bool, MessageResult> {
        if !self.is_user(receiver) || !self.is_user(sender) {
            return Err(MessageResult::NoSuchUser);
        }

        let receiver = self.clients.get(receiver).unwrap();
        receiver.send(msg);

        Ok(true)
    }

    /// Registra a un usuario al Hub.
    ///
    /// # Arguments
    ///
    /// * `client` - El cliente a agregar al hub.
    ///
    /// # Returns
    ///
    /// * `Ok(true)` si el cliente fue registrado exitosamente.
    /// * `Ok(false)` si el usuario reemplazó a uno previo con la misma clave.
    ///
    /// # Errors
    ///
    /// Retorna `Err(MessageResult::UserAlreadyExists)` si el usuario ya existe en el hub.
    pub fn register(&mut self, client: Client) -> Result<bool, MessageResult> {
        if self.clients.contains_key(client.username()) {
            return Err(MessageResult::UserAlreadyExists);
        }

        Ok(self
            .clients
            .insert(client.username().to_string(), client)
            .is_none())
    }

    /// Elimina a un usuario del Hub.
    ///
    /// # Arguments
    ///
    /// * `username` - El nombre de usuario del cliente a eliminar.
    ///
    pub fn unregister(&mut self, username: &str) {
        self.clients.remove(username);
    }

    /// Registra una sala y agrega a su creador como miembro.
    ///
    /// # Arguments
    ///
    /// * `room` - La sala a registrar.
    /// * `username` - El nombre de usuario del creador de la sala.
    ///
    /// # Returns
    ///
    /// Retorna `Ok(true)` si la sala fue registrada exitosamente.
    ///
    /// # Errors
    ///
    /// Retorna [`MessageResult::RoomAlreadyExists`] si ya existe una sala registrada
    /// con el mismo nombre.
    ///
    pub fn register_room(&mut self, roomname: &str, username: &str) -> Result<bool, MessageResult> {
        if self.rooms.contains_key(roomname) {
            return Err(MessageResult::RoomAlreadyExists);
        }

        if !self.is_user(username) {
            return Err(MessageResult::NoSuchUser);
        }

        let room = Room::new(roomname.to_string(), username.to_string());
        let room = self.rooms.insert(room.roomname().into(), room);

        let client = self.clients.get_mut(username).unwrap();
        client.add_membership(roomname);

        Ok(room.is_none())
    }

    /// Agrega a una lista de clientes a la lista de invitados de la sala.
    ///
    /// # Arguments
    ///
    /// * `room` - La sala a la que invitaron a los usuarios.
    /// * `usernames` - La lista de usuarios invitados.
    ///
    /// # Returns
    ///
    /// Retorna `Ok(true)` si se invitaron a todos los usuarios a la sala.
    ///
    /// # Errors
    ///
    /// Retorna [`MessageResult::NoSuchRoom`] si no existe la sala.
    /// Retorna [`MessageResult::NoSuchUser`] si al menos uno de los usuarios no existe.
    ///
    pub fn invite(
        &mut self,
        roomname: &str,
        usernames: Vec<String>,
    ) -> Result<bool, MessageResult> {
        for user in &usernames {
            if !self.is_user(user) {
                return Err(MessageResult::NoSuchUser);
            }
        }

        let Some(room) = self.rooms.get_mut(roomname) else {
            return Err(MessageResult::NoSuchRoom);
        };

        for user in usernames {
            room.invitate(user.as_str());
        }
        return Ok(true);
    }

    /// Agrega a una lista de clientes a la lista de invitados de la sala.
    ///
    /// # Arguments
    ///
    /// * `room` - La sala de invitación.
    /// * `usernames` - El usuario que aceptó la invitación.
    ///
    /// # Returns
    ///
    /// Retorna `Ok(true)` si fue posible unirse a la sala.
    ///
    /// # Errors
    ///
    /// Retorna [`MessageResult::NoSuchRoom`] si no existe la sala.
    /// Retorna [`MessageResult::NoSuchUser`] si el usuario no pertenece al hub.
    /// Retorna [`MessageResult::NotInvited`] si el usuario no fue invitado a la sala.
    ///
    pub fn be_member_of(&mut self, roomname: &str, username: &str) -> Result<bool, MessageResult> {
        if !self.is_user(username) {
            return Err(MessageResult::NoSuchUser);
        }

        let Some(room) = self.rooms.get_mut(roomname) else {
            return Err(MessageResult::NoSuchRoom);
        };

        if !room.is_invited(username) {
            return Err(MessageResult::NotInvited);
        }

        room.add_member(username).unwrap();
        let client = self.clients.get_mut(username).unwrap();
        client.add_membership(roomname);

        Ok(true)
    }

    /// Retorna los usernames y status de todos los clientes del Hub.
    ///
    /// # Returns
    ///
    /// Un `HashMap` donde la clave es el nombre de usuario (`String`)
    /// y el valor es su estado actual (`Status`).
    pub fn usernames(&self) -> HashMap<String, Status> {
        self.clients
            .iter()
            .map(|(username, client)| (username.clone(), *client.status()))
            .collect()
    }

    /// Cambia el status de un cliente en el Hub.
    ///
    /// # Arguments
    ///
    /// * `client` - El cliente al cual cambiar el status.
    /// * `new_status` - El nuevo status del cliente.
    ///
    /// # Returns
    ///
    /// * `true` si el status era diferente.
    /// * `false` si el status era igual.
    pub fn change_status(&mut self, username: &str, new_status: Status) -> bool {
        if let Some(client) = self.clients.get_mut(username) {
            client.set_status(new_status)
        } else {
            eprintln!("The user {} doesn't exists.", username);
            false
        }
    }

    /// Otorga una lista de los clientes en una sala.
    ///
    /// # Arguments
    ///
    /// * `roomname` - La sala de la que se solicita la lista.
    /// * `username` - El usuario que solicita la lista.
    ///
    /// # Returns
    ///
    /// Un `HashMap` donde la clave es el nombre de usuario (`String`)
    /// y el valor es su estado actual (`Status`).
    ///
    /// # Errors
    ///
    /// Retorna [`MessageResult::NoSuchRoom`] si no existe la sala.
    /// Retorna [`MessageResult::NotJoined`] si el usuario no pertenece a la sala.
    pub fn room_usernames(
        &self,
        roomname: &str,
        username: &str,
    ) -> Result<HashMap<String, Status>, MessageResult> {
        let Some(room) = self.rooms.get(roomname) else {
            return Err(MessageResult::NoSuchRoom);
        };

        if !room.is_member(username) {
            return Err(MessageResult::NotJoined);
        }

        Ok(room
            .members()
            .iter()
            .filter_map(|u| self.clients.get(u).map(|c| (u.to_string(), *c.status())))
            .collect())
    }

    /// Elimina a un cliente de una sala.
    ///
    /// # Arguments
    ///
    /// * `roomname` - La sala a abandonar.
    /// * `username` - El usuario que solicita abandonar la sala.
    ///
    /// # Returns
    ///
    /// Retorna `Ok(true)` si fue posible abandonar a la sala.
    ///
    /// # Errors
    ///
    /// Retorna [`MessageResult::NoSuchRoom`] si no existe la sala.
    /// Retorna [`MessageResult::NotJoined`] si el usuario no pertenece a la sala.
    pub fn leave_room(&mut self, roomname: &str, username: &str) -> Result<bool, MessageResult> {
        let Some(room) = self.rooms.get_mut(roomname) else {
            return Err(MessageResult::NoSuchRoom);
        };

        if !room.is_member(username) {
            return Err(MessageResult::NotJoined);
        }

        room.remove_member(username);

        if room.is_empty() {
            self.rooms.remove(roomname);
        }

        Ok(true)
    }
}

#[cfg(test)]
mod test {
    use tokio::sync::mpsc;

    use crate::{
        client::Client,
        hub::Hub,
        protocol::{outcoming::TypeS2C, result::MessageResult, status::Status},
    };

    #[test]
    fn test_register() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        let was_registered = hub.register(alice).unwrap();
        assert!(was_registered)
    }

    #[test]
    fn test_user_already_exists() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        let was_registered = hub.register(alice).unwrap();
        assert!(was_registered);

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        let error = hub.register(alice).unwrap_err();
        assert_eq!(MessageResult::UserAlreadyExists, error);
    }

    #[test]
    fn test_unregister() {
        let mut hub = Hub::new();

        assert!(hub.usernames().is_empty());

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        let _ = hub.register(alice).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        let _ = hub.register(bob);

        assert!(hub.is_user("bob".into()));
        assert!(hub.is_user("alice".into()));
        assert_eq!(hub.connected_users(), 2);

        hub.unregister("alice".into());
        assert!(!hub.is_user("alice".into()));
        assert_eq!(hub.connected_users(), 1);

        hub.unregister("bob".into());
        assert!(!hub.is_user("bob".into()));
        assert_eq!(hub.connected_users(), 0);
    }

    #[test]
    fn test_change_status() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);
        let _ = hub.register(alice);

        // Cambia porque por defecto status = Status::Active
        let does_change = hub.change_status("alice".into(), Status::Busy);
        assert!(does_change);

        let does_change = hub.change_status("alice".into(), Status::Busy);
        assert!(!does_change)
    }

    #[test]
    fn test_change_status_with_non_registered_client() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let _ = Client::new("alice".to_string(), tx_alice);

        assert!(!hub.change_status("alice".into(), Status::Away));
    }

    #[test]
    fn test_usernames() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        let _ = hub.register(alice).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        let _ = hub.register(bob);

        let (tx_charlie, _) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        let _ = hub.register(charlie);

        let usernames = hub.usernames();
        assert!(usernames.contains_key("alice"));
        assert!(usernames.contains_key("bob"));
        assert!(usernames.contains_key("charlie"));
    }

    #[test]
    fn test_broadcast() {
        let mut hub = Hub::new();

        let (tx_alice, mut rx_alice) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        let _ = hub.register(alice).unwrap();

        let (tx_bob, mut rx_bob) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        let _ = hub.register(bob);

        let msg = &TypeS2C::NewUser {
            username: "charlie".into(),
        };

        hub.broadcast(&msg.clone(), "alice".into());

        assert_eq!(rx_bob.try_recv().unwrap(), msg.clone());
        assert!(rx_alice.try_recv().is_err());
    }

    #[test]
    fn test_register_room() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        assert!(hub.is_room("Room 1".into()));
        assert!(hub.is_member_of("alice", "Room 1"));
    }

    #[test]
    fn test_register_room_with_non_registered_user() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let _ = Client::new("alice".to_string(), tx_alice);

        let error = hub.register_room("Room 1", "alice".into()).unwrap_err();
        assert_eq!(MessageResult::NoSuchUser, error);
    }

    #[test]
    fn test_invite_to_rooms() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        let (tx_charlie, _) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        hub.register(charlie).unwrap();

        let guests = vec!["bob".to_string(), "charlie".to_string()];

        assert!(hub.invite("Room 1", guests).unwrap());
    }

    #[test]
    fn test_join_to_room() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        let (tx_charlie, _) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        hub.register(charlie).unwrap();

        let guests = vec!["bob".to_string(), "charlie".to_string()];

        assert!(hub.invite("Room 1", guests.clone()).unwrap());
        for client in guests {
            assert!(hub.be_member_of("Room 1", client.as_str()).unwrap());
        }
    }

    #[test]
    fn test_cant_join_if_not_invited() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        let (tx_charlie, _) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        hub.register(charlie).unwrap();

        let guests = vec!["bob", "charlie"];

        for client in guests {
            assert_eq!(
                MessageResult::NotInvited,
                hub.be_member_of("Room 1", client).unwrap_err()
            );
        }
    }

    #[test]
    fn test_send_to_room() {
        let mut hub = Hub::new();

        let (tx_alice, mut rx_alice) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, mut rx_bob) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        let (tx_charlie, mut rx_charlie) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        hub.register(charlie).unwrap();

        let guests = vec!["bob".to_string()];

        hub.invite("Room 1", guests.clone()).unwrap();
        for client in guests {
            hub.be_member_of("Room 1", client.as_str()).unwrap();
        }

        let msg = TypeS2C::NewStatus {
            username: "alice".into(),
            status: Status::Away,
        };

        hub.to_room(&msg, "Room 1", "alice").unwrap();

        assert_eq!(rx_bob.try_recv().unwrap(), msg.clone());
        assert!(rx_alice.try_recv().is_err());
        assert!(rx_charlie.try_recv().is_err());
    }

    #[test]
    fn test_send_to_client() {
        let mut hub = Hub::new();

        let (tx_alice, mut rx_alice) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, mut rx_bob) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        let (tx_charlie, mut rx_charlie) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        hub.register(charlie).unwrap();

        let msg = TypeS2C::TextFrom {
            username: "alice".to_string(),
            text: "Hello".to_string(),
        };
        hub.send_to(&msg, "alice", "bob").unwrap();

        assert_eq!(rx_bob.try_recv().unwrap(), msg.clone());
        assert!(rx_alice.try_recv().is_err());
        assert!(rx_charlie.try_recv().is_err());

        let msg = TypeS2C::TextFrom {
            username: "bob".to_string(),
            text: "charlie".to_string(),
        };

        hub.send_to(&msg, "bob", "charlie").unwrap();

        assert_eq!(rx_charlie.try_recv().unwrap(), msg.clone());
        assert!(rx_bob.try_recv().is_err());
        assert!(rx_alice.try_recv().is_err());
    }

    #[test]
    fn test_room_userlist() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        let (tx_charlie, _) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        hub.register(charlie).unwrap();

        let (tx_diane, _) = mpsc::unbounded_channel();
        let diane = Client::new("diane".to_string(), tx_diane);
        hub.register(diane).unwrap();

        let guests = vec!["bob".to_string(), "charlie".to_string()];

        hub.invite("Room 1", guests.clone()).unwrap();
        for client in guests.clone() {
            hub.be_member_of("Room 1", client.as_str()).unwrap();
        }

        let room_usernames = hub.room_usernames("Room 1", "alice").unwrap();

        for client in guests {
            assert!(room_usernames.contains_key(client.as_str()));
        }
    }

    #[test]
    fn test_room_userlist_on_non_existent_room() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        assert_eq!(
            MessageResult::NoSuchRoom,
            hub.room_usernames("Room 2", "alice").unwrap_err()
        );
    }

    #[test]
    fn test_room_userlist_with_non_joined_client() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        assert_eq!(
            MessageResult::NotJoined,
            hub.room_usernames("Room 1", "bob").unwrap_err()
        );
    }

    #[test]
    fn test_leave_room() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        let (tx_bob, _) = mpsc::unbounded_channel();
        let bob = Client::new("bob".to_string(), tx_bob);
        hub.register(bob).unwrap();

        let (tx_charlie, _) = mpsc::unbounded_channel();
        let charlie = Client::new("charlie".to_string(), tx_charlie);
        hub.register(charlie).unwrap();

        let guests = vec!["bob".to_string(), "charlie".to_string()];

        assert!(hub.invite("Room 1", guests.clone()).unwrap());
        for client in guests {
            assert!(hub.be_member_of("Room 1", client.as_str()).unwrap());
        }

        assert!(hub.leave_room("Room 1", "bob").unwrap());
        assert!(!hub.is_member_of("Room 1", "bob"));
    }

    #[test]
    fn test_empty_room_is_removed() {
        let mut hub = Hub::new();

        let (tx_alice, _) = mpsc::unbounded_channel();
        let alice = Client::new("alice".to_string(), tx_alice);

        hub.register(alice).unwrap();
        hub.register_room("Room 1", "alice".into()).unwrap();

        assert!(hub.leave_room("Room 1", "alice").unwrap());
        assert!(!hub.is_room("Room 1"));
    }
}
