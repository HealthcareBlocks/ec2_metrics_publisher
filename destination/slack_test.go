package destination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/healthcareblocks/ec2_metrics_publisher/metadata"
	"github.com/healthcareblocks/ec2_metrics_publisher/metrics"
	"github.com/healthcareblocks/ec2_metrics_publisher/testutils"
	"github.com/stretchr/testify/assert"
)

// local test server
var slackServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}))

func TestSlackWithNilEndpoint(t *testing.T) {
	s := getSlackInstance()
	s.SetEndpoint(nil)
	err := s.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set service endpoint")
}

func TestSlackWithEmptyMachine(t *testing.T) {
	s := getSlackInstance()
	s.SetMachine(&metadata.Machine{})
	err := s.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set machine instance and region")
}

func TestSlackWithNilMachine(t *testing.T) {
	s := getSlackInstance()
	s.SetMachine(nil)
	err := s.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set machine instance and region")
}

func TestSlackWithoutAnyData(t *testing.T) {
	s := getSlackInstance()
	err := s.SendMessage(&metrics.System{})
	assert.EqualError(t, err, "must set at least one metric")
}

func TestSlackSendMessageSuccessfully(t *testing.T) {
	s := getSlackInstance()
	err := s.SendMessage(&testutils.Factory.System)
	assert.NoError(t, err)
}

func getSlackInstance() *Slack {
	s := &Slack{}
	s.SetMachine(&metadata.Machine{Instance: "i-abc123", Region: "us-west-2"})
	s.SetEndpoint(&ServiceEndpoint{slackServer.URL})
	return s
}
