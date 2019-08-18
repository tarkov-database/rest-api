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
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("connection error: %s", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("connection error: %s", err)
	}

	db = client.Database(cfg.Database)

	logger.Info("Successful connected to MongoDB server(s)\n")

	return nil
}

// Shutdown closes all sockets and shuts down the client gracefully
func Shutdown() error {
	logger.Info("Database client is shutting down...")

	client := db.Client()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("shutdown error: %s", err)
	}

	return nil
}

// GetDB returns a handle to the database
func GetDB() *mongo.Database {
	return db
}
