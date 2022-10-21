package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

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

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_inventory.listAwsSSMInventories", "connection_error", err)
		return nil, err
	}

	maxItems := int32(50)
	input := buildSSMInventoryFilter(ctx, d.Quals)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewGetInventoryPaginator(svc, &input, func(o *ssm.GetInventoryPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_inventory.listAwsSSMInventories", "api_error", err)
			return nil, err
		}

		for _, inventory := range output.Entities {
			if inventory.Data != nil {
				for _, data := range inventory.Data {
					d.StreamListItem(ctx, &InventoryInfo{
						Id:            inventory.Id,
						CaptureTime:   data.CaptureTime,
						SchemaVersion: data.SchemaVersion,
						TypeName:      data.TypeName,
						Content:       data.Content,
					})

					// Context may get cancelled due to manual cancellation or if the limit has been reached
					if d.QueryStatus.RowsRemaining(ctx) == 0 {
						return nil, nil
					}
				}
			}
		}
	}

	return nil, nil
}

func getAwsSSMInventorySchema(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_inventory.getAwsSSMInventorySchema", "connection_error", err)
		return nil, err
	}

	inventory := h.Item.(*InventoryInfo)

	input := &ssm.GetInventorySchemaInput{
		TypeName: inventory.TypeName,
	}

	paginator := ssm.NewGetInventorySchemaPaginator(svc, input, func(o *ssm.GetInventorySchemaPaginatorOptions) {
		o.Limit = 200
		o.StopOnDuplicateToken = true
	})

	var schemas []types.InventoryItemSchema

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_inventory.getAwsSSMInventorySchema", "api_error", err)
			return nil, err
		}
		schemas = append(schemas, output.Schemas...)
	}

	return schemas, nil
}

//// UTILITY FUNCTION

// Build SSM inventory list call input filter
func buildSSMInventoryFilter(ctx context.Context, quals plugin.KeyColumnQualMap) ssm.GetInventoryInput {

	input := ssm.GetInventoryInput{
		MaxResults: aws.Int32(50),
	}
	inventoryFilter := types.InventoryFilter{}

	filterQuals := []string{"id", "type_name"}

	for _, columnName := range filterQuals {
		if quals[columnName] != nil {
			value := getQualsValueByColumn(quals, columnName, "string")
			for _, q := range quals[columnName].Quals {
				switch columnName {
				case "id":
					input.Filters = []types.InventoryFilter{
						{
							Key:    aws.String("AWS:InstanceInformation.InstanceId"),
							Values: []string{value.(string)},
						},
					}
					if q.Operator == "=" {
						input.Filters[0].Type = types.InventoryQueryOperatorTypeEqual
					} else if q.Operator == "<>" {
						input.Filters[0].Type = types.InventoryQueryOperatorTypeNotEqual
					}
					input.Filters = append(input.Filters, inventoryFilter)

				case "type_name":
					if q.Operator == "=" {
						input.ResultAttributes = []types.ResultAttribute{{TypeName: aws.String(value.(string))}}
					}
				}
			}
		}
	}

	return input
}
