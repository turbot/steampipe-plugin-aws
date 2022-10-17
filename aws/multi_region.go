package aws

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

const matrixKeyRegion = "region"
const matrixKeyServiceCode = "serviceCode"

type RegionsData struct {
	AllRegions      []string
	ActiveRegions   []string
	NotOptedRegions []string
	APIRetrivedList bool
}

var (
	awsCommercialRegions = []string{
		"af-south-1", "ap-east-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ap-southeast-3", "ca-central-1", "eu-central-1", "eu-north-1", "eu-south-1", "eu-west-1", "eu-west-2", "eu-west-3", "me-central-1", "me-south-1", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2"}

	awsUsGovRegions  = []string{"us-gov-east-1", "us-gov-west-1"}
	awsChinaRegions  = []string{"cn-north-1", "cn-northwest-1"}
	awsUsIsoRegions  = []string{"us-iso-east-1", "us-iso-west-1"}
	awsUsIsobRegions = []string{"us-isob-east-1"}
)

func getAllAwsRegions() []string {
	awsRegions := append(awsCommercialRegions, awsUsGovRegions...)
	awsRegions = append(awsRegions, awsChinaRegions...)
	awsRegions = append(awsRegions, awsUsIsoRegions...)
	awsRegions = append(awsRegions, awsUsIsobRegions...)
	return awsRegions
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

	// Cache region list matrix
	cacheKey := "RegionListMatrix"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// Retrieve regions list from the AWS plugin steampipe connection config
	awsConfig := GetConfig(d.Connection)
	defaultAwsRegion := GetDefaultAwsRegion(d)

	// If regions are not mentioned in the plugin steampipe connection config
	// get the default aws region and prepare the matrix
	if awsConfig.Regions == nil {
		plugin.Logger(ctx).Debug("BuildRegionList", "connection_name", d.Connection.Name, "region", defaultAwsRegion)
		matrix := []map[string]interface{}{
			{matrixKeyRegion: defaultAwsRegion},
		}

		// set cache
		d.ConnectionManager.Cache.Set(cacheKey, matrix)
		return matrix
	}

	// If regions are mentioned in the connection config
	// Get the list of AWS region from EC2 DescribeRegions API
	regionData, _ := listRegions(ctx, d)
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

	plugin.Logger(ctx).Debug("BuildRegionList", "connection_name", d.Connection.Name, "regions", strings.Join(finalRegions, ", "))
	matrix := make([]map[string]interface{}, len(finalRegions))
	for i, region := range finalRegions {
		matrix[i] = map[string]interface{}{matrixKeyRegion: region}
	}

	// set cache
	d.ConnectionManager.Cache.Set(cacheKey, matrix)
	return matrix
}

func getInvalidRegions(regions []string, d *plugin.QueryData) []string {
	awsRegions := getAllAwsRegions()

	invalidRegions := []string{}
	for _, region := range regions {
		if !helpers.StringSliceContains(awsRegions, region) {
			invalidRegions = append(invalidRegions, region)
		}
	}
	return invalidRegions
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

func listRegions(ctx context.Context, d *plugin.QueryData) (RegionsData, error) {
	cacheKey := "listRegions"

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(RegionsData), nil
	}

	defaultRegions := awsCommercialRegions

	defaultRegion := GetDefaultAwsRegion(d)
	if strings.HasPrefix(defaultRegion, "us-gov") {
		defaultRegions = awsUsGovRegions
	} else if strings.HasPrefix(defaultRegion, "cn") {
		defaultRegions = awsChinaRegions
	} else if strings.HasPrefix(defaultRegion, "us-isob") {
		defaultRegions = awsUsIsobRegions
	} else if strings.HasPrefix(defaultRegion, "us-iso") {
		defaultRegions = awsUsIsoRegions
	}

	data := RegionsData{
		AllRegions:      defaultRegions,
		ActiveRegions:   defaultRegions,
		APIRetrivedList: false,
	}

	// We can query EC2 for the list of supported regions. If credentials
	// are insufficient this query will retry many times, so we create
	// a special client with a small number of retries to prevent hangs.
	svc, err := Ec2RegionsService(ctx, d, defaultRegion)
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
	resp, err := svc.DescribeRegions(params)
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

func SupportedRegionsForService(ctx context.Context, d *plugin.QueryData, serviceId string) ([]string, error) {
	cacheKey := fmt.Sprintf("supported-regions-%s", serviceId)
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]string), nil
	}

	var validRegions []string
	services := endpoints.AwsPartition().Services()
	serviceInfo, ok := services[serviceId]
	if !ok {
		return nil, fmt.Errorf("SupportedRegionsForService called with invalid service ID: " + serviceId)
	}

	regions := serviceInfo.Regions()
	for rs := range regions {
		validRegions = append(validRegions, rs)
	}

	// set cache
	d.ConnectionManager.Cache.Set(cacheKey, validRegions)
	return validRegions, nil
}

// BuildServiceQuotasServicesRegionList :: return a list of matrix items, one per region-services specified in the connection config
func BuildServiceQuotasServicesRegionList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// cache servicequotas services region matrix
	cacheKey := "ServiceQuotasServicesRegionList"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// get all the services
	services, err := listServiceQuotasServices(ctx, d)
	if err != nil {
		panic(err)
	}

	defaultAwsRegion := GetDefaultAwsRegion(d)
	regionData, _ := listRegions(ctx, d)
	var allRegions []string

	// retrieve regions from connection config
	awsConfig := GetConfig(d.Connection)
	// Get only the regions as required by config file
	if awsConfig.Regions != nil {
		for _, pattern := range awsConfig.Regions {
			for _, validRegion := range regionData.AllRegions {
				if ok, _ := path.Match(pattern, validRegion); ok {
					allRegions = append(allRegions, validRegion)
				}
			}
		}
	}

	if len(allRegions) > 0 {
		uniqueRegions := unique(allRegions)

		if len(getInvalidRegions(uniqueRegions, nil)) > 0 {
			panic("\n\nConnection config has invalid regions: " + strings.Join(getInvalidRegions(uniqueRegions, nil), ", "))
		}

		// Remove inactive regions from the list
		finalRegions := helpers.StringSliceDiff(uniqueRegions, regionData.NotOptedRegions)

		matrix := make([]map[string]interface{}, len(finalRegions)*len(services))
		for i, region := range finalRegions {
			for j, service := range services {
				matrix[len(services)*i+j] = map[string]interface{}{
					matrixKeyRegion:      region,
					matrixKeyServiceCode: *service.ServiceCode,
				}
				plugin.Logger(ctx).Debug("listServiceQuotasServices Matrix", (len(services)*i)+j, matrix[len(services)*i+j])
			}
		}

		// set ServiceQuotasServicesRegionList cache
		d.ConnectionManager.Cache.Set(cacheKey, matrix)

		return matrix
	}

	defaultMatrix := make([]map[string]interface{}, len(services))
	for j, service := range services {
		defaultMatrix[j] = map[string]interface{}{
			matrixKeyRegion:      defaultAwsRegion,
			matrixKeyServiceCode: *service.ServiceCode,
		}
		plugin.Logger(ctx).Debug("listServiceQuotasServices Matrix", j, defaultMatrix[j])
	}

	// set ServiceQuotasServicesRegionList cache
	d.ConnectionManager.Cache.Set(cacheKey, defaultMatrix)

	return defaultMatrix
}

func listServiceQuotasServices(ctx context.Context, d *plugin.QueryData) ([]*servicequotas.ServiceInfo, error) {
	plugin.Logger(ctx).Trace("listServiceQuotasServices")

	serviceCacheKey := "listServiceQuotasServices"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.([]*servicequotas.ServiceInfo), nil
	}

	// Create Session
	svc, err := ServiceQuotasService(ctx, d)
	if err != nil {
		return nil, err
	}

	services := []*servicequotas.ServiceInfo{}
	input := &servicequotas.ListServicesInput{
		MaxResults: aws.Int64(100),
	}

	// List call
	err = svc.ListServicesPages(
		input,
		func(page *servicequotas.ListServicesOutput, isLast bool) bool {
			services = append(services, page.Services...)
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	// save services in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, services)

	return services, err
}
