use std::collections::HashSet;

use tokio::sync::mpsc;

use crate::protocol::{outcoming::TypeS2C, status::Status};

#[derive(Clone)]
pub struct Client {
    memberships: HashSet<String>,
    username: String,
    status: Status,
    tx: mpsc::UnboundedSender<TypeS2C>,
}

impl Client {
    pub fn new(username: String, tx: mpsc::UnboundedSender<TypeS2C>) -> Self {
        Client {
            memberships: HashSet::new(),
            username,
            status: Status::Active,
            tx,
        }
    }

    /// Agrega una sala a las membresias del Cliente.
    ///
    /// # Arguments
    ///
    /// * `roomname` - El nombre de la sala.
    pub fn add_membership(&mut self, roomname: &str) {
        self.memberships.insert(roomname.to_string());
    }

    /// Agrega una sala a las membresias del Cliente.
    ///
    /// # Arguments
    ///
    /// * `roomname` - El nombre de la sala.
    ///
    /// # Returns
    ///
    /// Retorna `true` si es miembro de la sala.
    /// Retorna `false` si no es miembro de la sala.
    pub fn has_membership(&self, roomname: &str) -> bool {
        self.memberships.contains(&roomname.to_string())
    }

    pub fn username(&self) -> &str {
        &self.username
    }

    pub fn status(&self) -> &Status {
        &self.status
    }

    pub fn set_status(&mut self, status: Status) -> bool {
        if self.status != status {
            self.status = status;
            true
        } else {
            false
        }
    }

    pub fn send(&self, msg: &TypeS2C) {
        let _ = self.tx.send(msg.clone());
    }
}
