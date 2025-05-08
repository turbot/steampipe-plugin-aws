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
// - Last resort region: The default region is the "primary" / most common region within
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
// - `default_region`: Region that the Steampipe client (and likely AWS
//   CLI) is configured to use. This is the region that the client will use for
//   API calls that don't have a specific region specified.
// - `regions`: List of regions that the user has configured to use in
//   Steampipe.  Queries will combine results from these regions. But, the regions
//   config may include wildcard regions (e.g. `us-*`) that will be expanded to
//   include all enabled regions.
//
// Calculated for a connection at runtime:
// - Query regions: The set of regions that Steampipe will use for a given query.
//   This is the set of Enabled Regions that match the `regions` config.
//
// Query regions is based on these factors:
// - All regions in the partition
// - All regions available for service X in the partition
// - Hard-coded exclusions (sometimes the definition of service X is wrong)
// - Filter by regions enabled for this account (e.g. some might not be opted-in)
// - Filter by configured query `regions` in aws.spc
//
// How to guess the partition (and thus full region set):
// The basis for all decisions in the partition. We have two possible ways to determine
// the partition:
// 1. A guess from the default_region. (fast, no API call required)
// 2. GetCallerIdentity from the common columns. (cached, more accurate)
//
// Notes about region config & implementation:
// - It's complicated with multiple layers of lookups & filters. It's not just you.
// - Regions are not validated (other than to lower case).
// - We always try to make things work, even if the config is non-existent.
// - Try not to guess a default region, if guessing go to the last resort region
//   rather than using the `regions` list. This is not awesome, but consistency is
//   better than getting something "random" (e.g. alphabetical).
// - There is a tension between accurate data (e.g. via an API lookup) vs fast data
//   that can be obtained without needing permissions. For example, without being
//   given a region we have no idea what partition the account is in or what regions
//   are enabled. Much of the code is trying to balance this in a consistent way.

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	cloudwatchv1 "github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/logging"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const matrixKeyRegion = "region"

type RegionsData struct {
	AllRegions      []string
	ActiveRegions   []string
	NotOptedRegions []string
	APIRetrivedList bool
}

// Return a matrix of all regions for tables that target every region.
// It's normally better to use SupportedRegionMatrix instead, as it will
// filter out regions that are not enabled for the specific service.
func AllRegionsMatrix(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	return SupportedRegionMatrixWithExclusions("", []string{})(ctx, d)
}

// _metric_ tables must all be limited to the CloudWatch service regions.
// This is a convenience function for them to use.
func CloudWatchRegionsMatrix(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	return SupportedRegionMatrixWithExclusions(cloudwatchv1.EndpointsID, []string{})(ctx, d)
}

func S3TablesRegionsMatrix(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	commonColumnData, err := getCommonColumns(ctx, d, nil)
	if err != nil {
		plugin.Logger(ctx).Error("S3TablesRegionsMatrix", "connection_name", d.Connection.Name, "unable to get partition name", err)
		panic(err)
	}
	partitionName := commonColumnData.(*awsCommonColumnData).Partition

	// Get AWS partition based on the partition name
	// Get supported service along with the endpoints for the partition
	partition, err := getPartitionValueByPartitionName(partitionName)
	if err != nil {
		panic(fmt.Errorf("S3TablesRegionsMatrix: failed to get the endpoint details for the partition '%s', %v", partitionName, err))
	}

	s3SupportedRegions := partition.Services[AWS_S3_SERVICE_ID].Endpoints
	var unsupportedRegionsForS3Tables []string
	for region, ed := range s3SupportedRegions {
		if !slices.Contains(ed.SignatureVersions, "s3v4") {
			unsupportedRegionsForS3Tables = append(unsupportedRegionsForS3Tables, region)
		}
	}

	return SupportedRegionMatrixWithExclusions(AWS_S3_SERVICE_ID, unsupportedRegionsForS3Tables)(ctx, d)
}

// Return a matrix of regions supported by serviceID, which will then be
// queried for the table in parallel. This result factors in things like
// regions that are opted-in, regions for the service and even the `regions`
// config in aws.spc.
func SupportedRegionMatrix(serviceID string) func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	return SupportedRegionMatrixWithExclusions(serviceID, []string{})
}

// This function is used in the GetMatrixItemFunc implementations within the table definitions.
// GetMatrixItemFunc is designed to accept a single return type `[]map[string]interface{}`.
// AWS regional tables make API calls based on the region matrix return by the SupportedRegionMatrixWithExclusions function.
// In cases of incorrect credential configurations, listQueryRegionsForConnection returns an error, such as: "Error: operation error STS: GetCallerIdentity, failed to sign request: failed to retrieve credentials: failed to refresh cached credentials, operation error STS: AssumeRole, https response error StatusCode: 403, RequestID: a1028f7b-cb77-4b9e-b1e5-ce96ea77150e, api error InvalidClientTokenId: The security token included in the request is invalid."
// When an error is encountered, it should trigger a panic with that error; otherwise, regional tables return an empty row.
// The reason regional tables return an empty row because the function(SupportedRegionMatrixWithExclusions) returns an empty `[]map[string]interface{}` upon encountering any error.

// Similar to SupportedRegionMatrix, but excludes the regions in excludeRegions
// for manual overrides if the service definition is incorrect.
func SupportedRegionMatrixWithExclusions(serviceID string, excludeRegions []string) func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	return func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
		logging.LogTime("SupportedRegionMatrixWithExlusions start")
		defer logging.LogTime("SupportedRegionMatrixWithExlusions end")
		// Default to an empty list of regions
		matrix := []map[string]interface{}{}
		// Get the regions enabled for this account
		queryRegions, err := listQueryRegionsForConnection(ctx, d)
		if err != nil {
			plugin.Logger(ctx).Error("SupportedRegionMatrixWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "query_regions_error", err)
			panic(err)
		}
		plugin.Logger(ctx).Debug("SupportedRegionMatrixWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "query_regions", queryRegions)
		// Get the possible regions for this service
		var serviceRegions []string
		if serviceID == "" {
			// No service given, assume all regions are in scope
			serviceRegions = queryRegions
			plugin.Logger(ctx).Debug("SupportedRegionMatrixWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "service_regions_using_query_regions", serviceRegions)
		} else {
			serviceRegions, err = listRegionsForServiceWithExclusions(ctx, d, serviceID, excludeRegions)
			if err != nil {
				plugin.Logger(ctx).Error("SupportedRegionMatrixWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "service_regions_error", err)
				panic(err)
			}
			plugin.Logger(ctx).Debug("SupportedRegionMatrixWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "service_regions", serviceRegions)
		}
		// Find all regions in both the query regions and the service regions
		for _, region := range queryRegions {
			if slices.Contains(serviceRegions, region) {
				obj := map[string]interface{}{matrixKeyRegion: region}
				matrix = append(matrix, obj)
			}
		}
		plugin.Logger(ctx).Debug("SupportedRegionMatrixWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "matrix", matrix)
		return matrix
	}
}

// Calculate the regions that the user has requested to query for this
// connection.  Basically, we generate a possible list of regions (enabled
// regions for the account, or all regions for the partition) and then filter it
// by the `regions` config setting in aws.spc.
//
// The list is unique for each region specified in the connection config.
// Plugin supports wildcards "*" and "?" in the connection config for the
// regions.
//
// This function will build the regions list dynamically based the activated
// region in the AWS account. For this it uses
// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeRegions.html
//
// Some scenarios
// When no regions mentioned in connection config
//
//	result = "us-east-1"
//
// When no regions mentioned in connection config and started steampipe as AWS_REGION=ap-south-1 steampipe query
//
//	result = "ap-south-1"
//
// Wildcard:
//
//	regions = ["*"]
//	regions = [af-south-1, eu-north-1, ap-south-1, eu-west-3, eu-west-2, eu-south-1, eu-west-1, ap-northeast-3, ap-northeast-2, me-south-1, ap-northeast-1, sa-east-1, ca-central-1, ap-southeast-1, ap-southeast-2, eu-central-1, us-east-1, us-east-2, us-west-1, us-west-2]
//
// Partial wildcard:
//
//	regions = ["me-*", "ap-*", "us-*"]
//	result = [me-south-1, ap-south-1, ap-northeast-3, ap-northeast-2, ap-northeast-1, ap-southeast-1, ap-southeast-2, us-east-1, us-east-2, us-west-1, us-west-2]
//
// Mismatch with default region will return zero results:
//
//	default_region = "us-east-1"
//	regions = ["us-gov-*"]
//	result = []
func listQueryRegionsForConnection(ctx context.Context, d *plugin.QueryData) ([]string, error) {

	// Retrieve regions list from the AWS plugin steampipe connection config
	awsSpcConfig := GetConfig(d.Connection)

	// If there is no regions defined in SPC, then we default to targeting
	// the default region only for queries.
	if awsSpcConfig.Regions == nil {
		region, err := getDefaultRegion(ctx, d, nil)
		if err != nil || region == "" {
			plugin.Logger(ctx).Error("listQueryRegionsForConnection", "connection_name", d.Connection.Name, "default_region_error", err)
			return nil, fmt.Errorf("regions or default_region must be defined")
		}
		// Return the default region as the only query region for this connection
		plugin.Logger(ctx).Debug("listQueryRegionsForConnection", "connection_name", d.Connection.Name, "using default region", region)
		return []string{region}, nil
	}

	// PRE: there is a list of regions in the config to match against

	// Get information about the regions for this account, considering
	// the partition, opt-ins, etc.
	iRegionData, err := listRegionsCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	regionData := iRegionData.(RegionsData)

	// If we have a list of regions from the API, then it's the most accurate.
	// Otherwise just assume all regions for the account.
	maxTargetRegions := regionData.AllRegions
	if regionData.APIRetrivedList {
		maxTargetRegions = regionData.ActiveRegions
	} else {
		plugin.Logger(ctx).Warn("listQueryRegionsForConnection", "connection_name", d.Connection.Name, "status", "target regions not available via EC2.DescribeRegions API, assuming all regions for partition are active", "targetRegions", maxTargetRegions)
	}

	// Filter to regions that match the patterns in the config.
	var targetRegions []string
	for _, pattern := range awsSpcConfig.Regions {
		for _, validRegion := range maxTargetRegions {
			if ok, _ := path.Match(pattern, validRegion); ok {
				targetRegions = append(targetRegions, validRegion)
			}
		}
	}
	targetRegions = helpers.StringSliceDistinct(targetRegions)

	plugin.Logger(ctx).Debug("listQueryRegionsForConnection", "connection_name", d.Connection.Name, "targetRegions", targetRegions)

	return targetRegions, nil
}

// WAFRegionMatrix returns the general region list, with a special region
// called "global" added. This is a specific region name used only by the WAF
// service.
// Note that the global region is always included in WAF results, even if the
// target region list is limited to specific regions. Currently, there is no
// way to exclude it except by filtering the results.
func WAFRegionMatrix(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	regionMatrix := CloudWatchRegionsMatrix(ctx, d)
	matrix := make([]map[string]interface{}, 1, len(regionMatrix)+1)
	matrix[0] = map[string]interface{}{matrixKeyRegion: "global"}
	matrix = append(matrix, regionMatrix...)
	return matrix
}

// List all regions for a given service in the partition for this connection.
func listRegionsForService(ctx context.Context, d *plugin.QueryData, serviceID string) ([]string, error) {
	return listRegionsForServiceWithExclusions(ctx, d, serviceID, []string{})
}

// List all regions for a given service, defined to work with Memoize().
// Call listRegionsForService() instead of using this directly.
var listRegionsForServiceCached = plugin.HydrateFunc(listRegionsForServiceUncached).Memoize(memoize.WithCacheKeyFunction(listRegionsForServiceCacheKey))

// List all regions for a given service in the partition for this connection, but
// manually exclude any regions in excludeRegions.
func listRegionsForServiceWithExclusions(ctx context.Context, d *plugin.QueryData, serviceID string, excludeRegions []string) ([]string, error) {
	h := &plugin.HydrateData{Item: serviceID}
	// Get all regions for the service
	iRegions, err := listRegionsForServiceCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listRegionsForServiceWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "error", err)
		return nil, err
	}
	serviceRegions := iRegions.([]string)
	// Remove the excluded regions from the valid list
	serviceRegions = helpers.RemoveFromStringSlice(serviceRegions, excludeRegions...)
	plugin.Logger(ctx).Debug("listRegionsForServiceWithExclusions", "connection_name", d.Connection.Name, "serviceID", serviceID, "excludeRegions", excludeRegions, "serviceRegions", serviceRegions)
	return serviceRegions, nil
}

// getClient is per-region, but Memoize() is per-connection, so a setup
// a custom cache key with region information in it.
func listRegionsForServiceCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Extract the region from the hydrate data. This is not per-row data,
	// but a clever pass through of context for our case.
	serviceID := h.Item.(string)
	key := fmt.Sprintf("listRegionsForService-%s", serviceID)
	return key, nil
}

// Use the AWS SDK to get a list of regions that the given service (in hydrate
// data) supports.
// Implementation notes:
//   - Use AWS SDK v1 because v2 does not expose this data (ugh).
//   - Use getCommonColumns to get the accurate partition for the account (via
//     GetCallerIdentity). This is more accurate than guessing from the default
//     region.
func listRegionsForServiceUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var partitionName string
	var partition *Partition

	// Service ID is passed through the hydrate data
	serviceID := h.Item.(string)

	// We use getCommonColumns to get the partition name, which is more accurate
	// than guessing from the default region. It does include an API call to
	// GetCallerIdentity under the hood, but that is cached and used for almost
	// all tables / query results anyway.
	commonColumnData, err := getCommonColumns(ctx, d, nil)
	if err != nil {
		plugin.Logger(ctx).Error("listRegionsForServiceUncached", "connection_name", d.Connection.Name, "unable to get partition name", err)
		return nil, err
	}
	partitionName = commonColumnData.(*awsCommonColumnData).Partition

	// Get AWS partition based on the partition name
	// Get supported service along with the endpoints for the partition
	partition, err = getPartitionValueByPartitionName(partitionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get the endpoint details for the partition '%s', %v", partitionName, err)
	}

	var regionsForService []string

	// https://raw.githubusercontent.com/aws/aws-sdk-go-v2/master/codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoints.json
	services := partition.Services
	serviceInfo, ok := services[serviceID]
	if !ok {
		err := fmt.Errorf("listRegionsForServiceUncached called with invalid service ID: %s", serviceID)
		plugin.Logger(ctx).Error("listRegionsForServiceUncached", "connection_name", d.Connection.Name, "partition", partition, "serviceID", serviceID, "error", err)
		return nil, err
	}

	for rs := range serviceInfo.Endpoints {
		re := regexp.MustCompile(partition.RegionRegex)
		if re.Match([]byte(rs)) {
			regionsForService = append(regionsForService, rs)
		}
	}

	plugin.Logger(ctx).Debug("listRegionsForServiceUncached", "connection_name", d.Connection.Name, "partition", partition, "serviceID", serviceID, "regionsForService", regionsForService)
	return regionsForService, nil
}

// The list of regions is constant on a per-connection basis, so we cache it.
var listRegionsCached = plugin.HydrateFunc(listRegionsUncached).Memoize()

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
	clientRegion, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Debug("listRegionsUncached", "status", "starting", "connection_name", d.Connection.Name, "region", clientRegion)

	// If the client region is not AWS commercial (our default) then update
	// the full region list from a best guess based on the client region.
	allRegionsForClientPartition := getRegionByPartition("aws")
	if strings.HasPrefix(clientRegion, "us-gov") {
		allRegionsForClientPartition = getRegionByPartition("aws-us-gov")
	} else if strings.HasPrefix(clientRegion, "cn") {
		allRegionsForClientPartition = getRegionByPartition("aws-cn")
	} else if strings.HasPrefix(clientRegion, "us-isob") {
		allRegionsForClientPartition = getRegionByPartition("aws-iso-b")
	} else if strings.HasPrefix(clientRegion, "us-iso") {
		allRegionsForClientPartition = getRegionByPartition("aws-iso")
	}

	// We try to get the accurate region list via an API call below, but as a
	// safe fallback assume all regions for the client partition

	// Default region data to use if everything else fails
	data := RegionsData{
		APIRetrivedList: false,
		AllRegions:      allRegionsForClientPartition,
		ActiveRegions:   allRegionsForClientPartition,
	}

	// Get the AWS region list from the EC2 API (via cache)
	iRegions, err := listRawAwsRegions(ctx, d, h)
	if err != nil {
		// handle in case user doesn't have access to ec2 service
		// save to extension cache
		plugin.Logger(ctx).Warn("listRegionsUncached", "connection_name", d.Connection.Name, "region", clientRegion, "regions_error", err)
		return data, nil
	}

	var activeRegions []string   // All enabled regions in the account.
	var notOptedRegions []string // All not enabled regions in the account.
	var allRegions []string      // All regions listed by the API DescribeRegions.

	for _, region := range iRegions.([]types.Region) {
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

	plugin.Logger(ctx).Debug("listRegionsUncached", "status", "finished", "connection_name", d.Connection.Name, "region", clientRegion, "data", data)

	// save to extension cache
	return data, nil
}

// Cached list of regions actually active in the AWS account.
var listRawAwsRegions = plugin.HydrateFunc(listRawAwsRegionsUncached).Memoize()

// List regions for this AWS account connection by calling the EC2
// DescribeRegions API. The call is done using the default region (hopefully
// close to the user) and then cached. This region list is used as an input to
// which regions are opted-in (active) for the account.
func listRawAwsRegionsUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	clientRegion, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		logger.Error("listRawAwsRegionsUncached", "connection_name", d.Connection.Name, "clientRegion", clientRegion, "region_error", err)
		return nil, err
	}

	// Create Session
	svc, err := EC2LowRetryClientForRegion(ctx, d, clientRegion)
	if err != nil {
		logger.Error("listRawAwsRegionsUncached", "connection_name", d.Connection.Name, "clientRegion", clientRegion, "connnection_error", err)
		return nil, err
	}

	params := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	}

	// execute list call
	resp, err := svc.DescribeRegions(ctx, params)
	if err != nil {
		logger.Error("listRawAwsRegionsUncached", "connection_name", d.Connection.Name, "clientRegion", clientRegion, "params", params, "api_error", err)
		return nil, err
	}

	logger.Debug("listRawAwsRegionsUncached", "connection_name", d.Connection.Name, "clientRegion", clientRegion, "len(resp.Regions)", len(resp.Regions))

	return resp.Regions, nil
}

// The "last resort" region is generally the oldest / best final failsafe region
// to use for a given partition. For example, in AWS Commercial it's us-east-1.
// This region is used for API calls that must go to the base endpoint. In general,
// it's better to use the client region (see getDefaultRegion) if possible, this
// should be the last resort.
func getLastResortRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (string, error) {

	region, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getLastResortRegion", "connection_name", d.Connection.Name, "default_region_error", err)
		return "", err
	}

	plugin.Logger(ctx).Debug("getLastResortRegion", "connection_name", d.Connection.Name, "region", region)

	// Get the last resort region for the partition
	lastResortRegion := awsLastResortRegionFromRegionWildcard(region)
	plugin.Logger(ctx).Debug("getLastResortRegion", "connection_name", d.Connection.Name, "lastResortRegion", lastResortRegion)
	if lastResortRegion != "" {
		return lastResortRegion, nil
	}

	// If the given region didn't match any known partition then we are
	// stuck.
	return "", fmt.Errorf("cannot calculate last resort region for default region %s", region)
}

// Get the default region for AWS API calls that need to go to a central /
// non-regional endpoint (e.g. describe EC2 regions). Typically this should be
// the region closest to the user.
func getDefaultRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (string, error) {
	regionInterface, err := getDefaultRegionCached(ctx, d, h)
	if err != nil {
		return "", err
	}
	region := regionInterface.(string)
	return region, nil
}

// The default region is cached on a per-connection basis to prevent re-lookup
// and recalculations from the configuration.
var getDefaultRegionCached = plugin.HydrateFunc(getDefaultRegionUncached).Memoize()

// Helper function to get the default region from the Steampipe config
func getAwsSpcConfigDefaultRegion(ctx context.Context, d *plugin.QueryData) string {
	awsSpcConfig := GetConfig(d.Connection)

	if awsSpcConfig.DefaultRegion == nil {
		return ""
	}

	region := *awsSpcConfig.DefaultRegion
	plugin.Logger(ctx).Debug("getAwsSpcConfigDefaultRegion", "connection_name", d.Connection.Name, "region", region)
	return region
}

// Helper function to get the region from the AWS SDK. This will use the region
// defined in the AWS config files, or the AWS_REGION environment variable,
// or the default region for the partition.
func getAwsSdkRegion(ctx context.Context, d *plugin.QueryData) string {
	cfg, err := getBaseClientForAccount(ctx, d)
	if cfg != nil && cfg.Region != "" && err == nil {
		region := cfg.Region
		plugin.Logger(ctx).Debug("getAwsSdkRegion", "connection_name", d.Connection.Name, "region", region)
		// The AWS SDK will return us-east-1 if it can't find a region. So, we
		// can only trust regions other than us-east-1 as being intentional from
		// the config. Return those regions immediately as the default.
		// If it is us-east-1, then fall through to check the regions config for
		// a more reliable indication. If it's from a commercial partition (or
		// not specified) then we fall through to us-east-1 at the end anyway.
		if region != "us-east-1" {
			return region
		}
	}

	return ""
}

// Helper function to get the last resort region for the partition based on the
// list of regions in the Steampipe config if any of them have enough
// information to indicate our preferred partition.
func awsLastResortRegionFromRegionsConfig(ctx context.Context, d *plugin.QueryData) string {
	awsSpcConfig := GetConfig(d.Connection)

	if awsSpcConfig.Regions != nil {
		for _, r := range awsSpcConfig.Regions {
			lastResort := awsLastResortRegionFromRegionWildcard(r)
			if lastResort != "" {
				plugin.Logger(ctx).Debug("awsLastResortRegionFromRegionsConfig", "connection_name", d.Connection.Name, "region", lastResort)
				return lastResort
			}
		}
	}

	return ""
}

// Calculate the region we want to use by default for the plugin.
// It's complicated, because there are many different configuration sources
// (spc files, environment variables, AWS config files, etc.) and some choices
// are filters (e.g. regions = [ "*" ]).
// Here is the order of precedence:
// 1. default_region in the aws.spc file.
// 2. The region as calculated by the AWS SDK (AWS_REGION env var, ~/.aws/config, etc).
// 3. The last resort region for the partition best matched by each region added to regions in the aws.spc file.
// 4. us-east-1 (last resort region for the most common partition).
func getDefaultRegionUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	var region string

	plugin.Logger(ctx).Debug("getDefaultRegionUncached", "connection_name", d.Connection.Name)

	// The user has defined a specific default_region in their config. We use
	// it without further review. For example, they can have a default_region
	// that is not in the regions list.
	// Notes on default_region:
	// - It can be a region that is not in the regions list.
	region = getAwsSpcConfigDefaultRegion(ctx, d)
	if region != "" {
		plugin.Logger(ctx).Debug("getDefaultRegionUncached", "connection_name", d.Connection.Name, "region", region, "source", "default_region in config file")
		return region, nil
	}

	// Get the region from the AWS SDK. This will use the region defined in the
	// AWS config files, the AWS_REGION environment variable, or the default
	// region for the partition.
	region = getAwsSdkRegion(ctx, d)
	if region != "" {
		plugin.Logger(ctx).Debug("getDefaultRegionUncached", "connection_name", d.Connection.Name, "region", region, "source", "AWS SDK resolution")
		return region, nil
	}

	// Look through the list of regions, checking if any of them have enough
	// information to indicate our preferred partition. If available, then we
	// use the last resort region for the partition.
	// We've decided that it's better to default to a last resort region rather
	// than make the order of the regions list significant. For example, someone
	// may make the list alphabetical, and suddenly their first region is in Asia
	// Pacific.
	region = awsLastResortRegionFromRegionsConfig(ctx, d)
	if region != "" {
		plugin.Logger(ctx).Debug("getDefaultRegionUncached", "connection_name", d.Connection.Name, "region", region, "source", "best guess from regions config")
		return region, nil
	}

	// If all else fails, and we just don't know what to do ... default to
	// us-east-1 (the last resort region for the most common partition).
	region = "us-east-1"
	plugin.Logger(ctx).Debug("getDefaultRegionUncached", "connection_name", d.Connection.Name, "region", region, "source", "last resort region in most common partition")
	return region, nil
}

// Calculate the region we want to use for the plugin based on the Steampipe
// config.
// Unlike getDefaultRegionUncached, do not attempt to get the region from the
// SDK to avoid a circular dependency, since this function will primarily be
// used in the getBaseClient function.
// Here is the order of precedence:
// 1. default_region in the aws.spc file.
// 2. The last resort region for the partition best matched by each region added to regions in the aws.spc file.
// 3. us-east-1 (last resort region for the most common partition).
func getDefaultRegionFromConfig(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (string, error) {

	var region string

	plugin.Logger(ctx).Debug("getDefaultRegionFromConfig", "connection_name", d.Connection.Name)

	// The user has defined a specific default_region in their config. We use
	// it without further review. For example, they can have a default_region
	// that is not in the regions list.
	// Notes on default_region:
	// - It can be a region that is not in the regions list.
	region = getAwsSpcConfigDefaultRegion(ctx, d)
	if region != "" {
		plugin.Logger(ctx).Debug("getDefaultRegionFromConfig", "connection_name", d.Connection.Name, "region", region, "source", "default_region in config file")
		return region, nil
	}

	// Look through the list of regions, checking if any of them have enough
	// information to indicate our preferred partition. If available, then we
	// use the last resort region for the partition.
	// We've decided that it's better to default to a last resort region rather
	// than make the order of the regions list significant. For example, someone
	// may make the list alphabetical, and suddenly their first region is in Asia
	// Pacific.
	region = awsLastResortRegionFromRegionsConfig(ctx, d)
	if region != "" {
		plugin.Logger(ctx).Debug("getDefaultRegionFromConfig", "connection_name", d.Connection.Name, "region", region, "source", "best guess from regions config")
		return region, nil
	}

	// If all else fails, and we just don't know what to do ... default to
	// us-east-1 (the last resort region for the most common partition).
	region = "us-east-1"
	plugin.Logger(ctx).Debug("getDefaultRegionFromConfig", "connection_name", d.Connection.Name, "region", region, "source", "last resort region in most common partition")
	return region, nil
}

func getRegionByPartition(partition string) []string {
	regionsByPartition := []string{}

	partitionInfo, err := getPartitionValueByPartitionName(partition)
	if err != nil {
		panic("failed to get the partition info with given partition '" + partition + "': " + err.Error())
	}

	if partitionInfo != nil {
		for region := range partitionInfo.Regions {
			regionsByPartition = append(regionsByPartition, region)
		}
	}

	return regionsByPartition
}

// Given a region (including wildcards), guess at the best last resort region
// based on the partition. Examples:
//
//	us-gov-* -> us-gov-west-1
//	cn* -> cn-northwest-1
//	us-west-2 -> us-east-1
//	* -> us-east-1
//	crap -> ""
func awsLastResortRegionFromRegionWildcard(regionWildcard string) string {

	// Check prefixes for obscure partitions
	if strings.HasPrefix(regionWildcard, "us-gov") {
		return "us-gov-west-1"
	} else if strings.HasPrefix(regionWildcard, "cn") {
		return "cn-northwest-1"
	} else if strings.HasPrefix(regionWildcard, "us-isob") {
		return "us-isob-east-1"
	} else if strings.HasPrefix(regionWildcard, "us-iso") {
		return "us-iso-east-1"
	}

	// Check if the prefix is for a commercial region.
	// Must be done after obscure partitions, because they have the same
	// prefixes but longer.
	for _, prefix := range awsCommercialRegionPrefixes() {
		if strings.HasPrefix(regionWildcard, prefix) {
			return "us-east-1"
		}
	}

	// Unknown partition
	return ""
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

func awsCommercialRegionPrefixes() []string {
	return []string{
		"af",
		"ap",
		"ca",
		"eu",
		"me",
		"sa",
		"us",
	}
}
