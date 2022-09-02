package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSSMInventory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_inventory",
		Description: "AWS SSM Inventory",
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMInventories,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "id", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "type_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "ID of the inventory result entity.",
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
				Name:        "content",
				Description: "Contains all the inventory data of the item type. Results include attribute names and values.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "schema",
				Description: "The inventory item schema definition. Users can use this to compose inventory query filters.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMInventorySchema,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
		}),
	}
}

type InventoryInfo struct {
	CaptureTime   *string
	Content       interface{}
	Id            *string
	SchemaVersion *string
	TypeName      *string
}

//// LIST FUNCTION

func listAwsSSMInventories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsSSMInventories")

	// Create session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := buildSsmInventoryFilter(ctx, d.Quals)

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

	// List call
	err = svc.GetInventoryPages(
		input,
		func(page *ssm.GetInventoryOutput, isLast bool) bool {
			for _, inventory := range page.Entities {
				if inventory.Data != nil {
					for _, v := range inventory.Data {
						d.StreamListItem(ctx, &InventoryInfo{
							Id:            inventory.Id,
							CaptureTime:   v.CaptureTime,
							SchemaVersion: v.SchemaVersion,
							TypeName:      v.TypeName,
							Content:       v.Content,
						})

						// Context may get cancelled due to manual cancellation or if the limit has been reached
						if d.QueryStatus.RowsRemaining(ctx) == 0 {
							return false
						}
					}
				}

			}
			return !isLast
		},
	)

	return nil, err
}

func getAwsSSMInventorySchema(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMInventorySchema")

	// Create session
	svc, err := SsmService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getAwsSSMInventorySchema", "connection_error", err)
		return nil, err
	}

	inventory := h.Item.(*InventoryInfo)

	input := &ssm.GetInventorySchemaInput{
		TypeName: inventory.TypeName,
	}

	var schemas []*ssm.InventoryItemSchema

	err = svc.GetInventorySchemaPages(
		input,
		func(page *ssm.GetInventorySchemaOutput, isLast bool) bool {

			schemas = append(schemas, page.Schemas...)
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("getAwsSSMInventorySchema", "api_error", err)
		return nil, err
	}

	return schemas, nil
}

//// UTILITY FUNCTION

// Build ssm inventory list call input filter
func buildSsmInventoryFilter(ctx context.Context, quals plugin.KeyColumnQualMap) *ssm.GetInventoryInput {

	input := &ssm.GetInventoryInput{
		MaxResults: aws.Int64(50),
	}
	inventoryFilter := &ssm.InventoryFilter{}
	resultAttribute := &ssm.ResultAttribute{}

	filterQuals := []string{"id", "type_name"}

	for _, columnName := range filterQuals {
		if quals[columnName] != nil {
			value := getQualsValueByColumn(quals, columnName, "string")
			for _, q := range quals[columnName].Quals {
				switch columnName {
				case "id":
					inventoryFilter.Key = aws.String("AWS:InstanceInformation.InstanceId")
					inventoryFilter.Values = []*string{aws.String(value.(string))}
					if q.Operator == "=" {
						inventoryFilter.Type = aws.String("Equal")
					} else if q.Operator == "<>" {
						inventoryFilter.Type = aws.String("NotEqual")
					}
					input.Filters = append(input.Filters, inventoryFilter)
				case "type_name":
					if q.Operator == "=" {
						resultAttribute.TypeName = aws.String(value.(string))
						input.ResultAttributes = append(input.ResultAttributes, resultAttribute)
					}
				}
			}
		}
	}
	return input
}
