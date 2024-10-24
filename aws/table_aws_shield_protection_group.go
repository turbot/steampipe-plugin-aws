package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/shield"
	"github.com/aws/aws-sdk-go-v2/service/shield/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldProtectionGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_protection_group",
		Description: "AWS Shield Protection Group",
		Get: &plugin.GetConfig{
			KeyColumns:   plugin.SingleColumn("protection_group_id"),
			Hydrate:      getAwsShieldProtectionGroup,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags:         map[string]string{"service": "shield", "action": "DescribeProtectionGroup"},
		},
		List: &plugin.ListConfig{
			Hydrate:      listAwsShieldProtectionGroups,
			KeyColumns:   plugin.OptionalColumns([]string{"protection_group_id", "pattern", "resource_type", "aggregation"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags:         map[string]string{"service": "shield", "action": "ListProtectionGroups"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listTagsForShieldProtectionGroup,
				Tags: map[string]string{"service": "shield", "action": "ListTagsForResource"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "protection_group_id",
				Description: "The name of the protection group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "aggregation",
				Description: "Defines how Shield combines resource data for the group in order to detect, mitigate, and report events.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pattern",
				Description: "The criteria to use to choose the protected resources for inclusion in the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The resource type to include in the protection group. All protected resources of this type are included in the protection group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceType").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the protection group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProtectionGroupArn"),
			},
			{
				Name:        "members",
				Description: "The ARNs (Amazon Resource Names) of the resources that are included in the protection group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Members"),
			},
			{
				Name:        "tags_src",
				Description: "The list of tags associated with the protection group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForShieldProtectionGroup,
				Transform:   transform.FromField("Tags"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProtectionGroupId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProtectionGroupArn").Transform(transform.EnsureStringArray),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForShieldProtectionGroup,
				Transform:   transform.From(handleShieldProtectionGroupTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTIONS

func listAwsShieldProtectionGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection_group.listAwsShieldProtectionGroups", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	queryResultLimit := int32(10000)
	if d.QueryContext.Limit != nil {
		queryResultLimit = min(queryResultLimit, int32(*d.QueryContext.Limit))
	}

	input := &shield.ListProtectionGroupsInput{
		MaxResults: aws.Int32(queryResultLimit),
		InclusionFilters: &types.InclusionProtectionGroupFilters{},
	}

	if d.Quals["protection_group_id"] != nil {
		for _, q := range d.Quals["protection_group_id"].Quals {
			input.InclusionFilters.ProtectionGroupIds = []string{}
			input.InclusionFilters.ProtectionGroupIds = append(input.InclusionFilters.ProtectionGroupIds, q.Value.GetStringValue())
		}
	}

	if d.Quals["pattern"] != nil {
		for _, q := range d.Quals["pattern"].Quals {
			input.InclusionFilters.Patterns = []types.ProtectionGroupPattern{}
			input.InclusionFilters.Patterns = append(input.InclusionFilters.Patterns, types.ProtectionGroupPattern(q.Value.GetStringValue()))
		}
	}

	if d.Quals["resource_type"] != nil {
		for _, q := range d.Quals["resource_type"].Quals {
			input.InclusionFilters.ResourceTypes = []types.ProtectedResourceType{}
			input.InclusionFilters.ResourceTypes = append(input.InclusionFilters.ResourceTypes, types.ProtectedResourceType(q.Value.GetStringValue()))
		}
	}

	if d.Quals["aggregation"] != nil {
		for _, q := range d.Quals["aggregation"].Quals {
			input.InclusionFilters.Aggregations = []types.ProtectionGroupAggregation{}
			input.InclusionFilters.Aggregations = append(input.InclusionFilters.Aggregations, types.ProtectionGroupAggregation(q.Value.GetStringValue()))
		}
	}

	paginator := shield.NewListProtectionGroupsPaginator(svc, input, func(o *shield.ListProtectionGroupsPaginatorOptions) {
		o.Limit = queryResultLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_shield_protection.listAwsShieldProtections", "api_error", err)
			return nil, err
		}

		for _, items := range output.ProtectionGroups {
			d.StreamListItem(ctx, &items)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsShieldProtectionGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection_group.getAwsShieldProtectionGroup", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var protectionGroupId string
	if h.Item != nil {
		protectionGroupId = *h.Item.(*types.ProtectionGroup).ProtectionGroupId
	} else {
		protectionGroupId = d.EqualsQualString("protection_group_id")
	}

	var params *shield.DescribeProtectionGroupInput
	if protectionGroupId != "" {
		params = &shield.DescribeProtectionGroupInput{
			ProtectionGroupId: aws.String(protectionGroupId),
		}
	} else {
		return nil, nil
	}

	data, err := svc.DescribeProtectionGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection_group.getAwsShieldProtectionGroup", "api_error", err)
		return nil, err
	}

	if data != nil {
		return data.ProtectionGroup, nil
	}

	return nil, nil
}

func listTagsForShieldProtectionGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	protectionGroupArn := h.Item.(*types.ProtectionGroup).ProtectionGroupArn

	// Get client
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection_group.listTagsForShieldProtectionGroup", "get_client_error", err)
		return nil, err
	}

	// Build param
	param := &shield.ListTagsForResourceInput{
		ResourceARN: protectionGroupArn,
	}

	protectionTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection_group.listTagsForShieldProtectionGroup", "api_error", err)
		return nil, err
	}
	return protectionTags, nil
}

//// TRANSFORM FUNCTIONS

func handleShieldProtectionGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*shield.ListTagsForResourceOutput)
	if len(tags.Tags) == 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range tags.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}