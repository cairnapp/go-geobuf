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
