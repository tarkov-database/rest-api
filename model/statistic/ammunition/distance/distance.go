package distance

import (
	"context"
	"errors"
	"time"

	"github.com/google/logger"
	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type objectID = model.ObjectID

type timestamp = model.Timestamp

// AmmoDistanceStatistics describes the entity of a ammo distance statistics
type AmmoDistanceStatistics struct {
	ID         objectID               `json:"_id" bson:"_id"`
	Reference  objectID               `json:"ammo" bson:"ammo"`
	Distance   uint64                 `json:"distance" bson:"distance"`
	Properties AmmoDistanceProperties `json:"properties" bson:"properties"`
	Modified   timestamp              `json:"_modified" bson:"_modified"`
}

// Validate validates the fields of a DistanceStatistics
func (d AmmoDistanceStatistics) Validate() error {
	if d.Reference.IsZero() {
		return errors.New("ammo reference id is missing")
	}

	return nil
}

// AmmoDistanceProperties describes the properties of a ammo distance statistics
type AmmoDistanceProperties struct {
	Velocity         float64 `json:"velocity" bson:"velocity"`
	Damage           float64 `json:"damage" bson:"damage"`
	PenetrationPower float64 `json:"penetrationPower" bson:"penetrationPower"`
	TimeOfFlight     float64 `json:"timeOfFlight" bson:"timeOfFlight"`
	Drop             float64 `json:"drop" bson:"drop"`
}

// Collection indicates the MongoDB feature collection
const Collection = "statistics.ammunition.distances"

func getOneByFilter(filter interface{}) (*AmmoDistanceStatistics, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := &AmmoDistanceStatistics{}

	if err := c.FindOne(ctx, filter).Decode(stats); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return stats, model.MongoToAPIError(err)
	}

	return stats, nil
}

// GetByID returns the entity of the given ID
func GetByID(id string) (*AmmoDistanceStatistics, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &AmmoDistanceStatistics{}, err
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
		stats := &AmmoDistanceStatistics{}

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

// GetAll returns a result based on filters
func GetAll(opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.D{}, opts)
}

// GetByRefsAndRange returns a result by given IDs and range
func GetByRefsAndRange(ids []string, gte, lte *uint64, opts *Options) (*model.Result, error) {
	opts.Sort = append(opts.Sort, bson.D{
		bson.E{Key: "ammo", Value: 1},
		bson.E{Key: "armor.id", Value: 1},
	}...)

	filter := bson.D{}

	if ids != nil {
		objIDs := make([]objectID, len(ids))
		for i, id := range ids {
			objID, err := model.ToObjectID(id)
			if err != nil {
				return &model.Result{}, err
			}

			objIDs[i] = objID
		}

		filter = append(filter, bson.E{Key: "ammo", Value: bson.D{{Key: "$in", Value: objIDs}}})
	}

	switch {
	case gte != nil && lte != nil:
		filter = append(filter, bson.E{Key: "distance", Value: bson.D{
			bson.E{Key: "$gte", Value: gte}, bson.E{Key: "$lte", Value: lte},
		}})
	case gte != nil:
		filter = append(filter, bson.E{Key: "distance", Value: bson.D{{Key: "$gte", Value: gte}}})
	case lte != nil:
		filter = append(filter, bson.E{Key: "distance", Value: bson.D{{Key: "$lte", Value: lte}}})
	}

	return getManyByFilter(filter, opts)
}

// Create creates a new entity
func Create(stats *AmmoDistanceStatistics) error {
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
func Replace(id string, stats *AmmoDistanceStatistics) error {
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
