package dto

type UserDto struct {
	ID              string `json:"id" validate:"uuid_optional"`
	Title           string `json:"title" validate:"required"`
	Firstname       string `json:"firstname" validate:"required"`
	Lastname        string `json:"lastname" validate:"required"`
	Nickname        string `json:"nickname" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	LineID          string `json:"line_id" validate:"required"`
	Email           string `json:"email" validate:"email"`
	AllergyFood     string `json:"allergy_food"`
	FoodRestriction string `json:"food_restriction"`
	AllergyMedicine string `json:"allergy_medicine"`
	Disease         string `json:"disease"`
	GroupId         string `json:"group_id"`
	CanSelectBaan   *bool  `json:"can_select_baan" validate:"required"`
}

type UpdateUserDto struct {
	Title           string `json:"title"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Nickname        string `json:"nickname"`
	Phone           string `json:"phone"`
	LineID          string `json:"line_id"`
	Email           string `json:"email"`
	AllergyFood     string `json:"allergy_food"`
	FoodRestriction string `json:"food_restriction"`
	AllergyMedicine string `json:"allergy_medicine"`
	Disease         string `json:"disease"`
}
