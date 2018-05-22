package encode

import (
	"github.com/cairnapp/go-geobuf/pkg/geojson"
	"github.com/cairnapp/go-geobuf/pkg/geometry"
	"github.com/cairnapp/go-geobuf/pkg/math"
)

type EncodingConfig struct {
	Dimension uint
	Precision uint
	Keys      KeyStore
}

type EncodingOption func(o *EncodingConfig)

func WithPrecision(precision uint) EncodingOption {
	return func(o *EncodingConfig) {
		o.Precision = uint(math.DecodePrecision(uint32(precision)))
	}
}

func WithDimension(dimension uint) EncodingOption {
	return func(o *EncodingConfig) {
		o.Dimension = dimension
	}
}

func WithKeyStore(store KeyStore) EncodingOption {
	return func(o *EncodingConfig) {
		o.Keys = store
	}
}

func FromAnalysis(obj interface{}) EncodingOption {
	return func(o *EncodingConfig) {
		analyze(obj, o)
	}
}

func analyze(obj interface{}, opts *EncodingConfig) {
	opts.Dimension = 2
	switch t := obj.(type) {
	case *geojson.FeatureCollection:
		for _, feature := range t.Features {
			analyze(feature, opts)
		}
	case *geojson.Feature:
		analyze(geojson.NewGeometry(t.Geometry), opts)
		for key, _ := range t.Properties {
			opts.Keys.Add(key)
		}
	case *geojson.Geometry:
		switch t.Type {
		case GeometryPoint:
			updatePrecision(t.Coordinates.(geometry.Point), opts)
		case GeometryMultiPoint:
			coords := t.Coordinates.(geometry.MultiPoint)
			for _, coord := range coords {
				updatePrecision(coord, opts)
			}
		case GeometryLineString:
			coords := t.Coordinates.(geometry.LineString)
			for _, coord := range coords {
				updatePrecision(coord, opts)
			}
		case GeometryMultiLineString:
			lines := t.Coordinates.(geometry.MultiLineString)
			for _, line := range lines {
				for _, coord := range line {
					updatePrecision(coord, opts)
				}
			}
		case GeometryPolygon:
			lines := t.Coordinates.(geometry.Polygon)
			for _, line := range lines {
				for _, coord := range line {
					updatePrecision(coord, opts)
				}
			}
		case GeometryMultiPolygon:
			polygons := t.Coordinates.(geometry.MultiPolygon)
			for _, rings := range polygons {
				for _, ring := range rings {
					for _, coord := range ring {
						updatePrecision(coord, opts)
					}
				}
			}
		}
	}

}

func updatePrecision(point geometry.Point, opt *EncodingConfig) {
	for _, val := range point {
		e := math.GetPrecision(val)
		if e > opt.Precision {
			opt.Precision = e
		}
	}
}
