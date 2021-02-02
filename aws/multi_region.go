package aws

var regions = []string{
	"us-east-1",
	"us-east-2",
	"eu-west-1",
	"eu-west-2",
}

const fetchMetdataKeyRegion = "region"

// BuildFetchMetadataList :: fds
func BuildFetchMetadataList() []map[string]interface{} {
	// build a list of fetchMetadata - one per region
	fetchMetadataList := make([]map[string]interface{}, len(regions))
	for i, region := range regions {
		fetchMetadataList[i] = map[string]interface{}{fetchMetdataKeyRegion: region}
	}
	return fetchMetadataList
}
