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
        crate::handler::file::upload,
        crate::handler::user::update,
    ),
    components(schemas(
        crate::dto::Validate,
        crate::dto::RedeemNewToken,
        crate::dto::VerifyTicket,
        crate::dto::User,
        crate::dto::FileResponse,
        crate::dto::FileUploadRequest,
        crate::dto::UpdateUser,
    )),
    info(
        title = "RPKM66",
        contact(name = "isd.team.sgcu@gmail.com"),
    ),
    servers(
        (url = "http://localhost:{port}", description = "Local server", variables(
            ("port" = (default= "3000", description = "port"))
        )),
        (url = "https://pbeta.freshmen2023.sgcu.in.th", description = "beta server"),
        (url = "https://pdev.freshmen2023.sgcu.in.th", description = "dev server")
    )
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
