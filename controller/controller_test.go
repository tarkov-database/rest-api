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
	"github.com/tarkov-database/rest-api/model/hideout/module"
	"github.com/tarkov-database/rest-api/model/hideout/production"
	"github.com/tarkov-database/rest-api/model/item"
	"github.com/tarkov-database/rest-api/model/location"
	"github.com/tarkov-database/rest-api/model/location/feature"
	"github.com/tarkov-database/rest-api/model/location/featuregroup"
	"github.com/tarkov-database/rest-api/model/statistic/ammunition/armor"
	"github.com/tarkov-database/rest-api/model/user"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const contentTypeJSON = "application/json"

var (
	itemIDs           []primitive.ObjectID
	userIDs           []primitive.ObjectID
	moduleIDs         []primitive.ObjectID
	productionIDs     []primitive.ObjectID
	locationIDs       []primitive.ObjectID
	featureIDs        []primitive.ObjectID
	featureGroupIDs   []primitive.ObjectID
	ammoArmorStatsIDs []primitive.ObjectID
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
	createModules()
	createProductions()
	createLocations()
	createFeatureGroups()
	createStatisticAmmoArmor()
	createFeatures()
}

func mongoCleanup() {
	removeItems()
	removeModules()
	removeProductions()
	removeLocations()
	removeFeatures()
	removeFeatureGroups()
	removeStatisticAmmoArmor()
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
	ids := make([]primitive.ObjectID, 0, len(itemIDs)-1)
	for _, k := range itemIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	itemIDs = ids
}

func createModuleID() primitive.ObjectID {
	id := primitive.NewObjectID()
	moduleIDs = append(moduleIDs, id)

	return id
}

func removeModuleID(id primitive.ObjectID) {
	ids := make([]primitive.ObjectID, 0, len(moduleIDs)-1)
	for _, k := range moduleIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	moduleIDs = ids
}

func createProductionID() primitive.ObjectID {
	id := primitive.NewObjectID()
	productionIDs = append(productionIDs, id)

	return id
}

func removeProductionID(id primitive.ObjectID) {
	ids := make([]primitive.ObjectID, 0, len(productionIDs)-1)
	for _, k := range productionIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	productionIDs = ids
}

func createLocationID() primitive.ObjectID {
	id := primitive.NewObjectID()
	locationIDs = append(locationIDs, id)

	return id
}

func removeLocationID(id primitive.ObjectID) {
	ids := make([]primitive.ObjectID, 0, len(locationIDs)-1)
	for _, k := range locationIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	locationIDs = ids
}

func createFeatureID() primitive.ObjectID {
	id := primitive.NewObjectID()
	featureIDs = append(featureIDs, id)

	return id
}

func removeFeatureID(id primitive.ObjectID) {
	ids := make([]primitive.ObjectID, 0, len(featureIDs)-1)
	for _, k := range featureIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	featureIDs = ids
}

func createFeatureGroupID() primitive.ObjectID {
	id := primitive.NewObjectID()
	featureGroupIDs = append(featureGroupIDs, id)

	return id
}

func removeFeatureGroupID(id primitive.ObjectID) {
	ids := make([]primitive.ObjectID, 0, len(featureGroupIDs)-1)
	for _, k := range featureGroupIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	featureGroupIDs = ids
}

func createStatisticAmmoArmorID() primitive.ObjectID {
	id := primitive.NewObjectID()
	ammoArmorStatsIDs = append(ammoArmorStatsIDs, id)

	return id
}

func removeStatisticAmmoArmorID(id primitive.ObjectID) {
	ids := make([]primitive.ObjectID, 0, len(ammoArmorStatsIDs)-1)
	for _, k := range ammoArmorStatsIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	ammoArmorStatsIDs = ids
}

func createUserID() primitive.ObjectID {
	id := primitive.NewObjectID()
	userIDs = append(userIDs, id)

	return id
}

func removeUserID(id primitive.ObjectID) {
	ids := make([]primitive.ObjectID, 0, len(userIDs)-1)
	for _, k := range userIDs {
		if k != id {
			ids = append(ids, k)
		}
	}

	userIDs = ids
}

func createItems() {
	c := database.GetDB().Collection(item.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	itemA := item.Item{
		ID:       createItemID(),
		Name:     "item a",
		Modified: model.Timestamp{Time: time.Now()},
		Kind:     item.KindCommon,
	}
	itemB := item.Item{
		ID:       createItemID(),
		Name:     "item b",
		Modified: model.Timestamp{Time: time.Now()},
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

func createModules() {
	c := database.GetDB().Collection(module.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	moduleA := module.Module{
		ID:       createModuleID(),
		Name:     "module a",
		Modified: model.Timestamp{Time: time.Now()},
	}
	moduleB := module.Module{
		ID:       createModuleID(),
		Name:     "module b",
		Modified: model.Timestamp{Time: time.Now()},
	}

	if _, err := c.InsertMany(ctx, bson.A{moduleA, moduleB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeModules() {
	c := database.GetDB().Collection(module.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": moduleIDs}}); err != nil {
		log.Fatalf("Database cleanup error: %s", err)
	}
}

func createProductions() {
	c := database.GetDB().Collection(production.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	prodA := production.Production{
		ID:       createProductionID(),
		Modified: model.Timestamp{Time: time.Now()},
	}
	prodB := production.Production{
		ID:       createProductionID(),
		Modified: model.Timestamp{Time: time.Now()},
	}

	if _, err := c.InsertMany(ctx, bson.A{prodA, prodB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeProductions() {
	c := database.GetDB().Collection(production.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": productionIDs}}); err != nil {
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
		Modified: model.Timestamp{Time: time.Now()},
	}
	locationB := location.Location{
		ID:       createLocationID(),
		Name:     "location b",
		Modified: model.Timestamp{Time: time.Now()},
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
		Modified: model.Timestamp{Time: time.Now()},
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
		Modified: model.Timestamp{Time: time.Now()},
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
		Modified:    model.Timestamp{Time: time.Now()},
	}
	groupB := featuregroup.Group{
		ID:          createFeatureGroupID(),
		Name:        "group b",
		Description: "description of b",
		Tags:        []string{"test"},
		Location:    locationIDs[0],
		Modified:    model.Timestamp{Time: time.Now()},
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

func createStatisticAmmoArmor() {
	c := database.GetDB().Collection(armor.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	statsA := armor.AmmoArmorStatistics{
		ID:   createStatisticAmmoArmorID(),
		Ammo: primitive.NewObjectID(),
		Armor: armor.ItemRef{
			ID:   primitive.NewObjectID(),
			Kind: item.KindTacticalrig,
		},
		Distance:                  100,
		PenetrationChance:         [4]float64{},
		AverageShotsToDestruction: armor.Statistics{},
		AverageShotsTo50Damage:    armor.Statistics{},
		Modified:                  model.Timestamp{Time: time.Now()},
	}
	statsB := armor.AmmoArmorStatistics{
		ID:   createStatisticAmmoArmorID(),
		Ammo: primitive.NewObjectID(),
		Armor: armor.ItemRef{
			ID:   primitive.NewObjectID(),
			Kind: item.KindArmor,
		},
		Distance:                  500,
		PenetrationChance:         [4]float64{},
		AverageShotsToDestruction: armor.Statistics{},
		AverageShotsTo50Damage:    armor.Statistics{},
		Modified:                  model.Timestamp{Time: time.Now()},
	}

	if _, err := c.InsertMany(ctx, bson.A{statsA, statsB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func removeStatisticAmmoArmor() {
	c := database.GetDB().Collection(armor.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ammoArmorStatsIDs}}); err != nil {
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
