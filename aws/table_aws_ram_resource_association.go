package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ram"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsRAMResourceAssociation(_ context.Context) *plugin.Table {
	associationType := "RESOURCE"
	return &plugin.Table{
		Name:        "aws_ram_resource_association",
		Description: "AWS RAM Resource Association",
		List: &plugin.ListConfig{
			Hydrate: listResourceShareAssociations(associationType),
		},
		GetMatrixItemFunc: BuildRegionList,
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
		svc, err := RAMService(ctx, d)
		if err != nil {
			return nil, err
		}

		// List call
		input := &ram.GetResourceShareAssociationsInput{
			AssociationType: aws.String(associationType),
			MaxResults:      aws.Int64(100),
		}

		// Reduce the basic request limit down if the user has only requested a small number of rows
		limit := d.QueryContext.Limit
		if d.QueryContext.Limit != nil {
			if *limit <= 100 {
				if *limit < 1 {
					input.MaxResults = aws.Int64(1)
				} else {
					input.MaxResults = limit
				}
			}
		}

		// List call
		err = svc.GetResourceShareAssociationsPages(
			input,
			func(page *ram.GetResourceShareAssociationsOutput, isLast bool) bool {
				for _, association := range page.ResourceShareAssociations {
					d.StreamListItem(ctx, association)

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
}

//// HYDRATE FUNCTIONS

func getResourceSharePermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getResourceSharePermissions")

	arn := *h.Item.(*ram.ResourceShareAssociation).ResourceShareArn

	// Create Session
	svc, err := RAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	input := &ram.ListResourceSharePermissionsInput{
		ResourceShareArn: &arn,
	}
	summaries := []*ram.ResourceSharePermissionSummary{}

	// Get call
	err = svc.ListResourceSharePermissionsPages(
		input,
		func(page *ram.ListResourceSharePermissionsOutput, isLast bool) bool {
			summaries = append(summaries, page.Permissions...)

			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}
