package libs

import (
	"math"
)

func FloatIsEqual(f1, f2, p float64) bool {
	return math.Dim(f1, f2) < p
}
