use axum::{response::IntoResponse, Json};
use serde::Serialize;

macro_rules! into_dto {
    ($proto:ty, $dto:ty, $($field:ident),*) => {
        const _: () = {
            use crate::dto::IntoDto;

            impl IntoDto for $proto {
                type Target = $dto;

                fn into_dto(self) -> Self::Target {
                    Self::Target {
                        $(
                            $field: self.$field,
                        )*
                    }
                }
            }
        };
    };
}

pub trait IntoDto: Sized {
    type Target: Serialize;

    fn into_dto(self) -> Self::Target;

    fn into_response(self) -> axum::response::Response {
        Json(self.into_dto()).into_response()
    }
}

mod auth;
mod baan;
mod file;
mod user;

pub use auth::*;
pub use baan::*;
pub use file::*;
pub use user::*;
