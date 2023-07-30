use utoipa::ToSchema;

use super::IntoDto;

#[derive(serde::Serialize, ToSchema)]
pub struct HasCheckinResponse {
    has_checkin: bool,
}

impl From<bool> for HasCheckinResponse {
    fn from(value: bool) -> Self {
        Self { has_checkin: value }
    }
}

impl IntoDto for HasCheckinResponse {
    type Target = HasCheckinResponse;

    fn into_dto(self) -> Self::Target {
        self
    }
}

#[derive(serde::Serialize, ToSchema)]
pub struct CheckinResponse {
    success: bool,
}

impl From<bool> for CheckinResponse {
    fn from(value: bool) -> Self {
        Self { success: value }
    }
}

impl IntoDto for CheckinResponse {
    type Target = CheckinResponse;

    fn into_dto(self) -> Self::Target {
        self
    }
}


