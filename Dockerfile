FROM healthcareblocks/alpine:latest

COPY bin/ec2_metrics_publisher-linux-amd64 /bin/ec2_metrics_publisher
ENTRYPOINT ["ec2_metrics_publisher"]
