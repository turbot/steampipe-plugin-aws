package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudtrailChannel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_channel",
		Description: "AWS CloudTrail Channel",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getCloudTrailChannel,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ChannelNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudTrailChannels,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the cloudtrail channel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of a channel.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChannelArn"),
			},
			{
				Name:        "apply_to_all_regions",
				Description: "Specifies whether the channel applies to a single region or to all regions.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudTrailChannel,
				Transform:   transform.FromField("SourceConfig.ApplyToAllRegions"),
			},
			{
				Name:        "source",
				Description: "The event source for the cloudtrail channel.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailChannel,
			},
			{
				Name:        "advanced_event_selectors",
				Description: "The advanced event selectors that are configured for the channel.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudTrailChannel,
				Transform:   transform.FromField("SourceConfig.AdvancedEventSelectors"),
			},
			{
				Name:        "destinations",
				Description: "The Amazon Web Services service that created the service-linked channel.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudTrailChannel,
			},
			{
				Name:        "source_config",
				Description: "Configuration information about the channel.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudTrailChannel,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ChannelArn").Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudTrailChannels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_channel.listCloudTrailChannels", "client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &cloudtrail.ListChannelsInput{
		// Default to the maximum allowed
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := cloudtrail.NewListChannelsPaginator(svc, input, func(o *cloudtrail.ListChannelsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudtrail_channel.listCloudTrailChannels", "api_error", err)
			return nil, err
		}

		for _, item := range output.Channels {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudTrailChannel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var channelArn string
	if h.Item != nil {
		channelArn = *h.Item.(types.Channel).ChannelArn
	} else {
		channelArn = d.KeyColumnQualString("arn")
	}

	// Empty Check
	if channelArn == "" {
		return nil, nil
	}

	// Create session
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_channel.getCloudTrailChannel", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.GetChannelInput{
		Channel: &channelArn,
	}

	// execute list call
	op, err := svc.GetChannel(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_channel.getCloudTrailChannel", "api_error", err)
		return nil, err
	}

	return op, nil
}
