package item

import (
	"context"
	"fmt"
	"time"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Index describes the entity of an the item root endpoint
type Index struct {
	Total    int64               `json:"total" bson:"total"`
	Modified timestamp           `json:"modified" bson:"modified"`
	Kinds    map[Kind]*KindStats `json:"kinds" bson:"kinds"`
}

// KindStats describes the statistics of a kind
type KindStats struct {
	Count    int64     `json:"count" bson:"count"`
	Modified timestamp `json:"modified" bson:"modified"`
}

// WithKinds fills Index with kind data
func (i *Index) WithKinds(c *mongo.Collection) error {
	var err error

	s1 := bson.M{"_modified": -1}
	s2 := make(map[string]bson.A)
	s3 := make(map[string]interface{})

	s2["allCount"] = bson.A{bson.M{"$count": "count"}}
	s2["allDate"] = bson.A{
		bson.M{"$limit": 1},
		bson.M{"$project": bson.M{"modified": "$_modified"}},
	}
	s3["total"] = bson.M{"$arrayElemAt": bson.A{"$allCount.count", 0}}
	s3["modified"] = bson.M{"$arrayElemAt": bson.A{"$allDate.modified", 0}}

	kinds := make(map[Kind]bson.M)
	for _, kind := range KindList {
		count := fmt.Sprintf("%sCount", kind)
		date := fmt.Sprintf("%sDate", kind)
		s2[count] = bson.A{
			bson.M{"$match": bson.M{"_kind": kind}},
			bson.M{"$count": "count"},
		}
		s2[date] = bson.A{
			bson.M{"$match": bson.M{"_kind": kind}},
			bson.M{"$limit": 1},
			bson.M{"$project": bson.M{"modified": "$_modified"}},
		}
		kinds[kind] = bson.M{
			"count":    bson.M{"$arrayElemAt": bson.A{fmt.Sprintf("$%s.count", count), 0}},
			"modified": bson.M{"$arrayElemAt": bson.A{fmt.Sprintf("$%s.modified", date), 0}},
		}
	}

	s3["kinds"] = kinds

	pipeline := mongo.Pipeline{{{"$sort", s1}}, {{"$facet", s2}}, {{"$project", s3}}}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := c.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		if err := cur.Decode(i); err != nil {
			return model.MongoToAPIError(err)
		}
	}

	if err := cur.Err(); err != nil {
		return model.MongoToAPIError(err)
	}

	return nil
}

// WithoutKinds fills Index without kind data
func (i *Index) WithoutKinds(c *mongo.Collection) error {
	findOpts := options.Find()
	findOpts.SetSort(bson.M{"_modified": -1})
	findOpts.SetLimit(1)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error

	i.Total, err = c.EstimatedDocumentCount(ctx)
	if err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	if i.Total == 0 {
		return nil
	}

	cur, err := c.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err)
		}
		return model.MongoToAPIError(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		item := &Item{}

		if err := cur.Decode(item); err != nil {
			logger.Error(err)
			return model.MongoToAPIError(err)
		}

		i.Modified = item.Modified
	}

	if err := cur.Err(); err != nil {
		logger.Error(err)
		return model.MongoToAPIError(err)
	}

	return nil
}

// GetIndex returns the data of the item root endpoint
func GetIndex(skipKinds bool) (*Index, error) {
	db := database.GetDB()
	c := db.Collection(Collection)

	var err error

	index := &Index{}

	if skipKinds {
		err = index.WithoutKinds(c)
	} else {
		err = index.WithKinds(c)
	}

	return index, err
}
