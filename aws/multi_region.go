package aws

import (
	"context"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

const matrixKeyRegion = "region"

type RegionsData struct {
	AllRegions      []string
	ActiveRegions   []string
	NotOptedRegions []string
	APIRetrivedList bool
}

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config.
// Plugin supports wildcards "*" and "?" in the connection config for the regions.
//
// This function will build the regions list dynamically based the activated region in the AWS account.
// For this it uses https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeRegions.html
//
// Some scenarios
// When no regions mentioned in connection config
// region = "us-east-1"
//
// When no regions mentioned in connection config and started steampipe as AWS_REGION=ap-south-1 steampipe query
// region = "ap-south-1"
//
// regions = ["*"]
// regions="af-south-1, eu-north-1, ap-south-1, eu-west-3, eu-west-2, eu-south-1, eu-west-1, ap-northeast-3, ap-northeast-2, me-south-1, ap-northeast-1, sa-east-1, ca-central-1, ap-southeast-1, ap-southeast-2, eu-central-1, us-east-1, us-east-2, us-west-1, us-west-2"
//
// regions = ["me-*", "ap-*", "us-*"]
// regions="me-south-1, ap-south-1, ap-northeast-3, ap-northeast-2, ap-northeast-1, ap-southeast-1, ap-southeast-2, us-east-1, us-east-2, us-west-1, us-west-2"
func BuildRegionList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// Trace logging to debug cache and execution flows
	logRegion := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("BuildRegionList", "status", "starting", "connection_name", d.Connection.Name, "region", logRegion)

	// Cache region list matrix
	cacheKey := "RegionListMatrix"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// Retrieve regions list from the AWS plugin steampipe connection config
	awsConfig := GetConfig(d.Connection)

	// If the regions are not set in the aws.spc, this should try to pick the
	// AWS region based on the SDK logic similar to the AWS CLI command
	// "aws configure list"
	if awsConfig.Regions == nil {
		region := ""
		session, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			panic(err)
		}
		if session != nil && session.Config != nil {
			region = *session.Config.Region
		}

		/*
			AWS CLI behaviour if region is not found through "aws configure list"

			plugins/steampipe-plugin-aws|main⚡ ⇒  aws configure list
			Name                    Value             Type    Location
			----                    -----             ----    --------
			profile                <not set>             None    None
			access_key     ###############ABCD shared-credentials-file
			secret_key     ###############/g/w shared-credentials-file
			region                <not set>             None    None

			Global services work without any issue, if region is not set

			plugins/steampipe-plugin-aws|main⚡ ⇒  aws s3 ls
			2022-02-15 12:24:02 aws-cloudtrail-111122223333
			2021-08-16 22:06:46 aws-cloudtrail-111122223333
			.....

			Regional services fail with an warning to set the region

			plugins/steampipe-plugin-aws|main⚡ ⇒  aws sns list-topics
			You must specify a region. You can also configure your region by running "aws configure".
		*/
		if region == "" {
			panic("you must specify a region in \"regions\" in ~/.steampipe/config/aws.spc. Edit your connection configuration file and then restart Steampipe.")
		}

		matrix := []map[string]interface{}{
			// TODO
			// If the region is a invalid region. It will lead to long retry with error message like
			// dial tcp: lookup <service>.<invalid-region>.amazonaws.com: no such host (SQLSTATE HV000)
			{matrixKeyRegion: region},
		}

		// set cache
		d.ConnectionManager.Cache.Set(cacheKey, matrix)
		return matrix
	}

	// If regions are mentioned in the connection config
	// Get the list of AWS region from EC2 DescribeRegions API
	regionData, err := listRegions(ctx, d, nil)
	if err != nil {
		panic(err)
	}
	var finalRegions []string

	// If the list was retrived from AWS API - just looks for Active Regions
	if regionData.APIRetrivedList {
		for _, pattern := range awsConfig.Regions {
			for _, validRegion := range regionData.ActiveRegions {
				if ok, _ := path.Match(pattern, validRegion); ok {
					finalRegions = append(finalRegions, validRegion)
				}
			}
			finalRegions = helpers.StringSliceDistinct(finalRegions)
		}
	} else {
		// If the list was not retrived from AWS API
		// match for regions from all regions list
		for _, pattern := range awsConfig.Regions {
			for _, validRegion := range regionData.AllRegions {
				if ok, _ := path.Match(pattern, validRegion); ok {
					finalRegions = append(finalRegions, validRegion)
				}
			}
		}
		finalRegions = helpers.StringSliceDistinct(finalRegions)
	}

	matrix := make([]map[string]interface{}, len(finalRegions))
	for i, region := range finalRegions {
		matrix[i] = map[string]interface{}{matrixKeyRegion: region}
	}

	plugin.Logger(ctx).Debug("BuildRegionList", "status", "result", "connection_name", d.Connection.Name, "region", logRegion, "matrix", matrix)

	// set cache
	d.ConnectionManager.Cache.Set(cacheKey, matrix)
	return matrix
}

// BuildWafRegionList :: return a list of matrix items for AWS WAF resources, one per region specified in the connection config
func BuildWafRegionList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	var regionMatrix []map[string]interface{}
	if cachedData, ok := d.ConnectionManager.Cache.Get("RegionListMatrix"); ok {
		regionMatrix = cachedData.([]map[string]interface{})
	} else {
		regionMatrix = BuildRegionList(ctx, d)
	}

	matrix := make([]map[string]interface{}, 1, len(regionMatrix)+1)
	matrix[0] = map[string]interface{}{matrixKeyRegion: "global"}
	matrix = append(matrix, regionMatrix...)

	return matrix
}

func listRegions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (RegionsData, error) {
	cacheKey := "listRegions"

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(RegionsData), nil
	}

	// Default region data to use if everything else fails
	data := RegionsData{
		APIRetrivedList: false,
	}

	// The preferred region is used for two things:
	// 1. To make API calls to get the region list (if possible).
	// 2. To guess the partition we want to list regions from (if #1 fails).
	preferredRegion, err := getPreferredRegion(ctx, d, h)
	if err != nil {
		return data, err
	}

	// If the preferred region is not AWS commercial (our default) then update
	// the full region list from a best guess based on the preferred region.
	allRegionsForPreferredPartition := awsCommercialRegions()
	if strings.HasPrefix(preferredRegion, "us-gov") {
		allRegionsForPreferredPartition = awsUsGovRegions()
	} else if strings.HasPrefix(preferredRegion, "cn") {
		allRegionsForPreferredPartition = awsChinaRegions()
	} else if strings.HasPrefix(preferredRegion, "us-isob") {
		allRegionsForPreferredPartition = awsUsIsobRegions()
	} else if strings.HasPrefix(preferredRegion, "us-iso") {
		allRegionsForPreferredPartition = awsUsIsoRegions()
	}

	// We try to get the accurate region list via an API call below, but as a
	// safe fallback assume all regions for the preferred partition
	data.AllRegions = allRegionsForPreferredPartition
	data.ActiveRegions = allRegionsForPreferredPartition

	// We can query EC2 for the list of supported regions. If credentials
	// are insufficient this query will retry many times, so we create
	// a special client with a small number of retries to prevent hangs.
	svc, err := EC2RegionsClient(ctx, d, preferredRegion)
	if err != nil {
		// handle in case user doesn't have access to ec2 service
		// save to extension cache
		d.ConnectionManager.Cache.Set(cacheKey, data)
		return data, nil
	}

	params := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	}

	// execute list call
	resp, err := svc.DescribeRegions(ctx, params)
	if err != nil {
		// handle in case user doesn't have access to ec2 service
		d.ConnectionManager.Cache.Set(cacheKey, data)
		return data, nil
	}

	var activeRegions []string   // All enabled regions in the account.
	var notOptedRegions []string // All not enabled regions in the account.
	var allRegions []string      // All regions listed by the API DescribeRegions.

	for _, region := range resp.Regions {
		allRegions = append(allRegions, *region.RegionName)

		if *region.OptInStatus != "not-opted-in" {
			activeRegions = append(activeRegions, *region.RegionName)
		} else {
			notOptedRegions = append(notOptedRegions, *region.RegionName)
		}
	}

	data = RegionsData{
		AllRegions:      allRegions,
		ActiveRegions:   activeRegions,
		NotOptedRegions: notOptedRegions,
		APIRetrivedList: true,
	}

	// save to extension cache
	d.ConnectionManager.Cache.Set(cacheKey, data)
	return data, err
}

// Get the preferred region for AWS API calls that need to go to a central /
// non-regional endpoint (e.g. list S3 buckets). Typically this should be the
// region closest to the user. We don't currently have a way to specify this
// region, so we assume the default region is closest. This is not always true
// (i.e. not everyone lives near us-east-1), but it's the best we can do.
// OPTIONS / TODO:
//  1. Add a way to specify the preferred region in the connection config.
//  2. Assume the first region in the regions config list. We currently use
//     this model to guess the best partition.
func getPreferredRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (string, error) {
	return getDefaultRegion(ctx, d, h)
}

// The default region is the "primary" / most common region wtihin the AWS partition.
// This region is used for API calls that must go to the base endpoint. In general,
// it's better to use the preferred region (see getPreferredRegion) if possible, this
// should be the last resort.
func getDefaultRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (string, error) {
	regionInterface, err := getDefaultRegionCached(ctx, d, h)
	if err != nil {
		return "", err
	}
	region := regionInterface.(string)
	return region, nil
}

var getDefaultRegionCached = plugin.HydrateFunc(getDefaultRegionUncached).WithCache()

func getDefaultRegionUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	awsConfig := GetConfig(d.Connection)

	var regions []string
	var region string

	if awsConfig.Regions != nil {
		regions = awsConfig.Regions
		// Pick the first region from the regions list as a best guess to determine
		// the default region for the AWS partition based on the region prefix.
		region = regions[0]
	} else {
		// If the regions are not set in the aws.spc, this will try to pick the
		// AWS region based on the SDK logic similar to the AWS CLI command
		// "aws configure list"
		session, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			panic(err)
		}
		if session != nil && session.Config != nil {
			region = *session.Config.Region
		}
	}

	// Most of the global services like IAM, S3, Route 53, target these regions
	if strings.HasPrefix(region, "us-gov") {
		region = "us-gov-west-1"
	} else if strings.HasPrefix(region, "cn") {
		region = "cn-northwest-1"
	} else if strings.HasPrefix(region, "us-isob") {
		region = "us-isob-east-1"
	} else if strings.HasPrefix(region, "us-iso") {
		region = "us-iso-east-1"
	} else {
		region = "us-east-1"
	}

	return region, nil
}

//
// AWS STANDARD REGIONS
//
// Source: https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints
//
// Maintain a hard coded list of regions to use when API calls to the region
// list endpoint are not possible. This list must be updated manually as new
// regions are announced.
//

func awsCommercialRegions() []string {
	return []string{
		"af-south-1",
		"ap-east-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ca-central-1",
		"eu-central-1",
		"eu-central-2",
		"eu-north-1",
		"eu-south-1",
		"eu-south-2",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"me-central-1",
		"me-south-1",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}
}

func awsUsGovRegions() []string {
	return []string{
		"us-gov-east-1",
		"us-gov-west-1",
	}
}

func awsChinaRegions() []string {
	return []string{
		"cn-north-1",
		"cn-northwest-1",
	}
}

func awsUsIsoRegions() []string {
	return []string{
		"us-iso-east-1",
		"us-iso-west-1",
	}
}

func awsUsIsobRegions() []string {
	return []string{
		"us-isob-east-1",
	}
}
