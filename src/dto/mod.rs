macro_rules! direct_map {
    ($dto:ty, $proto:ty; $($field:ident),*) => {
        impl From<$proto> for $dto {
            fn from(value: $proto) -> Self {
                Self {
                    $($field: value.$field),*
                }
            }
        }
        impl From<$dto> for $proto {
            fn from(value: $dto) -> Self {
                Self {
                    $($field: value.$field),*
                }
            }
        }
    };
}

mod auth;
mod user;

pub use auth::*;
pub use user::*;
