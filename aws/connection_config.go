package aws

import (
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type awsConfig struct {
	Regions               []string `hcl:"regions,optional"`
	DefaultRegion         *string  `hcl:"default_region"`
	Profile               *string  `hcl:"profile"`
	AccessKey             *string  `hcl:"access_key"`
	SecretKey             *string  `hcl:"secret_key"`
	SessionToken          *string  `hcl:"session_token"`
	MaxErrorRetryAttempts *int     `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int     `hcl:"min_error_retry_delay"`
	IgnoreErrorCodes      []string `hcl:"ignore_error_codes,optional"`
	EndpointUrl           *string  `hcl:"endpoint_url"`
	S3ForcePathStyle      *bool    `hcl:"s3_force_path_style"`
	MaxConnectionsPerHost *int     `hcl:"max_connections_per_host"`
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

	/*
		if config.MaxConnectionsPerHost == nil {
			defaultMaxConnectionsPerHost := 5
			config.MaxConnectionsPerHost = &defaultMaxConnectionsPerHost
		}
	*/

	return config
}

func NormalizeRegion(region string) string {
	// ensure regions are lower case, to work consistently in matching
	// and comparisons
	return strings.ToLower(region)
}
