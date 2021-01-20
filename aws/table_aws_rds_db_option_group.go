package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsRDSDBOptionGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_option_group",
		Description: "AWS RDS DB Option Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"OptionGroupNotFoundFault"}),
			ItemFromKey:       optionGroupNameFromKey,
			Hydrate:           getRDSDBOptionGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBOptionGroups,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the option group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the option group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides a description of the option group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupDescription"),
			},
			{
				Name:        "allows_vpc_and_non_vpc_instance_memberships",
				Description: "Specifies whether this option group can be applied to both VPC and non-VPC instances",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "engine_name",
				Description: "Indicates the name of the engine that this option group can be applied to",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "major_engine_version",
				Description: "Indicates the major engine version associated with this option group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "Indicates the ID of the VPC, option group can be applied",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "options",
				Description: "Indicates what options are available in the option group",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tag_list",
				Description: "A list of tags attached to the option group",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSOptionGroupTags,
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSOptionGroupTags,
				Transform:   transform.From(getRDSDBOptionGroupTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OptionGroupArn").Transform(arnToAkas),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func optionGroupNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &rds.OptionGroup{
		OptionGroupName: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listRDSDBOptionGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listRDSDBOptionGroups", "AWS_REGION", defaultRegion)

	// Create Session
	svc, err := RDSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeOptionGroupsPages(
		&rds.DescribeOptionGroupsInput{},
		func(page *rds.DescribeOptionGroupsOutput, isLast bool) bool {
			for _, optionGroup := range page.OptionGroupsList {
				d.StreamListItem(ctx, optionGroup)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBOptionGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	optionGroup := h.Item.(*rds.OptionGroup)

	// Create service
	svc, err := RDSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeOptionGroupsInput{
		OptionGroupName: aws.String(*optionGroup.OptionGroupName),
	}

	op, err := svc.DescribeOptionGroups(params)
	if err != nil {
		return nil, err
	}

	if op.OptionGroupsList != nil && len(op.OptionGroupsList) > 0 {
		return op.OptionGroupsList[0], nil
	}
	return nil, nil
}

func getAwsRDSOptionGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRDSOptionGroupTags")
	defaultRegion := GetDefaultRegion()
	optionGroup := h.Item.(*rds.OptionGroup)

	// Create service
	svc, err := RDSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: optionGroup.OptionGroupArn,
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS ////

func getRDSDBOptionGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	optionGroup := d.HydrateItem.(*rds.ListTagsForResourceOutput)

	if optionGroup.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range optionGroup.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
