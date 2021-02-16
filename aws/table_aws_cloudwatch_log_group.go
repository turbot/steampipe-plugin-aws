package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsCloudwatchLogGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_group",
		Description: "AWS CloudWatch Log Group",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("name"),
			ItemFromKey: logGroupFromKey,
			Hydrate:     getCloudwatchLogGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudwatchLogGroups,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the log group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the log group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the log group",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime").Transform(convertTimestamp),
			},
			{
				Name:        "kms_key_id",
				Description: "The Amazon Resource Name (ARN) of the CMK to use when encrypting log data",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_filter_count",
				Description: "The number of metric filters",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "retention_in_days",
				Description: "The number of days to retain the log events in the specified log group. Possible values are: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, and 3653",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "stored_bytes",
				Description: "The number of bytes stored",
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

//// BUILD HYDRATE INPUT

func logGroupFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &cloudwatchlogs.LogGroup{
		LogGroupName: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listCloudwatchLogGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listCloudwatchLogGroups", "AWS_REGION", region)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	err = svc.DescribeLogGroupsPages(
		&cloudwatchlogs.DescribeLogGroupsInput{},
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, isLast bool) bool {
			for _, logGroup := range page.LogGroups {
				d.StreamListItem(ctx, logGroup)
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogGroup")

	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logGroup := h.Item.(*cloudwatchlogs.LogGroup)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: logGroup.LogGroupName,
	}

	// execute list call
	item, err := svc.DescribeLogGroups(params)
	if err != nil {
		return nil, err
	}

	for _, logGroup := range item.LogGroups {
		if *logGroup.LogGroupName == *logGroup.LogGroupName {
			return logGroup, nil
		}
	}

	return nil, nil
}

func getLogGroupTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogGroup")
	logGroup := h.Item.(*cloudwatchlogs.LogGroup)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d, region)
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
