package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tarkov-database/api/core/database"
	"github.com/tarkov-database/api/model"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidEmail = errors.New("invalid e-mail address")
)

type objectID = primitive.ObjectID

type timestamp = model.Timestamp

type User struct {
	ID       objectID  `json:"_id" bson:"_id"`
	Email    string    `json:"email" bson:"email"`
	Locked   bool      `json:"locked" bson:"locked"`
	Modified timestamp `json:"_modified" bson:"_modified"`
}

func (u *User) Validate() error {
	if len(u.Email) < 8 || !strings.Contains(u.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

const collection = "users"

func getOneByFilter(filter interface{}) (*User, error) {
	db := database.GetDB()
	c := db.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	u := &User{}

	err := c.FindOne(ctx, filter).Decode(u)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return u, model.MongoToAPIError(err)
	}

	return u, nil
}

func GetByID(id string) (*User, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &User{}, err
	}

	return getOneByFilter(bson.M{"_id": objID})
}

type Options struct {
	Sort   map[string]int64
	Limit  int64
	Offset int64
}

func getManyByFilter(filter interface{}, opts *Options) (*model.Result, error) {
	db := database.GetDB()
	c := db.Collection(collection)

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

func GetAll(opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.D{}, opts)
}

func GetByEmail(addr string, opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.M{"email": addr}, opts)
}

func GetByLockedState(locked bool, opts *Options) (*model.Result, error) {
	return getManyByFilter(bson.M{"locked": locked}, opts)
}

func Create(user *User) error {
	db := database.GetDB()
	c := db.Collection(collection)

	user.ID = primitive.NewObjectID()

	user.Modified = timestamp{time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := c.InsertOne(ctx, user)
	if err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	go createIndexes(c)

	return nil
}

func Replace(id string, user *User) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if user.ID.IsZero() {
		user.ID = objID
	}

	user.Modified = timestamp{time.Now()}

	db := database.GetDB()
	c := db.Collection(collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, user, opts).Decode(user)
	if err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	go createIndexes(c)

	return nil
}

func Remove(id string) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	db := database.GetDB()
	c := db.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = c.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	go createIndexes(c)

	return nil
}

func createIndexes(c *mongo.Collection) {
	index := c.Indexes()

	indexModels := []mongo.IndexModel{}
	indexModels = append(indexModels, mongo.IndexModel{
		Keys: bson.D{{"_modified", -1}},
	})
	indexModels = append(indexModels, mongo.IndexModel{
		Keys: bson.M{"email": 1},
	})
	indexModels = append(indexModels, mongo.IndexModel{
		Keys: bson.M{"locked": 1},
	})

	_, err := index.CreateMany(context.Background(), indexModels)
	if err != nil {
		logger.Errorf("Error while creating indexes: %v", err)
	}
}
