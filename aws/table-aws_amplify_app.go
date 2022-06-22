package aws

import (
	"context"

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
	// region := d.KeyColumnQualString(matrixKeyRegion)

	// // AWS Workspaces is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// // Post "https://workspaces.eu-north-1.amazonaws.com/": dial tcp: lookup workspaces.eu-north-1.amazonaws.com: no such host
	// serviceId := endpoints.WorkspacesServiceID
	// validRegions := SupportedRegionsForService(ctx, d, serviceId)
	// if !helpers.StringSliceContains(validRegions, region) {
	// 	return nil, nil
	// }

	// Create Session
	// svc, err := AmplifyService(ctx, d)
	_, err := AmplifyService(ctx, d)
	if err != nil {
		return nil, err
	}

	// input := &workspaces.DescribeWorkspacesInput{
	// 	Limit: aws.Int64(25),
	// }

	// // Reduce the basic request limit down if the user has only requested a small number of rows
	// limit := d.QueryContext.Limit
	// if d.QueryContext.Limit != nil {
	// 	if *limit < *input.Limit {
	// 		if *limit < 1 {
	// 			input.Limit = aws.Int64(1)
	// 		} else {
	// 			input.Limit = limit
	// 		}
	// 	}
	// }

	// equalQuals := d.KeyColumnQuals
	// if equalQuals["bundle_id"] != nil {
	// 	if equalQuals["bundle_id"].GetStringValue() != "" {
	// 		input.BundleId = aws.String(equalQuals["bundle_id"].GetStringValue())
	// 	}
	// }
	// if equalQuals["directory_id"] != nil {
	// 	if equalQuals["directory_id"].GetStringValue() != "" {
	// 		input.DirectoryId = aws.String(equalQuals["directory_id"].GetStringValue())
	// 	}
	// }
	// if equalQuals["user_name"] != nil {
	// 	if equalQuals["user_name"].GetStringValue() != "" {
	// 		input.UserName = aws.String(equalQuals["user_name"].GetStringValue())
	// 	}
	// }

	// // List call
	// err = svc.DescribeWorkspacesPages(
	// 	input,
	// 	func(page *workspaces.DescribeWorkspacesOutput, isLast bool) bool {
	// 		for _, Workspace := range page.Workspaces {
	// 			d.StreamListItem(ctx, Workspace)

	// 			// Context may get cancelled due to manual cancellation or if the limit has been reached
	// 			if d.QueryStatus.RowsRemaining(ctx) == 0 {
	// 				return false
	// 			}
	// 		}
	// 		return !isLast
	// 	},
	// )
	// if err != nil {
	// 	plugin.Logger(ctx).Error("listWorkspaces", "list", err)
	// 	return nil, err
	// }

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getApp")

	// region := d.KeyColumnQualString(matrixKeyRegion)

	// // AWS Workspaces is not supported in all regions. For unsupported regions the API throws an error, e.g.,
	// // Post "https://workspaces.eu-north-1.amazonaws.com/": dial tcp: lookup workspaces.eu-north-1.amazonaws.com: no such host
	// serviceId := endpoints.WorkspacesServiceID
	// validRegions := SupportedRegionsForService(ctx, d, serviceId)
	// if !helpers.StringSliceContains(validRegions, region) {
	// 	return nil, nil
	// }

	// WorkspaceId := d.KeyColumnQuals["workspace_id"].GetStringValue()

	// // check if workspace id is empty
	// if WorkspaceId == "" {
	// 	return nil, nil
	// }

	// Create service
	// svc, err := AmplifyService(ctx, d)
	_, err := AmplifyService(ctx, d)
	if err != nil {
		return nil, err
	}

	// // Build the params
	// params := &workspaces.DescribeWorkspacesInput{
	// 	WorkspaceIds: aws.StringSlice([]string{WorkspaceId}),
	// }

	// // Get call
	// data, err := svc.DescribeWorkspaces(params)
	// if err != nil {
	// 	plugin.Logger(ctx).Error("DescribeWorkspaces", "ERROR", err)
	// 	return nil, err
	// }

	// if len(data.Workspaces) > 0 {
	// 	return data.Workspaces[0], nil
	// }

	return nil, nil
}

//// TRANSFORM FUNCTION
