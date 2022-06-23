package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAmplifyApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_amplify_app",
		Description: "AWS Amplify app",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("app_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException"}),
			},
			Hydrate: getApp,
		},
		List: &plugin.ListConfig{
			Hydrate: listApps,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "workspace_id",
				Description: "The id of the WorkSpace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the WorkSpace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ComputerName"),
			},
			{
				Name:        "arn",
				Description: "The arn of the WorkSpace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWorkspaceArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "bundle_id",
				Description: "The identifier of the bundle used to create the WorkSpace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "directory_id",
				Description: "The identifier of the AWS Directory Service directory for the WorkSpace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The operational state of the WorkSpace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "error_code",
				Description: "The error code that is returned if the WorkSpace cannot be created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "error_message",
				Description: "The text of the error message that is returned if the WorkSpace cannot be created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "The IP address of the WorkSpace.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "root_volume_encryption_enabled",
				Description: "Indicates whether the data stored on the root volume is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "subnet_id",
				Description: "The identifier of the subnet for the WorkSpace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_name",
				Description: "The user for the WorkSpace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_volume_encryption_enabled",
				Description: "Indicates whether the data stored on the user volume is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "volume_encryption_key",
				Description: "The symmetric AWS KMS customer master key (CMK) used to encrypt data stored on your WorkSpace. Amazon WorkSpaces does not support asymmetric CMKs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "modification_states",
				Description: "The modification states of the WorkSpace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "workspace_properties",
				Description: "The properties of the WorkSpace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the WorkSpace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listWorkspacesTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ComputerName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listWorkspacesTags,
				Transform:   transform.From(workspaceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWorkspaceArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listApps")
	region := d.KeyColumnQualString(matrixKeyRegion)

	// AWS Amplify is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// Post "https://amplify.ap-southeast-3.amazonaws.com/": dial tcp: lookup amplify.ap-southeast-3.amazonaws.com: no such host
	serviceId := amplify.EndpointsID
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	// Create Session
	svc, err := AmplifyService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &amplify.ListAppsInput{
		MaxResults: aws.Int64(100),
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

	pagesLeft := true

	for pagesLeft {
		result, err := svc.ListApps(input)
		if err != nil {
			return nil, err
		}

		for _, app := range result.Apps {
			d.StreamListItem(ctx, app)

			// Context can be cancelled due to manual cancellation or the limit has been hit
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

func getApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getApp")

	region := d.KeyColumnQualString(matrixKeyRegion)

	// AWS Amplify is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// Post "https://amplify.ap-southeast-3.amazonaws.com/": dial tcp: lookup amplify.ap-southeast-3.amazonaws.com: no such host
	serviceId := amplify.EndpointsID
	validRegions := SupportedRegionsForService(ctx, d, serviceId)
	if !helpers.StringSliceContains(validRegions, region) {
		return nil, nil
	}

	appId := d.KeyColumnQuals["app_id"].GetStringValue()

	// check if workspace id is empty
	if appId == "" {
		return nil, nil
	}

	// Create service
	svc, err := AmplifyService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &amplify.GetAppInput{
		AppId: aws.String(appId),
	}

	// Get call
	data, err := svc.GetApp(params)
	if err != nil {
		plugin.Logger(ctx).Error("DescribeWorkspaces", "ERROR", err)
		return nil, err
	}

	if data.App == nil {
		err = errors.New("Expected valid App object but none was returned from Amplify GetApp call")
		return nil, err
	}

	return data.App, nil
}

//// TRANSFORM FUNCTION
