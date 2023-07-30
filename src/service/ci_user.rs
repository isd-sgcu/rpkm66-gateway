use rpkm66_rust_proto::rpkm66::{backend::user::v1::user_service_client::UserServiceClient, checkin::user::v1::GetUserEventByEventIdRequest};
use tonic::transport::Channel;

use crate::Result;

#[derive(Clone)]
pub struct Service {
    client: UserServiceClient<Channel>,
}

impl Service {
    pub fn new(client: UserServiceClient<Channel>) -> Self {
        Service { client }
    }

    pub async fn is_freshy_night_ticket_redeemed(&self, user_id: String) -> Result<bool> {
        Ok(self
            .client
            .clone()
            .get_user_event_by_event_id(GetUserEventByEventIdRequest {
                event_id: "freshy_night".to_string(),
                user_id,
            })
            .await?
            .into_inner()
            .user_event
            .is_some())
    }
}
