package decode

import (
	"github.com/cairnapp/go-geobuf/pkg/geojson"
	"github.com/cairnapp/go-geobuf/pkg/geometry"
	"github.com/cairnapp/go-geobuf/pkg/math"
	"github.com/cairnapp/go-geobuf/proto"
)

func DecodeGeometry(geo *proto.Data_Geometry, precision, dimensions uint32) *geojson.Geometry {
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

func makePoint(inCords []int64, precision uint32) geometry.Point {
	return geometry.Point(makeCoords(inCords, precision))
}

func makeMultiPoint(inCords []int64, precision uint32, dimension uint32) geometry.MultiPoint {
	return geometry.MultiPoint(makeLine(inCords, precision, dimension, false))
}

func makeMultiPolygon(lengths []uint32, inCords []int64, precision uint32, dimension uint32) geometry.MultiPolygon {
	polyCount := int(lengths[0])
	polygons := make([]geometry.Polygon, polyCount)
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
	return geometry.MultiPolygon(polygons)
}

func makePolygon(lengths []uint32, inCords []int64, precision uint32, dimension uint32) geometry.Polygon {
	lines := make([]geometry.Ring, len(lengths))
	for i, length := range lengths {
		l := int(length * dimension)
		lines[i] = makeRing(inCords[:l], precision, dimension)
		inCords = inCords[l:]
	}
	poly := geometry.Polygon(lines)
	return poly
}

func makeMultiLineString(lengths []uint32, inCords []int64, precision uint32, dimension uint32) geometry.MultiLineString {
	lines := make([]geometry.LineString, len(lengths))
	for i, length := range lengths {
		l := int(length * dimension)
		lines[i] = makeLineString(inCords[:l], precision, dimension)
		inCords = inCords[l:]
	}
	return geometry.MultiLineString(lines)
}

func makeRing(inCords []int64, precision uint32, dimension uint32) geometry.Ring {
	points := makeLine(inCords, precision, dimension, true)
	points = append(points, points[0])
	return geometry.Ring(points)
}

func makeLineString(inCords []int64, precision uint32, dimension uint32) geometry.LineString {
	return geometry.LineString(makeLine(inCords, precision, dimension, false))
}

func makeLine(inCords []int64, precision uint32, dimension uint32, isClosed bool) []geometry.Point {
	points := make([]geometry.Point, len(inCords)/int(dimension))
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
