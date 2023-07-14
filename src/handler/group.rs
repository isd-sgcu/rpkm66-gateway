use axum::{
    extract::{Path, State},
    response::IntoResponse,
};

use crate::{dto::IntoDto, middleware::auth::Cred};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::group::Service,
}

impl Handler {
    pub fn new(service: crate::service::group::Service) -> Self {
        Self { service }
    }
}

/// Get user current group
#[utoipa::path(
    get,
    path = "/group",
    tag = "Group",
    responses(
        (status = 200, description = "Success", body = Group),
        (status = 401, description = "Unauthorized"),
        (status = 404, description = "Not found"),
    ),
    security(
        ("api_key" = [])
    ),
)]
pub async fn find_one(State(handler): State<Handler>, cred: Cred) -> impl IntoResponse {
    handler
        .service
        .find_one(cred.user_id)
        .await
        .map(IntoDto::into_response)
}

/// Get group by group's token
#[utoipa::path(
    get,
    path = "/group/{token}",
    tag = "Group",
    responses(
        (status = 200, description = "Success", body = GroupOverview),
        (status = 404, description = "Not found"),
    ),
)]
pub async fn find_by_token(
    State(handler): State<Handler>,
    Path(token): Path<String>,
) -> impl IntoResponse {
    handler
        .service
        .find_by_token(token)
        .await
        .map(IntoDto::into_response)
}

/// Join group
/// 
/// Join group by using group's token
#[utoipa::path(
    post,
    path = "/group/{token}",
    tag = "Group",
    responses(
        (status = 200, description = "Success", body = Group),
        (status = 401, description = "Unauthorized"),
        (status = 404, description = "Not found"),
    ),
    security(
        ("api_key" = [])
    )
)]
pub async fn join(
    State(handler): State<Handler>,
    Path(token): Path<String>,
    cred: Cred,
) -> impl IntoResponse {
    handler
        .service
        .join(token, cred.user_id)
        .await
        .map(IntoDto::into_response)
}

/// Delete member from group
/// 
/// must be leader to do this.
#[utoipa::path(
    delete,
    path = "/group/members/{member_id}",
    tag = "Group",
    responses(
        (status = 200, description = "Success", body = GroupOverview),
        (status = 401, description = "Unauthorized"),
        (status = 404, description = "Not found"),
    ),
    security(
        ("api_key" = [])
    )
)]
pub async fn delete_member(
    State(handler): State<Handler>,
    Path(member_id): Path<String>,
    cred: Cred,
) -> impl IntoResponse {
    handler
        .service
        .delete_member(cred.user_id, member_id)
        .await
        .map(IntoDto::into_response)
}

/// Leave group
///
/// Leader of the group cannot leave group
#[utoipa::path(
    delete,
    path = "/group/leave",
    tag = "Group",
    responses(
        (status = 200, description = "Success", body = Group),
        (status = 401, description = "Unauthorized"),
        (status = 403, description = "Forbidden"),
    ),
    security(
        ("api_key" = [])
    )
)]
pub async fn leave(State(handler): State<Handler>, cred: Cred) -> impl IntoResponse {
    handler
        .service
        .leave(cred.user_id)
        .await
        .map(IntoDto::into_response)
}
