package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     uint
	ReceiverID uint
	Message    string
}
