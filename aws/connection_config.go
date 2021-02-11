package aws

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type awsConfig struct {
	Regions []string `cty:"regions"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"regions": {
		Type:     schema.TypeList,
		Elem:     &schema.Attribute{Type: schema.TypeString},
		Optional: true,
	},
}

func ConfigInstance() interface{} {
	return &awsConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) awsConfig {
	if connection == nil || connection.Config == nil {
		return awsConfig{}
	}
	config, _ := connection.Config.(awsConfig)
	return config
}
