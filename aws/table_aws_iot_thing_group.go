package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"

	iotEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIoTThingGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iot_thing_group",
		Description: "AWS IoT Thing Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("group_name"),
			Hydrate:    getIoTThingGroup,
			Tags:       map[string]string{"service": "iot", "action": "DescribeThingGroup"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIoTThingGroups,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "iot", "action": "ListThingGroups"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "parent_group_name", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getIoTThingGroup,
				Tags: map[string]string{"service": "iot", "action": "DescribeThingGroup"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(iotEndpoint.IOTServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "group_name",
				Description: "The group name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ThingGroupName", "GroupName"),
			},
			{
				Name:        "thing_group_id",
				Description: "The thing group ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTThingGroup,
			},
			{
				Name:        "thing_group_description",
				Description: "The thing group description.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTThingGroup,
				Transform:   transform.FromField("ThingGroupProperties.ThingGroupDescription"),
			},
			{
				Name:        "creation_date",
				Description: "The UNIX timestamp of when the thing group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getIoTThingGroup,
				Transform:   transform.FromField("ThingGroupMetadata.CreationDate"),
			},
			{
				Name:        "parent_group_name",
				Description: "The parent thing group name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTThingGroup,
				Transform:   transform.FromField("ThingGroupMetadata.ParentGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the thing group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupArn", "ThingGroupArn"),
			},
			{
				Name:        "status",
				Description: "The dynamic thing group status.",
				Hydrate:     getIoTThingGroup,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "index_name",
				Description: "The dynamic thing group index name.",
				Hydrate:     getIoTThingGroup,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "query_string",
				Description: "The dynamic thing group search query string.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTThingGroup,
			},
			{
				Name:        "query_version",
				Description: "The dynamic thing group query version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTThingGroup,
			},
			{
				Name:        "version",
				Description: "The version of the thing group.",
				Hydrate:     getIoTThingGroup,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "attribute_payload",
				Description: "The thing group attributes in JSON format.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTThingGroup,
				Transform:   transform.FromField("ThingGroupProperties.AttributePayload"),
			},
			{
				Name:        "root_to_parent_thing_groups",
				Description: "The root parent thing group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTThingGroup,
				Transform:   transform.FromField("ThingGroupMetadata.RootToParentThingGroups"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the thing group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTThingGroupTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ThingGroupName", "GroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTThingGroupTags,
				Transform:   transform.From(iotThingGroupTagListToTagsMap),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("GroupArn", "ThingGroupArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listIoTThingGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_group.listIoTThingGroups", "connection_error", err)
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

	input := &iot.ListThingGroupsInput{
		MaxResults: aws.Int32(maxLimit),
		Recursive:  aws.Bool(true),
	}

	if d.EqualsQualString("parent_group_name") != "" {
		input.ParentGroup = aws.String(d.EqualsQualString("parent_group_name"))
	}

	paginator := iot.NewListThingGroupsPaginator(svc, input, func(o *iot.ListThingGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iot_thing_group.listIoTThingGroups", "api_error", err)
			return nil, err
		}

		for _, item := range output.ThingGroups {
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

func getIoTThingGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	grpName := ""
	if h.Item != nil {
		g := h.Item.(types.GroupNameAndArn)
		grpName = *g.GroupName
	} else {
		grpName = d.EqualsQualString("group_name")
	}

	if grpName == "" {
		return nil, nil
	}

	// Create service
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_group.getIoTThingGroup", "connection_error", err)
		return nil, err
	}

	params := &iot.DescribeThingGroupInput{
		ThingGroupName: aws.String(grpName),
	}

	resp, err := svc.DescribeThingGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_group.getIoTThingGroup", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getIoTThingGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	grpArn := ""
	switch item := h.Item.(type) {
	case *iot.DescribeThingGroupOutput:
		grpArn = *item.ThingGroupArn
	case types.GroupNameAndArn:
		grpArn = *item.GroupArn
	}

	// Create service
	svc, err := IoTClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_group.getIoTThingGroupTags", "connection_error", err)
		return nil, err
	}

	params := &iot.ListTagsForResourceInput{
		ResourceArn: aws.String(grpArn),
	}

	endpointTags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_thing_group.getIoTThingGroupTags", "api_error", err)
		return nil, err
	}

	return endpointTags, nil
}

//// TRANSFORM FUNCTIONS

func iotThingGroupTagListToTagsMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
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
