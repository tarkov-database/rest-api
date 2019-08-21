package controller

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/item"
	"github.com/tarkov-database/rest-api/model/user"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const contentTypeJSON = "application/json"

var (
	itemIDs []primitive.ObjectID
	userIDs []primitive.ObjectID
)

func init() {
	logger.Init("default", false, false, ioutil.Discard)
}

func mongoStartup() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}

	createItems()
	createUsers()
}

func mongoCleanup() {
	removeItems()
	removeUsers()

	if err := database.Shutdown(); err != nil {
		log.Fatalf("Database shutdown error: %s", err)
	}
}

func createItemID() primitive.ObjectID {
	id := primitive.NewObjectID()
	itemIDs = append(itemIDs, id)

	return id
}

func removeItemID(id primitive.ObjectID) {
	new := make([]primitive.ObjectID, 0, len(itemIDs)-1)
	for _, k := range itemIDs {
		if k != id {
			new = append(new, k)
		}
	}

	itemIDs = new
}

func createUserID() primitive.ObjectID {
	id := primitive.NewObjectID()
	userIDs = append(userIDs, id)

	return id
}

func removeUserID(id primitive.ObjectID) {
	new := make([]primitive.ObjectID, 0, len(userIDs)-1)
	for _, k := range userIDs {
		if k != id {
			new = append(new, k)
		}
	}

	userIDs = new
}

func createItems() {
	c := database.GetDB().Collection(item.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	itemA := item.Item{
		ID:       createItemID(),
		Name:     "item a",
		Modified: model.Timestamp{time.Now()},
		Kind:     "common",
	}
	itemB := item.Item{
		ID:       createItemID(),
		Name:     "item b",
		Modified: model.Timestamp{time.Now()},
		Kind:     "common",
	}

	if _, err := c.InsertMany(ctx, bson.A{itemA, itemB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeItems() {
	c := database.GetDB().Collection(item.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": itemIDs}}); err != nil {
		log.Fatalf("Database cleanup error: %s", err)
	}
}

func createUsers() {
	c := database.GetDB().Collection(user.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	userA := user.User{ID: createUserID()}
	userB := user.User{ID: createUserID()}

	if _, err := c.InsertMany(ctx, bson.A{userA, userB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeUsers() {
	c := database.GetDB().Collection(user.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": userIDs}}); err != nil {
		log.Fatalf("Database cleanup error: %s", err)
	}
}

func TestMain(m *testing.M) {
	mongoStartup()
	code := m.Run()
	mongoCleanup()
	os.Exit(code)
}
