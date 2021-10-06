package aws

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-sdk-go/service/cloudcontrolapi"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsCCCloudTrailTrail(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "awscc_cloudtrail_trail",
		Description: "AWS Cloud Control CloudTrail Trail",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("trail_name"),
			Hydrate:    getResourceCloudTrailTrail,
		},
		List: &plugin.ListConfig{
			Hydrate: listResourcesCloudTrailTrails,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "trail_name",
				Description: "The name of the trail. NOTE: Description was missing from schema.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the trail. NOTE: Missing description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cloudwatch_logs_log_group_arn",
				Description: "Specifies a log group name using an Amazon Resource Name (ARN), a unique identifier that represents the log group to which CloudTrail logs will be delivered. Not required unless you specify CloudWatchLogsRoleArn.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getResourceCloudTrailTrail,
				Transform:   transform.FromField("CloudWatchLogsLogGroupArn"),
			},
			{
				Name:        "cloudwatch_logs_role_arn",
				Description: "Specifies the role for the CloudWatch Logs endpoint to assume to write to a user's log group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getResourceCloudTrailTrail,
				Transform:   transform.FromField("CloudWatchLogsRoleArn"),
			},
			{
				Name:        "enable_log_file_validation",
				Description: "Specifies whether log file validation is enabled. The default is false.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "event_selectors",
				Description: "Use event selectors to further specify the management and data event settings for your trail. By default, trails created without specific event selectors will be configured to log all read and write management events, and no data events. When an event occurs in your account, CloudTrail evaluates the event selector for all trails. For each trail, if the event matches any event selector, the trail processes and logs the event. If the event doesn't match any event selector, the trail doesn't log the event. You can configure up to five event selectors for a trail.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "include_global_service_events",
				Description: "Specifies whether the trail is publishing events from global services such as IAM to the log files.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "insight_selectors",
				Description: "Lets you enable Insights event logging by specifying the Insights selectors that you want to enable on an existing trail.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "is_logging",
				Description: "Specifies whether the CloudTrail is currently logging AWS API calls.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_multi_region_trail",
				Description: "Specifies whether the trail applies only to the current region or to all regions. The default is false. If the trail exists only in the current region and this value is set to true, shadow trails (replications of the trail) will be created in the other regions. If the trail exists in all regions and this value is set to false, the trail will remain in the region where it was created, and its shadow trails in other regions will be deleted. As a best practice, consider using trails that log events in all regions.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_organization_trail",
				Description: "Specifies whether the trail is created for all accounts in an organization in AWS Organizations, or only for the current AWS account. The default is false, and cannot be true unless the call is made on behalf of an AWS account that is the master account for an organization in AWS Organizations.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kms_key_id",
				Description: "Specifies the KMS key ID to use to encrypt the logs delivered by CloudTrail. The value can be an alias name prefixed by 'alias/', a fully specified ARN to an alias, a fully specified ARN to a key, or a globally unique identifier.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getResourceCloudTrailTrail,
				Transform:   transform.FromField("KMSKeyId"),
			},
			{
				Name:        "s3_bucket_name",
				Description: "Specifies the name of the Amazon S3 bucket designated for publishing log files. See Amazon S3 Bucket Naming Requirements.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "s3_key_prefix",
				Description: "Specifies the Amazon S3 key prefix that comes after the name of the bucket you have designated for log file delivery. For more information, see Finding Your CloudTrail Log Files. The maximum length is 200 characters.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sns_topic_arn",
				Description: "Specifies the ARN of the Amazon SNS topic that CloudTrail uses to send notifications when log files are delivered. NOTE: Missing description.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sns_topic_name",
				Description: "Specifies the name of the Amazon SNS topic defined for notification of log file delivery. The maximum length is 256 characters.",
				Hydrate:     getResourceCloudTrailTrail,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the trail. NOTE: Missing description.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getResourceCloudTrailTrail,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "raw",
				Description: "Raw data.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				//Transform:   transform.FromField("Tags").Transform(getCloudTrailTrailTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrailName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listResourcesCloudTrailTrails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listResourcesCloudTrailTrails")

	// Create session
	svc, err := CloudControlService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := cloudcontrolapi.ListResourcesInput{TypeName: types.String("AWS::CloudTrail::Trail")}

	err = svc.ListResourcesPages(&input,
		func(page *cloudcontrolapi.ListResourcesOutput, isLast bool) bool {
			for _, trail := range page.ResourceDescriptions {
				properties := trail.Properties
				var jsonMap map[string]interface{}
				json.Unmarshal([]byte(*properties), &jsonMap)

				d.StreamListItem(ctx, jsonMap)
				// This will return zero if context has been cancelled (i.e due to manual cancellation) or
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return true
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getResourceCloudTrailTrail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getResourceCloudTrailTrail")

	// Create session
	svc, err := CloudControlService(ctx, d)
	if err != nil {
		return nil, err
	}

	var trailName string

	if h.Item != nil {
		result := h.Item
		trailData := reflect.ValueOf(result).MapIndex(reflect.ValueOf("TrailName"))
		trailName = trailData.Interface().(string)
	} else {
		trailName = d.KeyColumnQuals["trail_name"].GetStringValue()
	}

	input := &cloudcontrolapi.GetResourceInput{
		Identifier: types.String(trailName),
		TypeName:   types.String("AWS::CloudTrail::Trail"),
	}

	item, err := svc.GetResource(input)
	if err != nil {
		return nil, err
	}

	properties := item.ResourceDescription.Properties
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(*properties), &jsonMap)

	return jsonMap, nil
}

//// TRANSFORM FUNCTIONS

/*
type resourceTag struct {
	Key   string
	Value string
}

func getCloudTrailTrailTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var tags []resourceTag
	//tags := d.HydrateItem.([]*cloudtrail.Tag)
	tags := d.HydrateItem.(map[string]interface{})
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
*/
