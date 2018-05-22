package geojson

type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

const FeatureCollectionType = "FeatureCollection"

func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     FeatureCollectionType,
		Features: []*Feature{},
	}
}

// Append appends a feature to the collection.
func (fc *FeatureCollection) Append(feature *Feature) *FeatureCollection {
	fc.Features = append(fc.Features, feature)
	return fc
}
