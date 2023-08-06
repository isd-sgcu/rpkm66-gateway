use axum::{
    extract::{Path, State},
    response::IntoResponse,
};

use crate::{dto::IntoDto, error::Error, middleware::auth::Cred};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::staff::Service,
    user_service: crate::service::user::Service,
}

impl Handler {
    pub fn new(
        service: crate::service::staff::Service,
        user_service: crate::service::user::Service,
    ) -> Self {
        Self {
            service,
            user_service,
        }
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

#[utoipa::path(
    get,
    path = "/staff/user/{user_id}",
    tag = "Staff",
    responses(
        (status = 200, description = "Success", body = User),
        (status = 400, description = "Bad request"),
        (status = 401, description = "Unauthorized"),
        (status = 401, description = "Not Found"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn get_user(
    State(hdr): State<Handler>,
    cred: Cred,
    Path(user_id): Path<String>,
) -> impl IntoResponse {
    if !hdr.service.is_staff(cred.user_id).await? {
        return Err(Error::Forbidden);
    }

    hdr.user_service
        .find_one(user_id)
        .await
        .map(IntoDto::into_response)
}
