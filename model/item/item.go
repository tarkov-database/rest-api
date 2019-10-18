package item

import (
	"context"
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

// Collection indicates the MongoDB item collection
const Collection = "items"

func getOneByFilter(filter interface{}, k Kind) (Entity, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	e, err := k.GetEntity()
	if err != nil {
		return e, err
	}

	if err = c.FindOne(ctx, filter).Decode(e); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return e, model.MongoToAPIError(err)
	}

	return e, nil
}

// GetByID returns the entity of the given ID
func GetByID(id string, k Kind) (Entity, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	return getOneByFilter(bson.M{"_id": objID, "_kind": k}, k)
}

// Options represents the options for a database operation
type Options struct {
	Sort   map[string]int64
	Limit  int64
	Offset int64
}

func getManyByFilter(filter interface{}, k Kind, opts *Options) (*model.Result, error) {
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
		e, err := k.GetEntity()
		if err != nil {
			return r, err
		}

		if err := cur.Decode(e); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, e)
	}

	if err := cur.Err(); err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	return r, nil
}

// GetAll returns a result based on filters
func GetAll(qs map[string]interface{}, k Kind, opts *Options) (*model.Result, error) {
	qs["_kind"] = k
	return getManyByFilter(qs, k, opts)
}

// GetByIDs returns a result by given IDs
func GetByIDs(ids []string, k Kind, opts *Options) (*model.Result, error) {
	objIDs := make([]objectID, 0, len(ids))
	for _, id := range ids {
		objID, err := model.ToObjectID(id)
		if err != nil {
			return &model.Result{}, err
		}

		objIDs = append(objIDs, objID)
	}

	return getManyByFilter(bson.M{"_id": bson.M{"$in": objIDs}, "_kind": k}, k, opts)
}

// GetByText returns a result based on given keyword
func GetByText(q string, opts *Options, k ...Kind) (*model.Result, error) {
	c := database.GetDB().Collection(Collection)

	findOpts := options.Find()
	findOpts.SetLimit(opts.Limit)
	findOpts.SetSort(opts.Sort)

	var kind Kind
	if len(k) > 0 {
		kind = k[0]
	}

	if kind.IsEmpty() {
		findOpts.SetProjection(bson.M{
			"name":        1,
			"shortName":   1,
			"description": 1,
			"price":       1,
			"weight":      1,
			"maxStack":    1,
			"rarity":      1,
			"grid":        1,
			"_modified":   1,
			"_kind":       1,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r := &model.Result{}

	q = regexp.QuoteMeta(q)
	re := strings.Join(strings.Split(q, " "), ".")

	var filter bson.D

	if !kind.IsEmpty() {
		filter = bson.D{
			{"_kind", kind},
			{"$or", bson.A{
				bson.M{"shortName": primitive.Regex{fmt.Sprintf("%s", re), "gi"}},
				bson.M{"name": primitive.Regex{fmt.Sprintf("%s", re), "gi"}},
			}},
		}
	} else {
		filter = bson.D{
			{"$or", bson.A{
				bson.M{"shortName": primitive.Regex{fmt.Sprintf("%s", re), "gi"}},
				bson.M{"name": primitive.Regex{fmt.Sprintf("%s", re), "gi"}},
			}},
		}
	}

	count, err := c.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	re = strings.Join(strings.Split(q, " "), "|")

	if count == 0 {
		if !kind.IsEmpty() {
			filter = bson.D{
				{"_kind", kind},
				{"$and", bson.A{
					bson.M{"$text": bson.M{"$search": q}},
					bson.M{"description": primitive.Regex{fmt.Sprintf("(%s)", re), "gim"}},
				}},
			}
		} else {
			filter = bson.D{
				{"$and", bson.A{
					bson.M{"$text": bson.M{"$search": q}},
					bson.M{"description": primitive.Regex{fmt.Sprintf("(%s)", re), "gim"}},
				}},
			}
		}
	}

	cur, err := c.Find(ctx, filter, findOpts)
	if err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var item Entity

		if !kind.IsEmpty() {
			item, err = kind.GetEntity()
			if err != nil {
				return r, err
			}
		} else {
			item = &Item{}
		}

		if err := cur.Decode(item); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, item)
	}

	if err := cur.Err(); err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	r.Count = int64(len(r.Items))

	return r, nil
}

// Create creates a new entity
func Create(e Entity) error {
	c := database.GetDB().Collection(Collection)

	if e.GetID().IsZero() {
		e.SetID(primitive.NewObjectID())
	}

	e.SetModified(timestamp{time.Now()})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, e); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, e Entity) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if e.GetID().IsZero() {
		e.SetID(objID)
	}

	e.SetModified(timestamp{time.Now()})

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_kind": e.GetKind(), "_id": objID}
	if err := c.FindOneAndReplace(ctx, filter, e, opts).Decode(e); err != nil {
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

	if _, err := c.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}
