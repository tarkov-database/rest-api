package database

import (
	"context"
	"time"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

// Init initiate the MongoDB connection
func Init() {
	logger.Info("Initiate MongoDB connection\n")

	config := newConfig()
	clientOptions := config.getOptions()

	var err error

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal(err)
	}
	cancel()

	db = client.Database(config.Database)

	logger.Info("Successful connected to MongoDB server(s)\n")
}

// GetDB returns a handle to the database
func GetDB() *mongo.Database {
	return db
}
