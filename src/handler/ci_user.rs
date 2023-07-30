use axum::{extract::State, response::IntoResponse};

use crate::{dto::IntoDto, middleware::auth::Cred};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::ci_user::Service,
}

impl Handler {
    pub fn new(service: crate::service::ci_user::Service) -> Self {
        Self { service }
    }
}

#[utoipa::path(
    get,
    path = "/freshy_night",
    tag = "Freshy Night",
    responses(
        (status = 200, description = "Success"),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn is_freshy_night_ticket_redeemed(
    State(handler): State<Handler>,
    cred: Cred,
) -> impl IntoResponse {
    handler
        .service
        .is_freshy_night_ticket_redeemed(cred.user_id)
        .await
        .map(crate::dto::IsFreshyNightTicketRedeemedResponse::from)
        .map(IntoDto::into_response)
}
