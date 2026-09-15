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

    pub fn is_user(&self, username: &str) -> bool {
        return self.clients.contains_key(username);
    }

    pub fn connected_users(&self) -> usize {
        return self.clients.len();
    }

    pub fn is_room(&self, roomname: &str) -> bool {
        return self.rooms.contains_key(roomname);
    }

    pub fn broadcast(&self, msg: &TypeS2C, except: &str) {
        for (id, client) in &self.clients {
            if id != except {
                client.send(msg);
            }
        }
    }

    pub fn register(&mut self, client: Client) -> Result<bool, MessageResult> {
        if self.clients.contains_key(client.username()) {
            return Err(MessageResult::UserAlreadyExists);
        }

        Ok(self
            .clients
            .insert(client.username().to_string(), client)
            .is_none())
    }

    pub fn unregister(&mut self, username: &str) {
        self.clients.remove(username);
    }

    pub fn usernames(&self) -> HashMap<String, Status> {
        self.clients
            .iter()
            .map(|(username, client)| (username.clone(), *client.status()))
            .collect()
    }

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
