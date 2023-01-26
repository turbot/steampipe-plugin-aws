package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	cloudwatchlogsv1 "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudwatchLogGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_group",
		Description: "AWS CloudWatch Log Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudwatchLogGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudwatchLogGroups,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchlogsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the log group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the log group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the log group.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime").Transform(convertTimestamp),
			},
			{
				Name:        "kms_key_id",
				Description: "The Amazon Resource Name (ARN) of the CMK to use when encrypting log data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_filter_count",
				Description: "The number of metric filters.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "retention_in_days",
				Description: "The number of days to retain the log events in the specified log group. Possible values are: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, and 3653.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "stored_bytes",
				Description: "The number of bytes stored.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogGroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLogGroupTagging,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudwatchLogGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Info("aws_cloudwatch_log_group.listCloudwatchLogGroups", "client_error", err)
		return nil, err
	}

	maxItems := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input := &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: &maxItems,
	}

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(svc, input, func(o *cloudwatchlogs.DescribeLogGroupsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		input.LogGroupNamePrefix = aws.String(equalQuals["name"].GetStringValue())
	}

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Info("aws_cloudwatch_log_group.listCloudwatchLogGroups", "api_error", err)
			return nil, err
		}

		for _, logGroup := range output.LogGroups {
			d.StreamListItem(ctx, logGroup)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Info("aws_cloudwatch_log_group.getCloudwatchLogGroup", "client_error", err)
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}

	params := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(name),
	}

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(svc, params, func(o *cloudwatchlogs.DescribeLogGroupsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Info("aws_cloudwatch_log_group.getCloudwatchLogGroup", "api_error", err)
			return nil, err
		}

		for _, logGroup := range output.LogGroups {
			if *logGroup.LogGroupName == name {
				return logGroup, nil
			}
		}
	}

	return nil, nil
}

func getLogGroupTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logGroup := h.Item.(types.LogGroup)

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Info("aws_cloudwatch_log_group.getLogGroupTagging", "client_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.ListTagsLogGroupInput{
		LogGroupName: logGroup.LogGroupName,
	}

	// List resource tags
	logGroupData, err := svc.ListTagsLogGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Info("aws_cloudwatch_log_group.getLogGroupTagging", "api_error", err)
		return nil, err
	}
	return logGroupData, nil
}
