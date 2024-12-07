package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddReview implements interfaces.ChatRepoInter.
func (c *ChatRepo) AddReview(review *model.Review) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.VideoCollection.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	return nil
}

// FindReviewMaterialID implements interfaces.ChatRepoInter.
func (c *ChatRepo) FindReviewMaterialID(materialID uint) ([]model.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query filter for material_id
	filter := bson.M{"material_id": materialID}

	// Slice to store retrieved reviews
	var reviews []model.Review
	cursor, err := c.VideoCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode reviews
	for cursor.Next(ctx) {
		var review model.Review
		if err := cursor.Decode(&review); err != nil {
			return nil, fmt.Errorf("failed to decode review: %w", err)
		}
		reviews = append(reviews, review)
	}

	// Check for cursor errors after iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return reviews, nil
}

// AddVideo adds a new video document (metadata) to the database
func (r *ChatRepo) AddVideo(video *model.Video) error {
	_, err := r.VideoCollection.InsertOne(context.TODO(), video)
	if err != nil {
		return fmt.Errorf("failed to add video metadata: %v", err)
	}
	return nil
}

// FindVideoByMaterialID fetches video documents by material ID
func (r *ChatRepo) FindVideoByMaterialID(materialID uint) ([]model.Video, error) {
	var videos []model.Video
	cursor, err := r.VideoCollection.Find(context.TODO(), bson.M{"material_id": materialID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch videos by material ID: %v", err)
	}
	if err := cursor.All(context.TODO(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode video documents: %v", err)
	}
	return videos, nil
}

// AddVideoChunk stores a chunk of the video in the database
func (c *ChatRepo) AddVideoChunk(videoID primitive.ObjectID, chunk *model.VideoChunk) error {
	// Append chunk to an existing video document
	filter := bson.M{"_id": videoID}
	update := bson.M{
		"$push": bson.M{"chunks": chunk},
	}

	_, err := c.VideoCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add video chunk: %v", err)
	}
	return nil
}

// FetchVideoChunks retrieves all chunks for a specific video
func (c *ChatRepo) FetchVideoChunks(videoID primitive.ObjectID) ([]model.VideoChunk, error) {
	video := &model.Video{}
	err := c.VideoCollection.FindOne(context.TODO(), bson.M{"_id": videoID}).Decode(video)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video chunks: %v", err)
	}
	return video.Chunks, nil
}

// FinalizeVideo updates metadata after all chunks are uploaded
func (c *ChatRepo) FinalizeVideo(videoID primitive.ObjectID, fileName string, userID uint32) error {
	filter := bson.M{"_id": videoID}
	update := bson.M{
		"$set": bson.M{
			"file_name": fileName,
			"user_id":   userID,
			"timestamp": time.Now(),
		},
	}

	_, err := c.VideoCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to finalize video metadata: %v", err)
	}
	return nil
}
