use axum::{
    extract::{Path, State},
    response::IntoResponse,
};

use crate::{dto::IntoDto, middleware::auth::Cred};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::staff::Service,
}

impl Handler {
    pub fn new(service: crate::service::staff::Service) -> Self {
        Self { service }
    }
}

#[utoipa::path(
    get,
    path = "/staff/check",
    tag = "Staff",
    responses(
        (status = 200, description = "Success", body = IsStaffResponse),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn is_staff(State(hdr): State<Handler>, cred: Cred) -> impl IntoResponse {
    hdr.service
        .is_staff(cred.user_id)
        .await
        .map(crate::dto::IsStaffResponse::from)
        .map(IntoDto::into_response)
}

#[utoipa::path(
    post,
    path = "/staff/checkin_freshy_night/{user_id}",
    tag = "Staff",
    responses(
        (status = 200, description = "Success", body = CheckingFreshyNightResponse),
        (status = 400, description = "Bad request"),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn checkin_freshy_night(
    State(hdr): State<Handler>,
    cred: Cred,
    Path(user_id): Path<String>,
) -> impl IntoResponse {
    hdr.service
        .checkin_freshy_night(cred.user_id, user_id)
        .await
        .map(crate::dto::CheckingFreshyNightResponse::from)
        .map(IntoDto::into_response)
}
