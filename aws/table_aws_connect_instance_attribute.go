package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/aws/aws-sdk-go-v2/service/connect/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Custom struct to hold instance attribute data with instance ID
type connectInstanceAttributeData struct {
	types.Attribute
	InstanceId string
}

//// TABLE DEFINITION

func tableAwsConnectInstanceAttribute(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_connect_instance_attribute",
		Description: "AWS Connect Instance Attribute",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"instance_id", "attribute_type"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getConnectInstanceAttribute,
			Tags:    map[string]string{"service": "connect", "action": "DescribeInstanceAttribute"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listConnectInstances,
			Hydrate:       listConnectInstanceAttributes,
			KeyColumns:    plugin.OptionalColumns([]string{"instance_id"}),
			Tags:          map[string]string{"service": "connect", "action": "ListInstanceAttributes"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The identifier of the Amazon Connect instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attribute_type",
				Description: "The type of attribute.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The value of the attribute.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AttributeType"),
			},
		}),
	}
}

//// LIST FUNCTION

func listConnectInstanceAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the parent instance data
	instance := h.Item.(types.InstanceSummary)
	instanceId := *instance.Id

	if d.EqualsQualString("instance_id") != "" && d.EqualsQualString("attribute_type") != instanceId {
		return nil, nil
	}

	// Create service
	svc, err := ConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_connect_instance_attribute.listConnectInstanceAttributes", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params for ListInstanceAttributes
	params := &connect.ListInstanceAttributesInput{
		InstanceId: aws.String(instanceId),
		MaxResults: aws.Int32(10),
	}

	// Create a paginator for ListInstanceAttributes
	paginator := connect.NewListInstanceAttributesPaginator(svc, params, func(o *connect.ListInstanceAttributesPaginatorOptions) {
		o.Limit = 10
	})

	// Iterate through all pages
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_connect_instance_attribute.listConnectInstanceAttributes", "api_error", err)
			return nil, err
		}

		if output != nil && output.Attributes != nil {
			for _, attr := range output.Attributes {
				// Create a custom struct to include instance_id
				attributeWithInstance := connectInstanceAttributeData{
					Attribute:  attr,
					InstanceId: instanceId,
				}

				d.StreamListItem(ctx, attributeWithInstance)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getConnectInstanceAttribute(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceId := d.EqualsQualString("instance_id")
	attributeType := d.EqualsQualString("attribute_type")

	if instanceId == "" || attributeType == "" {
		return nil, nil
	}

	// Create service
	svc, err := ConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_connect_instance_attribute.getConnectInstanceAttribute", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params for ListInstanceAttributes
	params := &connect.DescribeInstanceAttributeInput{
		InstanceId:    aws.String(instanceId),
		AttributeType: types.InstanceAttributeType(attributeType),
	}

	op, err := svc.DescribeInstanceAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_connect_instance_attribute.getConnectInstanceAttribute", "api_error", err)
		return nil, err
	}

	return connectInstanceAttributeData{
		Attribute:  *op.Attribute,
		InstanceId: instanceId,
	}, nil
}
