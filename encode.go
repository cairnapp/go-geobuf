package geobuf

import (
	"github.com/cairnapp/go-geobuf/pkg/encode"
	"github.com/cairnapp/go-geobuf/pkg/math"
	"github.com/cairnapp/go-geobuf/proto"

	"github.com/paulmach/orb/geojson"
)

func Encode(obj interface{}) *proto.Data {
	data, err := EncodeWithOptions(obj, encode.FromAnalysis(obj))
	if err != nil {
		panic(err)
	}
	return data
}

func EncodeWithOptions(obj interface{}, opts ...encode.EncodingOption) (*proto.Data, error) {
	cfg := &encode.EncodingConfig{
		Dimension: 2,
		Precision: 1,
		Keys:      encode.NewKeyStore(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	data := &proto.Data{
		Keys:       cfg.Keys.Keys(),
		Dimensions: uint32(cfg.Dimension),
		Precision:  math.EncodePrecision(cfg.Precision),
	}

	switch t := obj.(type) {
	case *geojson.FeatureCollection:
		collection, err := encode.EncodeFeatureCollection(t, cfg)
		if err != nil {
			return nil, err
		}
		data.DataType = &proto.Data_FeatureCollection_{
			FeatureCollection: collection,
		}
	case *geojson.Feature:
		feature, err := encode.EncodeFeature(t, cfg)
		if err != nil {
			return nil, err
		}
		data.DataType = &proto.Data_Feature_{
			Feature: feature,
		}
	case *geojson.Geometry:
		data.DataType = &proto.Data_Geometry_{
			Geometry: encode.EncodeGeometry(t, cfg),
		}
	}

	return data, nil
}
