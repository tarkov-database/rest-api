package user

import (
	"context"
	"errors"
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

var (
	// ErrInvalidEmail indicates that a email address is not valid
	ErrInvalidEmail = errors.New("invalid e-mail address")
)

type objectID = model.ObjectID

type timestamp = model.Timestamp

// User describes the entity of a user
type User struct {
	ID       objectID  `json:"_id" bson:"_id"`
	Email    string    `json:"email" bson:"email"`
	Locked   bool      `json:"locked" bson:"locked"`
	Modified timestamp `json:"_modified" bson:"_modified"`
}

// Validate validates the fields of a user
func (u *User) Validate() error {
	if len(u.Email) < 8 || !strings.Contains(u.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

// Collection indicates the MongoDB user collection
const Collection = "users"

func getOneByFilter(filter interface{}) (*User, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	u := &User{}

	if err := c.FindOne(ctx, filter).Decode(u); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return u, model.MongoToAPIError(err)
	}

	return u, nil
}

// GetByID returns the entity of the given ID
func GetByID(id string) (*User, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &User{}, err
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
		user := &User{}

		if err := cur.Decode(user); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, user)
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

// GetByEmail returns a result based on email address
func GetByEmail(addr string, opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.M{"email": addr}, opts)
}

// GetByLockedState returns a result based on lock state
func GetByLockedState(locked bool, opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.M{"locked": locked}, opts)
}

// Create creates a new entity
func Create(user *User) error {
	c := database.GetDB().Collection(Collection)

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	user.Modified = timestamp{Time: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, user); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, user *User) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if user.ID.IsZero() {
		user.ID = objID
	}

	user.Modified = timestamp{Time: time.Now()}

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, user, opts).Decode(user); err != nil {
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
