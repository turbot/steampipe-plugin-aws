package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsSSMInventory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_inventory",
		Description: "AWS SSM Inventory",
		// Get: &plugin.GetConfig{
		// 	KeyColumns:        plugin.SingleColumn("name"),
		// 	ShouldIgnoreError: isNotFoundError([]string{"ValidationException", "InvalidDocument"}),
		// 	Hydrate:           getAwsSSMInventory,
		// },
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMInventories,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "type_name",
				Description: "The type of inventory item returned by the request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "The managed node ID targeted by the request to query inventory information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capture_time",
				Description: "The time that inventory information was collected for the managed node(s).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema_version",
				Description: "The inventory schema version used by the managed node(s).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "entry",
				Description: "An inventory item on the managed node(s).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "data",
				Description: "An inventory item on the managed node(s).",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceId"),
			},
		}),
	}
}

type InventoryInfo struct {
	CaptureTime   *string
	Entry         map[string]*string
	InstanceId    *string
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

	input := &ssm.GetInventoryInput{
		MaxResults: aws.Int64(50),
	}

	// filters := buildSsmDocumentFilter(d.Quals)
	// if len(filters) > 0 {
	// 	input.Filters = filters
	// }

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
				d.StreamListItem(ctx, inventory)

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

//// HYDRATE FUNCTIONS

// func getAwsSSMInventory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	logger.Trace("getAwsSSMInventory")

// 	var name string
// 	if h.Item != nil {
// 		name = documentName(h.Item)
// 	} else {
// 		name = d.KeyColumnQuals["name"].GetStringValue()
// 	}

// 	// Create Session
// 	svc, err := SsmService(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build the params
// 	params := &ssm.GetInventoryInput{

// 	}

// 	// Get call
// 	data, err := svc.GetInventory(params)
// 	if err != nil {
// 		logger.Debug("getAwsSSMDocument", "ERROR", err)
// 		return nil, err
// 	}

// 	return data.Document, nil
// }

// func getAwsSSMDocumentPermissionDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	logger.Trace("getAwsSSMDocumentPermissionDetail")

// 	var name string
// 	if h.Item != nil {
// 		name = documentName(h.Item)
// 	} else {
// 		name = d.KeyColumnQuals["name"].GetStringValue()
// 	}

// 	// Create Session
// 	svc, err := SsmService(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build the params
// 	params := &ssm.DescribeDocumentPermissionInput{
// 		Name:           &name,
// 		PermissionType: aws.String("Share"),
// 	}

// 	// Get call
// 	data, err := svc.DescribeDocumentPermission(params)
// 	if err != nil {
// 		logger.Debug("getAwsSSMDocumentPermissionDetail", "ERROR", err)
// 		return nil, err
// 	}

// 	return data, nil
// }

// func getAwsSSMDocumentAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getAwsSSMDocumentAkas")
// 	region := d.KeyColumnQualString(matrixKeyRegion)
// 	name := documentName(h.Item)
// 	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
// 	c, err := getCommonColumnsCached(ctx, d, h)
// 	if err != nil {
// 		return nil, err
// 	}
// 	commonColumnData := c.(*awsCommonColumnData)
// 	aka := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":document"

// 	if strings.HasPrefix(name, "/") {
// 		aka = aka + name
// 	} else {
// 		aka = aka + "/" + name
// 	}

// 	return []string{aka}, nil
// }

// func ssmDocumentTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("ssmDocumentTagListToTurbotTags")
// 	data := resourceTags(d.HydrateItem)

// 	if data == nil {
// 		return nil, nil
// 	}
// 	// Mapping the resource tags inside turbotTags
// 	var turbotTagsMap map[string]string
// 	if data != nil {
// 		turbotTagsMap = map[string]string{}
// 		for _, i := range data {
// 			turbotTagsMap[*i.Key] = *i.Value
// 		}
// 	}

// 	return turbotTagsMap, nil
// }

// func documentName(item interface{}) string {
// 	switch item := item.(type) {
// 	case *ssm.DocumentDescription:
// 		return *item.Name
// 	case *ssm.DocumentIdentifier:
// 		return *item.Name
// 	}
// 	return ""
// }

// func resourceTags(item interface{}) []*ssm.Tag {
// 	switch item := item.(type) {
// 	case *ssm.DocumentDescription:
// 		return item.Tags
// 	case *ssm.DocumentIdentifier:
// 		return item.Tags
// 	}
// 	return nil
// }

//// UTILITY FUNCTION

// Build ssm documant list call input filter
// func buildSsmDocumentFilter(quals plugin.KeyColumnQualMap) []*ssm.DocumentKeyValuesFilter {
// 	filters := make([]*ssm.DocumentKeyValuesFilter, 0)

// 	filterQuals := map[string]string{
// 		"owner":         "Owner",
// 		"document_type": "DocumentType",
// 	}

// 	for columnName, filterName := range filterQuals {
// 		if quals[columnName] != nil {
// 			filter := ssm.DocumentKeyValuesFilter{
// 				Key: aws.String(filterName),
// 			}

// 			value := getQualsValueByColumn(quals, columnName, "string")
// 			val, ok := value.(string)
// 			if ok {
// 				filter.Values = []*string{&val}
// 			} else {
// 				filter.Values = value.([]*string)
// 			}
// 			filters = append(filters, &filter)
// 		}
// 	}
// 	return filters
// }
