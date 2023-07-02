use axum::extract::FromRef;
use axum::response::IntoResponse;
use axum::routing::{get, post};
use axum::{Router, Server};
use tonic::transport::Channel;
use tower_http::cors::Any;
use utoipa_swagger_ui::SwaggerUi;

pub mod config;
pub mod doc;
pub mod dto;
pub mod error;
pub mod handler;
pub mod middleware;
pub mod result;
pub mod service;

pub(crate) use error::Error;
pub(crate) use result::Result;

#[utoipa::path(
    get,
    path = "/",
    responses(
        (status = 200, description = "Success")
    )
)]
pub async fn health_check() -> impl IntoResponse {
    "OK"
}

#[derive(FromRef, Clone)]
pub struct AppState {
    pub auth_hdr: handler::auth::Handler,
    pub auth_svc: service::auth::Service,
}

#[tokio::main]
async fn main() {
    let config = config::get_config();

    let cors = tower_http::cors::CorsLayer::new()
        .allow_headers(Any)
        .allow_methods(Any)
        .allow_origin(Any);

    let auth_conn = Channel::from_shared(format!("http://{}", config.service.auth))
        .expect("Unable to connect to auth service")
        .connect_lazy();
    let backend_conn = Channel::from_shared(format!("http://{}", config.service.backend))
        .expect("Unable to connect to backend service")
        .connect_lazy();

    let auth_client =
        rpkm66_rust_proto::rpkm66::auth::auth::v1::auth_service_client::AuthServiceClient::new(
            auth_conn,
        );
    let backend_client =
        rpkm66_rust_proto::rpkm66::backend::user::v1::user_service_client::UserServiceClient::new(
            backend_conn,
        );

    let auth_svc = service::auth::Service::new(auth_client);
    let user_svc = service::user::Service::new(backend_client);

    let auth_handler = handler::auth::Handler::new(auth_svc.clone(), user_svc.clone());

    let state = AppState {
        auth_hdr: auth_handler.clone(),
        auth_svc: auth_svc.clone(),
    };

    let app = Router::new()
        .merge(SwaggerUi::new("/swagger-ui").url("/api-docs/openapi.json", doc::get_doc()))
        .route("/", get(health_check))
        .route("/auth/verify", post(handler::auth::verify_ticket))
        .route("/auth/me", get(handler::auth::validate))
        .route("/auth/refreshToken", post(handler::auth::refresh_token))
        .layer(cors)
        .with_state(state);

    let addr = format!("0.0.0.0:{}", config.app.port);

    Server::bind(&addr.parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
