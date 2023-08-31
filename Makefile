install:
	go build -tags netgo -o ~/.steampipe/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin *.go
	# go build -o ~/.steampipe/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin *.go
