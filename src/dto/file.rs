use utoipa::ToSchema;

#[derive(serde::Deserialize, serde::Serialize, ToSchema)]
pub struct FileResponse {
    pub url: String,
}

#[derive(ToSchema)]
pub struct FileUploadRequest {
    pub file: Vec<u8>,
    pub tag: i32,
    /// 1 for file
    /// 2 for image
    pub r#type: i32,
}
