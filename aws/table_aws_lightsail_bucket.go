package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	lightsailEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLightsailBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lightsail_bucket",
		Description: "AWS Lightsail Bucket",
		List: &plugin.ListConfig{
			Hydrate: listLightsailBuckets,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "name", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "lightsail", "action": "GetBuckets"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lightsailEndpoint.LIGHTSAILServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The timestamp when the bucket was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state_code",
				Description: "The state code of the bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Code"),
			},
			{
				Name:        "state_message",
				Description: "A message that describes the state of the bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Message"),
			},
			{
				Name:        "able_to_update_bundle",
				Description: "Indicates whether the bundle that is currently applied to a bucket can be changed to another bundle. You can update a bucket's bundle only one time within a monthly Amazon Web Services billing cycle.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "bundle_id",
				Description: "The ID of the bundle currently applied to the bucket. A bucket bundle specifies the monthly cost, storage space, and data transfer quota for a bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_versioning",
				Description: "Indicates whether object versioning is enabled for the bucket. The following options can be configured: Enabled, Suspended, NeverEnabled.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The Lightsail resource type of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "support_code",
				Description: "The support code for a bucket. Include this code in your email to support when you have questions about a Lightsail bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "The URL of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "An object that describes the location of the bucket, such as the Amazon Web Services Region and Availability Zone.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "readonly_access_accounts",
				Description: "An array of strings that specify the Amazon Web Services account IDs that have read-only access to the bucket.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "access_log_config",
				Description: "An object that describes the access log configuration for the bucket.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "access_rules",
				Description: "An object that describes the access rules of the bucket.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resources_receiving_access",
				Description: "An array of objects that describe Lightsail instances that have access to the bucket.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
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
				Transform:   transform.FromField("Tags").Transform(getLightsailBucketTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listLightsailBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_bucket.listLightsailBuckets", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &lightsail.GetBucketsInput{
		IncludeConnectedResources: aws.Bool(true),
	}

	if d.EqualsQualString("name") != "" {
		input.BucketName = aws.String(d.EqualsQualString("name"))
	}

	// List call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := svc.GetBuckets(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_lightsail_bucket.listLightsailBuckets", "query_error", err)
			return nil, nil
		}

		for _, item := range resp.Buckets {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPageToken != nil {
			input.PageToken = resp.NextPageToken
		} else {
			break
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getLightsailBucketTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)
	var turbotTagsMap map[string]string
	if tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}
