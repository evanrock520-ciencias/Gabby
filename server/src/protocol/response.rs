use crate::protocol::result::Result;
use serde::Serialize;

use crate::protocol::{incoming::TypeC2S, outcoming::TypeS2C};

pub struct Response {
    type_s2c: TypeS2C,
    type_c2s: TypeC2S,
    result: Result,
    extra: Option<String>,
}

// TODO: Implementar la serialización de respuestas.
impl Response {}
