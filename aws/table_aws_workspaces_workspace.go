package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	"github.com/aws/aws-sdk-go-v2/service/workspaces/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_workspaces_workspace",
		Description: "AWS Workspaces",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("workspace_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException"}),
			},
			Hydrate: getWorkspace,
		},
		List: &plugin.ListConfig{
			Hydrate: listWorkspaces,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bundle_id", Require: plugin.Optional},
				{Name: "directory_id", Require: plugin.Optional},
				{Name: "user_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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

func listWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := WorkspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_workspace.listWorkspaces", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(25)
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

	input := &workspaces.DescribeWorkspacesInput{
		Limit: aws.Int32(maxLimit),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["bundle_id"] != nil {
		if equalQuals["bundle_id"].GetStringValue() != "" {
			input.BundleId = aws.String(equalQuals["bundle_id"].GetStringValue())
		}
	}
	if equalQuals["directory_id"] != nil {
		if equalQuals["directory_id"].GetStringValue() != "" {
			input.DirectoryId = aws.String(equalQuals["directory_id"].GetStringValue())
		}
	}
	if equalQuals["user_name"] != nil {
		if equalQuals["user_name"].GetStringValue() != "" {
			input.UserName = aws.String(equalQuals["user_name"].GetStringValue())
		}
	}

	paginator := workspaces.NewDescribeWorkspacesPaginator(svc, input, func(o *workspaces.DescribeWorkspacesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_workspaces_workspace.listWorkspaces", "api_error", err)
			return nil, err
		}

		for _, items := range output.Workspaces {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	WorkspaceId := d.KeyColumnQuals["workspace_id"].GetStringValue()

	// check if workspace id is empty
	if WorkspaceId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := WorkspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_workspace.getWorkspace", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &workspaces.DescribeWorkspacesInput{
		WorkspaceIds: []string{WorkspaceId},
	}

	// Get call
	data, err := svc.DescribeWorkspaces(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_workspace.getWorkspace", "api_error", err)
		return nil, err
	}

	if len(data.Workspaces) > 0 {
		return data.Workspaces[0], nil
	}

	return nil, nil
}

func listWorkspacesTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspaceId := h.Item.(types.Workspace).WorkspaceId

	// Create Session
	svc, err := WorkspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_workspace.listWorkspacesTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &workspaces.DescribeTagsInput{
		ResourceId: workspaceId,
	}

	tags, err := svc.DescribeTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_workspace.listWorkspacesTags", "api_error", err)
		return nil, err
	}

	return tags, nil
}

// https://docs.aws.amazon.com/workspaces/latest/adminguide/workspaces-access-control.html
func getWorkspaceArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getWorkspaceArn")
	region := d.KeyColumnQualString(matrixKeyRegion)
	workspaceId := h.Item.(types.Workspace).WorkspaceId

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
