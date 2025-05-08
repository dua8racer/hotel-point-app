package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"-"`
	PointBalance int                `bson:"point_balance" json:"point_balance"`
	Role         string             `bson:"role" json:"role"` // "user" atau "admin"
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type PointTransaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount    int                `bson:"amount" json:"amount"`
	Type      string             `bson:"type" json:"type"`           // "annual_grant", "booking_deduction"
	Reference string             `bson:"reference" json:"reference"` // e.g., booking ID
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
