package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DateRule struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date      time.Time          `bson:"date" json:"date"`
	Type      string             `bson:"type" json:"type"`             // "regular", "weekend", "holiday"
	PointCost int                `bson:"point_cost" json:"point_cost"` // 1, 2, or 3
	Name      string             `bson:"name" json:"name,omitempty"`   // For holidays
}
