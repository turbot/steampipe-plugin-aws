package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsKinesisAnalyticsV2Application(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesisanalyticsv2_application",
		Description: "AWS Kinesis Analytics V2 Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("application_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getKinesisAnalyticsV2Application,
		},
		List: &plugin.ListConfig{
			Hydrate: listKinesisAnalyticsV2Applications,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := KinesisAnalyticsV2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	pagesLeft := true
	params := &kinesisanalyticsv2.ListApplicationsInput{
		Limit: aws.Int64(50),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.Limit {
			if *limit < 1 {
				params.Limit = aws.Int64(1)
			} else {
				params.Limit = limit
			}
		}
	}

	for pagesLeft {
		result, err := svc.ListApplications(params)
		if err != nil {
			return nil, err
		}

		for _, application := range result.ApplicationSummaries {
			d.StreamListItem(ctx, application)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
		plugin.Logger(ctx).Error("listKinesisAnalyticsV2Applications", "ListApplications_error", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKinesisAnalyticsV2Application(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getKinesisAnalyticsV2Application")

	var applicationName string
	if h.Item != nil {
		i := h.Item.(*kinesisanalyticsv2.ApplicationSummary)
		applicationName = *i.ApplicationName
	} else {
		applicationName = d.KeyColumnQuals["application_name"].GetStringValue()
	}

	// Create Session
	svc, err := KinesisAnalyticsV2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesisanalyticsv2.DescribeApplicationInput{
		ApplicationName: &applicationName,
	}

	// Get call
	data, err := svc.DescribeApplication(params)
	if err != nil {
		logger.Debug("getKinesisAnalyticsV2Application", "DescribeApplication_error", err)
		return nil, err
	}

	return data.ApplicationDetail, nil
}

func getKinesisAnalyticsV2ApplicationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getKinesisAnalyticsV2ApplicationTags")

	arn := applicationArn(h.Item)

	// Create Session
	svc, err := KinesisAnalyticsV2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesisanalyticsv2.ListTagsForResourceInput{
		ResourceARN: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getKinesisAnalyticsV2ApplicationTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func kinesisAnalyticsV2ApplicationTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("kinesisAnalyticsV2ApplicationTagListToTurbotTags")
	tagList := d.Value.([]*kinesisanalyticsv2.Tag)

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
	case *kinesisanalyticsv2.ApplicationDetail:
		return item.ApplicationARN
	case *kinesisanalyticsv2.ApplicationSummary:
		return item.ApplicationARN
	}
	return nil
}
