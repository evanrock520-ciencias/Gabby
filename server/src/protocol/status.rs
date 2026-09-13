use serde::{Deserialize, Serialize};

#[derive(Serialize, Debug, PartialEq, Deserialize, Clone, Copy)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum Status {
    Active,
    Away,
    Busy,
}
