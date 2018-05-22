package geojson

import (
	"github.com/cairnapp/go-geobuf/pkg/geometry"
)

const (
	GeometryPointType           = "Point"
	GeometryMultiPointType      = "MultiPoint"
	GeometryLineStringType      = "LineString"
	GeometryMultiLineStringType = "MultiLineString"
	GeometryPolygonType         = "Polygon"
	GeometryMultiPolygonType    = "MultiPolygon"
	GeometryCollectionType      = "GeometryCollectionType"
)

type Geometry struct {
	Type        string            `json:"type"`
	Coordinates geometry.Geometry `json:"coordinates,omitempty"`
	Geometries  []*Geometry       `json:"geometries,omitempty"`
}

func NewGeometry(g geometry.Geometry) *Geometry {
	geo := &Geometry{}
	switch typed := g.(type) {
	case geometry.Point:
		geo.Type = GeometryPointType
		geo.Coordinates = g
	case geometry.MultiPoint:
		geo.Type = GeometryMultiPointType
		geo.Coordinates = g
	case geometry.LineString:
		geo.Type = GeometryLineStringType
		geo.Coordinates = g
	case geometry.MultiLineString:
		geo.Type = GeometryMultiLineStringType
		geo.Coordinates = g
	case geometry.Polygon:
		geo.Type = GeometryPolygonType
		geo.Coordinates = g
	case geometry.MultiPolygon:
		geo.Type = GeometryMultiPolygonType
		geo.Coordinates = g
	case geometry.Collection:
		geo.Type = GeometryCollectionType
		geo.Geometries = make([]*Geometry, len(typed))
		for i, child := range typed {
			geo.Geometries[i] = NewGeometry(child)
		}
	}
	return geo
}
