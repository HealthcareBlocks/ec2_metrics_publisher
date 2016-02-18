SHELL := /bin/bash -e
NAMESPACE ?= healthcareblocks
.DEFAULT_GOAL := build

build: clean
	docker run --rm -it -v $(PWD):/src healthcareblocks/gobuild -o linux -a amd64

build_all: clean
	docker run --rm -it -v $(PWD):/src healthcareblocks/gobuild

docker:
	docker build -t $(NAMESPACE)/ec2_metrics_publisher .
	@docker images -f "dangling=true" -q | xargs docker rmi

clean:
	rm -fr ./bin/*

push_to_docker: tag_version
	version=$(shell docker run --rm healthcareblocks/ec2_metrics_publisher -v); \
	docker push $(NAMESPACE)/ec2_metrics_publisher:latest; \
	docker push $(NAMESPACE)/ec2_metrics_publisher:$$version;

tag_version:
	version=$(shell docker run --rm healthcareblocks/ec2_metrics_publisher -v); \
	docker tag $(NAMESPACE)/ec2_metrics_publisher $(NAMESPACE)/ec2_metrics_publisher:$$version

.PHONY: build build_all docker clean push_to_docker
