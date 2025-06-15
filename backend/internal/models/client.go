package models

import "time"

type Client struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Phone     string    `json:"phone" bson:"phone"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
