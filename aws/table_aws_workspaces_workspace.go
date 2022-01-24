package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/workspaces"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_workspaces_workspace",
		Description: "AWS Workspaces",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("workspace_id"),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException"}),
			Hydrate:           getWorkspace,
		},
		List: &plugin.ListConfig{
			Hydrate: listWorkspaces,
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

func listWorkspaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listWorkspaces")

	// Create Session
	svc, err := WorkspacesService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeWorkspacesPages(
		&workspaces.DescribeWorkspacesInput{},
		func(page *workspaces.DescribeWorkspacesOutput, isLast bool) bool {
			for _, Workspace := range page.Workspaces {
				d.StreamListItem(ctx, Workspace)
			}
			return !isLast
		},
	)
	if err != nil {
		// AWS workspaces is not available in every region yet. This section of code handles the errors that we get when the API call tries to use unsupported regions endpoint (it throws "no such host" error message)
		if strings.Contains(err.Error(), "no such host") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listWorkspaces", "list", err)
		return nil, err
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getWorkspace")

	WorkspaceId := d.KeyColumnQuals["workspace_id"].GetStringValue()

	// check if workspace id is empty
	if WorkspaceId == "" {
		return nil, nil
	}

	// Create service
	svc, err := WorkspacesService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &workspaces.DescribeWorkspacesInput{
		WorkspaceIds: aws.StringSlice([]string{WorkspaceId}),
	}

	// Get call
	data, err := svc.DescribeWorkspaces(params)
	if err != nil {
		// AWS workspaces is not available in every region yet. This section of code handles the errors that we get when the API call tries to use unsupported regions endpoint (it throws "no such host" error message)
		if strings.Contains(err.Error(), "no such host") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("DescribeWorkspaces", "ERROR", err)
		return nil, err
	}

	if len(data.Workspaces) > 0 {
		return data.Workspaces[0], nil
	}

	return nil, nil
}

func listWorkspacesTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listWorkspacesTags")
	workspaceId := h.Item.(*workspaces.Workspace).WorkspaceId

	// Create Session
	svc, err := WorkspacesService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &workspaces.DescribeTagsInput{
		ResourceId: workspaceId,
	}

	tags, err := svc.DescribeTags(params)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspacesTags", "error", err)
		return nil, err
	}

	return tags, nil
}

// https://docs.aws.amazon.com/workspaces/latest/adminguide/workspaces-access-control.html
func getWorkspaceArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getWorkspaceArn")
	region := d.KeyColumnQualString(matrixKeyRegion)
	workspaceId := h.Item.(*workspaces.Workspace).WorkspaceId

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":workspaces:" + region + ":" + commonColumnData.AccountId + ":workspace/" + *workspaceId

	return arn, nil
}

//// TRANSFORM FUNCTION

// Transform function for workspaces resources tags
func workspaceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*workspaces.DescribeTagsOutput)

	if tags.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range tags.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}

	return nil, nil
}
