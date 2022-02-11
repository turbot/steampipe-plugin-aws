package aws

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type awsConfig struct {
	Regions                  []string `cty:"regions"`
	Profile                  *string  `cty:"profile"`
	AccessKey                *string  `cty:"access_key"`
	SecretKey                *string  `cty:"secret_key"`
	SessionToken             *string  `cty:"session_token"`
	IgnoreAccessDeniedErrors *string  `cty:"ignore_access_denied_errors"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"regions": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"profile": {
		Type: schema.TypeString,
	},
	"access_key": {
		Type: schema.TypeString,
	},
	"secret_key": {
		Type: schema.TypeString,
	},
	"session_token": {
		Type: schema.TypeString,
	},
	"ignore_access_denied_errors": {
		Type: schema.TypeBool,
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
