package metrics

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/healthcareblocks/ec2_metrics_publisher/mathlib"
	"github.com/healthcareblocks/ec2_metrics_publisher/slice"
)

// MemoryInfo stores memory metrics
type MemoryInfo struct {
	LookupPath string
	Metrics    map[string]float64
	Timestamp  time.Time
}

// NewMemoryInfoWithDefaultLookup returns MemoryInfo struct with
// Lookup field set to the system default (/proc/meminfo)
func NewMemoryInfoWithDefaultLookup() *MemoryInfo {
	return &MemoryInfo{LookupPath: "/proc/meminfo"}
}

// Collect populates a MemoryInfo object;
// filters is a slice of field names that is used to populate m.Metrics
func (m *MemoryInfo) Collect(filters []string) error {
	if m.LookupPath == "" {
		return errors.New("must set LookupPath to location of memory file")
	}

	content, err := ioutil.ReadFile(m.LookupPath)
	if err != nil {
		return err
	}

	m.Metrics = make(map[string]float64)
	m.Timestamp = time.Now()

	for _, element := range strings.Split(string(content), "\n") {
		if element == "" {
			continue
		}
		data := strings.Split(element, ":")
		key := data[0]
		value := strings.Trim(data[1], "kB ")

		switch key {
		// we only care about a subset of all the memory metrics
		case "MemFree", "MemTotal", "SwapFree", "SwapTotal":
			v, err := stringToFloat(value)
			if err != nil {
				return err
			}
			m.Metrics[key] = v
		}
	}

	m.Metrics["Used"] = m.Metrics["MemTotal"] - m.Metrics["MemAvailable"]
	m.Metrics["UsedPercent"] = mathlib.RoundWithPrecision((m.Metrics["Used"]/m.Metrics["MemTotal"])*100, 5)
	m.Metrics["SwapUsed"] = m.Metrics["SwapTotal"] - m.Metrics["SwapFree"]

	if m.Metrics["SwapTotal"] == 0 {
		m.Metrics["SwapUsedPercent"] = 0
	} else {
		m.Metrics["SwapUsedPercent"] = mathlib.RoundWithPrecision((m.Metrics["SwapUsed"]/m.Metrics["SwapTotal"])*100, 5)
	}

	// intentionally filter out specific metrics after the above calculations have been performed
	if len(filters) > 0 && filters[0] != "" {
		for key := range m.Metrics {
			if !slice.ContainsString(filters, key) {
				delete(m.Metrics, key)
			}
		}
	}

	return nil
}

func stringToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
