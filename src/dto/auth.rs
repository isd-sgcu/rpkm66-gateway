use utoipa::ToSchema;

#[derive(serde::Deserialize, ToSchema)]
pub struct VerifyTicket {
    #[schema(
        example = "AAAAAAAAAAo2NTQwMDAwMDIxAAAAAAAAAARKb2huAAAAAAAAAANEb2UAAAAAAAAAFAu8DXGVMGlZFIP0MlQSHIOkOxLE"
    )]
    pub ticket: String,
}

#[derive(serde::Serialize, ToSchema)]
pub struct TokenPayloadAuth {
    pub user_id: String,
    pub role: String,
}

#[derive(serde::Deserialize, ToSchema)]
pub struct RedeemNewToken {
    pub refresh_token: String,
}

#[derive(serde::Deserialize, ToSchema)]
pub struct Validate {
    pub token: String,
}

#[derive(serde::Deserialize, serde::Serialize, ToSchema)]
pub struct Credential {
    #[schema(example = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL3BiZX...")]
    pub access_token: String,
    pub refresh_token: String,
    #[schema(example = 3600)]
    pub expires_in: i32,
}

into_dto!(
    rpkm66_rust_proto::rpkm66::auth::auth::v1::Credential,
    Credential,
    access_token,
    refresh_token,
    expires_in
);
