use axum::{extract::State, response::IntoResponse};

use crate::{dto::IntoDto, middleware::auth::Cred};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::checkin::Service,
}

impl Handler {
    pub fn new(service: crate::service::checkin::Service) -> Self {
        Self { service }
    }
}

/// Has checkin
///
/// Check whether user has checkin or not
#[utoipa::path(
    get,
    path = "/checkin",
    tag = "Check in",
    responses(
        (status = 200, description = "Success", body = HasCheckinResponse),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn has_checkin(State(hdr): State<Handler>, cred: Cred) -> impl IntoResponse {
    hdr.service
        .has_checkin(cred.user_id)
        .await
        .map(crate::dto::HasCheckinResponse::from)
        .map(IntoDto::into_response)
}

#[utoipa::path(
    post,
    path = "/checkin",
    tag = "Check in",
    responses(
        (status = 200, description = "Success", body = CheckinResponse),
        (status = 401, description = "Unauthorized"),
        (status = 409, description = "Duplicated"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn checkin(State(hdr): State<Handler>, cred: Cred) -> impl IntoResponse {
    hdr.service
        .checkin(cred.user_id)
        .await
        .map(crate::dto::CheckinResponse::from)
        .map(IntoDto::into_response)
}
