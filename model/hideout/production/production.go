package production

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/item"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type objectID = model.ObjectID

type timestamp = model.Timestamp

// Production describes the entity of a production
type Production struct {
	ID              objectID    `json:"_id" bson:"_id"`
	Module          objectID    `json:"module" bson:"module"`
	RequiredModules []ModuleRef `json:"requiredMods" bson:"requiredMods"`
	Materials       []ItemRef   `json:"materials" bson:"materials"`
	Tools           []ItemRef   `json:"tools" bson:"tools"`
	Outcome         []ItemRef   `json:"outcome" bson:"outcome"`
	Duration        int64       `json:"duration" bson:"duration"`
	Modified        timestamp   `json:"_modified" bson:"_modified"`
}

// Validate validates the fields of a production
func (p Production) Validate() error {
	if p.Module.IsZero() {
		return errors.New("module id is missing")
	}
	if len(p.Outcome) == 0 {
		return errors.New("outcome missing")
	}

	for i, v := range p.Tools {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error in tools index \"%v\": %s", i, err)
		}
	}

	for i, v := range p.Materials {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error in materials index \"%v\": %s", i, err)
		}
	}

	for i, v := range p.Outcome {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error in outcome index \"%v\": %s", i, err)
		}
	}

	return nil
}

// ModuleRef refers to a module and its stage
type ModuleRef struct {
	ID    objectID `json:"id" bson:"id"`
	Stage uint8    `json:"stage" bson:"stage"`
}

// Validate validates the fields of a module
func (m ModuleRef) Validate() error {
	if m.ID.IsZero() {
		return errors.New("id is missing")
	}

	return nil
}

// ItemRef refers to an item and specifies its quantity
type ItemRef struct {
	ID        objectID  `json:"id" bson:"id"`
	Count     uint64    `json:"count,omitempty" bson:"count,omitempty"`
	Resources uint64    `json:"resources,omitempty" bson:"resources,omitempty"`
	Kind      item.Kind `json:"kind" bson:"kind"`
}

// Validate validates the fields of an item
func (i ItemRef) Validate() error {
	if i.ID.IsZero() {
		return errors.New("id is missing")
	}
	if i.Kind.IsEmpty() {
		return errors.New("kind is missing")
	}
	if i.Count == 0 && i.Resources == 0 {
		return errors.New("count and resources is zero")
	}

	return nil
}

// Collection indicates the MongoDB production collection
const Collection = "production"

func getOneByFilter(filter interface{}) (*Production, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prod := &Production{}

	if err := c.FindOne(ctx, filter).Decode(prod); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return prod, model.MongoToAPIError(err)
	}

	return prod, nil
}

// GetByID returns the entity of the given ID
func GetByID(id string) (*Production, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &Production{}, err
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
		prod := &Production{}

		if err := cur.Decode(prod); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, prod)
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

// GetByIDs returns a result by given IDs
func GetByIDs(ids []string, opts *Options) (*model.Result, error) {
	objIDs := make([]objectID, len(ids))
	for i, id := range ids {
		objID, err := model.ToObjectID(id)
		if err != nil {
			return &model.Result{}, err
		}

		objIDs[i] = objID
	}

	return getManyByFilter(bson.M{"_id": bson.M{"$in": objIDs}}, opts)
}

// GetByModule returns a result based on module
func GetByModule(id string, opts *Options) (*model.Result, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.M{"module": objID}, opts)
}

// GetByMaterial returns a result based on material
func GetByMaterial(id string, opts *Options) (*model.Result, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.D{{Key: "materials.id", Value: objID}}, opts)
}

// GetByOutcome returns a result based on outcome
func GetByOutcome(id string, opts *Options) (*model.Result, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.D{{Key: "outcome.id", Value: objID}}, opts)
}

// Create creates a new entity
func Create(prod *Production) error {
	c := database.GetDB().Collection(Collection)

	if prod.ID.IsZero() {
		prod.ID = primitive.NewObjectID()
	}

	prod.Modified = timestamp{Time: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, prod); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, prod *Production) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if prod.ID.IsZero() {
		prod.ID = objID
	}

	prod.Modified = timestamp{Time: time.Now()}

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, prod, opts).Decode(prod); err != nil {
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
