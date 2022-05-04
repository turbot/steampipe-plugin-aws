package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_app",
		Description: "AWS Sagemaker App",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AnyColumn([]string{"name", "app_type", "domain_id", "user_profile_name"}),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
			Hydrate:           getAwsSageMakerApp,
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"domain_id", "user_profile_name"}),
			Hydrate:    listAwsSageMakerApps,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The app name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppName"),
			},
			{
				Name:        "app_type",
				Description: "The type of app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The app's Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerApp,
				Transform:   transform.FromField("AppArn"),
			},
			{
				Name:        "creation_time",
				Description: "A timestamp that indicates when the app was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "domain_id",
				Description: "The domain ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "failure_reason",
				Description: "The failure reason.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerApp,
			},
			{
				Name:        "last_health_check_timestamp",
				Description: "The timestamp of the last health check.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSageMakerApp,
			},
			{
				Name:        "last_user_activity_timestamp",
				Description: "The timestamp of the last health check.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSageMakerApp,
			},
			{
				Name:        "status",
				Description: "The app's status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the app.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerAppTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "user_profile_name",
				Description: "The user profile name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_spec",
				Description: "The instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerApp,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerAppTags,
				Transform:   transform.FromValue().Transform(sageMakerTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSageMakerApps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsSageMakerApps")

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &sagemaker.ListAppsInput{
		MaxResults: aws.Int64(100),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["domain_id"] != nil {
		if equalQuals["domain_id"].GetStringValue() != "" {
			input.DomainIdEquals = aws.String(equalQuals["domain_id"].GetStringValue())
		}
	} else if equalQuals["user_profile_name"] != nil {
		if equalQuals["user_profile_name"].GetStringValue() != "" {
			input.UserProfileNameEquals = aws.String(equalQuals["user_profile_name"].GetStringValue())
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListAppsPages(
		input,
		func(page *sagemaker.ListAppsOutput, isLast bool) bool {
			for _, app := range page.Apps {
				d.StreamListItem(ctx, app)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
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

func getAwsSageMakerApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = sageMakerAppId(h.Item)
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.DescribeAppInput{
		AppName: aws.String(id),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["app_type"] != nil {
		if equalQuals["app_type"].GetStringValue() != "" {
			params.AppType = aws.String(equalQuals["app_type"].GetStringValue())
		}
	} else if equalQuals["domain_id"] != nil {
		if equalQuals["domain_id"].GetStringValue() != "" {
			params.DomainId = aws.String(equalQuals["domain_id"].GetStringValue())
		}
	} else if equalQuals["user_profile_name"] != nil {
		if equalQuals["user_profile_name"].GetStringValue() != "" {
			params.DomainId = aws.String(equalQuals["user_profile_name"].GetStringValue())
		}
	}

	// Get call
	data, err := svc.DescribeApp(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsSageMakerApp", "ERROR", err)
		return nil, err
	}
	return data, nil
}

func listAwsSageMakerAppTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsSageMakerAppTags")

	var appArn string
	if h.Item != nil {
		appArn = sageMakerAppArn(h.Item)
	}

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(appArn),
	}

	pagesLeft := true
	tags := []*sagemaker.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListTags(params)
		if err != nil {
			plugin.Logger(ctx).Error("listAwsSageMakerAppTags", "ListTags_error", err)
			return nil, err
		}
		tags = append(tags, keyTags.Tags...)

		if keyTags.NextToken != nil {
			params.NextToken = keyTags.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}

//// TRANSFORM FUNCTION

func sageMakerAppId(item interface{}) string {
	switch item := item.(type) {
	case *sagemaker.AppDetails:
		return *item.AppName
	case *sagemaker.DescribeAppOutput:
		return *item.AppName
	}
	return ""
}

func sageMakerAppArn(item interface{}) string {
	switch item := item.(type) {
	// case *sagemaker.AppDetails:
	// 	return *item.AppArn
	case *sagemaker.DescribeAppOutput:
		return *item.AppArn
	}
	return ""
}
