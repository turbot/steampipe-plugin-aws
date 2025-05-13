package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEmrStudio(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_studio",
		Description: "AWS EMR Studio",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("studio_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidRequestException"}),
			},
			Hydrate: getEmrStudio,
			Tags:    map[string]string{"service": "elasticmapreduce", "action": "DescribeStudio"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEmrStudios,
			Tags:    map[string]string{"service": "elasticmapreduce", "action": "ListStudios"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ELASTICMAPREDUCE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the EMR Studio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "studio_id",
				Description: "The ID of the EMR Studio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "studio_arn",
				Description: "The Amazon Resource Name (ARN) of the EMR Studio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the EMR Studio.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StudioArn"),
			},
			{
				Name:        "auth_mode",
				Description: "Specifies whether the Studio authenticates users using IAM or IAM Identity Center.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_s3_location",
				Description: "The Amazon S3 location to back up EMR Studio Workspaces and notebook files.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The detailed description of the EMR Studio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_security_group_id",
				Description: "The ID of the Engine security group associated with the Studio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "idp_auth_url",
				Description: "The authentication endpoint of your identity provider (IdP).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "idp_relay_state_parameter_name",
				Description: "The name that your identity provider (IdP) uses for its RelayState parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_role",
				Description: "The IAM role that will be assumed by the Amazon EMR Studio.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEmrStudioServiceRole,
			},
			{
				Name:        "subnet_ids",
				Description: "A list of subnet IDs associated with the EMR Studio.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SubnetIds"),
				Hydrate:     getEmrStudioSubnetIds,
			},
			{
				Name:        "url",
				Description: "The unique access URL of the EMR Studio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_role",
				Description: "The IAM user role that will be assumed by users and groups logged in to a Studio.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEmrStudioUserRole,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC associated with the Studio.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEmrStudioVpcId,
			},
			{
				Name:        "workspace_security_group_id",
				Description: "The ID of the Workspace security group associated with the Studio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the EMR Studio was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StudioArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEmrStudios(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EMRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_studio.listEmrStudios", "connection_error", err)
		return nil, err
	}

	// List all items
	paginator := emr.NewListStudiosPaginator(svc, &emr.ListStudiosInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_emr_studio.listEmrStudios", "api_error", err)
			return nil, err
		}

		for _, item := range output.Studios {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEmrStudio(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := EMRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_studio.getEmrStudio", "client_error", err)
		return nil, err
	}

	// Get item ID from query data
	var id string
	if h.Item != nil {
		id = *h.Item.(types.StudioSummary).StudioId
	} else {
		id = d.EqualsQuals["studio_id"].GetStringValue()
		if id == "" {
			return nil, nil
		}
	}

	// Get item details
	params := &emr.DescribeStudioInput{
		StudioId: &id,
	}

	// Get call
	op, err := svc.DescribeStudio(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_studio.getEmrStudio", "api_error", err)
		return nil, err
	}

	return op.Studio, nil
}

func getEmrStudioSubnetIds(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch item := h.Item.(type) {
	case types.Studio:
		return item.SubnetIds, nil
	case types.StudioSummary:
		studio, err := getEmrStudio(ctx, d, h)
		if err != nil {
			return nil, err
		}
		return studio.(types.Studio).SubnetIds, nil
	default:
		return nil, nil
	}
}

func getEmrStudioVpcId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch item := h.Item.(type) {
	case types.Studio:
		return item.VpcId, nil
	case types.StudioSummary:
		studio, err := getEmrStudio(ctx, d, h)
		if err != nil {
			return nil, err
		}
		return studio.(types.Studio).VpcId, nil
	default:
		return nil, nil
	}
}

func getEmrStudioUserRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch item := h.Item.(type) {
	case types.Studio:
		return item.UserRole, nil
	case types.StudioSummary:
		studio, err := getEmrStudio(ctx, d, h)
		if err != nil {
			return nil, err
		}
		return studio.(types.Studio).UserRole, nil
	default:
		return nil, nil
	}
}

func getEmrStudioServiceRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch item := h.Item.(type) {
	case types.Studio:
		return item.ServiceRole, nil
	case types.StudioSummary:
		studio, err := getEmrStudio(ctx, d, h)
		if err != nil {
			return nil, err
		}
		return studio.(types.Studio).ServiceRole, nil
	default:
		return nil, nil
	}
}
