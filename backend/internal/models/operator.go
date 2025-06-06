package models

import (
	"course_project/internal/constants/roles"
	"time"
)

type Operator struct {
	ID        string     `json:"id" bson:"id" example:"0001M2PVBD5Q1DAMYJ0S2HADD6"`
	Username  string     `json:"username" bson:"username" example:"John Doe"`
	Email     string     `json:"email" bson:"email" example:"test@test.com"`
	PwdHash   string     `json:"pwd_hash" bson:"pwd_hash" example:"$2a$10$pikzoSYzIs1GRRPi0vermeY1mPH4"`
	Role      roles.Role `json:"role" bson:"role" example:"operator"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at" example:"2020-01-01T00:00:00+09:00"`
}
