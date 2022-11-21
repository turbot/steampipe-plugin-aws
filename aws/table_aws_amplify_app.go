package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amplify"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAmplifyApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_amplify_app",
		Description: "AWS Amplify App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("app_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "NotFoundException"}),
			},
			Hydrate: getAmplifyApp,
		},
		List: &plugin.ListConfig{
			Hydrate: listAmplifyApps,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "app_id",
				Description: "The unique ID of the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the Amplify app.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppArn"),
			},
			{
				Name:        "name",
				Description: "The name for the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description for the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "Creates a date and time for the Amplify app.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "Updates the date and time for the Amplify app.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "basic_auth_credentials",
				Description: "The basic authorization credentials for branches for the Amplify app. You must base64-encode the authorization credentials and provide them in the format user:password.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "custom_headers",
				Description: "Describes the custom HTTP headers for the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_domain",
				Description: "The default domain for the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_auto_branch_creation",
				Description: "Enables automated branch creation for the Amplify app.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "enable_basic_auth",
				Description: "Enables basic authorization for the Amplify app's branches.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "enable_branch_auto_build",
				Description: "Enables the auto-building of branches for the Amplify app.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "enable_branch_auto_deletion",
				Description: "Automatically disconnect a branch in the Amplify Console when you delete a branch from your Git repository.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "iam_service_role_arn",
				Description: "The AWS Identity and Access Management (IAM) service role for the Amazon Resource Name (ARN) of the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform",
				Description: "The platform for the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository",
				Description: "The Git repository for the Amplify app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_clone_method",
				Description: "The Amplify service uses this parameter to specify the authentication protocol to use to access the Git repository for an Amplify app. Amplify specifies TOKEN for a GitHub repository, SIGV4 for an AWS CodeCommit repository, and SSH for GitLab and Bitbucket repositories.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_branch_creation_config",
				Description: "Describes the automated branch creation configuration for the Amplify app.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "auto_branch_creation_patterns",
				Description: "Describes the automated branch creation glob patterns for the Amplify app.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "build_spec",
				Description: "Describes the content of the build specification (build spec) for the Amplify app.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BuildSpec").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "custom_rules",
				Description: "Describes the custom redirect and rewrite rules for the Amplify app.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "environment_variables",
				Description: "The environment variables for the Amplify app.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "production_branch",
				Description: "Describes the information about a production branch of the Amplify app.",
				Type:        proto.ColumnType_JSON,
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
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AppArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAmplifyApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := AmplifyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_amplify_app.listAmplifyApps", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &amplify.ListAppsInput{
		MaxResults: int32(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < input.MaxResults {
			if limit < 20 {
				input.MaxResults = int32(20)
			} else {
				input.MaxResults = int32(limit)
			}
		}
	}

	// API doesn't support aws-sdk-go-v2 paginator as of date.
	pagesLeft := true

	for pagesLeft {
		result, err := svc.ListApps(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_amplify_app.listAmplifyApps", "api_error", err)
			return nil, err
		}

		for _, item := range result.Apps {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			pagesLeft = true
			input.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAmplifyApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	appId := d.KeyColumnQuals["app_id"].GetStringValue()
	if appId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := AmplifyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_amplify_app.getAmplifyApp", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &amplify.GetAppInput{
		AppId: aws.String(appId),
	}

	// Get call
	data, err := svc.GetApp(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_amplify_app.getAmplifyApp", "api_error", err)
		return nil, err
	}

	return data.App, nil
}
