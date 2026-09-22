package aws

import (
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
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
	IgnoreErrorMessages   []string `hcl:"ignore_error_messages,optional"`
	IgnoreErrorCodes      []string `hcl:"ignore_error_codes,optional"`
	EndpointUrl           *string  `hcl:"endpoint_url"`
	S3ForcePathStyle      *bool    `hcl:"s3_force_path_style"`
	// S3UseDefaultRegionForBucketList opts out of always signing/routing the
	// aws_s3_bucket table's ListBuckets call through us-east-1 (the standard
	// behavior, kept for accurate `creation_date` -- see
	// https://www.marksayson.com/blog/s3-bucket-creation-dates-s3-master-regions/).
	// Set to true in networks that cannot reach us-east-1 (e.g. a single
	// region's VPC endpoint only) so ListBuckets is instead signed for the
	// connection's own default_region/regions. When enabled, `creation_date`
	// reflects the bucket's last-modified time rather than its true creation
	// time, matching ListBuckets' documented behavior for non-us-east-1 calls.
	S3UseDefaultRegionForBucketList *bool `hcl:"s3_use_default_region_for_bucket_list"`
}

func ConfigInstance() interface{} {
	return &awsConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) awsConfig {
	if connection == nil {
		return awsConfig{}
	}
	raw := connection.GetConfig()
	if raw == nil {
		return awsConfig{}
	}
	config, _ := raw.(awsConfig)

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
