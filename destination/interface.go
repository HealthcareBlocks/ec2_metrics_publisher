package destination

import (
	"github.com/healthcareblocks/ec2_metrics_publisher/metadata"
	"github.com/healthcareblocks/ec2_metrics_publisher/metrics"
)

// Service defines the interface for objects that send system metrics
type Service interface {
	SetEndpoint(*ServiceEndpoint)
	SetMachine(*metadata.Machine)
	SendMessage(*metrics.System) error
	String() string
}
