package metadata

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/healthcareblocks/ec2_metrics_publisher/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMachineIsEmpty(t *testing.T) {
	assert.True(t, new(Machine).IsEmpty())

	machineWithoutInstance := &Machine{Region: "us-west-2"}
	assert.True(t, machineWithoutInstance.IsEmpty())

	machineWithoutRegion := &Machine{Instance: "i-abc123"}
	assert.True(t, machineWithoutRegion.IsEmpty())
}

func TestGetMachineRetrievesEC2Metadata(t *testing.T) {
	server := testutils.EC2MetadataResponseStub()
	defer server.Close()

	machine := &Machine{}
	err := machine.LoadFromMetadata(&aws.Config{Endpoint: aws.String(server.URL + "/latest")})
	assert.NoError(t, err)
	assert.Equal(t, "i-12345", machine.Instance)
	assert.Equal(t, "us-west-2", machine.Region)
}

func TestGetMachineReturnsErrorOutsideOfEC2(t *testing.T) {
	server := testutils.InvalidEC2MetadataResponseStub()
	defer server.Close()

	machine := &Machine{}
	err := machine.LoadFromMetadata(&aws.Config{Endpoint: aws.String(server.URL + "/latest")})
	assert.Error(t, err)
}
