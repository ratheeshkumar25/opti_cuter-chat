package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ratheeshkumar25/opti_cut_chat_service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB(config *config.Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.DBurl)
	log.Printf("Attempting to connect to MongoDB at: %s", config.DBurl)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoclient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("MongoDB connection established")
	return mongoclient, nil
}
