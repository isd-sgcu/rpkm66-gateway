use crate::dto;
use crate::dto::RedeemNewToken;
use crate::middleware::auth::Cred;
use crate::service::auth::Service;
use axum::{
    extract::{Json, State},
    response::IntoResponse,
};

#[derive(Clone)]
pub struct Handler {
    service: Service,
}

impl Handler {
    pub fn new(service: Service) -> Self {
        Self { service }
    }
}

#[utoipa::path(
    post,
    path = "/auth/verify",
    request_body = VerifyTicket,
    responses(
        (status = 200, description = "Success"),
    ),
    security(
        (),
    )
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
        (),
        ("api_key" = []),
    )
)]
pub async fn validate(State(_handler): State<Handler>, cred: Cred) -> impl IntoResponse {
    println!("{}: {}", cred.user_id, cred.role);

    "TODO"
}

#[utoipa::path(
    post,
    path = "/auth/refreshToken",
    request_body = RedeemNewToken,
    responses(
        (status = 200, description = "Success"),
    ),
    security(
        (),
    )
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
