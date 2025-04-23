package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudWatchLogDestination(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_destination",
		Description: "AWS CloudWatch Log Destination",
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchLogDestinations,
			Tags:    map[string]string{"service": "logs", "action": "ListDestinations"},
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: []*plugin.Column{
			{
				Name:        "destination_name",
				Description: "The name of the destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of this destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the destination.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "role_arn",
				Description: "A role for impersonation, used when delivering log events to the target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_arn",
				Description: "The Amazon Resource Name (ARN) of the physical target where the log events are delivered (for example, a Kinesis stream).",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DestinationName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		},
	}
}

//// LIST FUNCTION

func listCloudWatchLogDestinations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_destination.listCloudWatchLogDestinations", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &cloudwatchlogs.DescribeDestinationsInput{
		Limit: &maxLimit,
	}

	paginator := cloudwatchlogs.NewDescribeDestinationsPaginator(svc, input, func(o *cloudwatchlogs.DescribeDestinationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_log_destination.listCloudWatchLogDestinations", "api_error", err)
			return nil, err
		}

		for _, destination := range output.Destinations {
			d.StreamListItem(ctx, destination)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
