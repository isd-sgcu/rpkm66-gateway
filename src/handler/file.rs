use axum::{
    extract::{Multipart, State},
    response::IntoResponse,
    Json,
};

use crate::{dto::FileResponse, middleware::auth::Cred, Error};

#[derive(Clone)]
pub struct Handler {
    service: crate::service::file::Service,
}

impl Handler {
    pub fn new(service: crate::service::file::Service) -> Self {
        Self { service }
    }
}

#[utoipa::path(
    post,
    path = "/file/upload",
    request_body(
        content = FileUploadRequest,
        content_type = "multipart/form-data",
    ),
    responses(
        (status = 200, description = "Success", body = FileResponse),
        (status = 400, description = "Bad request"),
        (status = 401, description = "Unauthorized"),
    ),
    security(
        ("api_key" = []),
    ),
)]
pub async fn upload(
    State(handler): State<Handler>,
    cred: Cred,
    mut multipart: Multipart,
) -> impl IntoResponse {
    let user_id = cred.user_id;

    let mut tag = -1;
    let mut ty = -1;
    let mut data = None;
    let mut filename = None;

    while let Ok(Some(field)) = multipart.next_field().await {
        let Some(name) = field.name() else {
            continue;
        };

        match name {
            "file" => {
                filename = Some(field.file_name()).flatten().map(ToString::to_string);
                if matches!(
                    field.content_type(),
                    Some("image/png" | "image/jpg" | "image/jpeg" | "image/gif")
                ) {
                    let bytes = field.bytes().await?;
                    data = Some(bytes.to_vec());
                }
            }
            "type" => {
                ty = field.text().await?.parse().unwrap_or(-1);
            }
            "tag" => {
                tag = field.text().await?.parse().unwrap_or(-1);
            }
            _ => {}
        }
    }

    let (Some(data), Some(filename)) = (data, filename) else {
        return Err(Error::BadRequest);
    };

    if tag == -1 || ty == -1 {
        return Err(Error::BadRequest);
    }

    let url = handler
        .service
        .upload(data, filename, user_id, tag, ty)
        .await?;

    Ok(Json(FileResponse { url }))
}
