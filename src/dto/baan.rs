use utoipa::ToSchema;

use super::IntoDto;

#[derive(ToSchema, serde::Deserialize, serde::Serialize)]
pub struct Baan {
    id: String,
    name_th: String,
    description_th: String,
    name_en: String,
    description_en: String,
    size: BaanSize,
    facebook: String,
    facebook_url: String,
    instagram: String,
    instagram_url: String,
    line: String,
    line_url: String,
    image_url: String,
}

#[derive(ToSchema, serde::Deserialize, serde::Serialize)]
pub struct BaanInfo {
    id: String,
    name_th: String,
    name_en: String,
    image_url: String,
    baan_size: BaanSize,
}

#[derive(serde::Serialize, serde::Deserialize, ToSchema)]
#[repr(i32)]
pub enum BaanSize {
    Unknown = 0,
    S = 1,
    M = 2,
    L = 3,
    Xl = 4,
    Xxl = 5,
}

impl From<i32> for BaanSize {
    fn from(value: i32) -> Self {
        match value {
            1 => BaanSize::S,
            2 => BaanSize::M,
            3 => BaanSize::L,
            4 => BaanSize::Xl,
            5 => BaanSize::Xxl,
            _ => BaanSize::Unknown,
        }
    }
}

impl IntoDto for rpkm66_rust_proto::rpkm66::backend::baan::v1::Baan {
    type Target = Baan;

    fn into_dto(self) -> Self::Target {
        Self::Target {
            id: self.id,
            name_th: self.name_th,
            description_th: self.description_th,
            name_en: self.name_en,
            description_en: self.description_en,
            size: self.size.into(),
            facebook: self.facebook,
            facebook_url: self.facebook_url,
            instagram: self.instagram,
            instagram_url: self.instagram_url,
            line: self.line,
            line_url: self.line_url,
            image_url: self.image_url,
        }
    }
}

impl IntoDto for rpkm66_rust_proto::rpkm66::backend::baan::v1::BaanInfo {
    type Target = BaanInfo;
    
    fn into_dto(self) -> Self::Target {
        BaanInfo {
            id: self.id,
            baan_size: self.size.into(),
            image_url: self.image_url,
            name_en: self.name_en,
            name_th: self.name_th,
        }
    }
}
