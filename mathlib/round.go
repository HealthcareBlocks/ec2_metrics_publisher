// Package mathlib provides additional mathematical functions
package mathlib

import (
	"math"
)

// Round is a rounding function for float64's
func Round(f float64) float64 {
	return math.Floor(f + .5)
}

// RoundWithPrecision is a rounding function for float64's with precision
func RoundWithPrecision(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}
