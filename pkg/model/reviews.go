package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	MaterialID uint32             `bson:"material_id"`
	UserID     uint32             `bson:"user_id"`
	Content    string             `bson:"content"`
	Rating     uint32             `bson:"rating"`
	Timestamp  time.Time          `bson:"timestamp"`
}

type Video struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     uint32             `bson:"user_id"`
	MaterialID uint32             `bson:"material_id"`
	FileName   string             `bson:"file_name"`
	VideoURL   string             `bson:"video_url"`
	Timestamp  time.Time          `bson:"timestamp"`
	Chunks     []VideoChunk       `bson:"chunks"`
}

type VideoChunk struct {
	ChunkID    primitive.ObjectID `bson:"_id,omitempty"`
	ChunkData  []byte             `bson:"chunk_data"`
	ChunkOrder int                `bson:"chunk_order"`
}

// type Video struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty"`
// 	MaterialID uint32             `bson:"material_id"`
// 	UserID     uint32             `bson:"user_id"`
// 	FileName   string             `bson:"file_name"`
// 	VideoURL   string             `bson:"video_url"`
// 	Timestamp  time.Time          `bson:"timestamp"`
// }
