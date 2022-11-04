package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRegion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_region",
		Description: "AWS Region",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"InvalidParameterValue"}),
			},
			Hydrate: getAwsRegion,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRegions,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the region",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionName"),
			},
			{
				Name:        "opt_in_status",
				Description: "The Region opt-in status. The possible values are opt-in-not-required, opted-in, and not-opted-in",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRegionAkas,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "partition",
				Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionName"),
			},
			{
				Name:        "account_id",
				Description: "The AWS Account ID in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listAwsRegions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	defaultRegion := GetDefaultAwsRegion(d)

	// Create Session
	svc, err := EC2RegionsClient(ctx, d, defaultRegion)
	if err != nil {
		logger.Error("aws_region.listAwsRegions", "connnection.error", err)
		return nil, err
	}

	params := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	}

	// execute list call
	resp, err := svc.DescribeRegions(ctx, params)
	if err != nil {
		logger.Error("aws_region.listAwsRegions", "api.error", err)
		return nil, err
	}

	for _, region := range resp.Regions {
		d.StreamListItem(ctx, region)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsRegion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)

	// Create service
	svc, err := EC2RegionsClient(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}
	regionName := d.KeyColumnQuals["name"].GetStringValue()

	params := &ec2.DescribeRegionsInput{
		AllRegions:  aws.Bool(true),
		RegionNames: []string{regionName},
	}

	// execute list call
	op, err := svc.DescribeRegions(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(op.Regions) > 0 {
		return op.Regions[0], nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getAwsRegionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := h.Item.(types.Region)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + "::" + *region.RegionName + ":" + commonColumnData.AccountId}
	return akas, nil
}
