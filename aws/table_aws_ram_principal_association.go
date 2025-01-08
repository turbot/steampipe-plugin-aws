package aws

import (
	"context"

	ramEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRAMPrincipalAssociation(_ context.Context) *plugin.Table {
	associationType := "PRINCIPAL"
	return &plugin.Table{
		Name:        "aws_ram_principal_association",
		Description: "AWS RAM Principal Association",
		List: &plugin.ListConfig{
			Hydrate: listResourceShareAssociations(associationType),
			Tags:    map[string]string{"service": "ram", "action": "GetResourceShareAssociations"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getResourceSharePermissions,
				Tags: map[string]string{"service": "ram", "action": "ListResourceSharePermissions"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ramEndpoint.RAMServiceID),
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
				Description: "The ID of an Amazon Web Services account/The Amazon Resoure Name (ARN) of an organization in Organizations/The ARN of an organizational unit (OU) in Organizations/The ARN of an IAM role The ARN of an IAM user.",
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
