package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_app",
		Description: "AWS Sagemaker App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "app_type", "domain_id", "user_profile_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "NotFoundException", "ResourceNotFound"}),
			},
			Hydrate: getSageMakerApp,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsSageMakerDomains,
			Hydrate:       listSageMakerApps,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_profile_name", Require: plugin.Optional},
				{Name: "domain_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Description: "The Amazon Resource Name (ARN) of the app.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSageMakerApp,
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
				Hydrate:     getSageMakerApp,
			},
			{
				Name:        "last_health_check_timestamp",
				Description: "The timestamp of the last health check.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getSageMakerApp,
			},
			{
				Name:        "last_user_activity_timestamp",
				Description: "The timestamp of the last user activity.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getSageMakerApp,
			},
			{
				Name:        "status",
				Description: "The status of the app.",
				Type:        proto.ColumnType_STRING,
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
				Hydrate:     getSageMakerApp,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     sageMakerAppArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSageMakerApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sageMakerDomain := h.Item.(*sagemaker.DomainDetails)
	plugin.Logger(ctx).Trace("listSageMakerApps")

	equalQuals := d.KeyColumnQuals

	if equalQuals["domain_id"].GetStringValue() != "" && equalQuals["domain_id"].GetStringValue() != *sageMakerDomain.DomainId {
		return nil, nil
	}
	input := &sagemaker.ListAppsInput{
		DomainIdEquals: sageMakerDomain.DomainId,
		MaxResults:     aws.Int64(100),
	}
	if equalQuals["user_profile_name"] != nil {
		if equalQuals["user_profile_name"].GetStringValue() != "" {
			input.UserProfileNameEquals = aws.String(equalQuals["user_profile_name"].GetStringValue())
		}
	}

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
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

func getSageMakerApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var params sagemaker.DescribeAppInput

	// Build the params
	if h.Item != nil {
		params = sageMakerAppParams(h.Item)
	} else {
		equalQuals := d.KeyColumnQuals
		appName := aws.String(equalQuals["name"].GetStringValue())
		appType := aws.String(equalQuals["app_type"].GetStringValue())
		userProfileName := aws.String(equalQuals["user_profile_name"].GetStringValue())
		domainId := aws.String(equalQuals["domain_id"].GetStringValue())
		if len(*appName) == 0 || len(*appType) == 0 || len(*userProfileName) == 0 || len(*domainId) == 0 {
			return nil, nil
		}

		params = sagemaker.DescribeAppInput{
			AppName:         appName,
			AppType:         appType,
			UserProfileName: userProfileName,
			DomainId:        domainId,
		}
	}

	// Create service
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get call
	data, err := svc.DescribeApp(&params)
	if err != nil {
		plugin.Logger(ctx).Error("getSageMakerApp", "ERROR", err)
		return nil, err
	}
	return data, nil
}

func sageMakerAppArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch item := h.Item.(type) {
	case *sagemaker.AppDetails:
		getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
		c, err := getCommonColumnsCached(ctx, d, h)
		if err != nil {
			return "", err
		}
		commonColumnData := c.(*awsCommonColumnData)
		return "arn:" + commonColumnData.Partition +
			":sagemaker:" + commonColumnData.Region +
			":" + commonColumnData.AccountId +
			":app/" + *item.DomainId +
			"/" + *item.UserProfileName +
			"/" + strings.ToLower(*item.AppType) +
			"/" + *item.AppName, nil
	case *sagemaker.DescribeAppOutput:
		return *item.AppArn, nil
	}
	return "", nil
}

func sageMakerAppParams(item interface{}) sagemaker.DescribeAppInput {
	switch item := item.(type) {
	case *sagemaker.AppDetails:
		return sagemaker.DescribeAppInput{
			AppName:         item.AppName,
			AppType:         item.AppType,
			UserProfileName: item.UserProfileName,
			DomainId:        item.DomainId,
		}
	case *sagemaker.DescribeAppOutput:
		return sagemaker.DescribeAppInput{
			AppName:         item.AppName,
			AppType:         item.AppType,
			UserProfileName: item.UserProfileName,
			DomainId:        item.DomainId,
		}
	}
	return sagemaker.DescribeAppInput{}
}
