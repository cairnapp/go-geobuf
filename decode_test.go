package geobuf_test

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	. "github.com/cairnapp/go-geobuf"
)

func TestDecodePoint(t *testing.T) {
	p := geojson.NewGeometry(orb.Point([2]float64{124.123, 234.456}))
	encoded := Encode(p)
	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiPoint(t *testing.T) {
	p := geojson.NewGeometry(orb.MultiPoint([]orb.Point{
		orb.Point([2]float64{124.123, 234.456}),
		orb.Point([2]float64{345.567, 456.678}),
	}))
	encoded := Encode(p)
	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeLineString(t *testing.T) {
	p := geojson.NewGeometry(orb.LineString([]orb.Point{
		orb.Point([2]float64{124.123, 234.456}),
		orb.Point([2]float64{345.567, 456.678}),
	}))
	encoded := Encode(p)
	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiLineString(t *testing.T) {
	p := geojson.NewGeometry(orb.MultiLineString([]orb.LineString{
		orb.LineString([]orb.Point{
			orb.Point([2]float64{124.123, 234.456}),
			orb.Point([2]float64{345.567, 456.678}),
		}),
		orb.LineString([]orb.Point{
			orb.Point([2]float64{224.123, 334.456}),
			orb.Point([2]float64{445.567, 556.678}),
		}),
	}))
	encoded := Encode(p)
	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodePolygon(t *testing.T) {
	p := geojson.NewGeometry(orb.Polygon([]orb.Ring{
		orb.Ring([]orb.Point{
			orb.Point([2]float64{124.123, 234.456}),
			orb.Point([2]float64{345.567, 456.678}),
			orb.Point([2]float64{124.123, 234.456}),
		}),
		orb.Ring([]orb.Point{
			orb.Point([2]float64{224.123, 334.456}),
			orb.Point([2]float64{445.567, 556.678}),
			orb.Point([2]float64{224.123, 334.456}),
		}),
	}))
	encoded := Encode(p)
	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiPolygon(t *testing.T) {
	p := geojson.NewGeometry(
		orb.MultiPolygon([]orb.Polygon{
			orb.Polygon([]orb.Ring{
				orb.Ring([]orb.Point{
					orb.Point([2]float64{124.123, 234.456}),
					orb.Point([2]float64{345.567, 456.678}),
					orb.Point([2]float64{124.123, 234.456}),
				}),
				orb.Ring([]orb.Point{
					orb.Point([2]float64{224.123, 334.456}),
					orb.Point([2]float64{445.567, 556.678}),
					orb.Point([2]float64{224.123, 334.456}),
				}),
			}),
			orb.Polygon([]orb.Ring{
				orb.Ring([]orb.Point{
					orb.Point([2]float64{124.123, 234.456}),
					orb.Point([2]float64{345.567, 456.678}),
					orb.Point([2]float64{124.123, 234.456}),
				}),
				orb.Ring([]orb.Point{
					orb.Point([2]float64{224.123, 334.456}),
					orb.Point([2]float64{445.567, 556.678}),
					orb.Point([2]float64{224.123, 334.456}),
				}),
			}),
		}))
	encoded := Encode(p)
	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeMultiPolygonEfficient(t *testing.T) {
	p := geojson.NewGeometry(
		orb.MultiPolygon([]orb.Polygon{
			orb.Polygon([]orb.Ring{
				orb.Ring([]orb.Point{
					orb.Point([2]float64{124.123, 234.456}),
					orb.Point([2]float64{345.567, 456.678}),
					orb.Point([2]float64{124.123, 234.456}),
				}),
				orb.Ring([]orb.Point{
					orb.Point([2]float64{224.123, 334.456}),
					orb.Point([2]float64{445.567, 556.678}),
					orb.Point([2]float64{224.123, 334.456}),
				}),
			}),
		}))
	encoded := Encode(p)
	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}

func TestDecodeFeature(t *testing.T) {
	p := geojson.NewFeature(orb.Polygon([]orb.Ring{
		orb.Ring([]orb.Point{
			orb.Point([2]float64{124.123, 234.456}),
			orb.Point([2]float64{345.567, 456.678}),
			orb.Point([2]float64{124.123, 234.456}),
		}),
		orb.Ring([]orb.Point{
			orb.Point([2]float64{224.123, 334.456}),
			orb.Point([2]float64{445.567, 556.678}),
			orb.Point([2]float64{224.123, 334.456}),
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

	decoded := Decode(*encoded)

	if !reflect.DeepEqual(p, decoded) {
		t.Errorf("Expected %+v, got %+v", p, decoded)
	}
}
