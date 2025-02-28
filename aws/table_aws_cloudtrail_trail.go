package aws

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	cloudtrailv1 "github.com/aws/aws-sdk-go/service/cloudtrail"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudtrailTrail(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_trail",
		Description: "AWS CloudTrail Trail",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name", "arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidTrailNameException", "TrailNotFoundException", "CloudTrailARNInvalidException"}),
			},
			Hydrate: getCloudtrailTrail,
			Tags:    map[string]string{"service": "cloudtrail", "action": "DescribeTrails"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudtrailTrails,
			Tags:    map[string]string{"service": "cloudtrail", "action": "DescribeTrails"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidTrailNameException", "TrailNotFoundException", "CloudTrailARNInvalidException"}),
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudtrailTrailStatus,
				Tags: map[string]string{"service": "cloudtrail", "action": "GetTrailStatus"},
			},
			{
				Func: getCloudtrailTrailEventSelector,
				Tags: map[string]string{"service": "cloudtrail", "action": "GetEventSelectors"},
			},
			{
				Func: getCloudtrailTrailInsightSelector,
				Tags: map[string]string{"service": "cloudtrail", "action": "GetInsightSelectors"},
			},
			{
				Func: getCloudtrailTrailTags,
				Tags: map[string]string{"service": "cloudtrail", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudtrailv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the trail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the trail.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrailARN"),
			},
			{
				Name:        "home_region",
				Description: "The region in which the trail was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_multi_region_trail",
				Description: "Specifies whether the trail exists only in one region or exists in all regions.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "log_file_validation_enabled",
				Description: "Specifies whether log file validation is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_logging",
				Description: "Specifies whether the CloudTrail is currently logging AWS API calls, or not.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "log_group_arn",
				Description: "Specifies an Amazon Resource Name (ARN), a unique identifier that represents the log group to which CloudTrail logs will be delivered.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CloudWatchLogsLogGroupArn"),
			},
			{
				Name:        "cloudwatch_logs_role_arn",
				Description: "Specifies the role for the CloudWatch Logs endpoint to assume to write to a user's log group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CloudWatchLogsRoleArn"),
			},
			{
				Name:        "has_custom_event_selectors",
				Description: "Specifies whether the trail has custom event selectors, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "has_insight_selectors",
				Description: "Specifies whether a trail has insight types specified in an InsightSelector list, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "include_global_service_events",
				Description: "Specifies whether to include AWS API calls from AWS global services, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_organization_trail",
				Description: "Specifies whether the trail is an organization trail, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kms_key_id",
				Description: "Specifies the KMS key ID that encrypts the logs delivered by CloudTrail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "s3_bucket_name",
				Description: "Name of the Amazon S3 bucket into which CloudTrail delivers your trail files.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "s3_key_prefix",
				Description: "Specifies the Amazon S3 key prefix that comes after the name of the bucket you have designated for log file delivery.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sns_topic_arn",
				Description: "Specifies the ARN of the Amazon SNS topic that CloudTrail uses to send notifications when log files are delivered.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnsTopicARN"),
			},

			// details of trail status
			{
				Name:        "latest_cloudwatch_logs_delivery_error",
				Description: "Displays any CloudWatch Logs error that CloudTrail encountered when attempting to deliver logs to CloudWatch Logs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudtrailTrailStatus,
				Transform:   transform.FromField("LatestCloudWatchLogsDeliveryError"),
			},
			{
				Name:        "latest_cloudwatch_logs_delivery_time",
				Description: "Displays the most recent date and time when CloudTrail delivered logs to CloudWatch Logs.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudtrailTrailStatus,
				Transform:   transform.FromField("LatestCloudWatchLogsDeliveryTime"),
			},
			{
				Name:        "latest_delivery_error",
				Description: "Displays any Amazon S3 error that CloudTrail encountered when attempting to deliver log files to the designated bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "latest_delivery_time",
				Description: "Specifies the date and time that CloudTrail last delivered log files to an account's Amazon S3 bucket.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "latest_digest_delivery_error",
				Description: "Displays any Amazon S3 error that CloudTrail encountered when attempting to deliver a digest file to the designated bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "latest_digest_delivery_time",
				Description: "Specifies the date and time that CloudTrail last delivered a digest file to an account's Amazon S3 bucket.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "latest_notification_error",
				Description: "Displays any Amazon SNS error that CloudTrail encountered when attempting to send a notification.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "latest_notification_time",
				Description: "Specifies the date and time of the most recent Amazon SNS notification that CloudTrail has written a new log file to an account's Amazon S3 bucket.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "start_logging_time",
				Description: "Specifies the most recent date and time when CloudTrail started recording API calls for an AWS account.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "stop_logging_time",
				Description: "Specifies the most recent date and time when CloudTrail stopped recording API calls for an AWS account.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudtrailTrailStatus,
			},
			{
				Name:        "advanced_event_selectors",
				Description: "Describes the advanced event selectors that are configured for the trail.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudtrailTrailEventSelector,
			},
			{
				Name:        "event_selectors",
				Description: "Describes the event selectors that are configured for the trail.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudtrailTrailEventSelector,
			},
			{
				Name:        "insight_selectors",
				Description: "A JSON string that contains the insight types you want to log on a trail.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudtrailTrailInsightSelector,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the trail.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudtrailTrailTags,
				Transform:   transform.FromValue(),
			},

			// steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudtrailTrailTags,
				Transform:   transform.From(getCloudtrailTrailTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TrailARN").Transform(arnToAkas),
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

func listCloudtrailTrails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.listCloudtrailTrails", "client_error", err)
		return nil, err
	}

	// Doesn't support paginator for CloudTrail DescribeTrails API
	resp, err := svc.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{})
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.listCloudtrailTrails", "api_error", err)
		return nil, err
	}

	for _, trail := range resp.TrailList {
		d.StreamListItem(ctx, trail)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudtrailTrail(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	arn := d.EqualsQuals["arn"].GetStringValue()
	if len(arn) > 0 {
		data := strings.Split(arn, "/")
		name = data[len(data)-1]
	}

	if d.EqualsQuals["name"] != nil && d.EqualsQuals["name"].GetStringValue() == "" {
		return nil, nil
	}

	// Create session
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailStatus", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.DescribeTrailsInput{
		TrailNameList: []string{name},
	}

	// execute list call
	item, err := svc.DescribeTrails(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrail", "api_error", err)
		return nil, err
	}

	if len(item.TrailList) > 0 {
		return item.TrailList[0], nil
	}

	return nil, nil
}

func getCloudtrailTrailStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailStatus", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	trail := h.Item.(types.Trail)

	// Avoid API call if Account ID of the client is not equal to the Account ID available in Trail ARN
	accountId := arnToAccountId(*trail.TrailARN)
	if commonColumnData.AccountId != accountId {
		return nil, nil
	}

	// Create session
	svc, err := CloudTrailRegionsClient(ctx, d, *trail.HomeRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailStatus", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.GetTrailStatusInput{
		Name: trail.TrailARN,
	}

	item, err := svc.GetTrailStatus(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if slices.Contains([]string{"TrailNotFoundException", "CloudTrailARNInvalidException"}, ae.ErrorCode()) {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailStatus", "api_error", err)
		return nil, err
	}

	return item, nil
}

func getCloudtrailTrailEventSelector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailEventSelector", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	trail := h.Item.(types.Trail)

	// Avoid api call if accountId is not equal to the accountId available in arn
	accountId := arnToAccountId(*trail.TrailARN)
	if commonColumnData.AccountId != accountId {
		return nil, nil
	}

	// Create session
	svc, err := CloudTrailRegionsClient(ctx, d, *trail.HomeRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailEventSelector", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.GetEventSelectorsInput{
		TrailName: trail.TrailARN,
	}

	// List resource tags
	item, err := svc.GetEventSelectors(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if slices.Contains([]string{"TrailNotFoundException", "CloudTrailARNInvalidException"}, ae.ErrorCode()) {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailEventSelector", "api_error", err)
		return nil, err
	}
	return item, nil
}

func getCloudtrailTrailInsightSelector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailInsightSelector", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	trail := h.Item.(types.Trail)

	// Avoid api call if accountId is not equal to the accountId available in arn
	accountId := arnToAccountId(*trail.TrailARN)
	if commonColumnData.AccountId != accountId {
		return nil, nil
	}

	// Create session
	svc, err := CloudTrailRegionsClient(ctx, d, *trail.HomeRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailInsightSelector", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.GetInsightSelectorsInput{
		TrailName: trail.TrailARN,
	}

	// List resource tags
	item, err := svc.GetInsightSelectors(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if slices.Contains([]string{"InsightNotEnabledException"}, ae.ErrorCode()) {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailInsightSelector", "api_error", err)
		return nil, err
	}
	return item, nil
}

func getCloudtrailTrailTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailTags", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	trail := h.Item.(types.Trail)

	var traiTag []types.Tag

	// Avoid api call if accountId is not equal to the accountId available in arn
	accountId := arnToAccountId(*trail.TrailARN)
	if commonColumnData.AccountId != accountId {
		return traiTag, nil
	}

	// Create session
	svc, err := CloudTrailRegionsClient(ctx, d, *trail.HomeRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailEventSelector", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.ListTagsInput{
		ResourceIdList: []string{*trail.TrailARN},
	}

	resp, err := svc.ListTags(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if slices.Contains([]string{"TrailNotFoundException", "CloudTrailARNInvalidException"}, ae.ErrorCode()) {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_cloudtrail_trail.getCloudtrailTrailEventSelector", "api_error", err)
		return nil, err
	}

	if len(resp.ResourceTagList) > 0 {
		return resp.ResourceTagList[0].TagsList, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getCloudtrailTrailTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)
	if tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}
	return turbotTagsMap, nil
}

func arnToAccountId(arn string) string {
	if arn != "" {
		return strings.Split(arn, ":")[4]
	}
	return ""
}
