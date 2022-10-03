package aws

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type awsConfig struct {
	Regions               []string `cty:"regions"`
	Profile               *string  `cty:"profile"`
	AccessKey             *string  `cty:"access_key"`
	SecretKey             *string  `cty:"secret_key"`
	SessionToken          *string  `cty:"session_token"`
	MaxErrorRetryAttempts *int     `cty:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int     `cty:"min_error_retry_delay"`
	IgnoreErrorCodes      []string `cty:"ignore_error_codes"`
	EndpointUrl           *string  `cty:"endpoint_url"`
	S3ForcePathStyle      *bool    `cty:"s3_force_path_style"`
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
	"ignore_error_codes": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"max_error_retry_attempts": {
		Type: schema.TypeInt,
	},
	"min_error_retry_delay": {
		Type: schema.TypeInt,
	},
	"endpoint_url": {
		Type: schema.TypeString,
	},
	"s3_force_path_style": {
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

	// Setting "regions = []" in the connection config is not valid
	if len(config.Regions) == 0 {
		errorMessage := fmt.Sprintf("\nconnection %s has invalid value for \"regions\", it must contain at least 1 region.", connection.Name)
		panic(errorMessage)
	}

	return config
}
