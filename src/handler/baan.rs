use crate::{dto::IntoDto, middleware::auth::Cred};
use axum::{
    extract::{Path, State},
    response::IntoResponse,
    Json,
};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::baan::Service,
    user_svc: crate::service::user::Service,
}

impl Handler {
    pub fn new(
        service: crate::service::baan::Service,
        user_svc: crate::service::user::Service,
    ) -> Self {
        Self { service, user_svc }
    }
}

#[utoipa::path(
    get,
    path = "/baan",
    tag = "Baan",
    responses(
        (status = 200, description = "Success", body = Vec<Baan>),
    ),
)]
pub async fn find_all(State(handler): State<Handler>) -> impl IntoResponse {
    handler.service.find_all().await.map(Json)
}

#[utoipa::path(
    get,
    path = "/baan/{id}",
    tag = "Baan",
    responses(
        (status = 200, description = "Success", body = Baan),
        (status = 400, description = "Bad request"),
        (status = 404, description = "Not found"),
    ),
)]
pub async fn find_one(State(handler): State<Handler>, Path(id): Path<String>) -> impl IntoResponse {
    handler.service.find_one(id).await.map(Json)
}

#[utoipa::path(
    get,
    path = "/baan/user",
    tag = "Baan",
    responses(
        (status = 200, description = "Success", body = Baan),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn get_user_baan(State(handler): State<Handler>, cred: Cred) -> impl IntoResponse {
    let user_id = cred.user_id;

    let user = handler.user_svc.find_one(user_id).await?;

    handler
        .service
        .find_one(user.baan_id)
        .await
        .map(IntoDto::into_response)
}
