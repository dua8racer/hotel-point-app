package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HotelID     primitive.ObjectID `bson:"hotel_id" json:"hotel_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Capacity    int                `bson:"capacity" json:"capacity"`
	Image       string             `bson:"image" json:"image"`
}
