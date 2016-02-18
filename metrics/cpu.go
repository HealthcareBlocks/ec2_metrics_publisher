package metrics

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/healthcareblocks/ec2_metrics_publisher/slice"
)

// CPUUsage returns usage for a map of cpu names
// Note: "cpu" represents all processors; cpu0, cpu1, etc. are individual ones
type CPUUsage struct {
	LookupPath string
	Metrics    map[string]float64
	Timestamp  time.Time
}

// CPUStats stores CPU statistics;
// Ref: /proc/stat (http://man7.org/linux/man-pages/man5/proc.5.html)
type CPUStats struct {
	// Time spent in user mode
	User uint64
	// Time spent in user mode with low priority (nice)
	Nice uint64
	// Time spent in system mode
	System uint64
	// Time spent in the idle task.  This value
	// should be USER_HZ times the second entry in the
	// /proc/uptime pseudo-file.
	Idle uint64
	// Time waiting for I/O to complete
	IOWait uint64
	// Time servicing interrupts
	IRQ uint64
	// Time servicing softirqs
	SoftIRQ uint64
	// Stolen time, which is the time spent in
	// other operating systems when running in a
	// virtualized environment
	Steal uint64
}

// NewCPUUsageWithDefaultLookup returns CPUUsage struct with
// Lookup field set to the system default (/proc/meminfo)
func NewCPUUsageWithDefaultLookup() *CPUUsage {
	return &CPUUsage{LookupPath: "/proc/stat"}
}

// Collect populates a CPUUsage object;
// filters is a slice of field names that is used to populate c.Metrics
func (c *CPUUsage) Collect(filters []string) error {
	if c.LookupPath == "" {
		return errors.New("must set LookupPath to location of CPU file")
	}

	c.Metrics = make(map[string]float64)
	c.Timestamp = time.Now()

	stats0, err := getCPUStats(c)
	if err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond)

	stats1, err := getCPUStats(c)
	if err != nil {
		return err
	}

	for cpuName, cpuStats1 := range stats1 {
		if len(filters) == 0 || filters[0] == "" || slice.ContainsString(filters, cpuName) {
			cpuStats0 := stats0[cpuName]
			total0 := cpuStats0.getRealIdleTime() + cpuStats0.getNonIdleTime()
			total1 := cpuStats1.getRealIdleTime() + cpuStats1.getNonIdleTime()
			deltaTotal := float64(total1 - total0)
			deltaIdle := float64(cpuStats1.getRealIdleTime() - cpuStats0.getRealIdleTime())

			c.Metrics[cpuName] = 0
			if deltaTotal > 0 {
				c.Metrics[cpuName] = 100 * (deltaTotal - deltaIdle) / deltaTotal
			}
		}
	}

	return nil
}

func getCPUStats(c *CPUUsage) (map[string]CPUStats, error) {
	content, err := ioutil.ReadFile(c.LookupPath)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]CPUStats)

	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if fields[0][:3] == "cpu" {
			cpu := CPUStats{}
			name := fields[0]

			for i := 1; i < len(fields); i++ {
				v, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					return nil, err
				}
				switch i {
				case 1:
					cpu.User = v
				case 2:
					cpu.Nice = v
				case 3:
					cpu.System = v
				case 4:
					cpu.Idle = v
				case 5:
					cpu.IOWait = v
				case 6:
					cpu.IRQ = v
				case 7:
					cpu.SoftIRQ = v
				case 8:
					cpu.Steal = v
				}
			}

			metrics[name] = cpu
		}
	}

	return metrics, nil
}

func (c CPUStats) getRealIdleTime() uint64 {
	return c.Idle + c.IOWait
}

func (c CPUStats) getNonIdleTime() uint64 {
	return c.User + c.Nice + c.System + c.IRQ + c.SoftIRQ + c.Steal
}
