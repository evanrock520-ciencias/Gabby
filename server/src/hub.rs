use std::collections::HashMap;
use uuid::Uuid;

use crate::{client::Client, protocol::result::MessageResult, room::Room};

#[derive(Clone)]
pub struct Hub {
    clients: HashMap<Uuid, Client>,
    rooms: HashMap<String, Room>,
}

impl Hub {
    pub fn new() -> Self {
        Hub {
            clients: HashMap::new(),
            rooms: HashMap::new(),
        }
    }

    pub fn register(&mut self, client: Client) -> Result<bool, MessageResult> {
        if self.clients.contains_key(client.id()) {
            return Err(MessageResult::UserAlreadyExists);
        }

        Ok(self.clients.insert(*client.id(), client).is_none())
    }

    pub fn unregister(&mut self, id: Uuid) {
        self.clients.remove(&id);
    }

    pub fn usernames(&self) -> Vec<String> {
        self.clients
            .values()
            .map(|client| client.username().to_string())
            .collect()
    }
}

#[cfg(test)]
mod test {}
