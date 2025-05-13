package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_app",
		Description: "AWS Sagemaker App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "app_type", "domain_id", "user_profile_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "NotFoundException", "ResourceNotFound"}),
			},
			Hydrate: getSageMakerApp,
			Tags:    map[string]string{"service": "sagemaker", "action": "DescribeApp"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsSageMakerDomains,
			Hydrate:       listSageMakerApps,
			Tags:          map[string]string{"service": "sagemaker", "action": "ListApps"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_profile_name", Require: plugin.Optional},
				{Name: "domain_id", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSageMakerApp,
				Tags: map[string]string{"service": "sagemaker", "action": "DescribeApp"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_API_SAGEMAKER_SERVICE_ID),
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
				Name:        "space_name",
				Description: "The name of the space. If this value is not set, then UserProfileName must be set.",
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
	sageMakerDomain := h.Item.(types.DomainDetails)

	equalQuals := d.EqualsQuals

	if equalQuals["domain_id"].GetStringValue() != "" && equalQuals["domain_id"].GetStringValue() != *sageMakerDomain.DomainId {
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
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

	input := &sagemaker.ListAppsInput{
		DomainIdEquals: sageMakerDomain.DomainId,
		MaxResults:     aws.Int32(maxLimit),
	}
	if equalQuals["user_profile_name"] != nil {
		if equalQuals["user_profile_name"].GetStringValue() != "" {
			input.UserProfileNameEquals = aws.String(equalQuals["user_profile_name"].GetStringValue())
		}
	}

	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_app.listSageMakerApps", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	paginator := sagemaker.NewListAppsPaginator(svc, input, func(o *sagemaker.ListAppsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_app.listSageMakerApps", "api_error", err)
			return nil, err
		}

		for _, items := range output.Apps {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSageMakerApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var params sagemaker.DescribeAppInput

	// Build the params
	if h.Item != nil {
		params = sageMakerAppParams(h.Item)
	} else {
		equalQuals := d.EqualsQuals
		appName := aws.String(equalQuals["name"].GetStringValue())
		appType := aws.String(equalQuals["app_type"].GetStringValue())
		userProfileName := aws.String(equalQuals["user_profile_name"].GetStringValue())
		domainId := aws.String(equalQuals["domain_id"].GetStringValue())
		if len(*appName) == 0 || len(*appType) == 0 || len(*userProfileName) == 0 || len(*domainId) == 0 {
			return nil, nil
		}

		params = sagemaker.DescribeAppInput{
			AppName:         appName,
			AppType:         types.AppType(*appType),
			UserProfileName: userProfileName,
			DomainId:        domainId,
		}
	}

	// Create service
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_app.getSageMakerApp", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Get call
	data, err := svc.DescribeApp(ctx, &params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_app.getSageMakerApp", "api_error", err)
		return nil, err
	}
	return data, nil
}

func sageMakerAppArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch item := h.Item.(type) {
	case types.AppDetails:

		c, err := getCommonColumns(ctx, d, h)
		if err != nil {
			return "", err
		}
		commonColumnData := c.(*awsCommonColumnData)
		return fmt.Sprintf("arn:%s:sagemaker:%s:%s:app/%s/%s/%s/%s", commonColumnData.Partition, commonColumnData.Region, commonColumnData.AccountId, *item.DomainId, *item.UserProfileName, strings.ToLower(string(item.AppType)), *item.AppName), nil
	case *sagemaker.DescribeAppOutput:
		return *item.AppArn, nil
	}
	return "", nil
}

func sageMakerAppParams(item interface{}) sagemaker.DescribeAppInput {
	switch item := item.(type) {
	case types.AppDetails:
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
