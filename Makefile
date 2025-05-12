STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
AWS_ACCOUNT ?= "891377056770"
AWS_REGION ?= us-west-2
IMAGE_TAG ?= latest
GIT_COMMIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin -tags "${BUILD_TAGS}" *.go

docker-build-steampipe:
	@docker build \
		-t $(AWS_ACCOUNT).dkr.ecr.$(AWS_REGION).amazonaws.com/steampipe:$(IMAGE_TAG) \
		-f docker/steampipe/Dockerfile \
		.