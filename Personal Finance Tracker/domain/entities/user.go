package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email      string             `bson:"email,omitempty" json:"email"`
	Password   string             `bson:"password,omitempty" json:"password"`
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	Role       string             `bson:"role,omitempty" json:"role,omitempty"`
	Currency   string             `bson:"currency,omitempty" json:"currency,omitempty"`
	Profile    Profile            `bson:"profile,omitempty" json:"profile,omitempty"`
	IsVerified bool               `bson:"is_verified,omitempty" json:"is_verified,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt  time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Profile struct {
	Address string `bson:"address,omitempty" json:"address,omitempty"`
	Phone   string `bson:"phone,omitempty" json:"phone,omitempty"`
}
