STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin -tags "${BUILD_TAGS}" *.go

# Exclude Parliament IAM permissions
dev:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin -tags "dev ${BUILD_TAGS}" *.go