package aws

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const matrixKeyRegion = "region"

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildRegionList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	// retrieve regions from connection config
	awsConfig := GetConfig(connection)

	if &awsConfig != nil && awsConfig.Regions != nil {
		regions := GetConfig(connection).Regions

		if len(getInvalidRegions(regions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(regions), ","))
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(regions))
		for i, region := range regions {
			matrix[i] = map[string]interface{}{matrixKeyRegion: region}
		}
		return matrix
	}

	return []map[string]interface{}{
		{matrixKeyRegion: GetDefaultRegion()},
	}
}

func getInvalidRegions(regions []string) []string {
	awsRegions := []string{
		"af-south-1", "ap-east-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "eu-central-1", "eu-north-1", "eu-south-1", "eu-west-1", "eu-west-2", "eu-west-3", "me-south-1", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2", "us-gov-east-1", "us-gov-west-1", "cn-north-1", "cn-northwest-1"}

	invalidRegions := []string{}
	for _, region := range regions {
		if !helpers.StringSliceContains(awsRegions, region) {
			invalidRegions = append(invalidRegions, region)
		}
	}
	return invalidRegions
}
