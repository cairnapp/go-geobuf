package geobuf

import (
	"github.com/paulmach/orb/geojson"

	"github.com/cairnapp/go-geobuf/pkg/decode"
	"github.com/cairnapp/go-geobuf/proto"
)

func Decode(msg proto.Data) interface{} {
	switch v := msg.DataType.(type) {
	case *proto.Data_Geometry_:
		geo := v.Geometry
		return decode.DecodeGeometry(geo, msg.Precision, msg.Dimensions)
	case *proto.Data_Feature_:
		return decode.DecodeFeature(msg, v.Feature, msg.Precision, msg.Dimensions)
	case *proto.Data_FeatureCollection_:
		collection := geojson.NewFeatureCollection()
		for _, feature := range v.FeatureCollection.Features {
			collection.Append(decode.DecodeFeature(msg, feature, msg.Precision, msg.Dimensions))
		}
		return collection
	}
	return struct{}{}
}
