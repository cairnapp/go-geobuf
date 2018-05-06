package encode_test

import (
	"fmt"
	"testing"

	. "github.com/cairnapp/go-geobuf/pkg/encode"
	"github.com/cairnapp/go-geobuf/proto"
)

func TestEncodeIntValue(t *testing.T) {
	testCases := []struct {
		Val      interface{}
		Expected uint
	}{
		{
			Val:      1,
			Expected: 1,
		},
		// Neg int values should be positively encoded
		{
			Val:      -1,
			Expected: 1,
		},
		{
			Val:      int8(1),
			Expected: 1,
		},
		{
			Val:      int16(10),
			Expected: 10,
		},
		{
			Val:      int32(-5),
			Expected: 5,
		},
		{
			Val:      int64(-128),
			Expected: 128,
		},
		{
			Val:      intPtr(8),
			Expected: 8,
		},
		{
			Val:      int8Ptr(8),
			Expected: 8,
		},
		{
			Val:      int16Ptr(8),
			Expected: 8,
		},
		{
			Val:      int32Ptr(-8),
			Expected: 8,
		},
		{
			Val:      int64Ptr(-128),
			Expected: 128,
		},
		{
			Val:      uint8(4),
			Expected: 4,
		},
		{
			Val:      uint16(16),
			Expected: 16,
		},
		{
			Val:      uint32(16),
			Expected: 16,
		},
		{
			Val:      uint64(36),
			Expected: 36,
		},
		{
			Val:      uint8Ptr(36),
			Expected: 36,
		},
		{
			Val:      uint16Ptr(36),
			Expected: 36,
		},
		{
			Val:      uint32Ptr(36),
			Expected: 36,
		},
		{
			Val:      uint64Ptr(36),
			Expected: 36,
		},
		{
			Val:      uintPtr(36),
			Expected: 36,
		},
		{
			Val:      uintPtr(18446744073709551615),
			Expected: uint(18446744073709551615),
		},
		{
			Val:      uint64(18446744073709551615),
			Expected: uint(18446744073709551615),
		},
		{
			Val:      int64(-9223372036854775808),
			Expected: uint(9223372036854775808),
		},
		{
			Val:      int64(9223372036854775807),
			Expected: uint(9223372036854775807),
		},
	}

	for i, test := range testCases {
		val, err := EncodeValue(test.Val)
		if err != nil {
			t.Fatalf("Case [%d]: Got unexpected error %s!", i, err)
		}

		switch cast := val.ValueType.(type) {
		case *proto.Data_Value_PosIntValue:
			if uint64(test.Expected) != cast.PosIntValue {
				t.Errorf("Case [%d]: Expected %d, got %d", i, test.Expected, cast.PosIntValue)
			}
		case *proto.Data_Value_NegIntValue:
			if uint64(test.Expected) != cast.NegIntValue {
				t.Errorf("Case [%d]: Expected %d, got %d", i, test.Expected, cast.NegIntValue)
			}
		default:
			t.Errorf("Case [%d]: Got unexpected type back %+v!", i, cast)
		}
	}
}

func TestEncodeStringValue(t *testing.T) {
	testCases := []struct {
		Val      interface{}
		Expected string
	}{
		{
			Val:      "Testing 123",
			Expected: "Testing 123",
		},
		{
			Val:      strPtr("Testing 123"),
			Expected: "Testing 123",
		},
	}

	for i, test := range testCases {
		val, err := EncodeValue(test.Val)
		if err != nil {
			t.Fatalf("Case [%d]: Got unexpected error %s!", i, err)
		}

		switch cast := val.ValueType.(type) {
		case *proto.Data_Value_StringValue:
			if test.Expected != cast.StringValue {
				t.Errorf("Case [%d]: Expected %s, got %s", i, test.Expected, cast.StringValue)
			}
		default:
			t.Errorf("Case [%d]: Got unexpected type back %+v!", i, cast)
		}
	}
}

func TestEncodeFloatValue(t *testing.T) {
	testCases := []struct {
		Val      interface{}
		Expected float64
	}{
		{
			Val:      float32(12.5),
			Expected: 12.5,
		},
		{
			Val:      float64(128.123),
			Expected: 128.123,
		},
		{
			Val:      float32Ptr(12.5),
			Expected: 12.5,
		},
		{
			Val:      float64Ptr(128.123),
			Expected: 128.123,
		},
	}

	for i, test := range testCases {
		val, err := EncodeValue(test.Val)
		if err != nil {
			t.Fatalf("Case [%d]: Got unexpected error %s!", i, err)
		}

		switch cast := val.ValueType.(type) {
		case *proto.Data_Value_DoubleValue:
			if test.Expected != cast.DoubleValue {
				t.Errorf("Case [%d]: Expected %f, got %f", i, test.Expected, cast.DoubleValue)
			}
		default:
			t.Errorf("Case [%d]: Got unexpected type back %+v!", i, cast)
		}
	}
}

func TestEncodeBoolValue(t *testing.T) {
	testCases := []struct {
		Val      interface{}
		Expected bool
	}{
		{
			Val:      true,
			Expected: true,
		},
		{
			Val:      false,
			Expected: false,
		},
		{
			Val:      boolPtr(true),
			Expected: true,
		},
		{
			Val:      boolPtr(false),
			Expected: false,
		},
	}

	for i, test := range testCases {
		val, err := EncodeValue(test.Val)
		if err != nil {
			t.Fatalf("Case [%d]: Got unexpected error %s!", i, err)
		}

		switch cast := val.ValueType.(type) {
		case *proto.Data_Value_BoolValue:
			if test.Expected != cast.BoolValue {
				t.Errorf("Case [%d]: Expected %t, got %t", i, test.Expected, cast.BoolValue)
			}
		default:
			t.Errorf("Case [%d]: Got unexpected type back %+v!", i, cast)
		}
	}
}

func TestEncodeJsonValue(t *testing.T) {
	testCases := []struct {
		Val      interface{}
		Expected string
	}{
		{
			Val:      []string{"A", "B", "C"},
			Expected: fmt.Sprintf("[%q,%q,%q]", "A", "B", "C"),
		},
		{
			Val:      map[string]int{"1": 1},
			Expected: "{\"1\":1}",
		},
	}

	for i, test := range testCases {
		val, err := EncodeValue(test.Val)
		if err != nil {
			t.Fatalf("Case [%d]: Got unexpected error %s!", i, err)
		}

		switch cast := val.ValueType.(type) {
		case *proto.Data_Value_JsonValue:
			if test.Expected != cast.JsonValue {
				t.Errorf("Case [%d]: Expected %s, got %s", i, test.Expected, cast.JsonValue)
			}
		default:
			t.Errorf("Case [%d]: Got unexpected type back %+v!", i, cast)
		}
	}
}

func boolPtr(val bool) *bool {
	return &val
}

func float32Ptr(val float32) *float32 {
	return &val
}

func float64Ptr(val float64) *float64 {
	return &val
}

func strPtr(val string) *string {
	return &val
}

func intPtr(val int) *int {
	return &val
}

func int8Ptr(val int8) *int8 {
	return &val
}

func int16Ptr(val int16) *int16 {
	return &val
}

func int32Ptr(val int32) *int32 {
	return &val
}

func int64Ptr(val int64) *int64 {
	return &val
}

func uintPtr(val uint) *uint {
	return &val
}

func uint8Ptr(val uint8) *uint8 {
	return &val
}

func uint16Ptr(val uint16) *uint16 {
	return &val
}

func uint32Ptr(val uint32) *uint32 {
	return &val
}

func uint64Ptr(val uint64) *uint64 {
	return &val
}
