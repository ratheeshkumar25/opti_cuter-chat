package repo

import (
	inter "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/repo/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepo struct {
	Collection          *mongo.Collection
	VideoCallCollection *mongo.Collection
	VideoCollection     *mongo.Collection
}

func NewChatRepository(mongo *mongo.Database) inter.ChatRepoInter {
	return &ChatRepo{
		Collection:          mongo.Collection("myNotification"),
		VideoCallCollection: mongo.Collection("video-call"),
		VideoCollection:     mongo.Collection("videos"),
	}
}
