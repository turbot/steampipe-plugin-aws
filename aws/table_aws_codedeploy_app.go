package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsCodeDeployApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codedeploy_app",
		Description: "AWS Code Deploy Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("application_name"),
			Hydrate:    getCodeDeployApplication,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeDeployApplications,
		},
		GetMatrixItem: BuildRegionList,
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

			// Steampipe standard columns
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
	svc, err := CodeDeployService(ctx, d)
	if err != nil {
		logger.Error("applicatioaws_codedeploy_application.listCodeDeployApplications", "service_creation_error", err)
		return nil, err
	}

	input := &codedeploy.ListApplicationsInput{}

	// List call
	err = svc.ListApplicationsPages(
		input,
		func(page *codedeploy.ListApplicationsOutput, isLast bool) bool {
			for _, application := range page.Applications {

				item := &codedeploy.ApplicationInfo{
					ApplicationName: application,
				}
				d.StreamListItem(ctx, item)

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

func getCodeDeployApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name string
	if h.Item != nil {
		name = *h.Item.(*codedeploy.ApplicationInfo).ApplicationName
	} else {
		name = d.KeyColumnQuals["application_name"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &codedeploy.GetApplicationInput{
		ApplicationName: &name,
	}

	// Create session
	svc, err := CodeDeployService(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_application.getCodeDeployApplication", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.GetApplication(params)
	if err != nil {
		logger.Error("aws_codedeploy_application.getCodeDeployApplication", "api_error", err)
		return nil, err
	}
	return data.Application, nil
}

func getCodeDeployApplicationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := *h.Item.(*codedeploy.ApplicationInfo).ApplicationName

	logger := plugin.Logger(ctx)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		logger.Error("applicatioaws_codedeploy_application.getCodeDeployApplicationArn", "caching_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	//arn:aws:codedeploy:region:account-id:application:application-name
	tableArn := "arn:" + commonColumnData.Partition + ":codedeploy:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":application/" + name
	return tableArn, nil
}
