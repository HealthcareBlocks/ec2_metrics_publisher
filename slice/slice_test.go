package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceContainsString(t *testing.T) {
	pets := []string{"gerbils", "hamsters", "puppies"}
	assert.True(t, ContainsString(pets, "hamsters"))
	assert.False(t, ContainsString(pets, "gila monsters"))
}
