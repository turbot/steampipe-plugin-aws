package aws

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type awsConfig struct {
	Regions      []string `cty:"regions"`
	Profile      *string  `cty:"profile"`
	AccessKey    *string  `cty:"access_key"`
	SecretKey    *string  `cty:"secret_key"`
	SessionToken *string  `cty:"session_token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"regions": {
		Type:        schema.TypeList,
		Elem:        &schema.Attribute{Type: schema.TypeString},
		Requirement: schema.AttributeOptional,
	},
	"profile": {
		Type:        schema.TypeString,
		Requirement: schema.AttributeOptional,
	},
	"access_key": {
		Type:        schema.TypeString,
		Requirement: schema.AttributeOptional,
	},
	"secret_key": {
		Type:        schema.TypeString,
		Requirement: schema.AttributeOptional,
	},
	"session_token": {
		Type:        schema.TypeString,
		Requirement: schema.AttributeOptional,
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
