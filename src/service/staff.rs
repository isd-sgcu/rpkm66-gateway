use crate::{Error, Result};
use rpkm66_rust_proto::rpkm66::checkin::staff::v1::{
    staff_service_client::StaffServiceClient, AddEventToUserRequest, IsStaffRequest,
};
use tonic::transport::Channel;

#[derive(Clone)]
pub struct Service {
    client: StaffServiceClient<Channel>,
}

impl Service {
    pub fn new(client: StaffServiceClient<Channel>) -> Self {
        Service { client }
    }

    pub async fn is_staff(&self, staff_id: String) -> Result<bool> {
        Ok(self
            .client
            .clone()
            .is_staff(IsStaffRequest { staff_id })
            .await?
            .into_inner()
            .is_staff)
    }

    pub async fn checkin_freshy_night(&self, staff_id: String, user_id: String) -> Result<bool> {
        Ok(self
            .client
            .clone()
            .add_event_to_user(AddEventToUserRequest {
                event_id: "freshy_night".to_string(),
                staff_user_id: staff_id,
                user_id,
            })
            .await?
            .into_inner()
            .success)
    }
}
