use utoipa::ToSchema;

#[derive(serde::Serialize, ToSchema)]
pub struct User {
    id: String,
    title: String,
    firstname: String,
    lastname: String,
    nickname: String,
    student_id: String,
    faculty: String,
    year: String,
    phone: String,
    line_id: String,
    email: String,
    allergy_food: String,
    food_restriction: String,
    allergy_medicine: String,
    disease: String,
    image_url: String,
    can_select_baan: bool,
    is_verify: bool,
    baan_id: String,
    is_got_ticket: bool,
}

into_dto!(
    rpkm66_rust_proto::rpkm66::backend::user::v1::User,
    User,
    id,
    title,
    firstname,
    lastname,
    nickname,
    student_id,
    faculty,
    year,
    phone,
    line_id,
    email,
    allergy_food,
    food_restriction,
    allergy_medicine,
    disease,
    image_url,
    can_select_baan,
    is_verify,
    baan_id,
    is_got_ticket
);

direct_map!(
    rpkm66_rust_proto::rpkm66::backend::user::v1::User,
    User,
    id,
    title,
    firstname,
    lastname,
    nickname,
    student_id,
    faculty,
    year,
    phone,
    line_id,
    email,
    allergy_food,
    food_restriction,
    allergy_medicine,
    disease,
    image_url,
    can_select_baan,
    is_verify,
    baan_id,
    is_got_ticket
);

#[derive(serde::Deserialize, ToSchema)]
pub struct UpdateUser {
    title: String,
    firstname: String,
    lastname: String,
    nickname: String,
    phone: String,
    line_id: String,
    email: String,
    allergy_food: String,
    food_restriction: String,
    allergy_medicine: String,
    disease: String,
}

impl UpdateUser {
    pub fn into_proto(
        self,
        user_id: String,
    ) -> rpkm66_rust_proto::rpkm66::backend::user::v1::UpdateUserRequest {
        rpkm66_rust_proto::rpkm66::backend::user::v1::UpdateUserRequest {
            id: user_id,
            title: self.title,
            firstname: self.firstname,
            lastname: self.lastname,
            nickname: self.nickname,
            phone: self.phone,
            line_id: self.line_id,
            email: self.email,
            allergy_food: self.allergy_food,
            food_restriction: self.food_restriction,
            allergy_medicine: self.allergy_medicine,
            disease: self.disease,
        }
    }
}
