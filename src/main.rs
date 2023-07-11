use axum::body::Body;
use axum::extract::FromRef;
use axum::response::IntoResponse;
use axum::routing::{get, patch, post};
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
    pub user_hdr: handler::user::Handler,
    pub auth_svc: service::auth::Service,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let config = config::get_config();

    let cors = tower_http::cors::CorsLayer::new()
        .allow_headers(Any)
        .allow_methods(Any)
        .allow_origin(Any);

    let body_limit_layer = axum::extract::DefaultBodyLimit::max((config.app.max_file_size * 1024 * 1024).try_into().expect("Unable to calculate max file size"));

    let trace = tower_http::trace::TraceLayer::new_for_http();

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
        rpkm66_rust_proto::rpkm66::file::file::v1::file_service_client::FileServiceClient::new(
            file_conn,
        );

    let auth_svc = service::auth::Service::new(auth_client);
    let user_svc = service::user::Service::new(user_client);
    let baan_svc = service::baan::Service::new(baan_client);
    let file_svc = service::file::Service::new(file_client);

    let auth_hdr = handler::auth::Handler::new(auth_svc.clone(), user_svc.clone());
    let baan_hdr = handler::baan::Handler::new(baan_svc.clone(), user_svc.clone());
    let file_hdr = handler::file::Handler::new(file_svc.clone());
    let user_hdr = handler::user::Handler::new(user_svc.clone());

    let state = AppState {
        auth_hdr: auth_hdr.clone(),
        auth_svc: auth_svc.clone(),
        baan_hdr: baan_hdr.clone(),
        user_hdr: user_hdr.clone(),
        file_hdr: file_hdr.clone(),
    };

    let mut non_state_app: Router<AppState, Body> = Router::new();

    if config.app.debug {
        non_state_app = non_state_app
            .merge(SwaggerUi::new("/swagger-ui").url("/api-docs/openapi.json", doc::get_doc()));
    }

    non_state_app = non_state_app
        .route("/", get(health_check))
        .route("/auth/verify", post(handler::auth::verify_ticket))
        .route("/auth/me", get(handler::auth::validate))
        .route("/auth/refreshToken", post(handler::auth::refresh_token))
        .route("/file/upload", post(handler::file::upload))
        .route("/user", patch(handler::user::update))
        .layer(body_limit_layer)
        .layer(trace)
        .layer(cors);

    let app = non_state_app.with_state(state);

    let addr = format!("0.0.0.0:{}", config.app.port);

    Server::bind(&addr.parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
