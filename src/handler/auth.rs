use crate::dto::{self, IntoDto};
use crate::dto::{GoogleAuthorizationCodeQuery, RedeemNewToken};
use crate::error::Error;
use crate::middleware::auth::Cred;
use crate::service::auth::Service;
use axum::extract::Query;
use axum::response::Redirect;
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
    tag = "Auth",
    request_body = VerifyTicket,
    responses(
        (status = 200, description = "Success", body = Credential),
        (status = 400, description = "Bad request"),
        (status = 401, description = "Unauthorized"),
    ),
)]
pub async fn verify_ticket(
    State(handler): State<Handler>,
    Json(ticket): Json<dto::VerifyTicket>,
) -> impl IntoResponse {
    handler
        .service
        .verify_ticket(ticket.ticket)
        .await
        .map(IntoDto::into_response)
}

#[utoipa::path(
    get,
    path = "/auth/me",
    tag = "Auth",
    responses(
        (status = 200, description = "Success", body = User),
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
    tag = "Auth",
    request_body = RedeemNewToken,
    responses(
        (status = 200, description = "Success", body = Credential),
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
        .map(IntoDto::into_response)
}

#[utoipa::path(
    get,
    path = "/auth/google",
    tag = "Auth",
    responses(
        (status = 200, description = "Success"),
    ),
)]
pub async fn get_google_oauth_redirect_uri(State(handler): State<Handler>) -> impl IntoResponse {
    let url = match handler.service.get_google_oauth_redirect_uri().await {
        Ok(x) => x,
        Err(_) => return Err(Error::InternalServer),
    };

    Ok(Redirect::temporary(&url))
}

#[utoipa::path(
    post,
    path = "/auth/google",
    tag = "Auth",
    responses(
        (status = 200, description = "Success", body = Credential),
        (status = 400, description = "Bad request"),
        (status = 401, description = "Unauthorized"),
    ),
)]
pub async fn get_token_from_google_oauth_code(
    State(handler): State<Handler>,
    Query(code): Query<GoogleAuthorizationCodeQuery>,
) -> impl IntoResponse {
    let code = code.code;

    handler
        .service
        .get_token_from_google_oauth_code(code)
        .await
        .map(IntoDto::into_response)
}
