package aws

import (
	"context"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const matrixKeyRegion = "region"

var pluginQueryData *plugin.QueryData

func init() {
	pluginQueryData = &plugin.QueryData{
		ConnectionManager: connection.NewManager(),
	}
}

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildRegionList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {
	// retrieve regions from connection config
	awsConfig := GetConfig(connection)
	pluginQueryData.Connection = connection

	validRegions, _ := listRegions(ctx, pluginQueryData)
	var allRegions []string

	if awsConfig.Regions != nil {
		for _, pattern := range awsConfig.Regions {
			for _, validRegion := range validRegions {
				if ok, _ := path.Match(pattern, validRegion); ok {
					allRegions = append(allRegions, validRegion)
				}
			}
		}
	}

	if allRegions != nil {
		uniqueRegions := unique(allRegions)
		// regions := GetConfig(connection).Regions

		if len(getInvalidRegions(uniqueRegions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(uniqueRegions), ","))
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(uniqueRegions))
		for i, region := range uniqueRegions {
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

// BuildWafRegionList :: return a list of matrix items for AWS WAF resources, one per region specified in the connection config
func BuildWafRegionList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	// retrieve regions from connection config
	awsConfig := GetConfig(connection)

	if awsConfig.Regions != nil {
		regions := awsConfig.Regions

		if len(getInvalidRegions(regions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(regions), ","))
		}
		regions = append(regions, "global")

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

func listRegions(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	cacheKey := "listRegions"

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]string), nil
	}

	awsCommercialRegions := []string{
		"af-south-1", "ap-east-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "eu-central-1", "eu-north-1", "eu-south-1", "eu-west-1", "eu-west-2", "eu-west-3", "me-south-1", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2"}

	awsUsGovRegions := []string{"us-gov-east-1", "us-gov-west-1"}

	awsChinaGovRegions := []string{"cn-north-1", "cn-northwest-1"}

	defaultRegion := GetDefaultAwsRegion(d)

	defaultRegions := awsCommercialRegions

	if strings.HasPrefix(defaultRegion, "us-gov") {
		defaultRegions = awsUsGovRegions
	} else if strings.HasPrefix(defaultRegion, "cn") {
		defaultRegions = awsChinaGovRegions
	}

	// Create Session
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		// handle in case user doesn't have access to ec2 service
		// save to extension cache
		d.ConnectionManager.Cache.Set(cacheKey, defaultRegions)
		return defaultRegions, nil
	}

	params := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	}

	// execute list call
	resp, err := svc.DescribeRegions(params)
	if err != nil {
		// handle in case user doesn't have access to ec2 service
		d.ConnectionManager.Cache.Set(cacheKey, defaultRegions)
		return defaultRegions, nil
	}

	var activeRegions []string
	for _, region := range resp.Regions {
		if *region.OptInStatus != "not-opted-in" {
			activeRegions = append(activeRegions, *region.RegionName)
		}
	}

	// save to extension cache
	d.ConnectionManager.Cache.Set(cacheKey, activeRegions)
	return activeRegions, err
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
