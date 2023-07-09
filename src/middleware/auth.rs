use axum::{
    async_trait,
    extract::{FromRef, FromRequestParts},
    headers::Authorization,
    http::request::Parts,
    TypedHeader,
};

use crate::{error::Error, service};

use axum::headers::authorization::Bearer;

pub struct Cred {
    pub user_id: String,
    pub role: String,
}

#[async_trait]
impl<S> FromRequestParts<S> for Cred
where
    service::auth::Service: FromRef<S>,
    S: Send + Sync,
{
    type Rejection = Error;

    async fn from_request_parts(
        parts: &mut Parts,
        state: &S,
    ) -> std::result::Result<Self, Self::Rejection> {
        let auth_svc = service::auth::Service::from_ref(state);

        if let Ok(TypedHeader(Authorization(token))) =
            TypedHeader::<Authorization<Bearer>>::from_request_parts(parts, state).await
        {
            if let Ok(x) = auth_svc.validate(token.token().to_owned()).await {
                Ok(Cred {
                    role: x.role,
                    user_id: x.user_id,
                })
            } else {
                Err(Error::InternalServer)
            }
        } else {
            Err(Error::Unauthorized)
        }
    }
}
