use rpkm66_rust_proto::rpkm66::backend::baan::v1::{baan_service_client::BaanServiceClient, *};
use tonic::transport::Channel;

use crate::{Error, Result};

#[derive(Clone)]
pub struct Service {
    client: BaanServiceClient<Channel>,
}

impl Service {
    pub fn new(client: BaanServiceClient<Channel>) -> Self {
        Self { client }
    }

    pub async fn find_all(&self) -> Result<Vec<Baan>> {
        Ok(self
            .client
            .clone()
            .find_all_baan(FindAllBaanRequest {})
            .await?
            .into_inner()
            .baans)
    }

    pub async fn find_one(&self, baan_id: String) -> Result<Baan> {
        self.client
            .clone()
            .find_one_baan(FindOneBaanRequest { id: baan_id })
            .await?
            .into_inner()
            .baan
            .ok_or(Error::NotFound)
    }
}
