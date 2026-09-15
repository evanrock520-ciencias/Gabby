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
mod test {}
