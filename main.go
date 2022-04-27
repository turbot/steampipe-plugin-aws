package main

import (
	"github.com/turbot/steampipe-plugin-aws/aws"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: aws.Plugin})
}
