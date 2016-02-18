package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVolumeCollectRequiresPath(t *testing.T) {
	err := new(Volume).Collect(nil)
	assert.EqualError(t, err, "must set Path")
}

func TestVolumeCollectUpdatesStats(t *testing.T) {
	v := Volume{LookupPath: "/"}
	err := v.Collect(nil)

	assert.NoError(t, err)
	assert.Len(t, v.Metrics, 8)
	assert.NotEmpty(t, v.Metrics["Available"])
	assert.NotEmpty(t, v.Metrics["Free"])
	assert.NotEmpty(t, v.Metrics["Size"])
	assert.NotEmpty(t, v.Metrics["Used"])
	assert.NotEmpty(t, v.Metrics["UsedPercent"])
	assert.NotEmpty(t, v.Metrics["INodesTotal"])
	assert.NotEmpty(t, v.Metrics["INodesUsed"])
	assert.NotEmpty(t, v.Metrics["INodesUsedPercent"])
	assert.NotNil(t, v.Timestamp)
}

func TestVolumeCollectWithFilters(t *testing.T) {
	v := Volume{LookupPath: "/"}
	err := v.Collect([]string{"Available", "INodesTotal"})

	assert.NoError(t, err)
	assert.Len(t, v.Metrics, 2)
	assert.NotContains(t, v.Metrics, "Free")
	assert.NotEmpty(t, v.Metrics["Available"])
	assert.NotEmpty(t, v.Metrics["INodesTotal"])
	assert.NotNil(t, v.Timestamp)
}

func TestVolumeCollectWithEmptyStringFilter(t *testing.T) {
	v := Volume{LookupPath: "/"}
	err := v.Collect([]string{""})

	assert.NoError(t, err)
	assert.Len(t, v.Metrics, 8)
}

func TestVolumeCollectWithBadFilter(t *testing.T) {
	v := Volume{LookupPath: "/"}
	err := v.Collect([]string{"Foo"})

	assert.NoError(t, err)
	assert.Len(t, v.Metrics, 0)
}
