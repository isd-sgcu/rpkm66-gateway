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

macro_rules! direct_map {
    ($a:ty, $b:ty, $($field:ident),*) => {
        const _: () = {
            impl From<$a> for $b {
                fn from(value: $a) -> Self {
                    Self {
                        $(
                            $field: value.$field,
                        )*
                    }
                }
            }

            impl From<$b> for $a {
                fn from(value: $b) -> Self {
                    Self {
                        $(
                            $field: value.$field,
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
mod group;
mod staff;
mod user;

pub use auth::*;
pub use baan::*;
pub use file::*;
pub use group::*;
pub use staff::*;
pub use user::*;
