package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

// Init initiate the MongoDB connection
func Init() error {
	logger.Info("Initiate MongoDB connection\n")

	clientOptions, err := cfg.getClientOptions()
	if err != nil {
		return fmt.Errorf("options error: %s", err)
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return fmt.Errorf("client error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connection error: %s", err)
	}
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("connection error: %s", err)
	}
	cancel()

	db = client.Database(cfg.Database)

	logger.Info("Successful connected to MongoDB server(s)\n")

	return nil
}

// GetDB returns a handle to the database
func GetDB() *mongo.Database {
	return db
}
