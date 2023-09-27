package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	cloudwatchlogsv1 "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudwatchLogSubscriptionFilter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_subscription_filter",
		Description: "AWS CloudWatch Log Subscription Filter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "log_group_name"}),
			Hydrate:    getCloudwatchLogSubscriptionFilter,
			Tags:       map[string]string{"service": "logs", "action": "DescribeSubscriptionFilters"},
		},
		List: &plugin.ListConfig{
			Hydrate:       listCloudwatchLogSubscriptionFilters,
			ParentHydrate: listCloudwatchLogGroups,
			Tags:          map[string]string{"service": "logs", "action": "DescribeSubscriptionFilters"},
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
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchlogsv1.EndpointsID),
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
	logGroup := h.Item.(types.LogGroup)

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_subscription_filter.listCloudwatchLogSubscriptionFilters", "service_creation_error", err)
		return nil, err
	}

	input := &cloudwatchlogs.DescribeSubscriptionFiltersInput{
		Limit:        aws.Int32(50),
		LogGroupName: logGroup.LogGroupName,
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
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
	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input.Limit = &maxLimit

	paginator := cloudwatchlogs.NewDescribeSubscriptionFiltersPaginator(svc, input, func(o *cloudwatchlogs.DescribeSubscriptionFiltersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_alarm.listCloudWatchAlarms", "api_error", err)
			return nil, err
		}
		for _, subscriptionFilter := range output.SubscriptionFilters {
			d.StreamListItem(ctx, subscriptionFilter)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogSubscriptionFilter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	if d.EqualsQuals["name"] == nil || d.EqualsQuals["log_group_name"] == nil {
		return nil, nil
	}
	name := d.EqualsQuals["name"].GetStringValue()
	logGroupName := d.EqualsQuals["log_group_name"].GetStringValue()

	// Empty input check
	if strings.TrimSpace(name) == "" || strings.TrimSpace(logGroupName) == "" {
		return nil, nil
	}

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_subscription_filter.getCloudwatchLogSubscriptionFilter", "service_creation_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.DescribeSubscriptionFiltersInput{
		FilterNamePrefix: &name,
		LogGroupName:     &logGroupName,
	}

	// execute list call
	op, err := svc.DescribeSubscriptionFilters(ctx, params)
	if err != nil {
		return nil, err
	}

	for _, subscriptionFilter := range op.SubscriptionFilters {
		if *subscriptionFilter.FilterName == name {
			return subscriptionFilter, nil
		}
	}
	return nil, nil
}

func getCloudwatchLogSubscriptionFilterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	subscriptionFilter := h.Item.(types.SubscriptionFilter)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_subscription_filter.getCloudwatchLogSubscriptionFilterAkas", "cache_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := fmt.Sprintf("arn:%s:logs:%s:%s:log-group:%s:subscription-filter:%s", commonColumnData.Partition, region, commonColumnData.AccountId, *subscriptionFilter.LogGroupName, *subscriptionFilter.FilterName)

	// Get data for turbot defined properties
	akas := []string{arn}

	return akas, nil
}
