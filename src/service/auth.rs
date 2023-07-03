use rpkm66_rust_proto::rpkm66::auth::auth::v1::{
    auth_service_client::AuthServiceClient, Credential, RefreshTokenRequest, ValidateRequest,
    VerifyTicketRequest,
};
use tonic::transport::Channel;

use crate::{dto::TokenPayloadAuth, Error, Result};

#[derive(Clone)]
pub struct Service {
    client: AuthServiceClient<Channel>,
}

impl Service {
    pub fn new(client: AuthServiceClient<Channel>) -> Self {
        Service { client }
    }

    pub async fn verify_ticket(&self, ticket: String) -> Result<Credential> {
        self.client
            .clone()
            .verify_ticket(VerifyTicketRequest {
                ticket,
                ..Default::default()
            })
            .await?
            .into_inner()
            .credential
            .ok_or(Error::InternalServer)
    }

    pub async fn validate(&self, token: String) -> Result<TokenPayloadAuth> {
        let response = self
            .client
            .clone()
            .validate(ValidateRequest {
                token,
                ..Default::default()
            })
            .await?
            .into_inner();

        Ok(TokenPayloadAuth {
            user_id: response.user_id,
            role: response.role,
        })
    }

    pub async fn refresh_token(&self, refresh_token: String) -> Result<Credential> {
        self.client
            .clone()
            .refresh_token(RefreshTokenRequest {
                refresh_token,
                ..Default::default()
            })
            .await?
            .into_inner()
            .credential
            .ok_or(Error::InternalServer)
    }
}
