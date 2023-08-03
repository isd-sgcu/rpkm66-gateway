use axum::{extract::State, response::IntoResponse, Json};

use crate::{
    dto::{UpdatePersonality, UpdateUser},
    middleware::auth::Cred,
};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::user::Service,
}

impl Handler {
    pub fn new(service: crate::service::user::Service) -> Self {
        Self { service }
    }
}

#[utoipa::path(
    patch,
    path = "/user",
    tag = "User",
    request_body = UpdateUser,
    responses(
        (status = 200, description = "Success", body = User),
        (status = 400, description = "Bad format"),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn update(
    State(handler): State<Handler>,
    cred: Cred,
    Json(updated_user): Json<UpdateUser>,
) -> impl IntoResponse {
    let user_id = cred.user_id;

    handler
        .service
        .update(user_id, updated_user)
        .await
        .map(Json)
}

#[utoipa::path(
    patch,
    path = "/user/personality",
    tag = "User",
    request_body = UpdatePersonality,
    responses(
        (status = 200, description = "Success", body = User),
        (status = 400, description = "Bad format"),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn update_personality(
    State(handler): State<Handler>,
    cred: Cred,
    Json(updated_personality): Json<UpdatePersonality>,
) -> impl IntoResponse {
    let user_id = cred.user_id;

    handler
        .service
        .update_personality(user_id, updated_personality.personality)
        .await
        .map(Json)
}
