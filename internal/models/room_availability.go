package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomAvailability struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	RoomID    primitive.ObjectID   `bson:"room_id" json:"room_id"`
	Date      time.Time            `bson:"date" json:"date"`
	Available bool                 `bson:"available" json:"available"`
	UserIDs   []primitive.ObjectID `bson:"user_ids,omitempty" json:"user_ids,omitempty"` // Jika kosong, semua user dapat memesan
}
