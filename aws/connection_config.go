package aws

import (
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/schema"
)

type awsConfig struct {
	Regions                  []string `cty:"regions"`
	Profile                  *string  `cty:"profile"`
	AccessKey                *string  `cty:"access_key"`
	SecretKey                *string  `cty:"secret_key"`
	SessionToken             *string  `cty:"session_token"`
	IgnoreAccessDeniedErrors *bool    `cty:"ignore_access_denied_errors"`
	AccessDeniedErrorsList   []string `cty:"access_denied_errors_list"` // To extend the list of access denied errors handling
	MaxErrorRetryAttempts    *int     `cty:"max_error_retry_attempts"`
	MinErrorRetryDelay       *int     `cty:"min_error_retry_delay"`
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
	"access_denied_errors_list": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"max_error_retry_attempts": {
		Type: schema.TypeInt,
	},
	"min_error_retry_delay": {
		Type: schema.TypeInt,
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
