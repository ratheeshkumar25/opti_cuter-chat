package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoCall struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     uint               `bson:"user_id"`
	ReceiverID uint               `bson:"receiver_id"`
	RoomURL    string             `bson:"room_url"`
	Timestamp  primitive.DateTime `bson:"timestamp"`
}
