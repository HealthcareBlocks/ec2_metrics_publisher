package destination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/healthcareblocks/ec2_metrics_publisher/_testutils"
	"github.com/healthcareblocks/ec2_metrics_publisher/metadata"
	"github.com/healthcareblocks/ec2_metrics_publisher/metrics"
	"github.com/stretchr/testify/assert"
)

// local test server
var cloudwatchServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}))

func TestCloudwatchWithNilEndpoint(t *testing.T) {
	cw := getCloudwatchInstance()
	cw.SetEndpoint(nil)
	err := cw.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set service endpoint")
}

func TestCloudwatchWithEmptyMachine(t *testing.T) {
	cw := getCloudwatchInstance()
	cw.SetMachine(&metadata.Machine{})
	err := cw.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set machine instance and region")
}

func TestCloudwatchWithNilMachine(t *testing.T) {
	cw := getCloudwatchInstance()
	cw.SetMachine(nil)
	err := cw.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set machine instance and region")
}

func TestCloudwatchWithoutAnyData(t *testing.T) {
	cw := getCloudwatchInstance()
	err := cw.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set at least one metric")
}

func TestCloudwatchSendMessageSuccessfully(t *testing.T) {
	cw := getCloudwatchInstance()
	err := cw.SendMessage(&testutils.Factory.System)
	assert.NoError(t, err)
}

func getCloudwatchInstance() *Cloudwatch {
	cw := &Cloudwatch{}
	cw.SetMachine(&metadata.Machine{Instance: "i-abc123", Region: "us-west-2"})
	cw.SetEndpoint(&ServiceEndpoint{cloudwatchServer.URL})
	return cw
}
