package encode

import (
	"github.com/cairnapp/go-geobuf/pkg/geojson"
	"github.com/cairnapp/go-geobuf/proto"
)

func EncodeFeature(feature *geojson.Feature, opts *EncodingConfig) (*proto.Data_Feature, error) {
	oldGeo := geojson.NewGeometry(feature.Geometry)
	geo := EncodeGeometry(oldGeo, opts)
	f := &proto.Data_Feature{
		Geometry: geo,
	}

	id, err := EncodeIntId(feature.ID)
	if err == nil {
		f.IdType = id
	} else {
		newId, newErr := EncodeId(feature.ID)
		if newErr != nil {
			return nil, newErr
		}
		f.IdType = newId
	}

	properties := make([]uint32, 0, 2*len(feature.Properties))
	values := make([]*proto.Data_Value, 0, len(feature.Properties))
	for key, val := range feature.Properties {
		encoded, err := EncodeValue(val)
		if err != nil {
			return f, err
		}

		idx := opts.Keys.IndexOf(key)
		values = append(values, encoded)
		properties = append(properties, uint32(idx))
		properties = append(properties, uint32(len(values)-1))
	}

	f.Values = values
	f.Properties = properties
	return f, nil
}
