use rpkm66_rust_proto::rpkm66::{
    checkin::event::v1::event_service_client::EventServiceClient,
    checkin::user::v1::{user_service_client::UserServiceClient, GetUserEventByEventIdRequest},
    checkin::{
        event::v1::GetEventsByNamespaceIdRequest,
        user::v1::{AddEventRequest, GetAllUserEventsByNamespaceIdRequest},
    },
};
use tonic::transport::Channel;

use crate::{error::Error, handler::user, Result};

#[derive(Clone)]
pub struct Service {
    event_client: EventServiceClient<Channel>,
    user_client: UserServiceClient<Channel>,
}

impl Service {
    pub fn new(
        event_client: EventServiceClient<Channel>,
        user_client: UserServiceClient<Channel>,
    ) -> Self {
        Service {
            event_client,
            user_client,
        }
    }

    pub async fn get_all_estamps(
        &self,
    ) -> Result<Vec<rpkm66_rust_proto::rpkm66::checkin::event::v1::Event>> {
        Ok(self
            .event_client
            .clone()
            .get_events_by_namespace_id(GetEventsByNamespaceIdRequest {
                namespace_id: "estamp".to_string(),
            })
            .await?
            .into_inner()
            .events)
    }

    pub async fn claim_estamp(
        &self,
        user_id: String,
        token: String,
    ) -> Result<rpkm66_rust_proto::rpkm66::checkin::event::v1::Event> {
        self.user_client
            .clone()
            .add_event(AddEventRequest { token, user_id })
            .await?
            .into_inner()
            .event
            .ok_or(Error::NotFound)
    }

    pub async fn get_user_estamp(
        &self,
        user_id: String,
    ) -> Result<Vec<rpkm66_rust_proto::rpkm66::checkin::event::v1::UserEvent>> {
        Ok(self
            .user_client
            .clone()
            .get_all_user_events_by_namespace_id(GetAllUserEventsByNamespaceIdRequest {
                user_id,
                namespace_id: "estamp".to_string(),
            })
            .await?
            .into_inner()
            .event)
    }

    pub async fn redeem_item(&self, user_id: String) -> Result<bool> {
        let has_redeem = self.has_redeem_item(user_id.clone()).await?;

        if !has_redeem {
            let stamps = self.get_user_estamp(user_id.clone()).await?;

            if stamps.len() != 4 {
                return Err(Error::Forbidden);
            }

            self.user_client
                .clone()
                .add_event(AddEventRequest {
                    token: "redeem".to_string(),
                    user_id,
                })
                .await?
                .into_inner()
                .event
                .ok_or(Error::NotFound)
                .map(|_| true)
        } else {
            Err(Error::Duplicated)
        }
    }

    pub async fn has_redeem_item(&self, user_id: String) -> Result<bool> {
        self.user_client
            .clone()
            .get_user_event_by_event_id(GetUserEventByEventIdRequest {
                event_id: "redeem-item".to_string(),
                user_id,
            })
            .await?
            .into_inner()
            .user_event
            .ok_or(Error::NotFound)
            .map(|x| x.is_taken)
    }
}
