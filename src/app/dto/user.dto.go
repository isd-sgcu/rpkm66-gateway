package dto

type UserDto struct {
	ID                    string `json:"id" validate:"uuid_optional"`
	Firstname             string `json:"firstname" validate:"required"`
	Lastname              string `json:"lastname" validate:"required"`
	Nickname              string `json:"nickname" validate:"required"`
	StudentID             string `json:"student_id" validate:"required"`
	Faculty               string `json:"faculty" validate:"required"`
	Year                  string `json:"year" validate:"required"`
	Phone                 string `json:"phone" validate:"required"`
	LineID                string `json:"line_id" validate:"required"`
	Email                 string `json:"email" validate:"email"`
	AllergyFood           string `json:"allergy_food"`
	FoodRestriction       string `json:"food_restriction"`
	AllergyMedicine       string `json:"allergy_medicine"`
	Disease               string `json:"disease"`
	VaccineCertificateUrl string `json:"vaccine_certificate_url" validate:"url"`
	ImageUrl              string `json:"image_url" validate:"url"`
}
