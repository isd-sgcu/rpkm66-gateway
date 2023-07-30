use axum::{
    extract::{Path, State},
    response::IntoResponse,
};

use crate::{
    dto::{GetAllEstampResponse, GetUserEstampsResponse, IntoDto},
    middleware::auth::Cred,
};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::estamp::Service,
}

impl Handler {
    pub fn new(service: crate::service::estamp::Service) -> Self {
        Self { service }
    }
}

#[utoipa::path(
    get,
    path = "/estamp",
    tag = "Estamp",
    responses(
        (status = 200, description = "Success", body = User),
    ),
)]
pub async fn get_all_estamps(State(handler): State<Handler>) -> impl IntoResponse {
    handler.service.get_all_estamps().await.map(|vec| {
        let estamp_events = vec.into_iter().map(crate::dto::EstampEvent::from);
        let response = GetAllEstampResponse::from(estamp_events);
        response.into_response()
    })
}

#[utoipa::path(
    post,
    path = "/estamp/:token",
    tag = "Estamp",
    responses(
        (status = 200, description = "Success"),
        (status = 401, description = "Unauthorized"),
        (status = 404, description = "Not found"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn claim_estamp(
    State(handler): State<Handler>,
    cred: Cred,
    Path(token): Path<String>,
) -> impl IntoResponse {
    handler
        .service
        .claim_estamp(cred.user_id, token)
        .await
        .map(crate::dto::EstampEvent::from)
        .map(IntoDto::into_response)
}

#[utoipa::path(
    get,
    path = "/estamp/my",
    tag = "Estamp",
    responses(
        (status = 200, description = "Success"),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn get_user_estamps(State(handler): State<Handler>, cred: Cred) -> impl IntoResponse {
    handler
        .service
        .get_user_estamp(cred.user_id)
        .await
        .map(|proto_events| {
            let user_events = proto_events
                .into_iter()
                .map(crate::dto::UserEstampEvent::from);
            let response = GetUserEstampsResponse::from(user_events);
            response.into_response()
        })
}
