package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudWatchLogDelivery(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_delivery",
		Description: "AWS CloudWatch Log Delivery",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCloudWatchLogDelivery,
			Tags:       map[string]string{"service": "logs", "action": "GetDelivery"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchLogDeliveries,
			Tags:    map[string]string{"service": "logs", "action": "ListDeliveries"},
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique ID that identifies this delivery in your account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies this delivery.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delivery_source_name",
				Description: "The name of the delivery source that is associated with this delivery.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delivery_destination_arn",
				Description: "The ARN of the delivery destination that is associated with this delivery.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delivery_destination_type",
				Description: "Displays whether the delivery destination associated with this delivery is CloudWatch Logs, Amazon S3, or Firehose.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
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

func listCloudWatchLogDeliveries(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery.listCloudWatchLogDeliveries", "connection_error", err)
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

	input := &cloudwatchlogs.DescribeDeliveriesInput{
		Limit: &maxLimit,
	}

	paginator := cloudwatchlogs.NewDescribeDeliveriesPaginator(svc, input, func(o *cloudwatchlogs.DescribeDeliveriesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery.listCloudWatchLogDeliveries", "api_error", err)
			return nil, err
		}

		for _, delivery := range output.Deliveries {
			d.StreamListItem(ctx, delivery)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudWatchLogDelivery(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	deliveryId := d.EqualsQuals["id"].GetStringValue()

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery.getCloudWatchLogDelivery", "connection_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.GetDeliveryInput{
		Id: &deliveryId,
	}

	op, err := svc.GetDelivery(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery.getCloudWatchLogDelivery", "api_error", err)
		return nil, err
	}

	if op.Delivery != nil {
		return *op.Delivery, nil
	}

	return nil, nil
}
