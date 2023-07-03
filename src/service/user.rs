use rpkm66_rust_proto::rpkm66::backend::user::v1::{user_service_client::UserServiceClient, *};
use tonic::transport::Channel;

use crate::{dto::UpdateUser, Error, Result};

#[derive(Clone)]
pub struct Service {
    client: UserServiceClient<Channel>,
}

impl Service {
    pub fn new(client: UserServiceClient<Channel>) -> Self {
        Service { client }
    }

    pub async fn find_one(&self, user_id: String) -> Result<User> {
        self.client
            .clone()
            .find_one(FindOneUserRequest {
                id: user_id,
                ..Default::default()
            })
            .await?
            .into_inner()
            .user
            .ok_or(Error::InternalServer)
    }

    pub async fn delete(&self, user_id: String) -> Result<()> {
        self.client
            .clone()
            .delete(DeleteUserRequest {
                id: user_id,
                ..Default::default()
            })
            .await?;

        Ok(())
    }

    pub async fn update(&self, user_id: String, updated_user: UpdateUser) -> Result<User> {
        self.client
            .clone()
            .update(updated_user.into_proto(user_id))
            .await?
            .into_inner()
            .user
            .ok_or(Error::InternalServer)
    }
}
