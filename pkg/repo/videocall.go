package repo

import (
	"context"
	"time"

	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

// GetVideoCallHistory implements interfaces.ChatRepoInter.
func (c *ChatRepo) GetVideoCallHistory(userID uint, receiverID uint) ([]model.VideoCall, error) {
	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":     userID,
		"receiver_id": receiverID,
	}

	cursor, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var videos []model.VideoCall
	err = cursor.All(ctx, &videos)
	if err != nil {
		return nil, err
	}
	return videos, nil

}

// LogVideoCall implements interfaces.ChatRepoInter.
func (c *ChatRepo) LogVideoCall(videoCall *model.VideoCall) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.Collection.InsertOne(ctx, videoCall)
	if err != nil {
		return err
	}
	return nil
}
