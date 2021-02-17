package controller

import (
	"context"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/item"
	"github.com/tarkov-database/rest-api/model/location"
	"github.com/tarkov-database/rest-api/model/location/feature"
	"github.com/tarkov-database/rest-api/model/location/featuregroup"
	"github.com/tarkov-database/rest-api/model/user"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const contentTypeJSON = "application/json"

var (
	itemIDs         []primitive.ObjectID
	userIDs         []primitive.ObjectID
	locationIDs     []primitive.ObjectID
	featureIDs      []primitive.ObjectID
	featureGroupIDs []primitive.ObjectID
)

func init() {
	logger.Init("default", false, false, io.Discard)
}

func mongoStartup() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}

	createUsers()
	createItems()
	createLocations()
	createFeatureGroups()
	createFeatures()
}

func mongoCleanup() {
	removeItems()
	removeLocations()
	removeFeatures()
	removeFeatureGroups()
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

func createLocationID() primitive.ObjectID {
	id := primitive.NewObjectID()
	locationIDs = append(locationIDs, id)

	return id
}

func removeLocationID(id primitive.ObjectID) {
	new := make([]primitive.ObjectID, 0, len(locationIDs)-1)
	for _, k := range locationIDs {
		if k != id {
			new = append(new, k)
		}
	}

	locationIDs = new
}

func createFeatureID() primitive.ObjectID {
	id := primitive.NewObjectID()
	featureIDs = append(featureIDs, id)

	return id
}

func removeFeatureID(id primitive.ObjectID) {
	new := make([]primitive.ObjectID, 0, len(featureIDs)-1)
	for _, k := range featureIDs {
		if k != id {
			new = append(new, k)
		}
	}

	featureIDs = new
}

func createFeatureGroupID() primitive.ObjectID {
	id := primitive.NewObjectID()
	featureGroupIDs = append(featureGroupIDs, id)

	return id
}

func removeFeatureGroupID(id primitive.ObjectID) {
	new := make([]primitive.ObjectID, 0, len(featureGroupIDs)-1)
	for _, k := range featureGroupIDs {
		if k != id {
			new = append(new, k)
		}
	}

	featureGroupIDs = new
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
		Kind:     item.KindCommon,
	}
	itemB := item.Item{
		ID:       createItemID(),
		Name:     "item b",
		Modified: model.Timestamp{time.Now()},
		Kind:     item.KindCommon,
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

func createLocations() {
	c := database.GetDB().Collection(location.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	locationA := location.Location{
		ID:       createLocationID(),
		Name:     "location a",
		Modified: model.Timestamp{time.Now()},
	}
	locationB := location.Location{
		ID:       createLocationID(),
		Name:     "location b",
		Modified: model.Timestamp{time.Now()},
	}

	if _, err := c.InsertMany(ctx, bson.A{locationA, locationB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeLocations() {
	c := database.GetDB().Collection(location.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": locationIDs}}); err != nil {
		log.Fatalf("Database cleanup error: %s", err)
	}
}

func createFeatures() {
	c := database.GetDB().Collection(feature.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	featureA := feature.Feature{
		ID:    createFeatureID(),
		Name:  "feature a",
		Group: featureGroupIDs[0],
		Geometry: feature.Geometry{
			Type:        feature.Point,
			Coordinates: createFeatureCoords(),
		},
		Location: locationIDs[0],
		Modified: model.Timestamp{time.Now()},
	}
	featureB := feature.Feature{
		ID:    createFeatureID(),
		Name:  "feature b",
		Group: featureGroupIDs[0],
		Geometry: feature.Geometry{
			Type:        feature.Point,
			Coordinates: createFeatureCoords(),
		},
		Location: locationIDs[0],
		Modified: model.Timestamp{time.Now()},
	}

	if _, err := c.InsertMany(ctx, bson.A{featureA, featureB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeFeatures() {
	c := database.GetDB().Collection(feature.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": featureIDs}}); err != nil {
		log.Fatalf("Database cleanup error: %s", err)
	}
}

func createFeatureCoords() feature.Coordinates {
	c := []float64{0, 0}
	fc := make(feature.Coordinates, len(c))
	for i, v := range c {
		fc[i] = v
	}

	return fc
}

func createFeatureGroups() {
	c := database.GetDB().Collection(featuregroup.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	groupA := featuregroup.Group{
		ID:          createFeatureGroupID(),
		Name:        "group a",
		Description: "description of a",
		Tags:        []string{"test"},
		Location:    locationIDs[0],
		Modified:    model.Timestamp{time.Now()},
	}
	groupB := featuregroup.Group{
		ID:          createFeatureGroupID(),
		Name:        "group b",
		Description: "description of b",
		Tags:        []string{"test"},
		Location:    locationIDs[0],
		Modified:    model.Timestamp{time.Now()},
	}

	if _, err := c.InsertMany(ctx, bson.A{groupA, groupB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeFeatureGroups() {
	c := database.GetDB().Collection(featuregroup.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": featureGroupIDs}}); err != nil {
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
