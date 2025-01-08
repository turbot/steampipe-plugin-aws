package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2/types"

	kinesisanalyticsv2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisAnalyticsV2Application(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesisanalyticsv2_application",
		Description: "AWS Kinesis Analytics V2 Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("application_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getKinesisAnalyticsV2Application,
			Tags:    map[string]string{"service": "kinesisanalytics", "action": "DescribeApplication"},
		},
		List: &plugin.ListConfig{
			Hydrate: listKinesisAnalyticsV2Applications,
			Tags:    map[string]string{"service": "kinesisanalytics", "action": "ListApplications"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(kinesisanalyticsv2Endpoint.KINESISANALYTICSServiceID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getKinesisAnalyticsV2Application,
				Tags: map[string]string{"service": "kinesisanalytics", "action": "DescribeApplication"},
			},
			{
				Func: getKinesisAnalyticsV2ApplicationTags,
				Tags: map[string]string{"service": "kinesisanalytics", "action": "ListTagsForResource"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "application_name",
				Description: "The name of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_version_id",
				Description: "Provides the current application version.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "application_arn",
				Description: "The ARN of the application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationARN"),
			},
			{
				Name:        "application_mode",
				Description: "To create a Managed Service for Apache Flink Studio notebook, you must set the mode to INTERACTIVE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_version_rolled_back_from",
				Description: "If you reverted the application using RollbackApplication , the application version when RollbackApplication was called.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "application_version_rolled_back_to",
				Description: "The version to which you want to roll back the application.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "application_version_updated_from",
				Description: "The previous application version before the latest application update.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "conditional_token",
				Description: "A value you use to implement strong concurrency for application updates.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "application_status",
				Description: "The status of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_timestamp",
				Description: "The current timestamp when the application was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "application_description",
				Description: "The description of the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "last_update_timestamp",
				Description: "The current timestamp when the application was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "runtime_environment",
				Description: "The runtime environment for the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_execution_role",
				Description: "Specifies the IAM role that the application uses to access external resources.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "application_configuration_description",
				Description: "Provides details about the application's Java, SQL, or Scala code and starting parameters.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "application_maintenance_configuration_description",
				Description: "The details of the maintenance configuration for the application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "cloud_watch_logging_option_descriptions",
				Description: "Describes the application Amazon CloudWatch logging options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKinesisAnalyticsV2Application,
			},
			{
				Name:        "tags_src",
				Description: "The key-value tags assigned to the application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKinesisAnalyticsV2ApplicationTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKinesisAnalyticsV2ApplicationTags,
				Transform:   transform.FromField("Tags").Transform(kinesisAnalyticsV2ApplicationTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listKinesisAnalyticsV2Applications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := KinesisAnalyticsV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesisanalyticsv2_application.listKinesisAnalyticsV2Applications", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// List call
	pagesLeft := true
	maxLimit := int32(50)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	params := &kinesisanalyticsv2.ListApplicationsInput{
		Limit: aws.Int32(maxLimit),
	}
	// API doesn't support aws-sdk-go-v2 paginator as of data

	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		result, err := svc.ListApplications(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kinesisanalyticsv2_application.listKinesisAnalyticsV2Applications", "api_error", err)
			return nil, err
		}

		for _, application := range result.ApplicationSummaries {
			d.StreamListItem(ctx, application)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesisanalyticsv2_application.listKinesisAnalyticsV2Applications", "api_error", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKinesisAnalyticsV2Application(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var applicationName string
	if h.Item != nil {
		i := h.Item.(types.ApplicationSummary)
		applicationName = *i.ApplicationName
	} else {
		applicationName = d.EqualsQuals["application_name"].GetStringValue()
	}

	// Create Session
	svc, err := KinesisAnalyticsV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesisanalyticsv2_application.getKinesisAnalyticsV2Application", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &kinesisanalyticsv2.DescribeApplicationInput{
		ApplicationName: &applicationName,
	}

	// Get call
	data, err := svc.DescribeApplication(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesisanalyticsv2_application.getKinesisAnalyticsV2Application", "api_error", err)
		return nil, err
	}

	return data.ApplicationDetail, nil
}

func getKinesisAnalyticsV2ApplicationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := applicationArn(h.Item)
	// Create Session
	svc, err := KinesisAnalyticsV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesisanalyticsv2_application.getKinesisAnalyticsV2ApplicationTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &kinesisanalyticsv2.ListTagsForResourceInput{
		ResourceARN: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesisanalyticsv2_application.getKinesisAnalyticsV2ApplicationTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func kinesisAnalyticsV2ApplicationTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

	if tagList == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func applicationArn(item interface{}) *string {
	switch item := item.(type) {
	case *types.ApplicationDetail:
		return item.ApplicationARN
	case types.ApplicationSummary:
		return item.ApplicationARN
	}
	return nil
}
