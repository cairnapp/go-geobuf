package geobuf_test

import (
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
