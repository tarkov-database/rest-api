package feature

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var (
	// ErrBadGeometrySemantic indicates that a GeoJSON geometry is semantically invalid
	ErrBadGeometrySemantic = errors.New("bad geometry semantic")

	// ErrUnknownGeometryType indicates that a GeoJSON geometry type is invalid
	ErrUnknownGeometryType = errors.New("unknown geometry type")

	// ErrBadGeometryCoords indicates that GeoJSON geometry coordinates are ivnalid
	ErrBadGeometryCoords = errors.New("bad geometry coordinates")

	// ErrBadGeometryCollection indicates that a GeoJSON geometry collection are invalid
	ErrBadGeometryCollection = errors.New("empty geometry collection")
)

// GeometryType represents an GeoJSON geometry type
type GeometryType int

// Represents the GeoJSON geometry types
const (
	UnknownGeometry GeometryType = iota
	Point
	MultiPoint
	LineString
	MultiLineString
	Polygon
	MultiPolygon
	GeometryCollection
)

var geometryStrings = [...]string{
	"",
	"Point",
	"MultiPoint",
	"LineString",
	"MultiLineString",
	"Polygon",
	"MultiPolygon",
	"GeometryCollection",
}

// String returns a string representing the GeometryType
func (g GeometryType) String() string {
	return geometryStrings[g]
}

// MarshalJSON implements the JSON marshaler
func (g GeometryType) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// UnmarshalJSON implements the JSON unmarshaler
func (g *GeometryType) UnmarshalJSON(b []byte) error {
	var t string

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	for i, k := range geometryStrings {
		if k == t {
			*g = GeometryType(i)
			return nil
		}
	}

	return ErrUnknownGeometryType
}

// MarshalBSONValue implements the BSON value marshaler
func (g GeometryType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, g.String()), nil
}

// UnmarshalBSONValue implements the BSON value unmarshaler
func (g *GeometryType) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	var s string

	if t == bsontype.String {
		if str, _, ok := bsoncore.ReadString(b); ok {
			s = str
		} else {
			return bsoncore.InsufficientBytesError{}
		}
	} else {
		return errors.New("wrong bson type")
	}

	for i, k := range geometryStrings {
		if k == s {
			*g = GeometryType(i)
			return nil
		}
	}

	return ErrUnknownGeometryType
}

// Coordinates ...
type Coordinates []interface{}

func isLongLat(a []interface{}) bool {
	if len(a) != 2 {
		return false
	}

	return true
}

func isCoords(a []interface{}) bool {
	for _, v := range a {
		var ok bool

		switch v := v.(type) {
		case []interface{}:
			ok = isCoords(v)
		case float64:
			ok = isLongLat(a)
		}

		if !ok {
			return false
		}
	}

	return true
}

// Geometry describes a GeoJSON geometry object
type Geometry struct {
	Type        GeometryType `json:"type" bson:"type"`
	Coordinates Coordinates  `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
	Geometries  []Geometry   `json:"geometries,omitempty" bson:"geometries,omitempty"`
}

// Validate validates a GeoJSON geometry object
func (g Geometry) Validate() error {
	var ok bool

	switch g.Type {
	case Point:
		ok = isCoords(g.Coordinates)
	case MultiPoint:
		ok = isCoords(g.Coordinates)
	case LineString:
		ok = isCoords(g.Coordinates)
	case MultiLineString:
		ok = isCoords(g.Coordinates)
	case Polygon:
		ok = isCoords(g.Coordinates)
	case MultiPolygon:
		ok = isCoords(g.Coordinates)
	case GeometryCollection:
		if len(g.Geometries) == 0 {
			return ErrBadGeometryCollection
		}
	default:
		return ErrUnknownGeometryType
	}

	if !ok {
		return ErrBadGeometryCoords
	}

	if g.Type == GeometryCollection {
		if g.Coordinates != nil {
			return ErrBadGeometrySemantic
		}
	} else {
		if len(g.Geometries) > 0 {
			return ErrBadGeometrySemantic
		}
	}

	return nil
}
