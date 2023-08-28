package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	"github.com/aws/aws-sdk-go-v2/service/workspaces/types"

	workspacesv1 "github.com/aws/aws-sdk-go/service/workspaces"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWorkspacesDirectory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_workspaces_directory",
		Description: "AWS Workspaces Directory",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("directory_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				 ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException"}),
			},
			Hydrate: getWorkspacesDirectory,
		},
		List: &plugin.ListConfig{
			Hydrate: listWorkspacesDirectories,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(workspacesv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "directory_id",
				Description: "The directory identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the directory.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DirectoryName"),
			},
			{
				Name:        "arn",
				Description: "The arn of the directory.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWorkspaceDirectoryArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "alias",
				Description: "The directory alias.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_based_auth_properties",
				Description: "The certificate-based authentication properties used to authenticate SAML 2.0 Identity Provider (IdP) user identities to Active Directory for WorkSpaces login.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "customer_user_name",
				Description: "The user name for the service account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "directory_type",
				Description: "The directory type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_ip_addresses",
				Description: "The IP addresses of the DNS servers for the directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_role_id",
				Description: "The identifier of the IAM role. This is the role that allows Amazon WorkSpaces to make calls to other services, such as Amazon EC2, on your behalf.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_group_ids",
				Description: "The identifiers of the IP access control groups associated with the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "registration_code",
				Description: "The registration code for the directory. This is the code that users enter in their Amazon WorkSpaces client application to connect to the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "saml_properties",
				Description: "Describes the enablement status, user access URL, and relay state parameter name that are used for configuring federation with an SAML 2.0 identity provider.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "selfservice_permissions",
				Description: "The default self-service permissions for WorkSpaces in the directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state",
				Description: "The state of the directory's registration with Amazon WorkSpaces.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_ids",
				Description: "The identifiers of the subnets used with the directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tenancy",
				Description: "Specifies whether the directory is dedicated or shared.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workspace_access_properties",
				Description: "The devices and operating systems that users can use to access WorkSpaces.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "workspace_creation_properties",
				Description: "The default creation properties for all WorkSpaces in the directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "workspace_security_group_id",
				Description: "The identifier of the security group that is assigned to new WorkSpaces.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceSecurityGroupId"),
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the directory.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listWorkspacesDirectoriesTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DirectoryName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listWorkspacesDirectoriesTags,
				Transform:   transform.From(workspaceDirectoryTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWorkspaceDirectoryArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listWorkspacesDirectories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := WorkspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_directory.listWorkspacesDirectories", "connection_error", err)
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

	input := &workspaces.DescribeWorkspaceDirectoriesInput{
		Limit: aws.Int32(maxLimit),
	}

	paginator := workspaces.NewDescribeWorkspaceDirectoriesPaginator(svc, input, func(o *workspaces.DescribeWorkspaceDirectoriesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_workspaces_directory.listWorkspacesDirectories", "api_error", err)
			return nil, err
		}

		for _, items := range output.Directories {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWorkspacesDirectory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	DirectoryId := d.EqualsQualString("directory_id")

	// check if workspace id is empty
	if DirectoryId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := WorkspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_directory.getWorkspacesDirectory", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &workspaces.DescribeWorkspaceDirectoriesInput{
		DirectoryIds: []string{DirectoryId},
	}

	// Get call
	data, err := svc.DescribeWorkspaceDirectories(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_directory.getWorkspacesDirectory", "api_error", err)
		return nil, err
	}

	if len(data.Directories) > 0 {
		return data.Directories[0], nil
	}

	return nil, nil
}

func listWorkspacesDirectoriesTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	DirectoryId := h.Item.(types.WorkspaceDirectory).DirectoryId

	// Create Session
	svc, err := WorkspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_directory.listWorkspacesDirectoriesTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &workspaces.DescribeTagsInput{
		ResourceId: DirectoryId,
	}

	tags, err := svc.DescribeTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_workspaces_directory.listWorkspacesDirectoriesTags", "api_error", err)
		return nil, err
	}

	return tags, nil
}

// https://docs.aws.amazon.com/workspaces/latest/adminguide/workspaces-access-control.html
func getWorkspaceDirectoryArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDirectoryArn")
	region := d.EqualsQualString(matrixKeyRegion)
	DirectoryId := h.Item.(types.WorkspaceDirectory).DirectoryId

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":workspaces:" + region + ":" + commonColumnData.AccountId + ":directory/" + *DirectoryId

	return arn, nil
}

//// TRANSFORM FUNCTION

// Transform function for workspaces resources tags
func workspaceDirectoryTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
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
