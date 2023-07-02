use crate::dto::RedeemNewToken;
use crate::dto::{self, IntoDto};
use crate::middleware::auth::Cred;
use crate::service::auth::Service;
use axum::{
    extract::{Json, State},
    response::IntoResponse,
};

#[derive(Clone)]
pub struct Handler {
    service: Service,
    user_service: crate::service::user::Service,
}

impl Handler {
    pub fn new(service: Service, user_service: crate::service::user::Service) -> Self {
        Self {
            service,
            user_service,
        }
    }
}

#[utoipa::path(
    post,
    path = "/auth/verify",
    request_body = VerifyTicket,
    responses(
        (status = 200, description = "Success"),
    ),
)]
pub async fn verify_ticket(
    State(handler): State<Handler>,
    Json(ticket): Json<dto::VerifyTicket>,
) -> impl IntoResponse {
    println!("Sending: {}", ticket.ticket);
    handler.service.verify_ticket(ticket.ticket).await.map(Json)
}

#[utoipa::path(
    get,
    path = "/auth/me",
    responses(
        (status = 200, description = "Success"),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn validate(State(handler): State<Handler>, cred: Cred) -> impl IntoResponse {
    handler
        .user_service
        .find_one(cred.user_id)
        .await
        .map(IntoDto::into_response)
}

#[utoipa::path(
    post,
    path = "/auth/refreshToken",
    request_body = RedeemNewToken,
    responses(
        (status = 200, description = "Success"),
    ),
)]
pub async fn refresh_token(
    State(handler): State<Handler>,
    Json(token): Json<RedeemNewToken>,
) -> impl IntoResponse {
    handler
        .service
        .refresh_token(token.refresh_token)
        .await
        .map(Json)
}
