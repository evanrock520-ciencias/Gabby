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
    /// * `username` - El nombre de la sala.
    ///
    /// # Returns
    ///
    /// `true` si es una sala del hub.
    /// `false` si no es sala del hub.
    pub fn is_room(&self, roomname: &str) -> bool {
        return self.rooms.contains_key(roomname);
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
}
