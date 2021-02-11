package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const matrixKeyRegion = "region"

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildRegionList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	// retrieve regions from connection config
	regions := GetConfig(connection).Regions
	matrix := make([]map[string]interface{}, len(regions))
	for i, region := range regions {
		matrix[i] = map[string]interface{}{matrixKeyRegion: region}
	}
	return matrix
}
