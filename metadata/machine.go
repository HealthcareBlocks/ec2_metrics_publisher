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
	// Internal awsConfig reference
	awsConfig *aws.Config
}

// NewMachine returns a new Machine pointer. The input parameters are based
// on the Machine struct fields. Passing in a nil aws.Config object uses the AWS SDK defaults.
//
//	   m := ebs.NewMachine(nil)
//
func NewMachine(awsConfig *aws.Config) *Machine {
	return &Machine{awsConfig: awsConfig}
}

// IsEmpty returns true if Instance or Region have not been set
func (machine *Machine) IsEmpty() bool {
	return machine.Instance == "" || machine.Region == ""
}

// IsEC2 returns true if this instance is represents an EC2 machine
func (machine *Machine) IsEC2() bool {
	if err := machine.LoadFromMetadata(); err != nil {
		return false
	}

	return true
}

// LoadFromMetadata reloads a Machine struct with EC2 metadata
func (machine *Machine) LoadFromMetadata() error {
	svc := ec2metadata.New(session.New(), machine.awsConfig)

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

// WithInstance returns machine with its Instance field set
func (machine *Machine) WithInstance(instance string) *Machine {
	machine.Instance = instance
	return machine
}

// WithRegion returns machine with its Region field set
func (machine *Machine) WithRegion(region string) *Machine {
	machine.Region = region
	return machine
}
