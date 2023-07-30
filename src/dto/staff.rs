use utoipa::ToSchema;

use super::IntoDto;

#[derive(serde::Serialize, ToSchema)]
pub struct IsStaffResponse {
    is_staff: bool,
}

impl From<bool> for IsStaffResponse {
    fn from(value: bool) -> Self {
        Self { is_staff: value }
    }
}

impl IntoDto for IsStaffResponse {
    type Target = IsStaffResponse;

    fn into_dto(self) -> Self::Target {
        self
    }
}

#[derive(serde::Serialize, ToSchema)]
pub struct CheckingFreshyNightResponse {
    success: bool,
}

impl From<bool> for CheckingFreshyNightResponse {
    fn from(value: bool) -> Self {
        Self { success: value }
    }
}

impl IntoDto for CheckingFreshyNightResponse {
    type Target = CheckingFreshyNightResponse;

    fn into_dto(self) -> Self::Target {
        self
    }
}
