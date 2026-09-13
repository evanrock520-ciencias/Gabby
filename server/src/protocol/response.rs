use crate::protocol::{incoming::TypeC2S, outcoming::TypeS2C, result::MessageResult};

/// Construye un mensaje de respuesta del servidor (S2C) a una operación de cliente.
pub fn new_response(operation: TypeC2S, result: MessageResult, extra: Option<String>) -> TypeS2C {
    TypeS2C::response_message(operation, result, extra)
}
