package dto

import (
	"course_project/internal/constants/roles"
)

type CreateOperatorDTO struct {
	Username string     `json:"username" bson:"username" example:"John Doe"`
	Email    string     `json:"email" bson:"email" example:"test@test.com"`
	PwdHash  string     `json:"pwd_hash" bson:"pwd_hash" example:"$2a$10$pikzoSYzIs1GRRPi0vermeY1mPH4"`
	Role     roles.Role `json:"role" bson:"role" example:"operator"`
}
