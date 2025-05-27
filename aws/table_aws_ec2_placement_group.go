package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsEc2PlacementGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_placement_group",
		Description: "AWS EC2 Placement Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"group_id", "group_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidPlacementGroup.Unknown"}),
			},
			Hydrate: getEc2PlacementGroup,
			Tags:    map[string]string{"service": "ec2", "action": "DescribePlacementGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2PlacementGroups,
			Tags:    map[string]string{"service": "ec2", "action": "DescribePlacementGroups"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "spread_level", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "strategy", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EC2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The ARN of the placement group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupArn"),
			},
			{
				Name:        "group_id",
				Description: "The ID of the placement group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_name",
				Description: "The name of the placement group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "partition_count",
				Description: "The number of partitions (for partition strategy).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "spread_level",
				Description: "The spread level for the placement group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the placement group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "strategy",
				Description: "The placement strategy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the placement group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2PlacementGroupTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("GroupArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2PlacementGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_placement_group.listEc2PlacementGroups", "connection_error", err)
		return nil, err
	}
	input := &ec2.DescribePlacementGroupsInput{}

	filters := buildEc2PlacementGroupFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Pagination does not support
	resp, err := svc.DescribePlacementGroups(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_placement_group.listEc2PlacementGroups", "api_error", err)
		return nil, err
	}
	for _, group := range resp.PlacementGroups {
		d.StreamListItem(ctx, group)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTION

func getEc2PlacementGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var input *ec2.DescribePlacementGroupsInput
	if d.EqualsQuals["group_id"] != nil && d.EqualsQuals["group_id"].GetStringValue() != "" {
		input = &ec2.DescribePlacementGroupsInput{
			GroupIds: []string{d.EqualsQuals["group_id"].GetStringValue()},
		}
	} else if d.EqualsQuals["group_name"] != nil && d.EqualsQuals["group_name"].GetStringValue() != "" {
		input = &ec2.DescribePlacementGroupsInput{
			GroupNames: []string{d.EqualsQuals["group_name"].GetStringValue()},
		}
	} else {
		return nil, nil
	}

	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_placement_group.getEc2PlacementGroup", "connection_error", err)
		return nil, err
	}

	op, err := svc.DescribePlacementGroups(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_placement_group.getEc2PlacementGroup", "api_error", err)
		return nil, err
	}
	if len(op.PlacementGroups) > 0 {
		return op.PlacementGroups[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTION

func getEc2PlacementGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(types.PlacementGroup)
	if group.Tags == nil {
		return nil, nil
	}
	turbotTagsMap := map[string]string{}
	for _, i := range group.Tags {
		if i.Key != nil && i.Value != nil {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return &turbotTagsMap, nil
}

// Build input parameter for list function
func buildEc2PlacementGroupFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := []types.Filter{}
	filterQuals := map[string]string{
		"spread_level": "spread-level",
		"state":        "state",
		"strategy":     "strategy",
	}
	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			value := getQualsValueByColumn(quals, columnName, "string")
			if val, ok := value.(string); ok && val != "" {
				filters = append(filters, types.Filter{
					Name:   aws.String(filterName),
					Values: []string{val},
				})
			}
		}
	}
	return filters
}
