package destination

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/healthcareblocks/ec2_metrics_publisher/metadata"
	"github.com/healthcareblocks/ec2_metrics_publisher/metrics"
)

// Cloudwatch implements the Service interface, wrapping AWS Cloudwatch functionality
type Cloudwatch struct {
	*metadata.Machine
	*ServiceEndpoint
}

// MetricData is a slice of cloudwatch.MetricDatum
type MetricData []*cloudwatch.MetricDatum

// String returns the name of the service
func (cw Cloudwatch) String() string {
	return "Cloudwatch"
}

// SetEndpoint sets the service endpoint
func (cw *Cloudwatch) SetEndpoint(endpoint *ServiceEndpoint) {
	cw.ServiceEndpoint = endpoint
}

// SetMachine sets the metadata.Machine object
func (cw *Cloudwatch) SetMachine(machine *metadata.Machine) {
	cw.Machine = machine
}

// SendMessage posts message to Cloudwatch API
func (cw *Cloudwatch) SendMessage(sys *metrics.System) error {
	if cw.ServiceEndpoint == nil {
		return errors.New("must set service endpoint")
	}

	if cw.Machine == nil || cw.Machine.IsEmpty() {
		return errors.New("must set machine instance and region")
	}

	var data MetricData

	data.add(sys.CPU.Metrics, cw.Instance, sys.CPU.Timestamp,
		cloudwatch.StandardUnitPercent, "CPUUtilization", "")

	data.add(sys.Memory.Metrics, cw.Instance, sys.Memory.Timestamp,
		cloudwatch.StandardUnitKilobytes, "Memory", "")

	for _, v := range sys.Volumes {
		data.add(v.Metrics, cw.Instance, v.Timestamp,
			cloudwatch.StandardUnitKilobytes, "Volume", v.LookupPath)
	}

	if len(data) == 0 {
		return errors.New("must set at least one metric")
	}

	metric := &cloudwatch.PutMetricDataInput{
		MetricData: data,
		Namespace:  aws.String("System/Linux"),
	}

	awsConfig := &aws.Config{Endpoint: aws.String(cw.URL), Region: aws.String(cw.Region)}
	api := cloudwatch.New(session.New(), awsConfig)
	if _, err := api.PutMetricData(metric); err != nil {
		return err
	}

	return nil
}

func (data *MetricData) add(metrics map[string]float64, instance string,
	timestamp time.Time, defaultUnit string, title string, optionalTitleSuffix string) {

	var metricName, unit string

	for key, value := range metrics {
		if optionalTitleSuffix != "" {
			// e.g. Volume_Used_/data
			metricName = fmt.Sprintf("%s_%s_%s", title, key, optionalTitleSuffix)
		} else {
			// e.g. Memory_Used
			metricName = fmt.Sprintf("%s_%s", title, key)
		}

		switch key {
		case "INodesUsedPercent", "SwapUsedPercent", "UsedPercent":
			unit = cloudwatch.StandardUnitPercent
		case "INodesTotal", "INodesUsed":
			unit = cloudwatch.StandardUnitCount
		default:
			unit = defaultUnit
		}

		datum := &cloudwatch.MetricDatum{
			MetricName: aws.String(metricName),
			Value:      aws.Float64(value),
			Unit:       aws.String(unit),
			Dimensions: []*cloudwatch.Dimension{
				{
					Name:  aws.String("InstanceId"),
					Value: aws.String(instance),
				},
			},
			Timestamp: aws.Time(timestamp),
		}
		*data = append(*data, datum)
	}
}
