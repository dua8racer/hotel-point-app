package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Address     string             `bson:"address" json:"address"`
	City        string             `bson:"city" json:"city"`
	Image       string             `bson:"image" json:"image"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
