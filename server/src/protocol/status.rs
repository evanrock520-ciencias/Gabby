use serde::Serialize;

#[derive(Serialize, Debug, PartialEq)]
pub enum Status {
    #[serde(rename = "ACTIVE")]
    Active,
    #[serde(rename = "AWAY")]
    Away,
    #[serde(rename = "BUSY")]
    Busy,
}
