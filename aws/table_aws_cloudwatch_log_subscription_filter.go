package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsCloudwatchLogSubscriptionFilter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_subscription_filter",
		Description: "AWS CloudWatch Log Subscription Filter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "log_group_name"}),
			Hydrate:    getCloudwatchLogSubscriptionFilter,
		},
		List: &plugin.ListConfig{
			Hydrate:       listCloudwatchLogSubscriptionFilters,
			ParentHydrate: listCloudwatchLogGroups,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "log_group_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the subscription filter.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FilterName"),
			},
			{
				Name:        "log_group_name",
				Description: "The name of the log group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the subscription filter.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime").Transform(convertTimestamp),
			},
			{
				Name:        "destination_arn",
				Description: "The Amazon Resource Name (ARN) of the destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "distribution",
				Description: "The method used to distribute log data to the destination.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter_pattern",
				Description: "A symbolic description of how CloudWatch Logs should interpret the data in each log event.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard column
			{
				Name:        "role_arn",
				Description: "The role associated to the filter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FilterName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudwatchLogSubscriptionFilterAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudwatchLogSubscriptionFilters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logGroup := h.Item.(*cloudwatchlogs.LogGroup)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_subscription_filter.listCloudwatchLogSubscriptionFilters", "service_creation_error", err)
		return nil, err
	}

	input := &cloudwatchlogs.DescribeSubscriptionFiltersInput{
		Limit:        aws.Int64(50),
		LogGroupName: logGroup.LogGroupName,
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.FilterNamePrefix = aws.String(equalQuals["name"].GetStringValue())
	}
	if equalQuals["log_group_name"] != nil {
		if *logGroup.LogGroupName != equalQuals["log_group_name"].GetStringValue() {
			return nil, nil
		}
		input.LogGroupName = aws.String(equalQuals["log_group_name"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = aws.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	err = svc.DescribeSubscriptionFiltersPages(
		input,
		func(page *cloudwatchlogs.DescribeSubscriptionFiltersOutput, isLast bool) bool {
			for _, subscriptionFilter := range page.SubscriptionFilters {
				d.StreamListItem(ctx, subscriptionFilter)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogSubscriptionFilter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	if d.KeyColumnQuals["name"] == nil || d.KeyColumnQuals["log_group_name"] == nil {
		return nil, nil
	}
	name := d.KeyColumnQuals["name"].GetStringValue()
	logGroupName := d.KeyColumnQuals["log_group_name"].GetStringValue()

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_subscription_filter.getCloudwatchLogSubscriptionFilter", "service_creation_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.DescribeSubscriptionFiltersInput{
		FilterNamePrefix: &name,
		LogGroupName:     &logGroupName,
	}

	// execute list call
	op, err := svc.DescribeSubscriptionFilters(params)
	if err != nil {
		return nil, err
	}

	for _, subscriptionFilter := range op.SubscriptionFilters {
		if *subscriptionFilter.FilterName == name {
			plugin.Logger(ctx).Error("aws_cloudwatch_log_subscription_filter.getCloudwatchLogSubscriptionFilter", "api_error", err)
			return subscriptionFilter, nil
		}
	}
	return nil, nil
}

func getCloudwatchLogSubscriptionFilterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	subscriptionFilter := h.Item.(*cloudwatchlogs.SubscriptionFilter)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_subscription_filter.getCloudwatchLogSubscriptionFilterAkas", "cache_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":logs:" + region + ":" + commonColumnData.AccountId + ":log-group:" + *subscriptionFilter.LogGroupName + ":subscription-filter:" + *subscriptionFilter.FilterName}

	return akas, nil
}
