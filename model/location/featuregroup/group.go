package featuregroup

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

// Group describes the entity of a feature group
type Group struct {
	ID          objectID  `json:"_id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Tags        []string  `json:"tags" bson:"tags"`
	Location    objectID  `json:"_location" bson:"_location"`
	Modified    timestamp `json:"_modified" bson:"_modified"`
}

// Validate validates the fields of a feature group
func (g Group) Validate() error {
	if len(g.Name) < 3 {
		return errors.New("name is too short or not set")
	}
	if len(g.Description) < 8 {
		return errors.New("description is too short or not set")
	}
	if len(g.Tags) == 0 {
		return errors.New("no tags specified")
	}

	return nil
}

// Collection indicates the MongoDB feature group collection
const Collection = "featureGroups"

func getOneByFilter(filter interface{}) (*Group, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ft := &Group{}

	if err := c.FindOne(ctx, filter).Decode(ft); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return ft, model.MongoToAPIError(err)
	}

	return ft, nil
}

// GetByID returns the entity of the given ID
func GetByID(id, loc string) (*Group, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &Group{}, err
	}

	lID, err := model.ToObjectID(loc)
	if err != nil {
		return &Group{}, err
	}

	return getOneByFilter(bson.M{"_id": objID, "_location": lID})
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
		ft := &Group{}

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
	lID, err := model.ToObjectID(loc)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.M{"_location": lID}, opts)
}

// GetByTags returns a result based on tag
func GetByTags(tags []string, loc string, opts *Options) (*model.Result, error) {
	lID, err := model.ToObjectID(loc)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.M{"tags": bson.M{"$all": tags}, "_location": lID}, opts)
}

// GetByText returns a result based on given keyword
func GetByText(q, loc string, opts *Options) (*model.Result, error) {
	c := database.GetDB().Collection(Collection)

	lID, err := model.ToObjectID(loc)
	if err != nil {
		return &model.Result{}, err
	}

	findOpts := options.Find()
	findOpts.SetLimit(opts.Limit)
	findOpts.SetSort(opts.Sort)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r := &model.Result{}

	q = regexp.QuoteMeta(q)
	re := strings.Join(strings.Split(q, " "), ".")

	var filter bson.D

	filter = bson.D{
		{Key: "_location", Value: lID},
		{Key: "name", Value: primitive.Regex{Pattern: fmt.Sprintf("%s", re), Options: "i"}},
	}

	count, err := c.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	re = strings.Join(strings.Split(q, " "), "|")

	if count == 0 {
		filter = bson.D{
			{Key: "_location", Value: lID},
			{Key: "$and", Value: bson.A{
				bson.M{"$text": bson.M{"$search": q}},
				bson.M{"description": primitive.Regex{Pattern: fmt.Sprintf("(%s)", re), Options: "gim"}},
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
		ft := &Group{}

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
func Create(ft *Group) error {
	c := database.GetDB().Collection(Collection)

	if ft.ID.IsZero() {
		ft.ID = primitive.NewObjectID()
	}

	ft.Modified = timestamp{Time: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, ft); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, fg *Group) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if fg.ID.IsZero() {
		fg.ID = objID
	}

	fg.Modified = timestamp{Time: time.Now()}

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, fg, opts).Decode(fg); err != nil {
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
