use rpkm66_rust_proto::rpkm66::file::file::v1::{file_service_client::FileServiceClient, *};
use tonic::transport::Channel;

use crate::Result;

#[derive(Clone)]
pub struct Service {
    client: FileServiceClient<Channel>,
}

impl Service {
    pub fn new(client: FileServiceClient<Channel>) -> Self {
        Self { client }
    }

    pub async fn upload(
        &self,
        data: Vec<u8>,
        filename: String,
        user_id: String,
        tag: i32,
        r#type: i32,
    ) -> Result<String> {
        Ok(self
            .client
            .clone()
            .upload(UploadRequest {
                data,
                filename,
                user_id,
                tag,
                r#type,
                ..Default::default()
            })
            .await?
            .into_inner()
            .url)
    }
}
