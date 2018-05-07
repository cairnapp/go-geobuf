package geobuf

import (
	"github.com/cairnapp/go-geobuf/pkg/math"
	"github.com/cairnapp/go-geobuf/proto"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func Decode(msg proto.Data) interface{} {
	switch v := msg.DataType.(type) {
	case *proto.Data_Geometry_:
		geo := v.Geometry
		return decodeGeometry(geo, msg.Precision, msg.Dimensions)
	case *proto.Data_Feature_:
		geo := v.Feature.Geometry
		decodedGeo := decodeGeometry(geo, msg.Precision, msg.Dimensions)
		feature := geojson.NewFeature(decodedGeo.Geometry())
		for i := 0; i < len(v.Feature.Properties); i = i + 2 {
			keyIdx := v.Feature.Properties[i]
			valIdx := v.Feature.Properties[i+1]
			val := v.Feature.Values[valIdx]
			switch actualVal := val.ValueType.(type) {
			case *proto.Data_Value_BoolValue:
				feature.Properties[msg.Keys[keyIdx]] = actualVal.BoolValue
			case *proto.Data_Value_DoubleValue:
				feature.Properties[msg.Keys[keyIdx]] = actualVal.DoubleValue
			case *proto.Data_Value_StringValue:
				feature.Properties[msg.Keys[keyIdx]] = actualVal.StringValue
			case *proto.Data_Value_PosIntValue:
				feature.Properties[msg.Keys[keyIdx]] = uint(actualVal.PosIntValue)
			case *proto.Data_Value_NegIntValue:
				feature.Properties[msg.Keys[keyIdx]] = int(actualVal.NegIntValue) * -1
			case *proto.Data_Value_JsonValue:
				feature.Properties[msg.Keys[keyIdx]] = actualVal.JsonValue
			}
		}
		switch id := v.Feature.IdType.(type) {
		case *proto.Data_Feature_Id:
			feature.ID = id.Id
		case *proto.Data_Feature_IntId:
			feature.ID = id.IntId
		}
		return feature
	case *proto.Data_FeatureCollection_:
	}
	return struct{}{}
}

func decodeGeometry(geo *proto.Data_Geometry, precision, dimensions uint32) *geojson.Geometry {
	switch geo.Type {
	case proto.Data_Geometry_POINT:
		return geojson.NewGeometry(makePoint(geo.Coords, precision))
	case proto.Data_Geometry_MULTIPOINT:
		return geojson.NewGeometry(makeMultiPoint(geo.Coords, precision, dimensions))
	case proto.Data_Geometry_LINESTRING:
		return geojson.NewGeometry(makeLineString(geo.Coords, precision, dimensions))
	case proto.Data_Geometry_MULTILINESTRING:
		return geojson.NewGeometry(makeMultiLineString(geo.Lengths, geo.Coords, precision, dimensions))
	case proto.Data_Geometry_POLYGON:
		return geojson.NewGeometry(makePolygon(geo.Lengths, geo.Coords, precision, dimensions))
	case proto.Data_Geometry_MULTIPOLYGON:
		return geojson.NewGeometry(makeMultiPolygon(geo.Lengths, geo.Coords, precision, dimensions))
	}
	return &geojson.Geometry{}
}

func makePoint(inCords []int64, precision uint32) orb.Point {
	point := [2]float64{}
	converted := makeCoords(inCords, precision)
	copy(point[:], converted[:])
	return orb.Point(point)
}

func makeMultiPoint(inCords []int64, precision uint32, dimension uint32) orb.MultiPoint {
	return orb.MultiPoint(makeLine(inCords, precision, dimension, false))
}

func makeMultiPolygon(lengths []uint32, inCords []int64, precision uint32, dimension uint32) orb.MultiPolygon {
	polyCount := int(lengths[0])
	polygons := make([]orb.Polygon, polyCount)
	lengths = lengths[1:]
	for i := 0; i < polyCount; i += 1 {
		ringCount := lengths[0]
		polygons[i] = makePolygon(lengths[1:ringCount+1], inCords, precision, dimension)
		skip := 0
		for i := 0; i < int(ringCount); i += 1 {
			skip += int(lengths[i]) * int(dimension)
		}

		lengths = lengths[ringCount:]
		inCords = inCords[skip:]
	}
	return orb.MultiPolygon(polygons)
}

func makePolygon(lengths []uint32, inCords []int64, precision uint32, dimension uint32) orb.Polygon {
	lines := make([]orb.Ring, len(lengths))
	for i, length := range lengths {
		l := int(length * dimension)
		lines[i] = makeRing(inCords[:l], precision, dimension)
		inCords = inCords[l:]
	}
	poly := orb.Polygon(lines)
	return poly
}

func makeMultiLineString(lengths []uint32, inCords []int64, precision uint32, dimension uint32) orb.MultiLineString {
	lines := make([]orb.LineString, len(lengths))
	for i, length := range lengths {
		l := int(length * dimension)
		lines[i] = makeLineString(inCords[:l], precision, dimension)
		inCords = inCords[l:]
	}
	return orb.MultiLineString(lines)
}

func makeRing(inCords []int64, precision uint32, dimension uint32) orb.Ring {
	points := makeLine(inCords, precision, dimension, true)
	points = append(points, points[0])
	return orb.Ring(points)
}

func makeLineString(inCords []int64, precision uint32, dimension uint32) orb.LineString {
	return orb.LineString(makeLine(inCords, precision, dimension, false))
}

func makeLine(inCords []int64, precision uint32, dimension uint32, isClosed bool) []orb.Point {
	points := make([]orb.Point, len(inCords)/int(dimension))
	prevCords := [2]int64{}
	for i, j := 0, 1; j < len(inCords); i, j = i+2, j+2 {
		prevCords[0] += inCords[i]
		prevCords[1] += inCords[j]
		points[i/2] = makePoint(prevCords[:], precision)
	}
	return points
}

func makeCoords(inCords []int64, precision uint32) []float64 {
	ret := make([]float64, len(inCords))
	e := math.DecodePrecision(precision)

	for i, val := range inCords {
		ret[i] = math.FloatWithPrecision(val, uint32(e))
	}
	return ret
}
