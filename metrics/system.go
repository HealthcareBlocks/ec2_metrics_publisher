package metrics

import (
	"errors"
	"strings"
)

// System is a container for system-wide metrics
type System struct {
	CPU     CPUUsage
	Memory  MemoryInfo
	Volumes []Volume
}

// Collectable defines the interface for objects that collect system metrics
type Collectable interface {
	Collect(filters string) error
}

const metricsNotFound = "no metrics were found, check your filters"

// CollectCPUInfo collects and populates CPU objects
func (sys *System) CollectCPUInfo(filters []string) error {
	cpuUsage := NewCPUUsageWithDefaultLookup()
	if err := cpuUsage.Collect(filters); err != nil {
		return err
	}
	if len(cpuUsage.Metrics) == 0 {
		return errors.New(metricsNotFound)
	}
	sys.CPU = *cpuUsage
	return nil
}

// CollectMemoryInfo collects and populates Memory object
func (sys *System) CollectMemoryInfo(filters []string) error {
	memoryInfo := NewMemoryInfoWithDefaultLookup()
	if err := memoryInfo.Collect(filters); err != nil {
		return err
	}
	if len(memoryInfo.Metrics) == 0 {
		return errors.New(metricsNotFound)
	}
	sys.Memory = *memoryInfo
	return nil
}

// CollectVolumesInfo collects and populates Volume objects
func (sys *System) CollectVolumesInfo(volumePaths string, filters []string) error {
	for _, path := range strings.Split(volumePaths, ",") {
		vol := &Volume{LookupPath: path}
		if err := vol.Collect(filters); err != nil {
			return err
		}
		if len(vol.Metrics) == 0 {
			return errors.New(metricsNotFound)
		}
		sys.Volumes = append(sys.Volumes, *vol)
	}
	return nil
}
