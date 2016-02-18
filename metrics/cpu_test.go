package metrics

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPUUsageCollectRequiresLookupPath(t *testing.T) {
	err := new(CPUUsage).Collect(nil)
	assert.EqualError(t, err, "must set LookupPath to location of CPU file")
}

func TestCPUUsageCollectUpdatesStats(t *testing.T) {
	curDir, _ := os.Getwd()
	fakeStatsFile := curDir + "/../testutils/fake_proc_stat"
	cpu := CPUUsage{LookupPath: fakeStatsFile}
	err := cpu.Collect(nil)

	assert.NoError(t, err)
	assert.Len(t, cpu.Metrics, 3)
	assert.EqualValues(t, 0, cpu.Metrics["cpu"])
	assert.NotNil(t, cpu.Timestamp)
}

func TestCPUUsageCollectWithFilters(t *testing.T) {
	curDir, _ := os.Getwd()
	fakeStatsFile := curDir + "/../testutils/fake_proc_stat"
	cpu := CPUUsage{LookupPath: fakeStatsFile}
	err := cpu.Collect([]string{"cpu1"})

	assert.NoError(t, err)
	assert.Len(t, cpu.Metrics, 1)
	assert.NotContains(t, cpu.Metrics, "cpu")
	assert.EqualValues(t, 0, cpu.Metrics["cpu1"])
	assert.NotNil(t, cpu.Timestamp)
}

func TestCPUUsageCollectWithEmptyStringFilter(t *testing.T) {
	curDir, _ := os.Getwd()
	fakeStatsFile := curDir + "/../testutils/fake_proc_stat"
	cpu := CPUUsage{LookupPath: fakeStatsFile}
	err := cpu.Collect([]string{""})

	assert.NoError(t, err)
	assert.Len(t, cpu.Metrics, 3)
}

func TestCPUUsageCollectWithBadFilter(t *testing.T) {
	curDir, _ := os.Getwd()
	fakeStatsFile := curDir + "/../testutils/fake_proc_stat"
	cpu := CPUUsage{LookupPath: fakeStatsFile}
	err := cpu.Collect([]string{"Foo"})

	assert.NoError(t, err)
	assert.Len(t, cpu.Metrics, 0)
}
