package testutils

import (
	"time"

	"github.com/healthcareblocks/ec2_metrics_publisher/metrics"
)

type dataFactory struct {
	System metrics.System
}

// Factory contains test data
var Factory = dataFactory{
	System: metrics.System{
		CPU: metrics.CPUUsage{
			Metrics:   map[string]float64{"cpu": 25.0, "cpu1": 50.0, "cpu2": 0},
			Timestamp: time.Now(),
		},
		Memory: metrics.MemoryInfo{
			Metrics: map[string]float64{
				"MemFree":         50.0,
				"MemTotal":        100.0,
				"Used":            50.0,
				"UsedPercent":     50.0,
				"SwapFree":        10.0,
				"SwapTotal":       5.0,
				"SwapUsed":        5.0,
				"SwapUsedPercent": 50.0,
			},
			Timestamp: time.Now(),
		},
		Volumes: []metrics.Volume{
			{
				LookupPath: "/",
				Metrics: map[string]float64{
					"Available":         100.0,
					"Free":              100.0,
					"Size":              200.0,
					"Used":              100.0,
					"UsedPercent":       50.0,
					"INodesTotal":       1000000.0,
					"INodesUsed":        500000.0,
					"INodesUsedPercent": 50.0,
				},
				Timestamp: time.Now(),
			},
			{
				LookupPath: "/data",
				Metrics: map[string]float64{
					"Available":         60.0,
					"Free":              60.0,
					"Size":              100.0,
					"Used":              40.0,
					"UsedPercent":       60.0,
					"INodesTotal":       1000000.0,
					"INodesUsed":        500000.0,
					"INodesUsedPercent": 50.0,
				},
				Timestamp: time.Now(),
			},
		},
	},
}
