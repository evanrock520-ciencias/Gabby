use serde::{Deserialize, Serialize};

#[derive(Serialize, Debug, PartialEq, Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum Status {
    Active,
    Away,
    Busy,
}
