package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	HotelID   primitive.ObjectID `bson:"hotel_id" json:"hotel_id"`
	RoomID    primitive.ObjectID `bson:"room_id" json:"room_id"`
	CheckIn   time.Time          `bson:"check_in" json:"check_in"`
	CheckOut  time.Time          `bson:"check_out" json:"check_out"`
	PointCost int                `bson:"point_cost" json:"point_cost"`
	Status    string             `bson:"status" json:"status"` // "pending", "confirmed", "cancelled"
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
