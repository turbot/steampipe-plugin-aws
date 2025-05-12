package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudWatchLogDeliverySource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_delivery_source",
		Description: "AWS CloudWatch Log Delivery Source",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudWatchLogDeliverySource,
			Tags:       map[string]string{"service": "logs", "action": "GetDeliverySource"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudWatchLogDeliverySources,
			Tags:    map[string]string{"service": "logs", "action": "ListDeliverySources"},
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the delivery source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies this delivery source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service",
				Description: "The Amazon Web Services service that is sending logs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_type",
				Description: "The type of log that the source is sending.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_arns",
				Description: "This array contains the ARN of the Amazon Web Services resource that sends logs and is represented by this delivery source",
				Type:        proto.ColumnType_JSON,
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

func listCloudWatchLogDeliverySources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_source.listCloudWatchLogDeliverySources", "connection_error", err)
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

	input := &cloudwatchlogs.DescribeDeliverySourcesInput{
		Limit: &maxLimit,
	}

	paginator := cloudwatchlogs.NewDescribeDeliverySourcesPaginator(svc, input, func(o *cloudwatchlogs.DescribeDeliverySourcesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_source.listCloudWatchLogDeliverySources", "api_error", err)
			return nil, err
		}

		for _, source := range output.DeliverySources {
			d.StreamListItem(ctx, source)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudWatchLogDeliverySource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	sourceName := d.EqualsQuals["name"].GetStringValue()

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_source.getCloudWatchLogDeliverySource", "connection_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.GetDeliverySourceInput{
		Name: &sourceName,
	}

	op, err := svc.GetDeliverySource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_delivery_source.getCloudWatchLogDeliverySource", "api_error", err)
		return nil, err
	}

	if op.DeliverySource != nil {
		return *op.DeliverySource, nil
	}

	return nil, nil
}
