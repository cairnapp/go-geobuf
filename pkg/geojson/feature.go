package geojson

import "github.com/cairnapp/go-geobuf/pkg/geometry"

type Properties map[string]interface{}

const FeatureType = "Feature"

// A Feature corresponds to GeoJSON feature object
type Feature struct {
	ID         interface{}       `json:"id,omitempty"`
	Type       string            `json:"type"`
	Geometry   geometry.Geometry `json:"geometry"`
	Properties Properties        `json:"properties"`
}

func NewFeature(geometry geometry.Geometry) *Feature {
	return &Feature{
		Type:       FeatureType,
		Geometry:   geometry,
		Properties: make(map[string]interface{}),
	}
}
