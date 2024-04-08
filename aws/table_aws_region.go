package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRegion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_region",
		Description: "AWS Region",
		List: &plugin.ListConfig{
			Hydrate: listAwsRegions,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeRegions"},
		},
		// Get is not implemented because the API is not paged anyway, so
		// the List has the same cost but better caching benefit.
		Columns: awsAccountColumns([]*plugin.Column{
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
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionName"),
			},
			{
				Name:        "connection_config_set",
				Description: "True if the region is specified in the connection config file.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAWSRegionsInConfig,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

func listAwsRegions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	iRegions, err := listRawAwsRegions(ctx, d, h)
	if err != nil {
		return nil, err
	}
	for _, region := range iRegions.([]types.Region) {
		d.StreamListItem(ctx, region)
	}
	return nil, nil
}

func getAwsRegionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := h.Item.(types.Region)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + "::" + *region.RegionName + ":" + commonColumnData.AccountId}
	return akas, nil
}

func getAWSRegionsInConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := h.Item.(types.Region)

	// Retrieve regions list from the AWS plugin steampipe connection config
	configRegions, err := listQueryRegionsForConnection(ctx, d)
	if err != nil {
		return nil, err
	}

	// check if the region is set in connection config
	if helpers.StringSliceContains(configRegions, *region.RegionName) {
		return true, nil
	} else {
		return false, nil
	}
}
