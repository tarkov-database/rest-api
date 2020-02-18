package feature

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type objectID = model.ObjectID

type timestamp = model.Timestamp

// Feature ...
type Feature struct {
	ID         objectID               `json:"_id" bson:"_id"`
	Name       string                 `json:"name" bson:"name"`
	Geometry   Geometry               `json:"geometry" bson:"geometry"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
	Location   objectID               `json:"_location" bson:"_location"`
	Group      objectID               `json:"group" bson:"group"`
	Modified   timestamp              `json:"_modified" bson:"_modified"`
}

// Validate ...
func (f Feature) Validate() error {
	if len(f.Name) < 3 {
		return errors.New("name is too short or not set")
	}

	return f.Geometry.Validate()
}

// Collection ...
const Collection = "features"

func getOneByFilter(filter interface{}) (*Feature, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ft := &Feature{}

	if err := c.FindOne(ctx, filter).Decode(ft); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return ft, model.MongoToAPIError(err)
	}

	return ft, nil
}

// GetByID returns the entity of the given ID
func GetByID(id, loc string) (*Feature, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &Feature{}, err
	}

	locID, err := model.ToObjectID(loc)
	if err != nil {
		return &Feature{}, err
	}

	return getOneByFilter(bson.M{"_id": objID, "_location": locID})
}

// Options represents the options for a database operation
type Options struct {
	Sort   map[string]int64
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
		ft := &Feature{}

		if err := cur.Decode(ft); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, ft)
	}

	if err := cur.Err(); err != nil {
		return r, model.MongoToAPIError(err)
	}

	return r, nil
}

// GetAll returns a result based on filters
func GetAll(loc string, opts *Options) (*model.Result, error) {
	locID, err := model.ToObjectID(loc)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.M{"_location": locID}, opts)
}

// GetByGroup returns a result based on location
func GetByGroup(id, loc string, opts *Options) (*model.Result, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &model.Result{}, err
	}

	locID, err := model.ToObjectID(loc)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.M{"group": objID, "_location": locID}, opts)
}

// GetByText returns a result based on given keyword
func GetByText(q string, opts *Options) (*model.Result, error) {
	c := database.GetDB().Collection(Collection)

	findOpts := options.Find()
	findOpts.SetLimit(opts.Limit)
	findOpts.SetSort(opts.Sort)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r := &model.Result{}

	q = regexp.QuoteMeta(q)
	re := strings.Join(strings.Split(q, " "), ".")

	var filter interface{}

	filter = bson.M{"name": primitive.Regex{fmt.Sprintf("%s", re), "gi"}}

	count, err := c.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	re = strings.Join(strings.Split(q, " "), "|")

	if count == 0 {
		filter = bson.D{
			{"$and", bson.A{
				bson.M{"$text": bson.M{"$search": q}},
			}},
		}
	}

	cur, err := c.Find(ctx, filter, findOpts)
	if err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		ft := &Feature{}

		if err := cur.Decode(ft); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, ft)
	}

	if err = cur.Err(); err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	r.Count = int64(len(r.Items))

	return r, nil
}

// Create creates a new entity
func Create(ft *Feature) error {
	c := database.GetDB().Collection(Collection)

	if ft.ID.IsZero() {
		ft.ID = primitive.NewObjectID()
	}

	ft.Modified = timestamp{time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, ft); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, ft *Feature) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if ft.ID.IsZero() {
		ft.ID = objID
	}

	ft.Modified = timestamp{time.Now()}

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, ft, opts).Decode(ft); err != nil {
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
