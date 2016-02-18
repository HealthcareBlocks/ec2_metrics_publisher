package metrics

import (
	"errors"
	"syscall"
	"time"

	"github.com/healthcareblocks/ec2_metrics_publisher/mathlib"
	"github.com/healthcareblocks/ec2_metrics_publisher/slice"
)

// Volume contains storage volume metrics
type Volume struct {
	LookupPath string
	Metrics    map[string]float64
	Timestamp  time.Time
}

// Collect populates a Volume object;
// filters is a slice of field names that is used to populate v.Metrics
func (v *Volume) Collect(filters []string) error {
	if v.LookupPath == "" {
		return errors.New("must set Path")
	}

	v.Metrics = make(map[string]float64)
	v.Timestamp = time.Now()

	var stat syscall.Statfs_t
	if err := syscall.Statfs(v.LookupPath, &stat); err != nil {
		return err
	}

	bsize := uint64(stat.Bsize)
	available := float64(stat.Bavail * bsize)
	free := float64(stat.Bfree * bsize)
	size := float64(stat.Blocks * bsize)
	inodesTotal := float64(stat.Files)
	inodesFree := float64(stat.Ffree)
	inodesUsed := inodesTotal - inodesFree

	const KB = 1024
	v.Metrics["Available"] = available / KB
	v.Metrics["Free"] = free / KB
	v.Metrics["Size"] = size / KB
	v.Metrics["Used"] = (size - free) / KB
	v.Metrics["UsedPercent"] = mathlib.RoundWithPrecision(((size-free)/size)*100, 5)
	v.Metrics["INodesTotal"] = inodesTotal
	v.Metrics["INodesUsed"] = inodesUsed
	v.Metrics["INodesUsedPercent"] = mathlib.RoundWithPrecision((inodesUsed/inodesTotal)*100, 5)

	// intentionally filter out specific metrics after the above calculations have been performed
	if len(filters) > 0 && filters[0] != "" {
		for key := range v.Metrics {
			if !slice.ContainsString(filters, key) {
				delete(v.Metrics, key)
			}
		}
	}

	return nil
}
