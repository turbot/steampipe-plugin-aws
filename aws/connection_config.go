package aws

import (
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type awsConfig struct {
	Regions               []string `cty:"regions"`
	ClientRegion          *string  `cty:"client_region"`
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
	"client_region": {
		Type: schema.TypeString,
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

	if config.Regions != nil {
		if len(config.Regions) == 0 {
			// Setting "regions = []" in the connection config is not valid
			errorMessage := fmt.Sprintf("connection %s has invalid value for \"regions\", it must contain at least 1 region.", connection.Name)
			panic(errorMessage)
		}

		for i, r := range config.Regions {
			config.Regions[i] = NormalizeRegion(r)
		}
	}

	return config
}

func NormalizeRegion(region string) string {
	// ensure regions are lower case, to work consistently in matching
	// and comparisons
	return strings.ToLower(region)
}
