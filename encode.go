package geobuf

import (
	"sort"

	"github.com/cairnapp/go-geobuf/pkg/encode"
	"github.com/cairnapp/go-geobuf/pkg/math"
	"github.com/cairnapp/go-geobuf/proto"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func Encode(obj interface{}) *proto.Data {
	builder := newBuilder()
	builder.Analyze(obj)
	b, err := builder.Build(obj)
	if err != nil {
		panic(err)
	}
	return b
}

type protoBuilder struct {
	keys      map[string]bool
	precision uint32
	dimension uint32
	data      *proto.Data
}

func newBuilder() *protoBuilder {
	pb := &protoBuilder{
		keys:      map[string]bool{},
		precision: 1,
		// Since Orb forces us into a 2 dimensional point, we'll have to use other ways to encode elevation + time
		// int(math.Max(float64(b.dimension), float64(len(point))))
		dimension: 2,
		data:      nil,
	}
	return pb
}

func (b *protoBuilder) Build(obj interface{}) (*proto.Data, error) {
	precision := math.EncodePrecision(float64(b.precision))
	keys := make([]string, 0, len(b.keys))
	for key, _ := range b.keys {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	b.data = &proto.Data{
		Keys:       keys,
		Dimensions: b.dimension,
		Precision:  precision,
	}

	switch t := obj.(type) {
	case *geojson.Feature:
		data, err := b.encodeFeature(t)
		if err != nil {
			return b.data, err
		}
		b.data.DataType = data
	case *geojson.Geometry:
		b.data.DataType = b.buildGeometry(t)
	}
	return b.data, nil
}

func (b protoBuilder) encodeFeature(feature *geojson.Feature) (*proto.Data_Feature_, error) {
	oldGeo := geojson.NewGeometry(feature.Geometry)
	geo := b.buildGeometry(oldGeo)
	f := &proto.Data_Feature_{
		Feature: &proto.Data_Feature{
			Geometry: geo.Geometry,
		},
	}

	id, err := encode.EncodeIntId(feature.ID)
	if err == nil {
		f.Feature.IdType = id
	} else {
		newId, newErr := encode.EncodeId(feature.ID)
		if newErr != nil {
			return nil, newErr
		}
		f.Feature.IdType = newId
	}

	for key, val := range feature.Properties {
		encoded, err := encode.EncodeValue(val)
		if err != nil {
			return f, err
		}
		idx := indexOf(b.data.Keys, key)
		f.Feature.Values = append(f.Feature.Values, encoded)
		f.Feature.Properties = append(f.Feature.Properties, idx)
		f.Feature.Properties = append(f.Feature.Properties, uint32(len(f.Feature.Values))-1)
	}
	return f, nil
}

func indexOf(values []string, key string) uint32 {
	return uint32(sort.SearchStrings(values, key))
}

func (b protoBuilder) buildGeometry(t *geojson.Geometry) *proto.Data_Geometry_ {
	switch t.Type {
	case "Point":
		p := t.Coordinates.(orb.Point)
		return &proto.Data_Geometry_{
			Geometry: &proto.Data_Geometry{
				Type:   proto.Data_Geometry_POINT,
				Coords: translateCoords(b.precision, p[:]),
			},
		}
	case "MultiPoint":
		p := t.Coordinates.(orb.MultiPoint)
		return &proto.Data_Geometry_{
			Geometry: &proto.Data_Geometry{
				Type:   proto.Data_Geometry_MULTIPOINT,
				Coords: translateLine(b.precision, b.dimension, p, false),
			},
		}
	case "LineString":
		p := t.Coordinates.(orb.LineString)
		return &proto.Data_Geometry_{
			Geometry: &proto.Data_Geometry{
				Type:   proto.Data_Geometry_LINESTRING,
				Coords: translateLine(b.precision, b.dimension, p, false),
			},
		}
	case "MultiLineString":
		p := t.Coordinates.(orb.MultiLineString)
		coords, lengths := translateMultiLine(b.precision, b.dimension, p)
		return &proto.Data_Geometry_{
			Geometry: &proto.Data_Geometry{
				Type:    proto.Data_Geometry_MULTILINESTRING,
				Coords:  coords,
				Lengths: lengths,
			},
		}
	case "Polygon":
		p := []orb.Ring(t.Coordinates.(orb.Polygon))
		coords, lengths := translateMultiRing(b.precision, b.dimension, p)
		return &proto.Data_Geometry_{
			Geometry: &proto.Data_Geometry{
				Type:    proto.Data_Geometry_POLYGON,
				Coords:  coords,
				Lengths: lengths,
			},
		}
	case "MultiPolygon":
		p := []orb.Polygon(t.Coordinates.(orb.MultiPolygon))
		coords, lengths := translateMultiPolygon(b.precision, b.dimension, p)
		return &proto.Data_Geometry_{
			Geometry: &proto.Data_Geometry{
				Type:    proto.Data_Geometry_MULTIPOLYGON,
				Coords:  coords,
				Lengths: lengths,
			},
		}
	}
	return nil
}

func (b *protoBuilder) Analyze(obj interface{}) {
	switch t := obj.(type) {
	case geojson.FeatureCollection:
		for _, feature := range t.Features {
			b.Analyze(feature)
		}
	case *geojson.Feature:
		b.Analyze(geojson.NewGeometry(t.Geometry))
		for key, _ := range t.Properties {
			_, ok := b.keys[key]
			if !ok {
				b.keys[key] = true
			}
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
	for _, val := range point {
		e := math.GetPrecision(val)
		if e > b.precision {
			b.precision = e
		}
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
			n := math.IntWithPrecision(p, e) - sums[j]
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
		ret[i] = math.IntWithPrecision(p, e)
	}
	return ret
}
