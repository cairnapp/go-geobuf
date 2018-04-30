package geobuf

import (
	"github.com/cairnapp/go-geobuf/proto"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"math"
)

func Decode(msg proto.Data) interface{} {
	switch v := msg.DataType.(type) {
	case *proto.Data_Geometry_:
		geo := v.Geometry
		switch geo.Type {
		case proto.Data_Geometry_POINT:
			return geojson.NewGeometry(makePoint(geo.Coords, msg.Precision))
		case proto.Data_Geometry_MULTIPOINT:
		case proto.Data_Geometry_LINESTRING:
		}
		return geojson.Geometry{}
	case *proto.Data_Feature_:
	case *proto.Data_FeatureCollection_:
	}
	return struct{}{}
}

func makePoint(inCords []int64, precision uint32) orb.Point {
	point := [2]float64{}
	converted := makeCoords(inCords, precision)
	copy(point[:], converted[:])
	return orb.Point(point)
}

func makeCoords(inCords []int64, precision uint32) []float64 {
	ret := make([]float64, len(inCords))
	e := math.Pow10(int(precision))

	for i, val := range inCords {
		ret[i] = float64(val) / e
	}
	return ret
}
