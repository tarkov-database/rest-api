package item

import (
	"errors"

	"github.com/tarkov-database/rest-api/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type objectID = primitive.ObjectID

type timestamp = model.Timestamp

type ItemList map[Kind][]objectID

const KindCommon Kind = "common"

type Item struct {
	ID          objectID  `json:"_id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	ShortName   string    `json:"shortName" bson:"shortName"`
	Description string    `json:"description" bson:"description"`
	Price       int64     `json:"price" bson:"price"`
	Weight      float64   `json:"weight" bson:"weight"`
	MaxStack    int64     `json:"maxStack" bson:"maxStack"`
	Rarity      string    `json:"rarity" bson:"rarity"`
	Grid        GridProps `json:"grid" bson:"grid"`
	Modified    timestamp `json:"_modified" bson:"_modified"`
	Kind        Kind      `json:"_kind" bson:"_kind"`
}

func (i *Item) GetID() objectID {
	return i.ID
}

func (i *Item) SetID(id objectID) {
	i.ID = id
}

func (i *Item) GetKind() Kind {
	return i.Kind
}

func (i *Item) SetKind(k Kind) {
	i.Kind = k
}

func (i *Item) GetModified() timestamp {
	return i.Modified
}

func (i *Item) SetModified(t timestamp) {
	i.Modified = t
}

func (i *Item) Validate() error {
	if len(i.Name) < 3 {
		return errors.New("name is too short or not set")
	}
	if len(i.ShortName) < 1 {
		return errors.New("short name is too short or not set")
	}
	if len(i.Description) < 8 {
		return errors.New("description is too short or not set")
	}
	if i.Price < 0 {
		return errors.New("no negative price allowed")
	}
	if i.Weight < 0 {
		return errors.New("weight is too low or not set")
	}
	if i.MaxStack < 1 {
		return errors.New("maximum stack is too low or not set")
	}
	if i.Rarity == "" {
		return errors.New("rarity is not set")
	}
	if !i.Kind.IsValid() {
		return model.ErrInvalidKind
	}

	return nil
}

type GridProps struct {
	Color  RGBA  `json:"color" bson:"color"`
	Height int64 `json:"height" bson:"height"`
	Width  int64 `json:"width" bson:"width"`
}

type RGBA struct {
	R int `json:"r" bson:"r"`
	G int `json:"g" bson:"g"`
	B int `json:"b" bson:"b"`
	A int `json:"a" bson:"a"`
}
