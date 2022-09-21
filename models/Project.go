package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project
type Project struct {
	ID           primitive.ObjectID `bson:"_id"`
	Project_name *string            `json:"project_name" validate:"required,min=2,max=100"`
	Project_type *string            `json:"project_type" validate:"required,eq=NORMAL|eq=VIP"`
	Logo         *string            `json:"logo"`
	Member       []RoleProject      `json:"member"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}

type RoleProject struct {
	ID      primitive.ObjectID `bson:"_id"`
	ID_User primitive.ObjectID `bson:"user_id"`
	Role    string             `bson:"role" validate:"required,eq=MANAGER|eq=MEMBER|eq=OWNER"`
}
