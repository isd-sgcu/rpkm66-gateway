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
    pub baan_hdr: handler::baan::Handler,
    pub file_hdr: handler::file::Handler,
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
    let file_conn = Channel::from_shared(format!("http://{}", config.service.file))
    .expect("Unable to connect to file service")
    .connect_lazy();

    let auth_client =
        rpkm66_rust_proto::rpkm66::auth::auth::v1::auth_service_client::AuthServiceClient::new(
            auth_conn,
        );
    let user_client =
        rpkm66_rust_proto::rpkm66::backend::user::v1::user_service_client::UserServiceClient::new(
            backend_conn.clone(),
        );
    let baan_client =
        rpkm66_rust_proto::rpkm66::backend::baan::v1::baan_service_client::BaanServiceClient::new(
            backend_conn,
        );
    let file_client = 
        rpkm66_rust_proto::rpkm66::file::file::v1::file_service_client::FileServiceClient::new(file_conn);

    let auth_svc = service::auth::Service::new(auth_client);
    let user_svc = service::user::Service::new(user_client);
    let baan_svc = service::baan::Service::new(baan_client);
    let file_svc = service::file::Service::new(file_client);

    let auth_hdr = handler::auth::Handler::new(auth_svc.clone(), user_svc.clone());
    let baan_hdr = handler::baan::Handler::new(baan_svc.clone(), user_svc.clone());
    let file_hdr = handler::file::Handler::new(file_svc.clone());

    let state = AppState {
        auth_hdr: auth_hdr.clone(),
        auth_svc: auth_svc.clone(),
        baan_hdr: baan_hdr.clone(),
        file_hdr: file_hdr.clone(),
    };

    let app = Router::new()
        .merge(SwaggerUi::new("/swagger-ui").url("/api-docs/openapi.json", doc::get_doc()))
        .route("/", get(health_check))
        .route("/auth/verify", post(handler::auth::verify_ticket))
        .route("/auth/me", get(handler::auth::validate))
        .route("/auth/refreshToken", post(handler::auth::refresh_token))
        .route("/baan", get(handler::baan::find_all))
        .route("/baan/:id", get(handler::baan::find_one))
        .route("/baan/user", get(handler::baan::get_user_baan))
        .route("/file/upload", post(handler::file::upload))
        .layer(cors)
        .with_state(state);

    let addr = format!("0.0.0.0:{}", config.app.port);

    Server::bind(&addr.parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
