use utoipa::{
    openapi::security::{Http, SecurityScheme},
    OpenApi,
};

#[derive(utoipa::OpenApi)]
#[openapi(
    paths(
        crate::health_check,
        crate::handler::auth::validate,
        crate::handler::auth::verify_ticket,
        crate::handler::auth::refresh_token,
    ),
    components(schemas(
        crate::dto::Validate,
        crate::dto::RedeemNewToken,
        crate::dto::VerifyTicket,
        crate::dto::User,
    ))
)]
pub struct ApiDoc;

pub fn get_doc() -> utoipa::openapi::OpenApi {
    let mut doc = ApiDoc::openapi();

    if doc.components.is_none() {
        doc.components = Some(Default::default());
    }

    doc.components.as_mut().unwrap().add_security_scheme(
        "api_key",
        SecurityScheme::Http(Http::new(utoipa::openapi::security::HttpAuthScheme::Bearer)),
    );

    doc
}
