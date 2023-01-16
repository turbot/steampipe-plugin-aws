package aws

// Calculating regions for Steampipe
//
// Working across multiple regions is a core capability of Steampipe. In the majority
// of cases it's really simple and intuitive to use, but there are a few edge cases
// that can be confusing. This section describes how regions are handled in Steampipe.
//
// AWS concepts for regions:
// - Regions: AWS regions are the geographic locations where AWS resources are hosted.
// - Partitions: AWS partitions are logical groupings of regions. There are three main
//   partitions: AWS Commercial, AWS GovCloud, and AWS China. Each partition has its own
//   set of regions.
// - Default region: The default region is the "primary" / most common region wtihin
//   the AWS partition. This region is used for API calls that must go to the base
//   endpoint. In general, it's better to use the client region (see below) if
//   possible, this should be the last resort.
// - Available Regions: The set of regions that are available in the partition.
//   This list includes optional regions, and does not change unless AWS announces
//   new regions.
// - Enabled Regions: The set of regions that are enabled in the account. This list
//   includes regions that have been opted-in (and excludes those that have
//   not), so it can change at any time.
//
// Steampipe config for regions:
// - `client_region`: Region that the Steampipe client (and likely AWS
//   CLI) is configured to use. This is the region that the client will use for
//   API calls that don't have a specific region specified.
// - `regions`: List of regions that the user has configured to use in
//   Steampipe.  Queries will combine results from these regions. But, the regions
//   config may include wildcard regions (e.g. `us-*`) that will be expanded to
//   include all enabled regions.
//
// Calculated for a connection at runtime:
// - query_regions: The set of regions that Steampipe will use for a given query.
//   This is the set of Enabled Regions that match the `regions` config.
// - query_client_region: The region that the Steampipe client will use for global
//   API calls (e.g. list S3 buckets).
// - query_default_region: The region the Steampipe client will use for primary
//   API calls on a default region (e.g. us-east-1).
//
// Bootstrapping and validating the region configuration is tricky because:
// * We calculate the partition from a region.
// * We can only determine enabled regions by running a query against a region.
// * The configuration may be invalid (a bad region name).
// * The configuration may be impossible (mix regions across partitions, e.g. us-east-1 and us-gov-east-1).
// * The configuration may be incomplete (no regions specified).
// * The configuration may have a client_region that is not in the regions list.
//
// First, we calculate the client_region:
// 1. `client_region` from the config.
// 2. Region returned by the AWS SDK (e.g. from the `AWS_REGION` environment variable).
// 3. First valid region in the regions list. This may require expansion of the
//    regions, and we don't know the correct partition yet.
// 4. Assume us-east-1, default region for the AWS Commercial partition
//
// Second, with a client_region, we can calculate:
// 1. The partition
// 2. Default region for the partition
// 3. Available regions for the partition
// 4. Enabled regions for the account (by querying the client_region). If this fails, default to the available regions.
//
// Third, to calculate the query_regions:
// 1. If the regions list is empty, default to the client_region.
// 2. Expand the query regions by matching it against the enabled regions.
// 3. If any specific regions in the config are not in the query regions then error, the config is invalid. This will catch partition mismatches and more.
// 4. If the client_region is not in the query regions then issue a warning, but continue.
//
// Fourth, calculate specific global regions:
// 1. The query_client_region is the client_region if it is in the query regions, otherwise the first region in the query regions.
// 2. The query_default_region is the default region for the partition if it is in the query regions, otherwise ????

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

	// Get information about the regions for this account, considering
	// the partition, opt-ins, etc.
	iRegionData, err := listRegionsCached(ctx, d, nil)
	if err != nil {
		panic(err)
	}
	regionData := iRegionData.(RegionsData)

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

// BuildWafRegionList returns the general region list, with a special region
// called "global" added. This is a specific region name used only by the WAF
// service.
// Note that the global region is always included in WAF results, even if the
// target region list is limited to specific regions. Currently, there is no
// way to exclude it except by filtering the results.
func BuildWafRegionList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	regionMatrix := BuildRegionList(ctx, d)
	matrix := make([]map[string]interface{}, 1, len(regionMatrix)+1)
	matrix[0] = map[string]interface{}{matrixKeyRegion: "global"}
	matrix = append(matrix, regionMatrix...)
	return matrix
}

// The list of regions is constant on a per-connection basis, so we cache it.
var listRegionsCached = plugin.HydrateFunc(listRegionsUncached).WithCache()

// Region data is complicated. Every account may have different combinations
// of regions based on opt-ins, partition, etc. This function is responsible
// for building a data structure representing region data for the account.
// The data is cached, and includes calls to APIs in it's construction when
// possible (e.g. region list API). But, it will also fall back to defaults
// like a full region list.
func listRegionsUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// The client region is used for two things:
	// 1. To make API calls to get the region list (if possible).
	// 2. To guess the partition we want to list regions from (if #1 fails).
	clientRegion, err := getClientRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("listRegions", "status", "starting", "connection_name", d.Connection.Name, "region", clientRegion)

	// If the client region is not AWS commercial (our default) then update
	// the full region list from a best guess based on the client region.
	allRegionsForClientPartition := awsCommercialRegions()
	if strings.HasPrefix(clientRegion, "us-gov") {
		allRegionsForClientPartition = awsUsGovRegions()
	} else if strings.HasPrefix(clientRegion, "cn") {
		allRegionsForClientPartition = awsChinaRegions()
	} else if strings.HasPrefix(clientRegion, "us-isob") {
		allRegionsForClientPartition = awsUsIsobRegions()
	} else if strings.HasPrefix(clientRegion, "us-iso") {
		allRegionsForClientPartition = awsUsIsoRegions()
	}

	// We try to get the accurate region list via an API call below, but as a
	// safe fallback assume all regions for the client partition

	// Default region data to use if everything else fails
	data := RegionsData{
		APIRetrivedList: false,
		AllRegions:      allRegionsForClientPartition,
		ActiveRegions:   allRegionsForClientPartition,
	}

	// We can query EC2 for the list of supported regions. If credentials
	// are insufficient this query will retry many times, so we create
	// a special client with a small number of retries to prevent hangs.
	svc, err := EC2RegionsClient(ctx, d, clientRegion)
	if err != nil {
		// handle in case user doesn't have access to ec2 service
		// save to extension cache
		return data, nil
	}

	params := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	}

	// execute list call
	resp, err := svc.DescribeRegions(ctx, params)
	if err != nil {
		// handle in case user doesn't have access to ec2 service
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

	plugin.Logger(ctx).Trace("listRegions", "status", "finished", "connection_name", d.Connection.Name, "region", clientRegion, "data", data)

	// save to extension cache
	return data, nil
}

// The default region is the "primary" / most common region wtihin the AWS partition.
// This region is used for API calls that must go to the base endpoint. In general,
// it's better to use the client region (see getClientRegion) if possible, this
// should be the last resort.
func getDefaultRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (string, error) {

	region, err := getClientRegion(ctx, d, h)
	if err != nil {
		return "", err
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

// Get the client region for AWS API calls that need to go to a central /
// non-regional endpoint (e.g. describe EC2 regions). Typically this should be
// the region closest to the user. Until we have a specific config option for
// it, we treat the first region in the regions list as the client region. If
// not given in the regions config, then try to read from the AWS config files.
// As a last resort, fall back to the default region for the partition.
func getClientRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (string, error) {
	regionInterface, err := getClientRegionCached(ctx, d, h)
	if err != nil {
		return "", err
	}
	region := regionInterface.(string)
	return region, nil
}

// The client region is cached on a per-connection basis to prevent re-lookup and
// recalculations from the configuration.
var getClientRegionCached = plugin.HydrateFunc(getClientRegionUncached).WithCache()

func getClientRegionUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	awsConfig := GetConfig(d.Connection)

	//var regions []string
	var region string

	plugin.Logger(ctx).Trace("getClientRegionUncached", "connection_name", d.Connection.Name, "awsConfig.Regions", awsConfig.Regions)

	if awsConfig.ClientRegion != nil {
		// The user has defined a specific home_region in their config. We use it
		// without further review. For example, they can have a home_region that is
		// not in the regions list.
		region = *awsConfig.ClientRegion
		plugin.Logger(ctx).Trace("getClientRegionUncached", "connection_name", d.Connection.Name, "region", region, "source", "home_region in config file")
		return region, nil
	}

	// Get the region from the AWS SDK. This will use the region defined in the
	// AWS config files, or the AWS_REGION environment variable, or the default
	// region for the partition.
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if session != nil && session.Config != nil && err == nil {
		// We have a home region from the SDK level
		region = *session.Config.Region
		plugin.Logger(ctx).Trace("getClientRegionUncached", "connection_name", d.Connection.Name, "region", region, "source", "AWS SDK lookup")
		return region, nil
	}

	/*

		TODO - fall back through other client region choices

		regions = awsConfig.Regions

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

	*/

	plugin.Logger(ctx).Trace("getClientRegionUncached", "connection_name", d.Connection.Name, "region", region)

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
		"ap-south-2",
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
