STEAMPIPE_INSTALL_DIR ?= ~/.steampipe

# Determine the operating system
OS := $(shell uname)

# Check if the OS is Mac OS/Darwin
ifeq ($(OS),Darwin)
  BUILD_TAGS = netgo
endif

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin -tags "$(BUILD_TAGS)" *.go