package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	ssmv1 "github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSSMInventory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_inventory",
		Description: "AWS SSM Inventory",
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMInventories,
			Tags:    map[string]string{"service": "ssm", "action": "GetInventory"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "filter_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "filter_value", Require: plugin.Optional, Operators: []string{"=", "<>", ">", "<", ">=", "<="}},
				{Name: "type_name", Require: plugin.Optional},
				{Name: "component_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "component_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "instance_detailed_information_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "instance_detailed_information_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "network_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "network_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "windows_role_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "windows_role_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "service_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "service_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "windows_registry_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "windows_registry_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "compliance_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "compliance_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "file_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "file_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "instance_information_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "instance_information_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "patch_compliance_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "patch_compliance_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "patch_summary_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "patch_summary_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "application_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "application_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "tag_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "tag_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "resource_group_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "resource_group_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "windows_update_attribute_key", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "windows_update_attribute_value", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsSSMInventorySchema,
				Tags: map[string]string{"service": "ssm", "action": "GetInventorySchema"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssmv1.EndpointsID),
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
				Name:        "filter_key",
				Description: "The name of the filter key. Example: inventory filter key where managed node ID 'AWS:InstanceInformation.InstanceId'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "filter_key"),
			},
			{
				Name:        "filter_value",
				Description: "Inventory filter values. Example: inventory filter where managed node IDs are specified as values 'i-a12b3c4d5e6g'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "filter_value"),
			},
			{
				Name:        "component_attribute_key",
				Description: "The attribute key that are supported for type name AWS:AWSComponent, Possible values are: Name,ApplicationType,Publisher,Version,InstalledTime,Architecture and URL.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "component_attribute_key"),
			},
			{
				Name:        "component_attribute_value",
				Description: "The value for the component attribute key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "component_attribute_value"),
			},
			{
				Name:        "application_attribute_key",
				Description: "The attribute key of the type name AWS:Application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "application_attribute_key"),
			},
			{
				Name:        "application_attribute_value",
				Description: "The value for the attribute key of the type name AWS:Application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "application_attribute_value"),
			},
			{
				Name:        "compliance_attribute_key",
				Description: "The attribute key of the type name AWS:ComplianceItem.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "compliance_attribute_key"),
			},
			{
				Name:        "compliance_attribute_value",
				Description: "The value for the attribute key of the type name AWS:ComplianceItem.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "compliance_attribute_value"),
			},
			{
				Name:        "file_attribute_key",
				Description: "The attribute key of the type name AWS:File.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "file_attribute_key"),
			},
			{
				Name:        "file_attribute_value",
				Description: "The value for the attribute key of the type name AWS:File.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "file_attribute_value"),
			},
			{
				Name:        "instance_detailed_information_attribute_key",
				Description: "The attribute key of the type name AWS:InstanceDetailedInformation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "instance_detailed_information_attribute_key"),
			},
			{
				Name:        "instance_detailed_information_attribute_value",
				Description: "The value for the attribute key of the type name AWS:InstanceDetailedInformation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "instance_detailed_information_attribute_value"),
			},
			{
				Name:        "instance_information_attribute_key",
				Description: "The attribute key of the type name AWS:InstanceInformation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "instance_information_attribute_key"),
			},
			{
				Name:        "instance_information_attribute_value",
				Description: "The value for the attribute key of the type name AWS:InstanceInformation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "instance_information_attribute_value"),
			},
			{
				Name:        "network_attribute_key",
				Description: "The attribute key of the type name AWS:Network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "network_attribute_key"),
			},
			{
				Name:        "network_attribute_value",
				Description: "The value for the attribute key of the type name AWS:Network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "network_attribute_value"),
			},
			{
				Name:        "windows_registry_attribute_key",
				Description: "The attribute key of the type name AWS:WindowsRegistry.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "windows_registry_attribute_key"),
			},
			{
				Name:        "windows_registry_attribute_value",
				Description: "The value for the attribute key of the type name AWS:WindowsRegistry.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "windows_registry_attribute_value"),
			},
			{
				Name:        "patch_compliance_attribute_key",
				Description: "The attribute key of the type name AWS:PatchCompliance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "patch_compliance_attribute_key"),
			},
			{
				Name:        "patch_compliance_attribute_value",
				Description: "The value for the attribute key of the type name AWS:PatchCompliance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "patch_compliance_attribute_value"),
			},
			{
				Name:        "patch_summary_attribute_key",
				Description: "The attribute key of the type name AWS:PatchSummary.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "patch_summary_attribute_key"),
			},
			{
				Name:        "patch_summary_attribute_value",
				Description: "The value for the attribute key of the type name AWS:PatchSummary.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "patch_summary_attribute_value"),
			},
			{
				Name:        "resource_group_attribute_key",
				Description: "The attribute key of the type name AWS:ResourceGroup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "resource_group_attribute_key"),
			},
			{
				Name:        "resource_group_attribute_value",
				Description: "The value for the attribute key of the type name AWS:ResourceGroup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "resource_group_attribute_value"),
			},
			{
				Name:        "service_attribute_key",
				Description: "The attribute key of the type name AWS:Service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "service_attribute_key"),
			},
			{
				Name:        "service_attribute_value",
				Description: "The value for the attribute key of the type name AWS:Service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "service_attribute_value"),
			},
			{
				Name:        "tag_attribute_key",
				Description: "The attribute key of the type name AWS:Tag.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "tag_attribute_key"),
			},
			{
				Name:        "tag_attribute_value",
				Description: "The value for the attribute key of the type name AWS:Tag.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "tag_attribute_value"),
			},
			{
				Name:        "windows_role_attribute_key",
				Description: "The attribute key of the type name AWS:WindowsRole.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "windows_role_attribute_key"),
			},
			{
				Name:        "windows_role_attribute_value",
				Description: "The value for the attribute key of the type name AWS:WindowsRole.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "windows_role_attribute_value"),
			},
			{
				Name:        "windows_update_attribute_key",
				Description: "The attribute key of the type name AWS:WindowsUpdate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "windows_update_attribute_key"),
			},
			{
				Name:        "windows_update_attribute_value",
				Description: "The value for the attribute key of the type name AWS:WindowsUpdate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getFilterKeyWithOperator, "windows_update_attribute_value"),
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
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(50)
	input := buildSSMInventoryFilter(ctx, d)

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
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

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
					if d.RowsRemaining(ctx) == 0 {
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
	if svc == nil {
		// Unsupported region check
		return nil, nil
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
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

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
func buildSSMInventoryFilter(ctx context.Context, quals *plugin.QueryData) ssm.GetInventoryInput {

	input := ssm.GetInventoryInput{
		MaxResults: aws.Int32(50),
	}
	inventoryFilter := types.InventoryFilter{}

	shouldFilterKeyValueApplied := quals.EqualsQualString("filter_key") != "" && quals.EqualsQualString("filter_value") != ""

	var filterValue, filterOperator = "", ""

	if quals.EqualsQualString("filter_value") != "" {
		for _, q := range quals.Quals["filter_value"].Quals {
			filterValue = q.Value.GetStringValue()
			filterOperator = q.Operator
		}
	}

	// Optimize query results by filtering with "instance_id" if provided as a query parameter in the child table "aws_ssm_inventory_entry".
	if quals.EqualsQualString("instance_id") != "" {
		for _, q := range quals.Quals["instance_id"].Quals {
			input.Filters = []types.InventoryFilter{
				{
					Key:    aws.String("AWS:InstanceInformation.InstanceId"),
					Values: []string{quals.EqualsQualString("instance_id")},
				},
			}
			if q.Operator == "=" {
				input.Filters[0].Type = types.InventoryQueryOperatorTypeEqual
			} else if q.Operator == "<>" {
				input.Filters[0].Type = types.InventoryQueryOperatorTypeNotEqual
			}
		}
	}

	filterQuals := []string{"id", "type_name", "filter_key"}

	// Optional qulas for all supported type name along with their attribute key/value.
	filterQualsByTypeName := map[string][]string{
		"AWS:AWSComponent":                {"component_attribute_key", "component_attribute_value"},
		"AWS:Application":                 {"application_attribute_key", "application_attribute_value"},
		"AWS:ComplianceItem":              {"compliance_attribute_key", "compliance_attribute_value"},
		"AWS:File":                        {"file_attribute_key", "file_attribute_value"},
		"AWS:InstanceDetailedInformation": {"instance_detailed_information_attribute_key", "instance_detailed_information_attribute_value"},
		"AWS:InstanceInformation":         {"instance_information_attribute_key", "instance_information_attribute_value"},
		"AWS:Network":                     {"network_attribute_key", "network_attribute_value"},
		"AWS:PatchCompliance":             {"patch_compliance_attribute_key", "patch_compliance_attribute_value"},
		"AWS:PatchSummary":                {"patch_summary_attribute_key", "patch_summary_attribute_value"},
		"AWS:ResourceGroup":               {"resource_group_attribute_key", "resource_group_attribute_value"},
		"AWS:Service":                     {"service_attribute_key", "service_attribute_value"},
		"AWS:Tag":                         {"tag_attribute_key", "tag_attribute_value"},
		"AWS:WindowsRegistry":             {"windows_registry_attribute_key", "windows_registry_attribute_value"},
		"AWS:WindowsRole":                 {"windows_role_attribute_key", "windows_role_attribute_value"},
		"AWS:WindowsUpdate":               {"windows_update_attribute_key", "windows_update_attribute_value"},
	}

	for _, columnName := range filterQuals {
		if quals.Quals[columnName] != nil {
			value := getQualsValueByColumn(quals.Quals, columnName, "string")
			for _, q := range quals.Quals[columnName].Quals {
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

				// If we want to get the inventory details for a custom type name.
				case "filter_key":
					if shouldFilterKeyValueApplied {
						inventoryFilter.Key = aws.String(value.(string))
						inventoryFilter.Values = []string{filterValue}
						if filterOperator == "=" {
							inventoryFilter.Type = types.InventoryQueryOperatorTypeEqual
						} else if filterOperator == "<>" {
							inventoryFilter.Type = types.InventoryQueryOperatorTypeNotEqual
						} else if filterOperator == "<" || filterOperator == "<=" {
							inventoryFilter.Type = types.InventoryQueryOperatorTypeLessThan
						} else if filterOperator == ">" || filterOperator == ">=" {
							inventoryFilter.Type = types.InventoryQueryOperatorTypeGreaterThan
						}
						input.Filters = append(input.Filters, inventoryFilter)
					}

				// Supported type names are: AWS:InstanceInformation, AWS:PatchSummary. Default result type name is AWS:InstanceInformation.Supported pattern is ^(AWS|Custom):.*$
				// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_ResultAttribute.html
				case "type_name":
					if q.Operator == "=" && quals.Table.Name == "aws_ssm_inventory" {
						if value.(string) == "AWS:InstanceInformation" || value.(string) == "AWS:PatchSummary" || strings.HasPrefix(value.(string), "Custom:") {
							input.ResultAttributes = []types.ResultAttribute{{TypeName: aws.String(value.(string))}}
						}
					}
				}
			}
		}
	}

	// Build the filter as per the supported attributes for inventory type names.
	// You can get the supported attributs by type name by running the command "aws ssm get-inventory-schema"
	for typeName, typeNameQualsColumn := range filterQualsByTypeName {
		f := types.InventoryFilter{}

		if typeName == "AWS:Tag" {
			if quals.Quals[typeNameQualsColumn[0]] != nil {
				f.Key = aws.String("AWS:Tag.Key")
				if v, ok := getQualsValueByColumn(quals.Quals, typeNameQualsColumn[0], "string").(string); ok {
					f.Values = []string{v}
				} else if v, ok := getQualsValueByColumn(quals.Quals, typeNameQualsColumn[0], "string").([]*string); ok {
					f.Values = getStringSliceWithoutPointer(v)
				}
				f.Type = types.InventoryQueryOperatorTypeEqual
				input.Filters = append(input.Filters, f)
			}
			if quals.Quals[typeNameQualsColumn[1]] != nil {
				f.Key = aws.String("AWS:Tag.Value")
				if v, ok := getQualsValueByColumn(quals.Quals, typeNameQualsColumn[1], "string").(string); ok {
					f.Values = []string{v}
				} else if v, ok := getQualsValueByColumn(quals.Quals, typeNameQualsColumn[1], "string").([]*string); ok {
					f.Values = getStringSliceWithoutPointer(v)
				}
				f.Type = types.InventoryQueryOperatorTypeEqual
				input.Filters = append(input.Filters, f)
			}
			continue
		}

		if quals.Quals[typeNameQualsColumn[0]] != nil && quals.Quals[typeNameQualsColumn[1]] != nil {
			attKey := getQualsValueByColumn(quals.Quals, typeNameQualsColumn[0], "string")
			attValue := getQualsValueByColumn(quals.Quals, typeNameQualsColumn[1], "string")

			if k, ok := attKey.(string); ok {
				f.Key = aws.String(typeName + "." + k)
				if v, ok := attValue.(string); ok {
					f.Values = []string{v}
				} else if v, ok := attValue.([]*string); ok {
					f.Values = getStringSliceWithoutPointer(v)
				}
				f.Type = types.InventoryQueryOperatorTypeEqual
				input.Filters = append(input.Filters, f)
			} else if k, ok := attKey.([]*string); ok {
				for _, aK := range k {
					f.Key = aws.String(typeName + "." + *aK)
					if v, ok := attValue.(string); ok {
						f.Values = []string{v}
					} else if v, ok := attValue.([]*string); ok {
						f.Values = getStringSliceWithoutPointer(v)
					}
					f.Type = types.InventoryQueryOperatorTypeEqual
					input.Filters = append(input.Filters, f)
				}
			}

		}
	}

	return input
}

func getFilterKeyWithOperator(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	data := d.KeyColumnQuals[param]
	for _, q := range data {
		if q.Operator == "=" {
			return q.Value.GetStringValue(), nil
		}
	}
	return "", nil
}

//// UTILITY FUNCTION

func getStringSliceWithoutPointer(originalSlice []*string) []string {
	convertedSlice := make([]string, 0, len(originalSlice))
	for _, strPtr := range originalSlice {
		if strPtr != nil {
			convertedSlice = append(convertedSlice, *strPtr)
		}
	}
	return convertedSlice
}
