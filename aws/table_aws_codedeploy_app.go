package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy/types"

	codedeployEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeDeployApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codedeploy_app",
		Description: "AWS CodeDeploy Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("application_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ApplicationDoesNotExistException"}),
			},
			Hydrate: getCodeDeployApplication,
			Tags:    map[string]string{"service": "codedeploy", "action": "GetApplication"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeDeployApplications,
			Tags:    map[string]string{"service": "codedeploy", "action": "ListApplications"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCodeDeployApplicationTags,
				Tags: map[string]string{"service": "codedeploy", "action": "ListTagsForResource"},
			},
			{
				Func: getCodeDeployApplication,
				Tags: map[string]string{"service": "codedeploy", "action": "GetApplication"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codedeployEndpoint.CODEDEPLOYServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "application_id",
				Description: "The application ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployApplication,
			},
			{
				Name:        "application_name",
				Description: "The application name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployApplicationArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "compute_platform",
				Description: "The destination platform type for deployment of the application (Lambda or Server).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployApplication,
			},
			{
				Name:        "create_time",
				Description: "The time at which the application was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeDeployApplication,
			},
			{
				Name:        "github_account_name",
				Description: "The name for a connection to a GitHub account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployApplication,
				Transform:   transform.FromField("GitHubAccountName"),
			},
			{
				Name:        "linked_to_github",
				Description: "True if the user has authenticated with GitHub for the specified application. Otherwise, false.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCodeDeployApplication,
				Transform:   transform.FromField("LinkedToGitHub"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tag key and value pairs associated with this application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployApplicationTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployApplicationTags,
				Transform:   transform.From(codeDeployApplicationTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployApplicationArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodeDeployApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create session
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_app.listCodeDeployApplications", "service_creation_error", err)
		return nil, err
	}

	input := codedeploy.ListApplicationsInput{}

	paginator := codedeploy.NewListApplicationsPaginator(svc, &input, func(o *codedeploy.ListApplicationsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codedeploy_app.listCodeDeployApplications", "api_error", err)
			return nil, err
		}

		for _, application := range output.Applications {
			item := &types.ApplicationInfo{
				ApplicationName: aws.String(application),
			}
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeDeployApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name string
	if h.Item != nil {
		name = *h.Item.(*types.ApplicationInfo).ApplicationName
	} else {
		name = d.EqualsQuals["application_name"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &codedeploy.GetApplicationInput{
		ApplicationName: &name,
	}

	// Create session
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_app.getCodeDeployApplication", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.GetApplication(ctx, params)
	if err != nil {
		logger.Error("aws_codedeploy_app.getCodeDeployApplication", "api_error", err)
		return nil, err
	}
	return data.Application, nil
}

func getCodeDeployApplicationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name string
	if h.Item != nil {
		name = *h.Item.(*types.ApplicationInfo).ApplicationName
	} else {
		name = d.EqualsQuals["application_name"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &codedeploy.ListTagsForResourceInput{
		ResourceArn: aws.String(codeDeployApplicationArn(ctx, d, h)),
	}

	// Create session
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_app.getCodeDeployApplicationTags", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		logger.Error("aws_codedeploy_app.getCodeDeployApplicationTags", "api_error", err)
		return nil, err
	}
	return data, nil
}

func getCodeDeployApplicationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return codeDeployApplicationArn(ctx, d, h), nil
}

func codeDeployApplicationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) string {
	name := *h.Item.(*types.ApplicationInfo).ApplicationName
	region := d.EqualsQualString(matrixKeyRegion)
	logger := plugin.Logger(ctx)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		logger.Error("aws_codedeploy_app.getCodeDeployApplicationArn", "caching_error", err)
		return ""
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	//arn:aws:codedeploy:region:account-id:application:application-name
	tableArn := "arn:" + commonColumnData.Partition + ":codedeploy:" + region + ":" + commonColumnData.AccountId + ":application:" + name
	return tableArn
}

//// TRANSFORM FUNCTIONS

func codeDeployApplicationTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*codedeploy.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
