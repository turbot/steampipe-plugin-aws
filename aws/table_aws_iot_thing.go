package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"

	iotv1 "github.com/aws/aws-sdk-go/service/iot"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIotThing(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iot_thing",
		Description: "AWS Iot Thing",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("thing_name"),
			Hydrate:    getIotThing,
			Tags:       map[string]string{"service": "iot", "action": "DescribeThing"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIotThings,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "iot", "action": "ListThings"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "attribute_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "attribute_value", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "thing_type_name", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getIotThing,
				Tags: map[string]string{"service": "iot", "action": "DescribeThing"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(iotv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "thing_name",
				Description: "The name of the thing.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "thing_id",
				Description: "The ID of the thing to describe.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIotThing,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the thing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ThingArn"),
			},
			{
				Name:        "thing_type_name",
				Description: "The name of the thing type, if the thing has been associated with a type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attribute_name",
				Description: "The attribute name of the thing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("attribute_name"),
			},
			{
				Name:        "attribute_value",
				Description: "The attribute value for the attribute name of the thing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("attribute_value"),
			},
			{
				Name:        "billing_group_name",
				Description: "The name of the billing group the thing belongs to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIotThing,
			},
			{
				Name:        "default_client_id",
				Description: "The default MQTT client ID. For a typical device, the thing name is also used as the default MQTT client ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIotThing,
			},
			{
				Name:        "version",
				Description: "The version of the thing record in the registry.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "attributes",
				Description: "A list of thing attributes which are name-value pairs.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ThingName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ThingArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listIotThings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IOTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing.listIotThings", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(250)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &iot.ListThingsInput{
		MaxResults:              aws.Int32(maxLimit),
		UsePrefixAttributeValue: false,
	}

	if d.EqualsQualString("attribute_name") != "" {
		input.AttributeName = aws.String(d.EqualsQualString("attribute_name"))
	}
	if d.EqualsQualString("attribute_value") != "" && d.EqualsQualString("attribute_name") != ""{
		input.AttributeValue = aws.String(d.EqualsQualString("attribute_value"))
	}
	if d.EqualsQualString("thing_type_name") != "" {
		input.ThingTypeName = aws.String(d.EqualsQualString("thing_type_name"))
	}

	paginator := iot.NewListThingsPaginator(svc, input, func(o *iot.ListThingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iot_thing.listIotThings", "api_error", err)
			return nil, err
		}

		for _, item := range output.Things {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIotThing(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	thingName := ""
	if h.Item != nil {
		thing := h.Item.(types.ThingAttribute)
		thingName = *thing.ThingName
	} else {
		thingName = d.EqualsQualString("thing_name")
	}

	if thingName == "" {
		return nil, nil
	}

	// Create service
	svc, err := IOTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing.getIotThing", "connection_error", err)
		return nil, err
	}

	params := &iot.DescribeThingInput{
		ThingName: aws.String(thingName),
	}

	thing, err := svc.DescribeThing(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing.getIotThing", "api_error", err)
		return nil, err
	}

	return thing, nil
}
