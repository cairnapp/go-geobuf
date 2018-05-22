package encode_test

import (
	"reflect"
	"testing"

	. "github.com/cairnapp/go-geobuf/pkg/encode"
	"github.com/cairnapp/go-geobuf/pkg/geojson"
	"github.com/cairnapp/go-geobuf/pkg/geometry"
	"github.com/cairnapp/go-geobuf/proto"
)

func TestEncodePoint(t *testing.T) {
	testCases := []struct {
		Precision uint
		Expected  []int64
	}{
		{
			Precision: 1000,
			Expected:  []int64{124123, 234456},
		},
		// Should round up when truncating precision
		{
			Precision: 100,
			Expected:  []int64{12412, 23446},
		},
		// Should round up (.5) when truncating precision
		{
			Precision: 10,
			Expected:  []int64{1241, 2345},
		},
		{
			Precision: 1,
			Expected:  []int64{124, 234},
		},
	}

	p := geojson.NewGeometry(geometry.Point([]float64{124.123, 234.456}))
	for i, test := range testCases {
		expected := &proto.Data_Geometry{
			Type:   proto.Data_Geometry_POINT,
			Coords: test.Expected,
		}
		encoded := EncodeGeometry(p, &EncodingConfig{
			Dimension: 2,
			Precision: test.Precision,
		})

		if !reflect.DeepEqual(encoded, expected) {
			t.Errorf("Case [%d]: Expected %+v, got %+v", i, expected, encoded)
		}
	}

}
