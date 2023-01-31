package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsRegion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_region",
		Description: "AWS Region",
		List: &plugin.ListConfig{
			Hydrate: listAwsRegions,
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
