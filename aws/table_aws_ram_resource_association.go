package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ram"
	"github.com/aws/aws-sdk-go-v2/service/ram/types"

	ramv1 "github.com/aws/aws-sdk-go/service/ram"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRAMResourceAssociation(_ context.Context) *plugin.Table {
	associationType := "RESOURCE"
	return &plugin.Table{
		Name:        "aws_ram_resource_association",
		Description: "AWS RAM Resource Association",
		List: &plugin.ListConfig{
			Hydrate: listResourceShareAssociations(associationType),
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ramv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "resource_share_name",
				Description: "The name of the resource share.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_share_arn",
				Description: "The Amazon Resoure Name (ARN) of the resource share.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_entity",
				Description: "The Amazon Resoure Name (ARN) of the associated resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_type",
				Description: "The type of entity included in this association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time when the association was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "external",
				Description: "Indicates whether the principal belongs to the same organization in Organizations as the Amazon Web Services account that owns the resource share.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_updated_time",
				Description: "The date and time when the association was last updated..",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status_message",
				Description: "A message about the status of the association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_share_permission",
				Description: "Information about an RAM permission that is associated with a resource share and any of its resources of a specified type.",
				Hydrate:     getResourceSharePermissions,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceShareName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceShareArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listResourceShareAssociations(associationType string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		// Create session
		svc, err := RAMClient(ctx, d)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ram_resource_association.listResourceShareAssociations", "connection_error", err)
			return nil, err
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

		input := &ram.GetResourceShareAssociationsInput{
			AssociationType: types.ResourceShareAssociationType(associationType),
			MaxResults:      aws.Int32(maxLimit),
		}

		paginator := ram.NewGetResourceShareAssociationsPaginator(svc, input, func(o *ram.GetResourceShareAssociationsPaginatorOptions) {
			o.Limit = maxLimit
			o.StopOnDuplicateToken = true
		})

		// List call
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("aws_ram_resource_association.listResourceShareAssociations", "api_error", err)
				return nil, err
			}

			for _, items := range output.ResourceShareAssociations {
				d.StreamListItem(ctx, items)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
		return nil, err
	}
}

//// HYDRATE FUNCTIONS

func getResourceSharePermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	arn := *h.Item.(types.ResourceShareAssociation).ResourceShareArn

	// Create Session
	svc, err := RAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ram_resource_association.getResourceSharePermissions", "connection_error", err)
		return nil, err
	}

	// Build the params
	input := &ram.ListResourceSharePermissionsInput{
		ResourceShareArn: &arn,
	}
	summaries := []types.ResourceSharePermissionSummary{}

	paginator := ram.NewListResourceSharePermissionsPaginator(svc, input, func(o *ram.ListResourceSharePermissionsPaginatorOptions) {
		o.Limit = int32(100)
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ram_resource_association.getResourceSharePermissions", "api_error", err)
			return nil, err
		}

		summaries = append(summaries, output.Permissions...)
	}

	return summaries, nil
}
