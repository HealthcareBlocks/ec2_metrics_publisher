// Amazon Web Services EC2 agent that collects CPU, memory, and disk metrics,
// publishing them to AWS CloudWatch and/or Slack.
//
// See https://github.com/HealthcareBlocks/ec2_metrics_publisher for usage details.
//
package main // import "github.com/healthcareblocks/ec2_metrics_publisher"

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/healthcareblocks/ec2_metrics_publisher/destination"
	"github.com/healthcareblocks/ec2_metrics_publisher/logger"
	"github.com/healthcareblocks/ec2_metrics_publisher/metadata"
	"github.com/healthcareblocks/ec2_metrics_publisher/metrics"
)

var (
	// Required flag
	metricsDestinations = flag.String("destinations", "", "Destinations to publish to")

	// Required if destinations includes slack
	slackHook = flag.String("slack-hook", "", "Slack hook URL. Required if -destinations includes slack.")

	// Optional flags with default values
	interval         = flag.Int("interval", 60, "Frequency of publishing metrics to destination(s)")
	metricsToCollect = flag.String("metrics", "cpu,memory,volume", "Metrics to collect")
	volumePaths      = flag.String("paths", "/", "Volume paths to calculate usage")
	appVersion       = flag.Bool("v", false, "Prints version of this app and exits")

	// Optional flags - setting these will prevent an EC2 metadata lookup (useful for testing)
	instanceID = flag.String("instance", "", "EC2 instance ID")
	region     = flag.String("region", "", "EC2 region of this instance")

	// Logrus logging object
	log = logger.NewLogger()

	// Used for registering active destinations for publishing metrics
	destinations []destination.Service
)

func main() {
	if err := parseFlags(nil); err != nil {
		log.Fatal(err.Error())
	}
	go periodicStatsCollection(time.Duration(*interval) * time.Second)
	select {}
}

func parseFlags(awsConfig *aws.Config) error {
	flag.Parse()

	if *appVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *metricsDestinations == "" {
		return errors.New("-destinations cannot be blank")
	}

	if *metricsToCollect == "" {
		return errors.New("-metrics cannot be blank")
	}

	if strings.Contains(*metricsToCollect, "volume") && *volumePaths == "" {
		return errors.New("-paths cannot be blank")
	}

	var acceptedMetrics = regexp.MustCompile(`(cpu|memory|volume)+(\[[\w,]+\])?`)
	if !acceptedMetrics.MatchString(*metricsToCollect) {
		return errors.New("-metrics contains an unrecognized item")
	}

	machine := &metadata.Machine{Instance: *instanceID, Region: *region}
	if machine.IsEmpty() {
		if err := machine.LoadFromMetadata(awsConfig); err != nil {
			log.Fatal(err.Error())
		}
	}

	if strings.Contains(*metricsDestinations, "cloudwatch") {
		// blank url results in the default aws service endpoint
		endpoint := &destination.ServiceEndpoint{URL: ""}

		cw := &destination.Cloudwatch{}
		cw.SetEndpoint(endpoint)
		cw.SetMachine(machine)

		destinations = append(destinations, cw)
		log.Info("Registered CloudWatch Service")
	}

	if strings.Contains(*metricsDestinations, "slack") {
		if *slackHook == "" {
			return errors.New("must set -slack-hook")
		}

		endpoint := &destination.ServiceEndpoint{URL: *slackHook}

		slack := &destination.Slack{}
		slack.SetEndpoint(endpoint)
		slack.SetMachine(machine)

		destinations = append(destinations, slack)
		log.Info("Registered Slack Service")
	}

	if len(destinations) == 0 {
		return errors.New("set at least one valid destination via -destinations")
	}

	return nil
}

func collectAndPublishStats() {
	sys := &metrics.System{}
	var err error

	var cpu = regexp.MustCompile(`cpu(\[([\w,]+)\])?`)
	var memory = regexp.MustCompile(`memory(\[([\w,]+)\])?`)
	var volume = regexp.MustCompile(`volume(\[([\w,]+)\])?`)

	for _, metric := range strings.Split(*metricsToCollect, ",") {
		if cpu.MatchString(metric) {
			filters := strings.Split(cpu.FindStringSubmatch(metric)[2], ",")
			err = sys.CollectCPUInfo(filters)
		}

		if memory.MatchString(metric) {
			filters := strings.Split(memory.FindStringSubmatch(metric)[2], ",")
			err = sys.CollectMemoryInfo(filters)
		}

		if volume.MatchString(metric) {
			filters := strings.Split(volume.FindStringSubmatch(metric)[2], ",")
			err = sys.CollectVolumesInfo(*volumePaths, filters)
		}

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	for _, item := range destinations {
		if err := item.SendMessage(sys); err != nil {
			log.Fatal(err.Error())
		}
		log.Infof("Message sent to %v", item)
	}
}

func periodicStatsCollection(n time.Duration) {
	// collect stats and publish immediately
	collectAndPublishStats()

	// now start a ticker for next iteration and beyond
	for range time.Tick(n) {
		collectAndPublishStats()
	}
}
