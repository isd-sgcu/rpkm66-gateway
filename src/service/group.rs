use crate::{Error, Result};
use rpkm66_rust_proto::rpkm66::backend::group::v1::{
    group_service_client::GroupServiceClient, DeleteMemberGroupRequest, FindByTokenGroupRequest,
    FindByTokenGroupResponse, FindOneGroupRequest, Group, JoinGroupRequest, LeaveGroupRequest,
    SelectBaanRequest,
};
use tonic::transport::Channel;

#[derive(Clone)]
pub struct Service {
    client: GroupServiceClient<Channel>,
}

impl Service {
    pub fn new(client: GroupServiceClient<Channel>) -> Self {
        Self { client }
    }

    pub async fn find_one(&self, user_id: String) -> Result<Group> {
        self.client
            .clone()
            .find_one(FindOneGroupRequest {
                user_id,
                ..Default::default()
            })
            .await?
            .into_inner()
            .group
            .ok_or(Error::NotFound)
    }

    pub async fn find_by_token(&self, token: String) -> Result<FindByTokenGroupResponse> {
        self.client
            .clone()
            .find_by_token(FindByTokenGroupRequest {
                token,
                ..Default::default()
            })
            .await
            .map(tonic::Response::into_inner)
            .map_err(Into::into)
    }

    pub async fn join(&self, token: String, user_id: String) -> Result<Group> {
        self.client
            .clone()
            .join(JoinGroupRequest {
                token,
                user_id,
                ..Default::default()
            })
            .await?
            .into_inner()
            .group
            .ok_or(Error::NotFound)
    }

    pub async fn delete_member(&self, leader_id: String, user_id: String) -> Result<Group> {
        self.client
            .clone()
            .delete_member(DeleteMemberGroupRequest {
                leader_id,
                user_id,
                ..Default::default()
            })
            .await?
            .into_inner()
            .group
            .ok_or(Error::NotFound)
    }

    pub async fn leave(&self, user_id: String) -> Result<Group> {
        self.client
            .clone()
            .leave(LeaveGroupRequest {
                user_id,
                ..Default::default()
            })
            .await?
            .into_inner()
            .group
            .ok_or(Error::NotFound)
    }

    pub async fn select_baans(&self, user_id: String, baans: Vec<String>) -> Result<bool> {
        Ok(self
            .client
            .clone()
            .select_baan(SelectBaanRequest {
                baans,
                user_id,
                ..Default::default()
            })
            .await?
            .into_inner()
            .success)
    }
}
