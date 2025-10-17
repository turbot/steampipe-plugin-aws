STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
GIT_COMMIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin -tags "${BUILD_TAGS}" *.go

docker-build-steampipe:
	@docker build \
		--platform $(PLATFORM) \
		-t local/steampipe:latest \
		-f docker/steampipe/Dockerfile \
		.