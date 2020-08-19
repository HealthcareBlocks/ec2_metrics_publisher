package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryInfoCollectRequiresLookupPath(t *testing.T) {
	err := new(MemoryInfo).Collect(nil)
	assert.EqualError(t, err, "must set LookupPath to location of memory file")
}

func TestMemoryInfoCollectUpdatesStats(t *testing.T) {
	fakeMemFile := "../testdata/fake_proc_meminfo"
	mem := MemoryInfo{LookupPath: fakeMemFile}
	err := mem.Collect(nil)

	assert.NoError(t, err)
	assert.Len(t, mem.Metrics, 8)
	assert.NotEmpty(t, mem.Metrics["MemFree"])
	assert.NotEmpty(t, mem.Metrics["MemTotal"])
	assert.NotEmpty(t, mem.Metrics["Used"])
	assert.NotEmpty(t, mem.Metrics["UsedPercent"])
	assert.NotEmpty(t, mem.Metrics["SwapFree"])
	assert.NotEmpty(t, mem.Metrics["SwapTotal"])
	assert.NotEmpty(t, mem.Metrics["SwapUsed"])
	assert.NotEmpty(t, mem.Metrics["SwapUsedPercent"])
	assert.NotNil(t, mem.Timestamp)
}

func TestMemoryInfoCollectWithFilters(t *testing.T) {
	fakeMemFile := "../testdata/fake_proc_meminfo"
	mem := MemoryInfo{LookupPath: fakeMemFile}
	err := mem.Collect([]string{"Used", "SwapUsed"})

	assert.NoError(t, err)
	assert.Len(t, mem.Metrics, 2)
	assert.NotContains(t, mem.Metrics, "MemFree")
	assert.NotEmpty(t, mem.Metrics["Used"])
	assert.NotEmpty(t, mem.Metrics["SwapUsed"])
	assert.NotNil(t, mem.Timestamp)
}

func TestMemoryInfoCollectWithEmptyStringFilter(t *testing.T) {
	fakeMemFile := "../testdata/fake_proc_meminfo"
	mem := MemoryInfo{LookupPath: fakeMemFile}
	err := mem.Collect([]string{""})

	assert.NoError(t, err)
	assert.Len(t, mem.Metrics, 8)
}

func TestMemoryInfoCollectWithBadFilter(t *testing.T) {
	fakeMemFile := "../testdata/fake_proc_meminfo"
	mem := MemoryInfo{LookupPath: fakeMemFile}
	err := mem.Collect([]string{"Foo"})

	assert.NoError(t, err)
	assert.Len(t, mem.Metrics, 0)
}
