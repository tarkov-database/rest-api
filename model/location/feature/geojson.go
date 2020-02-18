package feature

import (
	"encoding/json"
	"errors"
)

var (
	// ErrUnsupportedFeatureType ...
	ErrUnsupportedFeatureType = errors.New("unsupported feature type")

	// ErrBadFeatureSemantic ...
	ErrBadFeatureSemantic = errors.New("bad geometry semantic")

	// ErrBadGeometrySemantic ...
	ErrBadGeometrySemantic = errors.New("bad geometry semantic")

	// ErrUnknownGeometryType ...
	ErrUnknownGeometryType = errors.New("unknown geometry type")

	// ErrBadGeometryCoords ...
	ErrBadGeometryCoords = errors.New("bad geometry coordinates")

	// ErrEmtpyGeometryCollection ...
	ErrEmtpyGeometryCollection = errors.New("empty geometry collection")
)

// GeometryType ...
type GeometryType int

// ...
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

// MarshalJSON ...
func (g GeometryType) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// UnmarshalJSON ...
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

type (
	// Coordinates ...
	Coordinates []interface{}

	// LongLat ...
	LongLat []interface{}

	// CoordsPoint ...
	CoordsPoint LongLat

	// CoordsMultiPoint ...
	CoordsMultiPoint []LongLat

	// CoordsLineString ...
	CoordsLineString []LongLat

	// CoordsMultiLineString ...
	CoordsMultiLineString [][]LongLat

	// CoordsPolygon ...
	CoordsPolygon [][]LongLat

	// CoordsMultiPolygon ...
	CoordsMultiPolygon [][][]LongLat
)

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

// Geometry ...
type Geometry struct {
	Type        GeometryType `json:"type" bson:"type"`
	Coordinates Coordinates  `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
	Geometries  []Geometry   `json:"geometries,omitempty" bson:"geometries,omitempty"`
}

// Validate ...
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
			return ErrEmtpyGeometryCollection
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
