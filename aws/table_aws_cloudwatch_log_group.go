package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
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
		GetMatrixItemFunc: BuildRegionList,
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
	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int64(50),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.LogGroupNamePrefix = types.String(equalQuals["name"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = types.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	err = svc.DescribeLogGroupsPages(
		input,
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, isLast bool) bool {
			for _, logGroup := range page.LogGroups {
				d.StreamListItem(ctx, logGroup)

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

func getCloudwatchLogGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogGroup")

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	params := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(name),
	}

	// execute list call
	item, err := svc.DescribeLogGroups(params)
	if err != nil {
		return nil, err
	}

	for _, logGroup := range item.LogGroups {
		if types.SafeString(logGroup.LogGroupName) == name {
			return logGroup, nil
		}
	}

	return nil, nil
}

func getLogGroupTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogGroup")
	logGroup := h.Item.(*cloudwatchlogs.LogGroup)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &cloudwatchlogs.ListTagsLogGroupInput{
		LogGroupName: logGroup.LogGroupName,
	}

	// List resource tags
	logGroupData, err := svc.ListTagsLogGroup(params)
	if err != nil {
		return nil, err
	}
	return logGroupData, nil
}
