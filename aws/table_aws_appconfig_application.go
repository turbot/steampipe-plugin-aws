package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/aws/aws-sdk-go-v2/service/appconfig/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAppConfigApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appconfig_application",
		Description: "AWS AppConfig Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAppConfigApplication,
		},
		List: &plugin.ListConfig{
			Hydrate: listAppConfigApplication,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The application ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The application name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that identifies the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAppConfigApplicationArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "description",
				Description: "The description of the application.",
				Type:        proto.ColumnType_STRING,
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
				Hydrate:     getAppConfigTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppConfigApplicationArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAppConfigApplication(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AppConfigClient(ctx, d)
	if err != nil {
		logger.Error("aws_appconfig_application.listAppConfigApplication", "service_creation_error", err)
		return nil, err
	}

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

	params := &appconfig.ListApplicationsInput{
		MaxResults: *aws.Int32(maxLimit),
	}

	paginator := appconfig.NewListApplicationsPaginator(svc, params, func(o *appconfig.ListApplicationsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_appconfig_application.listAppConfigApplication", "api_error", err)
			return nil, err
		}
		for _, application := range output.Items {
			d.StreamListItem(ctx, application)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAppConfigApplication(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AppConfigClient(ctx, d)
	if err != nil {
		logger.Error("aws_appconfig_application.getAppConfigApplication", "service_creation_error", err)
		return nil, err
	}

	id := d.KeyColumnQuals["id"].GetStringValue()
	params := &appconfig.GetApplicationInput{
		ApplicationId: aws.String(id),
	}

	application, err := svc.GetApplication(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appconfig_application.getAppConfigApplication", "api_error", err)
		return nil, err
	}

	if application != nil {
		api := types.Application{
			Name:        application.Name,
			Id:          application.Id,
			Description: application.Description,
		}
		return api, nil
	}

	return nil, nil
}

func getAppConfigTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AppConfigClient(ctx, d)
	if err != nil {
		logger.Error("aws_appconfig_application.getAppConfigTags", "service_creation_error", err)
		return nil, err
	}

	arn := getArnFormat(ctx, d, h)
	params := &appconfig.ListTagsForResourceInput{
		ResourceArn: &arn,
	}

	tags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appconfig_application.getAppConfigTags", "api_error", err)
		return nil, err
	}

	return tags.Tags, nil
}

func getAppConfigApplicationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getArnFormat(ctx, d, h), nil
}

func getArnFormat(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) string {

	region := d.KeyColumnQualString(matrixKeyRegion)
	id := h.Item.(types.Application).Id

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appconfig_application.getArnFormat", "cache_error", err)
		return ""
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// arn:${Partition}:appconfig:${Region}:${Account}:application/${ApplicationId}
	arn := "arn:" + commonColumnData.Partition + ":appconfig:" + region + ":" + commonColumnData.AccountId + ":application/" + *id
	return arn
}
