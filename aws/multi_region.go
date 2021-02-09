package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const fetchMetdataKeyRegion = "region"

// BuildFetchMetadataList :: fds
func BuildFetchMetadataList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	// retrieve regions from connection config
	regions := GetConfig(connection).Regions
	// build a list of fetchMetadata - one per region
	fetchMetadataList := make([]map[string]interface{}, len(regions))
	for i, region := range regions {
		fetchMetadataList[i] = map[string]interface{}{fetchMetdataKeyRegion: region}
	}
	return fetchMetadataList
}
