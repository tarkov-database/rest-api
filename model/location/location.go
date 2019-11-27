package location

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

// Location describes the entity of a location
type Location struct {
	ID             objectID  `json:"_id" bson:"_id"`
	Name           string    `json:"name" bson:"name"`
	Description    string    `json:"description" bson:"description"`
	MinimumPlayers int64     `json:"minPlayers" bson:"minPlayers"`
	MaximumPlayers int64     `json:"maxPlayers" bson:"maxPlayers"`
	EscapeTime     int64     `json:"escapeTime" bson:"escapeTime"`
	Insurance      bool      `json:"insurance" bson:"insurance"`
	Available      bool      `json:"available" bson:"available"`
	Exits          []Exit    `json:"exits" bson:"exits"`
	Bosses         []Boss    `json:"bosses" bson:"bosses"`
	Modified       timestamp `json:"_modified" bson:"_modified"`
}

// Validate validates the fields of a location
func (l *Location) Validate() error {
	if len(l.Name) < 3 {
		return errors.New("name is too short or not set")
	}
	if len(l.Description) < 8 {
		return errors.New("description is too short or not set")
	}
	if l.MinimumPlayers < 1 {
		return errors.New("minimum player count is too low")
	}
	if l.MaximumPlayers < 2 {
		return errors.New("maximum player count is too low")
	}

	return nil
}

// Exit describes an exit of a location
type Exit struct {
	Name             string  `json:"name" bson:"name"`
	Description      string  `json:"description" bson:"description"`
	Chance           float64 `json:"chance" bson:"chance"`
	MinimumTime      int64   `json:"minTime" bson:"minTime"`
	MaximumTime      int64   `json:"maxTime" bson:"maxTime"`
	ExfiltrationTime int64   `json:"exfilTime" bson:"exfilTime"`
	Requirement      string  `json:"requirement,omitempty" bson:"requirement,omitempty"`
}

// Boss describes a boss of a location
type Boss struct {
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Chance      float64 `json:"chance" bson:"chance"`
	Followers   int64   `json:"followers" bson:"followers"`
}

// Collection indicates the MongoDB location collection
const Collection = "location"

func getOneByFilter(filter interface{}) (*Location, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	loc := &Location{}

	if err := c.FindOne(ctx, filter).Decode(loc); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return loc, model.MongoToAPIError(err)
	}

	return loc, nil
}

// GetByID returns the entity of the given ID
func GetByID(id string) (*Location, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &Location{}, err
	}

	return getOneByFilter(bson.M{"_id": objID})
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
		loc := &Location{}

		if err := cur.Decode(loc); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, loc)
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

// GetByAvailability returns a result based on availability
func GetByAvailability(a bool, opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.M{"available": a}, opts)
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
				bson.M{"description": primitive.Regex{fmt.Sprintf("(%s)", re), "gim"}},
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
		loc := &Location{}

		if err := cur.Decode(loc); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, loc)
	}

	if err := cur.Err(); err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	r.Count = int64(len(r.Items))

	return r, nil
}

// Create creates a new entity
func Create(loc *Location) error {
	c := database.GetDB().Collection(Collection)

	if loc.ID.IsZero() {
		loc.ID = primitive.NewObjectID()
	}

	loc.Modified = timestamp{time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, loc); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, loc *Location) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if loc.ID.IsZero() {
		loc.ID = objID
	}

	loc.Modified = timestamp{time.Now()}

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, loc, opts).Decode(loc); err != nil {
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
