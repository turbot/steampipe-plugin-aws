package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIoTThingType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iot_thing_type",
		Description: "AWS IoT Thing Type",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("thing_type_name"),
			Hydrate:    getIoTThingType,
			Tags:       map[string]string{"service": "iot", "action": "DescribeThingType"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIoTThingTypes,
			Tags:    map[string]string{"service": "iot", "action": "ListThingTypes"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getIoTThingType,
				Tags: map[string]string{"service": "iot", "action": "DescribeThingType"},
			},
			{
				Func: getIoTThingTypeTags,
				Tags: map[string]string{"service": "iot", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_IOT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "thing_type_name",
				Description: "The name of the thing type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the thing type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ThingTypeArn"),
			},
			{
				Name:        "thing_type_id",
				Description: "The thing type ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTThingType,
			},
			{
				Name:        "thing_type_description",
				Description: "The description of the thing type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ThingTypeProperties.ThingTypeDescription"),
			},
			{
				Name:        "creation_date",
				Description: "The UNIX timestamp of when the thing type was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ThingTypeMetadata.CreationDate"),
			},
			{
				Name:        "deprecated",
				Description: "Whether the thing type is deprecated. If true, no new things could be associated with this type.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ThingTypeMetadata.Deprecated"),
			},
			{
				Name:        "deprecation_date",
				Description: "The date and time when the thing type was deprecated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ThingTypeMetadata.DeprecationDate"),
			},
			{
				Name:        "searchable_attributes",
				Description: "A list of searchable thing attribute names.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTThingType,
				Transform:   transform.FromField("ThingTypeProperties.SearchableAttributes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the thing type.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTThingTypeTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ThingTypeName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTThingTypeTags,
				Transform:   transform.From(iotThingTypeTagListToTagsMap),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ThingTypeArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listIoTThingTypes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_type.listIoTThingTypes", "connection_error", err)
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

	input := &iot.ListThingTypesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := iot.NewListThingTypesPaginator(svc, input, func(o *iot.ListThingTypesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iot_thing_type.listIoTThingTypes", "api_error", err)
			return nil, err
		}

		for _, item := range output.ThingTypes {
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

func getIoTThingType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	typeName := ""
	if h.Item != nil {
		t := h.Item.(types.ThingTypeDefinition)
		typeName = *t.ThingTypeName
	} else {
		typeName = d.EqualsQualString("thing_type_name")
	}

	if typeName == "" {
		return nil, nil
	}

	// Create service
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_type.getIoTThingType", "connection_error", err)
		return nil, err
	}

	params := &iot.DescribeThingTypeInput{
		ThingTypeName: aws.String(typeName),
	}

	resp, err := svc.DescribeThingType(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_type.getIoTThingType", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getIoTThingTypeTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	typeArn := ""
	switch item := h.Item.(type) {
	case *iot.DescribeThingTypeOutput:
		typeArn = *item.ThingTypeArn
	case types.ThingTypeDefinition:
		typeArn = *item.ThingTypeArn
	}

	// Create service
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_type.getIoTThingTypeTags", "connection_error", err)
		return nil, err
	}

	params := &iot.ListTagsForResourceInput{
		ResourceArn: aws.String(typeArn),
	}

	endpointTags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_type.getIoTThingTypeTags", "api_error", err)
		return nil, err
	}

	return endpointTags, nil
}

//// TRANSFORM FUNCTIONS

func iotThingTypeTagListToTagsMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*iot.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	if data.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
