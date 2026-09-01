use serde::{Deserialize, Serialize};

#[derive(Serialize, Debug, PartialEq, Deserialize)]
#[serde(rename = "INVALID")]
pub enum Status {
    Active,
    Away,
    Busy,
}
