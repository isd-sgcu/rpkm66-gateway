use rpkm66_rust_proto::rpkm66::backend::group::v1::FindByTokenGroupResponse;
use utoipa::ToSchema;

use super::IntoDto;

#[derive(ToSchema, serde::Serialize, serde::Deserialize)]
pub struct Group {
    id: String,
    leader_id: String,
    token: String,
    members: Vec<UserInfo>,
    baans: Vec<crate::dto::BaanInfo>,
}

#[derive(ToSchema, serde::Serialize, serde::Deserialize, Default)]
pub struct UserInfo {
    id: String,
    firstname: String,
    lastname: String,
    image_url: String,
}

#[derive(ToSchema, serde::Serialize, serde::Deserialize)]
pub struct GroupOverview {
    id: String,
    token: String,
    leader: UserInfo,
}

impl IntoDto for FindByTokenGroupResponse {
    type Target = GroupOverview;

    fn into_dto(self) -> Self::Target {
        GroupOverview {
            id: self.id,
            leader: self
                .leader
                .map(IntoDto::into_dto)
                .unwrap_or_else(|| UserInfo::default()),
            token: self.token,
        }
    }
}

into_dto!(
    rpkm66_rust_proto::rpkm66::backend::group::v1::UserInfo,
    UserInfo,
    id,
    firstname,
    lastname,
    image_url
);

impl IntoDto for rpkm66_rust_proto::rpkm66::backend::group::v1::Group {
    type Target = Group;

    fn into_dto(self) -> Self::Target {
        Group {
            id: self.id,
            leader_id: self.leader_id,
            members: self.members.into_iter().map(IntoDto::into_dto).collect(),
            token: self.token,
            baans: self.baans.into_iter().map(IntoDto::into_dto).collect(),
        }
    }
}
