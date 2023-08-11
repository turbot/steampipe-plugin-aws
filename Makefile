STEAMPIPE_INSTALL_DIR?=~/.steampipe

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin *.go
