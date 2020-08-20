# EC2 Metrics Publisher

Standalone agent that collects *cpu*, *memory*, and *storage volume* metrics and publishes them to a destination, such as AWS CloudWatch and/or Slack. See **Usage Requirements** section for caveats.

## Example Usage

Publish All Metrics to CloudWatch Every Minute:
```
ec2_metrics_publisher -destinations=cloudwatch
```

Adding an ```-interval``` flag allows you to override the default 60 second ticker.

Publish All Metrics to CloudWatch and Slack every 5 minutes:
```
ec2_metrics_publisher -destinations=cloudwatch,slack -interval=300 -slack-hook=https://hooks.slack.com/etc...
```

Publish a Specific Metric to Specific Services:
```
ec2_metrics_publisher -destinations=cloudwatch -metrics=memory,volume
ec2_metrics_publisher -destinations=slack -metrics=cpu -slack-hook=https://hooks.slack.com/etc...
```

Publish Specific Metric Field Names:
```
ec2_metrics_publisher -destinations=cloudwatch -metrics=memory[SwapUsedPercent,UsedPercent],volume[UsedPercent]
```

Tracking Usage of More than One Volume:
```
ec2_metrics_publisher -destinations=cloudwatch -paths=/dev/xvda1,/dev/xvdf
```

## Usage Requirements

### Compatibility

#### Linux Kernel Required
CPU and memory metrics are retrieved via the /proc/ virtual file system, which only exists on Linux systems.

### AWS CloudWatch

* The runtime agent needs to have access to host's EC2 instance metadata.
* AWS credentials or IAM role associated with machine need to have the **cloudwatch:PutMetricData** permission set in an [IAM policy](http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/UsingIAM.html).
* See [AWS's pricing page](https://aws.amazon.com/cloudwatch/pricing/) to understand how custom metrics are priced

### Slack

* Configure an Incoming WebHook integration; the resulting URL will be used as the *-slack-hook* parameter.

### Running in Production

Using a Systemd, SystemV, Supervisor, or Upstart script is the recommended way for controlling this agent as a background process. If running as a Docker container, setting the ```restart='always'``` policy at runtime is recommended.

## Common Issues

*"NoCredentialProviders: no valid providers in chain."*

Posting data to the CloudWatch API requires that your EC2 host machine is either (a) associated with an IAM profile that has permission to post to the CloudWatch API; or (b) has the proper AWS credentials configured. One way is via environment variables:

* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY
* AWS_SESSION_TOKEN *(required if using session tokens)*
