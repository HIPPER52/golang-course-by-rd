package dto

type RegisterClientDTO struct {
	Name  string `json:"name" validate:"required,min=2" example:"John Doe"`
	Phone string `json:"phone" validate:"required" example:"+380931234567"`
}
