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

Publish a Specific Metric to Specific Channels:
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

#### Docker
Packaging the agent inside a Docker image can be useful for distribution needs or testing on platforms that don't support /proc (e.g. Mac OS X). A [sample Dockerfile](Dockerfile) is included, and this project's [makefile](Makefile) shows how to compile the Go app for other platforms.

To collect stats from the **host** system while running the agent inside a Docker container, there are two important caveats:
* Enable the container to collect host stats via *--pid=host*
* Mount the host's devices to the container if you are collecting volume usage

Example:
```
docker run -d --pid=host -v /dev/xvdf:/dev/xvdf \
healthcareblocks/ec2_metrics_publisher -destinations=cloudwatch -paths=/dev/xvdf
```

### AWS CloudWatch

* The runtime agent needs to have access to host's EC2 instance metadata. If the agent is running directly on the EC2 instance, this access is granted by default. Inside a Docker container, this is possible unless you've configured strict networking rules.
* AWS credentials or IAM role associated with machine need to have the **cloudwatch:PutMetricData** permission set in an [IAM policy](http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/UsingIAM.html).
* On your EC2 machine, you should avoid setting the Docker engine's DNS to an external service (e.g. --dns 8.8.8.8), as this will impact your container's ability to interact with EC2 metadata and IAM roles.
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

See [this post](http://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs) for details.

## Testing

Ensure your environment has ```GO15VENDOREXPERIMENT=1``` set and then run:
```
godep restore
```

Testing project only:
```
go test -v $(go list ./... | grep -v /vendor/)
```

Testing project and vendored dependencies:
```
go test -v ./...
```

Or use the included test helper script:
```
./test
./test -v
./test all
./test all -v
```

### Live Testing

When testing against your AWS or Slack account, you can test locally by using the included Docker image. This enables you to test the agent on systems that don't meet the compatibility requirements mentioned above.

To pull the existing Docker Hub image:
```
docker pull healthcareblocks/ec2_metrics_publisher
```

Or to build from source, run ```make build docker``` from the root of this project.

Now you can use Docker to test against different endpoints. Notice both the *-instance* and *-region* flags are set explicitly since a local non-EC2 machine does not have the concept of EC2 metadata.

**CloudWatch**
```
docker run -it --rm \
  -e AWS_SESSION_TOKEN=$AWS_SESSION_TOKEN \
  -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  healthcareblocks/ec2_metrics_publisher \
    -destinations=cloudwatch
    -interval 10 \
    -instance i-11223344 \
    -region us-west-2
```

**Slack**
```
docker run -it --rm healthcareblocks/ec2_metrics_publisher \
  -destinations=slack \
  -interval 10 \
  -slack-hook https://hooks.slack.com/services/YOUR-SLACKHOOK-ID-GOES-HERE \
  -instance i-11223344 \
  -region us-west-2
```

## Contributing

Feel free to submit a pull request containing additional destinations,  enhancements, and bug fixes.

Ensure all external packages are vendored:
```
GO15VENDOREXPERIMENT=1 godep save ./...
```
