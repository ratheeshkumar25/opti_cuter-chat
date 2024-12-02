package repo

import (
	"context"
	"time"

	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

// Createchat implements interfaces.ChatRepoInter.
func (c *ChatRepo) Createchat(chat *model.History) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.Collection.InsertOne(ctx, chat)
	if err != nil {
		return err
	}
	return nil
}

// Findchat implements interfaces.ChatRepoInter.
func (c *ChatRepo) Findchat(userID uint, receiverID uint) (*[]model.History, error) {
	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"user_id": userID, "receiver_id": receiverID},
			{"user_id": receiverID, "receiver_id": userID},
		},
	}
	// Use Find to get all the relevant chat history
	cursor, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []model.History
	// Loop through the cursor to decode the chat documents
	for cursor.Next(ctx) {
		var chat model.History
		if err := cursor.Decode(&chat); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return &chats, nil
}
