use tokio::sync::mpsc;
use uuid::Uuid;

use crate::protocol::{outcoming::TypeS2C, status::Status};

pub struct Client {
    id: Uuid,
    username: String,
    status: Status,
    tx: mpsc::UnboundedSender<TypeS2C>,
}

impl Client {
    pub fn new(
        id: Uuid,
        username: String,
        status: Status,
        tx: mpsc::UnboundedSender<TypeS2C>,
    ) -> Self {
        Client {
            id: id,
            username: username,
            status: status,
            tx: tx,
        }
    }

    pub fn id(&self) -> &Uuid {
        &self.id
    }

    pub fn username(&self) -> &String {
        &self.username
    }

    pub fn status(&self) -> &Status {
        &self.status
    }

    pub fn set_status(&mut self, status: Status) {
        self.status = status;
    }

    pub fn send(&self, msg: TypeS2C) {
        let _ = self.tx.send(msg);
    }
}
