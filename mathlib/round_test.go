package mathlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	cases := []struct {
		input  float64
		expect float64
	}{
		{1, 1},
		{1.1, 1.0},
		{1.5, 2.0},
		{1.6, 2.0},
		{2.0, 2.0},
	}

	for _, test := range cases {
		result := Round(test.input)
		assert.EqualValues(t, test.expect, result)
	}
}

func TestRoundWithPrecision(t *testing.T) {
	cases := []struct {
		input     float64
		precision int
		expect    float64
	}{
		{1, 1, 1},
		{1.0001, 1, 1.0},
		{1.0001, 3, 1.0},
		{1.0001, 4, 1.0001},
		{1.558, 2, 1.56},
	}

	for _, test := range cases {
		result := RoundWithPrecision(test.input, test.precision)
		assert.EqualValues(t, test.expect, result)
	}
}
