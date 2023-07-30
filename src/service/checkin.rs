use crate::Result;
use rpkm66_rust_proto::rpkm66::checkin::user::v1::{
    user_service_client::UserServiceClient, AddEventRequest, GetUserEventByEventIdRequest,
};
use tonic::transport::Channel;

#[derive(Clone)]
pub struct Service {
    client: UserServiceClient<Channel>,
    config: crate::config::AppConfig,
}

impl Service {
    pub fn new(client: UserServiceClient<Channel>, config: crate::config::AppConfig) -> Self {
        Service { client, config }
    }

    pub async fn has_checkin(&self, user_id: String) -> Result<bool> {
        Ok(self
            .client
            .clone()
            .get_user_event_by_event_id(GetUserEventByEventIdRequest {
                user_id,
                event_id: self.get_checkin_event_id(),
            })
            .await?
            .into_inner()
            .user_event
            .is_some())
    }

    pub async fn checkin(&self, user_id: String) -> Result<bool> {
        Ok(self
            .client
            .clone()
            .add_event(AddEventRequest {
                user_id,
                token: self.get_checkin_event_id(),
            })
            .await?
            .into_inner()
            .event
            .is_some())
    }

    fn get_checkin_event_id(&self) -> String {
        format!("checkin-day-{}", self.config.event_day)
    }
}
