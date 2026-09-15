use tokio::sync::mpsc;

use crate::protocol::{outcoming::TypeS2C, status::Status};

#[derive(Clone)]
pub struct Client {
    username: String,
    status: Status,
    tx: mpsc::UnboundedSender<TypeS2C>,
}

impl Client {
    pub fn new(username: String, tx: mpsc::UnboundedSender<TypeS2C>) -> Self {
        Client {
            username,
            status: Status::Active,
            tx,
        }
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
