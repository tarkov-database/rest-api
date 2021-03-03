package module

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
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

// Module describes the entity of a module
type Module struct {
	ID            objectID  `json:"_id" bson:"_id"`
	Name          string    `json:"name" bson:"name"`
	RequiresPower bool      `json:"requiresPower" bson:"requiresPower"`
	Stages        []Stage   `json:"stages" bson:"stages"`
	Modified      timestamp `json:"_modified" bson:"_modified"`
}

// Validate validates the fields of a module
func (m Module) Validate() error {
	if len(m.Stages) == 0 {
		return errors.New("stages missing")
	}

	for i, v := range m.Stages {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error in stages index \"%v\": %s", i, err)
		}
	}

	return nil
}

// Stage describes a stage of a module
type Stage struct {
	Description      string        `json:"description" bson:"description"`
	Bonuses          []Bonus       `json:"bonuses" bson:"bonuses"`
	Requirements     []Requirement `json:"requirements" bson:"requirements"`
	RequiredModules  []Ref         `json:"requiredMods" bson:"requiredMods"`
	Materials        []ItemRef     `json:"materials" bson:"materials"`
	ConstructionTime int64         `json:"constructionTime" bson:"constructionTime"`
}

// Validate validates the fields of a stage
func (s Stage) Validate() error {
	for i, v := range s.Requirements {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error in requirements index \"%v\": %s", i, err)
		}
	}

	for i, v := range s.RequiredModules {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error in required modules index \"%v\": %s", i, err)
		}
	}

	for i, v := range s.Materials {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error in materials index \"%v\": %s", i, err)
		}
	}

	return nil
}

// Bonus describes a bonus of a stage
type Bonus struct {
	Description string  `json:"description" bson:"description"`
	Value       float64 `json:"value,omitempty" bson:"value,omitempty"`
	Type        string  `json:"type" bson:"type"`
}

// Validate validates the fields of a module bonus
func (b Bonus) Validate() error {
	if b.Type == "" {
		return errors.New("type is missing")
	}
	if len(b.Description) < 3 {
		return errors.New("name is too short or not set")
	}

	return nil
}

// Requirement describes a requirement of different types of a stage
type Requirement struct {
	Name  string `json:"name" bson:"name"`
	Level uint8  `json:"level,omitempty" bson:"level,omitempty"`
	Type  string `json:"type" bson:"type"`
}

// Validate validates the fields of a module requirement
func (r Requirement) Validate() error {
	if r.Type == "" {
		return errors.New("type is missing")
	}
	if len(r.Name) < 3 {
		return errors.New("name is too short or not set")
	}

	return nil
}

// Ref refers to a module and its stage
type Ref struct {
	ID    objectID `json:"id" bson:"id"`
	Stage uint8    `json:"stage" bson:"stage"`
}

// Validate validates the fields of a referred module
func (m Ref) Validate() error {
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

// Validate validates the fields of an referred item
func (i ItemRef) Validate() error {
	if i.ID.IsZero() {
		return errors.New("id is missing")
	}
	if i.Kind.IsEmpty() {
		return errors.New("kind is missing")
	}
	if i.Count == 0 {
		return errors.New("count is zero")
	}

	return nil
}

// Collection indicates the MongoDB module collection
const Collection = "modules"

func getOneByFilter(filter interface{}) (*Module, error) {
	c := database.GetDB().Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mod := &Module{}

	if err := c.FindOne(ctx, filter).Decode(mod); err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return mod, model.MongoToAPIError(err)
	}

	return mod, nil
}

// GetByID returns the entity of the given ID
func GetByID(id string) (*Module, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &Module{}, err
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
		mod := &Module{}

		if err := cur.Decode(mod); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, mod)
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

// GetByMaterial returns a result based on stage materials
func GetByMaterial(id string, opts *Options) (*model.Result, error) {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return &model.Result{}, err
	}

	return getManyByFilter(bson.D{{Key: "stages.materials.id", Value: objID}}, opts)
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

	filter = bson.M{"name": primitive.Regex{Pattern: fmt.Sprintf("%s", re), Options: "gi"}}

	count, err := c.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	re = strings.Join(strings.Split(q, " "), "|")

	if count == 0 {
		filter = bson.D{
			{Key: "$and", Value: bson.A{
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
		mod := &Module{}

		if err := cur.Decode(mod); err != nil {
			logger.Error(err)
			return r, model.MongoToAPIError(err)
		}

		r.Items = append(r.Items, mod)
	}

	if err := cur.Err(); err != nil {
		logger.Error(err)
		return r, model.MongoToAPIError(err)
	}

	r.Count = int64(len(r.Items))

	return r, nil
}

// Create creates a new entity
func Create(mod *Module) error {
	c := database.GetDB().Collection(Collection)

	if mod.ID.IsZero() {
		mod.ID = primitive.NewObjectID()
	}

	mod.Modified = timestamp{Time: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := c.InsertOne(ctx, mod); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// Replace replaces the data of an existing entity
func Replace(id string, mod *Module) error {
	objID, err := model.ToObjectID(id)
	if err != nil {
		return err
	}

	if mod.ID.IsZero() {
		mod.ID = objID
	}

	mod.Modified = timestamp{Time: time.Now()}

	c := database.GetDB().Collection(Collection)

	opts := options.FindOneAndReplace()
	opts.SetUpsert(false)
	opts.SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = c.FindOneAndReplace(ctx, bson.M{"_id": objID}, mod, opts).Decode(mod); err != nil {
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
