package geobuf

import (
	"math"

	"github.com/cairnapp/go-geobuf/proto"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

const MaxPrecision = 100000000

func Encode(obj interface{}) *proto.Data {
	builder := newBuilder()
	builder.Analyze(obj)
	return builder.Build(obj)
}

type protoBuilder struct {
	keys      []string
	precision uint32
	dimension uint32
}

func newBuilder() *protoBuilder {
	pb := &protoBuilder{
		keys:      []string{},
		precision: 1,
		// Since Orb forces us into a 2 dimensional point, we'll have to use other ways to encode elevation + time
		// int(math.Max(float64(b.dimension), float64(len(point))))
		dimension: 2,
	}
	return pb
}

func (b *protoBuilder) Build(obj interface{}) *proto.Data {
	precision := math.Ceil(math.Log(float64(b.precision)) / math.Ln10)
	pbf := proto.Data{
		Keys:       b.keys,
		Dimensions: b.dimension,
		Precision:  uint32(precision),
	}

	switch t := obj.(type) {
	case *geojson.Geometry:
		switch t.Type {
		case "Point":
			p := t.Coordinates.(orb.Point)
			pbf.DataType = &proto.Data_Geometry_{
				Geometry: &proto.Data_Geometry{
					Type:   proto.Data_Geometry_POINT,
					Coords: translateCoords(b.precision, p[:]),
				},
			}
			return &pbf
		case "MultiPoint":
			p := t.Coordinates.(orb.MultiPoint)
			pbf.DataType = &proto.Data_Geometry_{
				Geometry: &proto.Data_Geometry{
					Type:   proto.Data_Geometry_MULTIPOINT,
					Coords: translateLine(b.precision, b.dimension, p, false),
				},
			}
			return &pbf
		case "LineString":
			p := t.Coordinates.(orb.LineString)
			pbf.DataType = &proto.Data_Geometry_{
				Geometry: &proto.Data_Geometry{
					Type:   proto.Data_Geometry_LINESTRING,
					Coords: translateLine(b.precision, b.dimension, p, false),
				},
			}
		case "MultiLineString":
			p := t.Coordinates.(orb.MultiLineString)
			coords, lengths := translateMultiLine(b.precision, b.dimension, p)
			pbf.DataType = &proto.Data_Geometry_{
				Geometry: &proto.Data_Geometry{
					Type:    proto.Data_Geometry_MULTILINESTRING,
					Coords:  coords,
					Lengths: lengths,
				},
			}
		case "Polygon":
			p := []orb.Ring(t.Coordinates.(orb.Polygon))
			coords, lengths := translateMultiRing(b.precision, b.dimension, p)
			pbf.DataType = &proto.Data_Geometry_{
				Geometry: &proto.Data_Geometry{
					Type:    proto.Data_Geometry_POLYGON,
					Coords:  coords,
					Lengths: lengths,
				},
			}
		case "MultiPolygon":
			p := []orb.Polygon(t.Coordinates.(orb.MultiPolygon))
			coords, lengths := translateMultiPolygon(b.precision, b.dimension, p)
			pbf.DataType = &proto.Data_Geometry_{
				Geometry: &proto.Data_Geometry{
					Type:    proto.Data_Geometry_MULTIPOLYGON,
					Coords:  coords,
					Lengths: lengths,
				},
			}

		}

	}
	return &pbf
}

func (b *protoBuilder) Analyze(obj interface{}) {
	switch t := obj.(type) {
	case geojson.FeatureCollection:
		for _, feature := range t.Features {
			b.Analyze(feature)
		}
	case geojson.Feature:
		b.Analyze(t.Geometry)
		for key, _ := range t.Properties {
			b.keys = append(b.keys, key)
		}
	case *geojson.Geometry:
		switch t.Type {
		case "Point":
			b.updatePrecision(t.Coordinates.(orb.Point))
		case "MultiPoint":
			coords := t.Coordinates.(orb.MultiPoint)
			for _, coord := range coords {
				b.updatePrecision(coord)
			}
		case "LineString":
			coords := t.Coordinates.(orb.LineString)
			for _, coord := range coords {
				b.updatePrecision(coord)
			}
		case "MultiLineString":
			lines := t.Coordinates.(orb.MultiLineString)
			for _, line := range lines {
				for _, coord := range line {
					b.updatePrecision(coord)
				}
			}
		case "Polygon":
			lines := t.Coordinates.(orb.Polygon)
			for _, line := range lines {
				for _, coord := range line {
					b.updatePrecision(coord)
				}
			}
		case "MultiPolygon":
			polygons := t.Coordinates.(orb.MultiPolygon)
			for _, rings := range polygons {
				for _, ring := range rings {
					for _, coord := range ring {
						b.updatePrecision(coord)
					}
				}
			}
		}
	// case geojson.GeometryCollection:
	case geojson.Polygon:
		for _, line := range t {
			for _, coord := range line {
				b.updatePrecision(coord)
			}
		}
	case geojson.MultiPolygon:
		for _, polygon := range t {
			for _, line := range polygon {
				for _, coord := range line {
					b.updatePrecision(coord)
				}
			}
		}
	}
}

func (b *protoBuilder) updatePrecision(point orb.Point) {
	e := getPrecision([2]float64(point))
	if e > b.precision {
		b.precision = e
	}
}

func translateMultiLine(e uint32, dim uint32, lines []orb.LineString) ([]int64, []uint32) {
	lengths := make([]uint32, len(lines))
	coords := []int64{}

	for i, line := range lines {
		lengths[i] = uint32(len(line))
		coords = append(coords, translateLine(e, dim, line, false)...)
	}
	return coords, lengths
}

func translateMultiPolygon(e uint32, dim uint32, polygons []orb.Polygon) ([]int64, []uint32) {
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

func translateMultiRing(e uint32, dim uint32, lines []orb.Ring) ([]int64, []uint32) {
	lengths := make([]uint32, len(lines))
	coords := []int64{}
	for i, line := range lines {
		lengths[i] = uint32(len(line) - 1)
		newLine := translateLine(e, dim, line, true)
		coords = append(coords, newLine...)
	}
	return coords, lengths
}

func translateLine(e uint32, dim uint32, points []orb.Point, isClosed bool) []int64 {
	sums := make([]int64, dim)
	ret := make([]int64, len(points)*int(dim))
	for i, point := range points {
		for j, p := range point {
			n := doTheMaths(e, p) - sums[j]
			ret[(int(dim)*i)+j] = n
			sums[j] = sums[j] + n
		}
	}
	if isClosed {
		return ret[:(len(ret) - int(dim))]
	}
	return ret
}

func translateCoords(e uint32, point []float64) []int64 {
	ret := make([]int64, len(point))
	for i, p := range point {
		ret[i] = doTheMaths(e, p)
	}
	return ret
}

func doTheMaths(e uint32, p float64) int64 {
	return int64(math.Round(p * float64(e)))
}

func getPrecision(point [2]float64) uint32 {
	var e uint32 = 1
	for _, val := range point {
		for {
			base := math.Round(float64(val * float64(e)))
			if (base/float64(e)) != val && float64(e) < MaxPrecision {
				e = e * 10
			} else {
				break
			}
		}
	}
	return e
}
