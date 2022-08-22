package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Project
type Project struct {
	ID           primitive.ObjectID `bson:"_id"`
	Project_name *string            `json:"project_name" validate:"required,min=2,max=100"`
	Project_type *string            `json:"project_type" validate:"required,eq=ADMIN|eq=USER"`
	Logo         *string            `json:"logo" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"email,required"`
	Member       *string            `json:"phone" validate:"required"`
	User_created *string            `json:"token"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	User_id      string             `json:"user_id"`
}
