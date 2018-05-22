package geobuf_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	. "github.com/cairnapp/go-geobuf"
	"github.com/cairnapp/go-geobuf/pkg/geojson"
	"github.com/cairnapp/go-geobuf/pkg/geometry"
)

func TestDecodePoint(t *testing.T) {
	p := geojson.NewGeometry(geometry.Point([]float64{124.123, 234.456}))
	encoded := Encode(p)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiPoint(t *testing.T) {
	p := geojson.NewGeometry(geometry.MultiPoint([]geometry.Point{
		geometry.Point([]float64{124.123, 234.456}),
		geometry.Point([]float64{345.567, 456.678}),
	}))
	encoded := Encode(p)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeLineString(t *testing.T) {
	p := geojson.NewGeometry(geometry.LineString([]geometry.Point{
		geometry.Point([]float64{124.123, 234.456}),
		geometry.Point([]float64{345.567, 456.678}),
	}))
	encoded := Encode(p)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiLineString(t *testing.T) {
	p := geojson.NewGeometry(geometry.MultiLineString([]geometry.LineString{
		geometry.LineString([]geometry.Point{
			geometry.Point([]float64{124.123, 234.456}),
			geometry.Point([]float64{345.567, 456.678}),
		}),
		geometry.LineString([]geometry.Point{
			geometry.Point([]float64{224.123, 334.456}),
			geometry.Point([]float64{445.567, 556.678}),
		}),
	}))
	encoded := Encode(p)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodePolygon(t *testing.T) {
	p := geojson.NewGeometry(geometry.Polygon([]geometry.Ring{
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{124.123, 234.456}),
			geometry.Point([]float64{345.567, 456.678}),
			geometry.Point([]float64{124.123, 234.456}),
		}),
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{224.123, 334.456}),
			geometry.Point([]float64{445.567, 556.678}),
			geometry.Point([]float64{224.123, 334.456}),
		}),
	}))
	encoded := Encode(p)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiPolygon(t *testing.T) {
	p := geojson.NewGeometry(
		geometry.MultiPolygon([]geometry.Polygon{
			geometry.Polygon([]geometry.Ring{
				geometry.Ring([]geometry.Point{
					geometry.Point([]float64{124.123, 234.456}),
					geometry.Point([]float64{345.567, 456.678}),
					geometry.Point([]float64{124.123, 234.456}),
				}),
				geometry.Ring([]geometry.Point{
					geometry.Point([]float64{224.123, 334.456}),
					geometry.Point([]float64{445.567, 556.678}),
					geometry.Point([]float64{224.123, 334.456}),
				}),
			}),
			geometry.Polygon([]geometry.Ring{
				geometry.Ring([]geometry.Point{
					geometry.Point([]float64{124.123, 234.456}),
					geometry.Point([]float64{345.567, 456.678}),
					geometry.Point([]float64{124.123, 234.456}),
				}),
				geometry.Ring([]geometry.Point{
					geometry.Point([]float64{224.123, 334.456}),
					geometry.Point([]float64{445.567, 556.678}),
					geometry.Point([]float64{224.123, 334.456}),
				}),
			}),
		}))
	encoded := Encode(p)
	spew.Dump(encoded)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiPolygonEfficient(t *testing.T) {
	p := geojson.NewGeometry(
		geometry.MultiPolygon([]geometry.Polygon{
			geometry.Polygon([]geometry.Ring{
				geometry.Ring([]geometry.Point{
					geometry.Point([]float64{124.123, 234.456}),
					geometry.Point([]float64{345.567, 456.678}),
					geometry.Point([]float64{124.123, 234.456}),
				}),
				geometry.Ring([]geometry.Point{
					geometry.Point([]float64{224.123, 334.456}),
					geometry.Point([]float64{445.567, 556.678}),
					geometry.Point([]float64{224.123, 334.456}),
				}),
			}),
		}))
	encoded := Encode(p)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeFeatureIntId(t *testing.T) {
	p := geojson.NewFeature(geometry.Polygon([]geometry.Ring{
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{124.123, 234.456}),
			geometry.Point([]float64{345.567, 456.678}),
			geometry.Point([]float64{124.123, 234.456}),
		}),
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{224.123, 334.456}),
			geometry.Point([]float64{445.567, 556.678}),
			geometry.Point([]float64{224.123, 334.456}),
		}),
	}))
	p.ID = int64(1)
	p.Properties["int"] = uint(4)
	p.Properties["float"] = float64(2.0)
	p.Properties["neg_int"] = -1
	p.Properties["string"] = "string"
	p.Properties["bool"] = true
	encoded := Encode(p)
	spew.Dump(encoded)
	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeFeatureStringId(t *testing.T) {
	p := geojson.NewFeature(geometry.Polygon([]geometry.Ring{
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{124.123, 234.456}),
			geometry.Point([]float64{345.567, 456.678}),
			geometry.Point([]float64{124.123, 234.456}),
		}),
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{224.123, 334.456}),
			geometry.Point([]float64{445.567, 556.678}),
			geometry.Point([]float64{224.123, 334.456}),
		}),
	}))
	p.ID = "1234"
	p.Properties["int"] = uint(4)
	p.Properties["float"] = float64(2.0)
	p.Properties["neg_int"] = -1
	p.Properties["string"] = "string"
	p.Properties["bool"] = true
	encoded := Encode(p)
	spew.Dump(encoded)

	decoded := Decode(encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeFeatureCollection(t *testing.T) {
	p := geojson.NewFeature(geometry.Polygon([]geometry.Ring{
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{124.123, 234.456}),
			geometry.Point([]float64{345.567, 456.678}),
			geometry.Point([]float64{124.123, 234.456}),
		}),
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{224.123, 334.456}),
			geometry.Point([]float64{445.567, 556.678}),
			geometry.Point([]float64{224.123, 334.456}),
		}),
	}))
	p.ID = "1234"
	p.Properties["int"] = uint(4)
	p.Properties["float"] = float64(2.0)
	p.Properties["neg_int"] = -1
	p.Properties["string"] = "string"
	p.Properties["bool"] = true

	p2 := geojson.NewFeature(geometry.Polygon([]geometry.Ring{
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{224.123, 334.456}),
			geometry.Point([]float64{445.567, 556.678}),
			geometry.Point([]float64{224.123, 334.456}),
		}),
		geometry.Ring([]geometry.Point{
			geometry.Point([]float64{124.123, 234.456}),
			geometry.Point([]float64{345.567, 456.678}),
			geometry.Point([]float64{124.123, 234.456}),
		}),
	}))
	p2.ID = "5679"
	p2.Properties["int"] = uint(4)
	p2.Properties["float"] = float64(2.0)
	p2.Properties["neg_int"] = -1
	p2.Properties["string"] = "string"
	p2.Properties["bool"] = true

	collection := geojson.NewFeatureCollection()
	collection.Append(p)
	collection.Append(p2)
	encoded := Encode(collection)

	decoded := Decode(encoded)

	if !reflect.DeepEqual(collection, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}
