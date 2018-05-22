package encode

import (
	"github.com/cairnapp/go-geobuf/pkg/geojson"
	"github.com/cairnapp/go-geobuf/pkg/geometry"
	"github.com/cairnapp/go-geobuf/pkg/math"
	"github.com/cairnapp/go-geobuf/proto"
)

const (
	GeometryPoint           = "Point"
	GeometryMultiPoint      = "MultiPoint"
	GeometryLineString      = "LineString"
	GeometryMultiLineString = "MultiLineString"
	GeometryPolygon         = "Polygon"
	GeometryMultiPolygon    = "MultiPolygon"
)

func EncodeGeometry(g *geojson.Geometry, opt *EncodingConfig) *proto.Data_Geometry {
	switch g.Type {
	case geojson.GeometryPointType:
		p := g.Coordinates.(geometry.Point)
		return &proto.Data_Geometry{
			Type:   proto.Data_Geometry_POINT,
			Coords: translateCoords(opt.Precision, p[:]),
		}
	case geojson.GeometryMultiPointType:
		p := g.Coordinates.(geometry.MultiPoint)
		return &proto.Data_Geometry{
			Type:   proto.Data_Geometry_MULTIPOINT,
			Coords: translateLine(opt.Precision, opt.Dimension, p, false),
		}
	case geojson.GeometryLineStringType:
		p := g.Coordinates.(geometry.LineString)
		return &proto.Data_Geometry{
			Type:   proto.Data_Geometry_LINESTRING,
			Coords: translateLine(opt.Precision, opt.Dimension, p, false),
		}
	case geojson.GeometryMultiLineStringType:
		p := g.Coordinates.(geometry.MultiLineString)
		coords, lengths := translateMultiLine(opt.Precision, opt.Dimension, p)
		return &proto.Data_Geometry{
			Type:    proto.Data_Geometry_MULTILINESTRING,
			Coords:  coords,
			Lengths: lengths,
		}
	case geojson.GeometryPolygonType:
		p := []geometry.Ring(g.Coordinates.(geometry.Polygon))
		coords, lengths := translateMultiRing(opt.Precision, opt.Dimension, p)
		return &proto.Data_Geometry{
			Type:    proto.Data_Geometry_POLYGON,
			Coords:  coords,
			Lengths: lengths,
		}
	case geojson.GeometryMultiPolygonType:
		p := []geometry.Polygon(g.Coordinates.(geometry.MultiPolygon))
		coords, lengths := translateMultiPolygon(opt.Precision, opt.Dimension, p)
		return &proto.Data_Geometry{
			Type:    proto.Data_Geometry_MULTIPOLYGON,
			Coords:  coords,
			Lengths: lengths,
		}
	}
	return nil
}

func translateMultiLine(e uint, dim uint, lines []geometry.LineString) ([]int64, []uint32) {
	lengths := make([]uint32, len(lines))
	coords := []int64{}

	for i, line := range lines {
		lengths[i] = uint32(len(line))
		coords = append(coords, translateLine(e, dim, line, false)...)
	}
	return coords, lengths
}

func translateMultiPolygon(e uint, dim uint, polygons []geometry.Polygon) ([]int64, []uint32) {
	lengths := []uint32{uint32(len(polygons))}
	coords := []int64{}
	for _, rings := range polygons {
		lengths = append(lengths, uint32(len(rings)))
		newLine, newLength := translateMultiRing(e, dim, rings)
		lengths = append(lengths, newLength...)
		coords = append(coords, newLine...)
	}
	return coords, lengths
}

func translateMultiRing(e uint, dim uint, lines []geometry.Ring) ([]int64, []uint32) {
	lengths := make([]uint32, len(lines))
	coords := []int64{}
	for i, line := range lines {
		lengths[i] = uint32(len(line) - 1)
		newLine := translateLine(e, dim, line, true)
		coords = append(coords, newLine...)
	}
	return coords, lengths
}

/*
Since we're converting floats to ints, we can get additional compression out of
how Google does varint encoding (#1). Smaller numbers can be packed into less bytes,
even when using large primitives (int64). To take advantage of this, we subtract
out the prior coordinate x/y value from the current coordinate x/y value to (hopefully)
reduce the number to a small integer.

For example: (123.123, 234.234), (123.134, 234.236) would be encoded out to
(123123, 234234), (11, 2). The first point takes the full penalty for encoding size,
while the remaining points become much smaller.

A further enhancement comes from the fact that lines that start and end in the same place,
such as with a polygon, we can skip the last point, and place it back when we decode.

1. https://developers.google.com/protocol-buffers/docs/encoding#varints
*/
func translateLine(precision uint, dim uint, points []geometry.Point, isClosed bool) []int64 {
	sums := make([]int64, dim)
	ret := make([]int64, len(points)*int(dim))
	for i, point := range points {
		for j, p := range point {
			n := math.IntWithPrecision(p, precision) - sums[j]
			ret[(int(dim)*i)+j] = n
			sums[j] = sums[j] + n
		}
	}
	if isClosed {
		return ret[:(len(ret) - int(dim))]
	}
	return ret
}

// Converts a floating point geojson point to int64 by multiplying it by a factor of 10,
// potentially truncating and rounding
func translateCoords(precision uint, point []float64) []int64 {
	ret := make([]int64, len(point))
	for i, p := range point {
		ret[i] = math.IntWithPrecision(p, precision)
	}
	return ret
}
