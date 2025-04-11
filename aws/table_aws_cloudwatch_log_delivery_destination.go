package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudWatchLogDeliveryDestination(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_delivery_destination",
		Description: "AWS CloudWatch Log Delivery Destination",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudWatchLogDeliveryDestination,
			Tags:       map[string]string{"service": "logs", "action": "GetDeliveryDestination"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchLogDeliveryDestinations,
			Tags:    map[string]string{"service": "logs", "action": "ListDeliveryDestinations"},
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the delivery destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies this delivery destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "destination_resource_arn",
				Description: "The ARN of the Amazon Web Services destination that this delivery destination represents.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeliveryDestinationConfiguration.DestinationResourceArn"),
			},
			{
				Name:        "delivery_destination_type",
				Description: "Displays whether this delivery destination is CloudWatch Logs, Amazon S3, or Firehose.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "output_format",
				Description: "The format of the logs that are sent to this delivery destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy",
				Description: "The policy of the delivery destination.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudWatchLogDeliveryDestinationPolicy,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
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

func listCloudWatchLogDeliveryDestinations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_destination.listCloudWatchLogDeliveryDestinations", "connection_error", err)
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

	input := &cloudwatchlogs.DescribeDeliveryDestinationsInput{
		Limit: &maxLimit,
	}

	paginator := cloudwatchlogs.NewDescribeDeliveryDestinationsPaginator(svc, input, func(o *cloudwatchlogs.DescribeDeliveryDestinationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_destination.listCloudWatchLogDeliveryDestinations", "api_error", err)
			return nil, err
		}

		for _, destination := range output.DeliveryDestinations {
			d.StreamListItem(ctx, destination)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudWatchLogDeliveryDestination(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	destinationName := d.EqualsQuals["name"].GetStringValue()

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_destination.getCloudWatchLogDeliveryDestination", "connection_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.GetDeliveryDestinationInput{
		Name: &destinationName,
	}

	op, err := svc.GetDeliveryDestination(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_destination.getCloudWatchLogDeliveryDestination", "api_error", err)
		return nil, err
	}

	if op.DeliveryDestination != nil {
		return *op.DeliveryDestination, nil
	}

	return nil, nil
}

func getCloudWatchLogDeliveryDestinationPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(types.DeliveryDestination)

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_destination.getCloudWatchLogDeliveryDestination", "connection_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.GetDeliveryDestinationPolicyInput{
		DeliveryDestinationName: item.Name,
	}

	op, err := svc.GetDeliveryDestinationPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_destination.getCloudWatchLogDeliveryDestination", "api_error", err)
		return nil, err
	}

	if op.Policy != nil {
		return op.Policy, nil
	}

	return nil, nil
}
