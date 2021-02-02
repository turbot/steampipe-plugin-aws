package aws

var regions = []string{
	"eu-north-1",
	"eu-south-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
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
