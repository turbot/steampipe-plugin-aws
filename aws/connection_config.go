package aws

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type awsConfig struct {
	Regions []string `cty:"regions"`
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) awsConfig {
	if connection == nil || connection.Config == nil {
		return awsConfig{}
	}
	config, _ := connection.Config.(awsConfig)
	return config
}
