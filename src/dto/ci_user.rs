use utoipa::ToSchema;

use super::IntoDto;

#[derive(serde::Serialize, ToSchema)]
pub struct IsFreshyNightTicketRedeemedResponse {
    redeemed: bool,
}

impl From<bool> for IsFreshyNightTicketRedeemedResponse {
    fn from(value: bool) -> Self {
        Self { redeemed: value }
    }
}

impl IntoDto for IsFreshyNightTicketRedeemedResponse {
    type Target = IsFreshyNightTicketRedeemedResponse;

    fn into_dto(self) -> Self::Target {
        self
    }
}
