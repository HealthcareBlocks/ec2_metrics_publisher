// Package metadata handles EC2 metadata
package metadata

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Machine represents the current host executing this program
type Machine struct {
	// EC2 instance ID
	Instance string
	// EC2 region
	Region string
}

// IsEmpty returns true if Instance or Region have not been set
func (machine *Machine) IsEmpty() bool {
	return machine.Instance == "" || machine.Region == ""
}

// LoadFromMetadata populates a Machine struct based on EC2 metadata;
// Passing in a nil aws.Config object will result in using AWS SDK defaults
func (machine *Machine) LoadFromMetadata(awsConfig *aws.Config) error {
	svc := ec2metadata.New(session.New(), awsConfig)

	instance, err := svc.GetMetadata("instance-id")
	if err != nil {
		return err
	}

	region, err := svc.Region()
	if err != nil {
		return err
	}

	machine.Instance = instance
	machine.Region = region

	return nil
}
