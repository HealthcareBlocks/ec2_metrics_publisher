SHELL := /bin/bash -e
.DEFAULT_GOAL := build

build_native:
	go build -o bin/ec2_metrics_publisher

build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/ec2_metrics_publisher-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o bin/ec2_metrics_publisher-linux-arm64

.PHONY: build_native build_linux
