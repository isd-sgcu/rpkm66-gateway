use utoipa::ToSchema;

#[derive(serde::Deserialize, ToSchema)]
pub struct VerifyTicket {
    #[schema(example = "AAAAAAAAAAo2NTQwMDAwMDIxAAAAAAAAAARKb2huAAAAAAAAAANEb2UAAAAAAAAAFAu8DXGVMGlZFIP0MlQSHIOkOxLE")]
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
