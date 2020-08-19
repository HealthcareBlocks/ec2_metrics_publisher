package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/healthcareblocks/ec2_metrics_publisher/_testutils"
	"github.com/stretchr/testify/assert"
)

var (
	defaultMetricsDestinations string
	defaultSlackHook           string
	defaultInterval            int
	defaultMetricsToCollect    string
	defaultVolumePaths         string
)

func TestFlags(t *testing.T) {
	cases := []struct {
		input  []string
		expect string
	}{
		{
			[]string{},
			"-destinations cannot be blank",
		},
		{
			[]string{"-destinations="},
			"-destinations cannot be blank",
		},
		{
			[]string{"-destinations=cloudwatch", "-metrics="},
			"-metrics cannot be blank",
		},
		{
			[]string{"-destinations=invalid_item"},
			"set at least one valid destination via -destinations",
		},
		{
			[]string{"-destinations=cloudwatch", "-paths="},
			"-paths cannot be blank",
		},
		{
			[]string{"-destinations=cloudwatch", "-metrics=volume", "-paths="},
			"-paths cannot be blank",
		},
		{
			[]string{"-destinations=cloudwatch", "-metrics=foo"},
			"-metrics contains an unrecognized item",
		},
		{
			[]string{"-destinations=slack", "-slack-hook="},
			"must set -slack-hook",
		},
	}

	defaultMetricsDestinations = *metricsDestinations
	defaultSlackHook = *slackHook
	defaultInterval = *interval
	defaultMetricsToCollect = *metricsToCollect
	defaultVolumePaths = *volumePaths

	originalArgs := os.Args

	server := testutils.EC2MetadataResponseStub()
	defer server.Close()

	for _, test := range cases {
		os.Args = append(originalArgs, test.input...)
		assert.EqualError(t, parseFlags(&aws.Config{Endpoint: aws.String(server.URL + "/latest")}), test.expect)
		resetFlags()
	}
}

func TestNestedFlag(t *testing.T) {
	server := testutils.EC2MetadataResponseStub()
	defer server.Close()

	os.Args = append(os.Args, "-destinations=cloudwatch", "-metrics=memory[UsedPercent]")
	err := parseFlags(&aws.Config{Endpoint: aws.String(server.URL + "/latest")})
	assert.NoError(t, err)
}

func TestMetricsRegex(t *testing.T) {
	cases := []struct {
		input  string
		expect bool
	}{
		{
			"cpu[Free]",
			true,
		},
		{
			"cpu[Free,Available]",
			true,
		},
		{
			"cpu",
			true,
		},
		{
			"memory",
			false,
		},
	}

	var cpu = regexp.MustCompile(`cpu(\[([\w,]+)\])?`)
	for _, test := range cases {
		actual := cpu.MatchString(test.input)
		assert.Equal(t, test.expect, actual)
	}
}

func resetFlags() {
	*metricsDestinations = defaultMetricsDestinations
	*slackHook = defaultSlackHook
	*interval = defaultInterval
	*metricsToCollect = defaultMetricsToCollect
	*volumePaths = defaultVolumePaths
}
