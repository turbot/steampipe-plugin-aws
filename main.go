package main

import (
	"github.com/turbot/steampipe-plugin-aws/aws"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: aws.Plugin})
}
