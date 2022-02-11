package armor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/logger"
	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/item"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type objectID = model.ObjectID

type timestamp = model.Timestamp

// AmmoArmorStatistics describes the entity of a ammo against armor statistics
type AmmoArmorStatistics struct {
	ID                        objectID   `json:"_id" bson:"_id"`
	Ammo                      objectID   `json:"ammo" bson:"ammo"`
	Armor                     ItemRef    `json:"armor" bson:"armor"`
	Distance                  uint64     `json:"distance" bson:"distance"`
	PenetrationChance         [4]float64 `json:"penetrationChance" bson:"penetrationChance"`
	AverageShotsToDestruction Statistics `json:"avgShotsToDestruct" bson:"avgShotsToDestruct"`
	AverageShotsTo50Damage    Statistics `json:"avgShotsTo50Damage" bson:"avgShotsTo50Damage"`
	Modified                  timestamp  `json:"_modified" bson:"_modified"`
}

// Validate validates the fields of a ArmorStatistics
func (d AmmoArmorStatistics) Validate() error {
	if d.Ammo.IsZero() {
		return errors.New("ammo id is missing")
	}
	if err := d.Armor.Validate(); err != nil {
		return fmt.Errorf("armor reference is invalid: %w", err)
	}

	return nil
}

// ItemRef refers to an item entity
type ItemRef struct {
	ID   objectID  `json:"id" bson:"id"`
	Kind item.Kind `json:"kind" bson:"kind"`
}

// Validate validates the fields of a ItemRef
func (d ItemRef) Validate() error {
	if d.ID.IsZero() {
		return errors.New("item reference id is missing")
	}
	if !d.Kind.IsValid() {
		return errors.New("item reference kind is missing")
	}

	return nil
}

// Statistics describes the statistical values
type Statistics struct {
	Min    float64 `json:"min" bson:"min"`
	Max    float64 `json:"max" bson:"max"`
	Mean   float64 `json:"mean" bson:"mean"`
	Median float64 `json:"median" bson:"median"`
	StdDev float64 `json:"stdDev" bson:"stdDev"`
}

// Collection indicates the MongoDB feature collection
const Collection = "statistics.ammunition.armor"

func getOneByFilter(filter interface{}) (*AmmoArmorStatistics, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := &AmmoArmorStatistics{}

	if err := c.FindOne(ctx, filter).Decode(stats); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return stats, model.MongoToAPIError(err)
	}

	return stats, nil
}

// GetByID returns the entity of the given ID
func GetByID(id string) (*AmmoArmorStatistics, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &AmmoArmorStatistics{}, err
	}

	return getOneByFilter(bson.M{"_id": objID})
}

// Options represents the options for a database operation
type Options struct {
	Sort   bson.D
	Limit  int64
	Offset int64
}

func getManyByFilter(filter interface{}, opts *Options) (*model.Result, error) {
	c := database.GetDB().Collection(Collection)

	findOpts := options.Find()
	findOpts.SetLimit(opts.Limit)
	findOpts.SetSkip(opts.Offset)
	findOpts.SetSort(opts.Sort)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	r := &model.Result{}

	r.Count, err = c.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	if r.Count == 0 {
		return r, nil
	}

	cur, err := c.Find(ctx, filter, findOpts)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return r, model.MongoToAPIError(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		stats := &AmmoArmorStatistics{}

		if err := cur.Decode(stats); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, stats)
	}

	if err := cur.Err(); err != nil {
		return r, model.MongoToAPIError(err)
	}

	return r, nil
}

// RangeOptions represents the range options of a query
type RangeOptions struct {
	GTE *uint64
	LTE *uint64
}

// GetAll returns a result based on filters
func GetAll(opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.D{}, opts)
}

// GetByRefs returns a result by given ammo and armor IDs
func GetByRefs(ammo, armor []string, r *RangeOptions, opts *Options) (*model.Result, error) {
	opts.Sort = append(opts.Sort, bson.D{
		bson.E{Key: "ammo", Value: 1},
		bson.E{Key: "armor.id", Value: 1},
	}...)

	filter := bson.D{}

	if ammo != nil {
		IDs := make([]objectID, len(ammo))
		for i, id := range ammo {
			objID, err := model.ToObjectID(id)
			if err != nil {
				return &model.Result{}, err
			}

			IDs[i] = objID
		}

		filter = append(filter, bson.E{Key: "ammo", Value: bson.D{{Key: "$in", Value: IDs}}})
	}
	if armor != nil {
		IDs := make([]objectID, len(armor))
		for i, id := range armor {
			objID, err := model.ToObjectID(id)
			if err != nil {
				return &model.Result{}, err
			}

			IDs[i] = objID
		}

		filter = append(filter, bson.E{Key: "armor.id", Value: bson.D{{Key: "$in", Value: IDs}}})
	}

	switch {
	case r.GTE != nil && r.LTE != nil:
		filter = append(filter, bson.E{Key: "distance", Value: bson.D{
			bson.E{Key: "$gte", Value: r.GTE}, bson.E{Key: "$lte", Value: r.LTE},
		}})
	case r.GTE != nil:
		filter = append(filter, bson.E{Key: "distance", Value: bson.D{{Key: "$gte", Value: r.GTE}}})
	case r.LTE != nil:
		filter = append(filter, bson.E{Key: "distance", Value: bson.D{{Key: "$lte", Value: r.LTE}}})
	}

	return getManyByFilter(filter, opts)
}

// Create creates a new entity
func Create(stats *AmmoArmorStatistics) error {
	c := database.GetDB().Collection(Collection)

	if stats.ID.IsZero() {
		stats.ID = primitive.NewObjectID()
	}

	stats.Modified = timestamp{Time: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, stats); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, stats *AmmoArmorStatistics) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if stats.ID.IsZero() {
		stats.ID = objID
	}

	stats.Modified = timestamp{Time: time.Now()}

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, stats, opts).Decode(stats); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Remove removes an entity
func Remove(id string) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err = c.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}
