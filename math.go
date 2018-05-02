package geobuf

import (
	"math"
)

var (
	MaxPrecision = uint32(math.Pow10(9))
)

func GetPrecision(point float64) uint32 {
	var e uint32 = 1
	for {
		base := math.Round(float64(point * float64(e)))
		if (base/float64(e)) != point && e < MaxPrecision {
			e = e * 10
		} else {
			break
		}
	}
	return e
}

func IntWithPrecision(point float64, precision uint32) int64 {
	return int64(math.Round(point * float64(precision)))
}

func FloatWithPrecision(point int64, precision uint32) float64 {
	return float64(point) / float64(precision)
}
