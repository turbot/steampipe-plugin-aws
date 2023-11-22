package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmv1 "github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSSMInventoryEntry(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_inventory_entry",
		Description: "AWS SSM Inventory Entry",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsSSMInventories,
			Hydrate:       listAwsSSMInventoryEntries,
			Tags:          map[string]string{"service": "ssm", "action": "ListInventoryEntries"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "type_name", Require: plugin.Required},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssmv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The managed node ID targeted by the request to query inventory information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type_name",
				Description: "The type of inventory item returned by the request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capture_time",
				Description: "The time that inventory information was collected for the managed node(s).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "schema_version",
				Description: "The inventory schema version used by the managed node(s).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "entries",
				Description: "The inventory items on the managed node(s).",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceId"),
			},
		}),
	}
}

type InventoryEntryInfo struct {
	CaptureTime   *string
	Entries       map[string]string
	InstanceId    *string
	SchemaVersion *string
	TypeName      *string
}

//// LIST FUNCTION

func listAwsSSMInventoryEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	inventory := h.Item.(*InventoryInfo)

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_inventory_entry.listAwsSSMInventoryEntry", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	if d.EqualsQualString("instance_id") != "" {
		if d.EqualsQualString("instance_id") != *inventory.Id {
			return nil, nil
		}
	}
	// if d.EqualsQualString("entry_type_name") != "" {
	// 	if d.EqualsQualString("entry_type_name") != *inventory.TypeName {
	// 		return nil, nil
	// 	}
	// }

	maxItems := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := &ssm.ListInventoryEntriesInput{
		InstanceId: inventory.Id,
		TypeName:   aws.String(d.EqualsQualString("type_name")),
		MaxResults: &maxItems,
	}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		op, err := svc.ListInventoryEntries(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_inventory_entry.listAwsSSMInventoryEntry", "api_error", err)
			return nil, err
		}

		for _, item := range op.Entries {
			d.StreamListItem(ctx, &InventoryEntryInfo{
				CaptureTime:   op.CaptureTime,
				InstanceId:    op.InstanceId,
				Entries:       item,
				SchemaVersion: op.SchemaVersion,
				TypeName:      op.TypeName,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if op.NextToken != nil {
			input.NextToken = op.NextToken
		} else {
			break
		}
	}

	return nil, nil
}
