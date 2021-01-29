package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// MultiRegionGet :: helper function to execute list call for all regions
func MultiRegionList(listFunc plugin.HydrateFunc) plugin.HydrateFunc {
	// get regions from the connection config
	regions := getRegions()

	return func(ctx context.Context, queryData *plugin.QueryData, hydrateData *plugin.HydrateData) (interface{}, error) {
		// build a list of parameter maps - one per region
		paramsList := make([]map[string]string, len(regions))
		for i, region := range regions {
			paramsList[i] = map[string]string{"region": region}
		}
		return plugin.ListForEach(ctx, queryData, hydrateData, listFunc, paramsList)
	}
}

// MultiRegionGet :: helper function to execute get for all regions
func MultiRegionGet(getFunc plugin.HydrateFunc) plugin.HydrateFunc {
	// get regions from the connection config
	regions := getRegions()

	return func(ctx context.Context, queryData *plugin.QueryData, hydrateData *plugin.HydrateData) (interface{}, error) {
		// build a list of parameter maps - one per region
		paramsList := make([]map[string]string, len(regions))
		for i, region := range regions {
			paramsList[i] = map[string]string{"region": region}
		}
		return plugin.GetForEach(ctx, queryData, hydrateData, getFunc, paramsList)
	}
}

// dummy function to get regions - this wil be replace with connection config read
func getRegions() []string {
	// TODO read the actual conneciton config
	return []string{"eu-west1", "eu-west2"}
}
