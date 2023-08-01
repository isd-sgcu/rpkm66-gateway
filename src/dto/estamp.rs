use utoipa::ToSchema;

use super::IntoDto;

#[derive(serde::Serialize, ToSchema)]
pub struct GetAllEstampResponse {
    events: Vec<EstampEvent>,
}
#[derive(serde::Serialize, ToSchema)]

pub struct GetUserEstampsResponse {
    events: Vec<UserEstampEvent>,
}

#[derive(serde::Serialize, ToSchema)]
pub struct EstampEvent {
    id: String,
    name: String,
    additional_info: String,
}

#[derive(serde::Serialize, ToSchema)]
pub struct UserEstampEvent {
    event: EstampEvent,
    is_taken: bool,
    taken_at: i64,
}

#[derive(serde::Serialize, ToSchema)]
pub struct RedeemItemResponse {
    success: bool,
}

#[derive(serde::Serialize, ToSchema)]
pub struct HasRedeemItemResponse {
    redeemed: bool,
}

impl From<rpkm66_rust_proto::rpkm66::checkin::event::v1::Event> for EstampEvent {
    fn from(value: rpkm66_rust_proto::rpkm66::checkin::event::v1::Event) -> Self {
        Self {
            id: value.event_id,
            name: value.event_name,
            additional_info: value.additional_info,
        }
    }
}

impl From<rpkm66_rust_proto::rpkm66::checkin::event::v1::UserEvent> for UserEstampEvent {
    fn from(value: rpkm66_rust_proto::rpkm66::checkin::event::v1::UserEvent) -> Self {
        Self {
            event: value.event.unwrap_or_default().into(),
            is_taken: value.is_taken,
            taken_at: value.taken_at,
        }
    }
}

impl From<bool> for RedeemItemResponse {
    fn from(value: bool) -> Self {
        Self { success: value }
    }
}

impl IntoDto for RedeemItemResponse {
    type Target = Self;

    fn into_dto(self) -> Self::Target {
        self
    }
}

impl From<bool> for HasRedeemItemResponse {
    fn from(value: bool) -> Self {
        Self { redeemed: value }
    }
}

impl IntoDto for HasRedeemItemResponse {
    type Target = Self;

    fn into_dto(self) -> Self::Target {
        self
    }
}

impl IntoDto for UserEstampEvent {
    type Target = Self;

    fn into_dto(self) -> Self::Target {
        self
    }
}

impl IntoDto for EstampEvent {
    type Target = Self;

    fn into_dto(self) -> Self::Target {
        self
    }
}

impl IntoDto for GetAllEstampResponse {
    type Target = GetAllEstampResponse;

    fn into_dto(self) -> Self::Target {
        self
    }
}

impl IntoDto for GetUserEstampsResponse {
    type Target = Self;

    fn into_dto(self) -> Self::Target {
        self
    }
}

impl<T> From<T> for GetAllEstampResponse
where
    T: Iterator<Item = EstampEvent>,
{
    fn from(value: T) -> Self {
        GetAllEstampResponse {
            events: value.collect(),
        }
    }
}

impl<T> From<T> for GetUserEstampsResponse
where
    T: Iterator<Item = UserEstampEvent>,
{
    fn from(value: T) -> Self {
        GetUserEstampsResponse {
            events: value.collect(),
        }
    }
}
