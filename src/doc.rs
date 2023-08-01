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
        crate::handler::group::find_one,
        crate::handler::group::find_by_token,
        crate::handler::group::join,
        crate::handler::group::delete_member,
        crate::handler::group::leave,
        crate::handler::group::select_baans,
        crate::handler::baan::find_one,
        crate::handler::baan::find_all,
        crate::handler::baan::get_user_baan,
        crate::handler::staff::is_staff,
        crate::handler::staff::checkin_freshy_night,
        crate::handler::ci_user::is_freshy_night_ticket_redeemed,
        crate::handler::estamp::get_all_estamps,
        crate::handler::estamp::get_user_estamps,
        crate::handler::estamp::claim_estamp,
        crate::handler::estamp::has_redeem_item,
        crate::handler::estamp::redeem_item,
        crate::handler::checkin::has_checkin,
        crate::handler::checkin::checkin,
    ),
    components(schemas(
        crate::dto::Validate,
        crate::dto::RedeemNewToken,
        crate::dto::VerifyTicket,
        crate::dto::Credential,
        crate::dto::User,
        crate::dto::FileResponse,
        crate::dto::FileUploadRequest,
        crate::dto::UpdateUser,
        crate::dto::UserInfo,
        crate::dto::Group,
        crate::dto::GroupOverview,
        crate::dto::BaanSize,
        crate::dto::Baan,
        crate::dto::BaanInfo,
        crate::dto::SelectBaan,
        crate::dto::IsStaffResponse,
        crate::dto::CheckingFreshyNightResponse,
        crate::dto::IsFreshyNightTicketRedeemedResponse,
        crate::dto::GetAllEstampResponse,
        crate::dto::GetUserEstampsResponse,
        crate::dto::EstampEvent,
        crate::dto::UserEstampEvent,
        crate::dto::HasCheckinResponse,
        crate::dto::CheckinResponse,
        crate::dto::RedeemItemResponse,
        crate::dto::HasRedeemItemResponse,
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
    ),
    tags(
        (name = "Health check"),
        (name = "Auth"),
        (name = "User"),
        (name = "File"),
        (name = "Group"),
        (name = "Baan"),
        (name = "Staff"),
        (name = "Freshy Night"),
        (name = "Estamp"),
        (name = "Check in")
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
