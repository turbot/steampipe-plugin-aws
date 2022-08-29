package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsRDSDBOptionGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_option_group",
		Description: "AWS RDS DB Option Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"OptionGroupNotFoundFault"}),
			},
			Hydrate: getRDSDBOptionGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBOptionGroups,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "engine_name", Require: plugin.Optional},
				{Name: "major_engine_version", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the option group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the option group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides a description of the option group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupDescription"),
			},
			{
				Name:        "allows_vpc_and_non_vpc_instance_memberships",
				Description: "Specifies whether this option group can be applied to both VPC and non-VPC instances.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "engine_name",
				Description: "Indicates the name of the engine that this option group can be applied to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "major_engine_version",
				Description: "Indicates the major engine version associated with this option group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "Indicates the ID of the VPC, option group can be applied.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "options",
				Description: "Indicates what options are available in the option group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the option group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSOptionGroupTags,
				Transform:   transform.FromField("TagList"),
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

//// LIST FUNCTION

func listRDSDBOptionGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRDSDBOptionGroups")

	// Create Session
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &rds.DescribeOptionGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxRecords {
			if *limit < 20 {
				input.MaxRecords = aws.Int64(20)
			} else {
				input.MaxRecords = limit
			}
		}
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["engine_name"] != nil {
		input.EngineName = aws.String(equalQuals["engine_name"].GetStringValue())
	}

	// We must need to pass engine name if we are passing the major engine version
	if equalQuals["engine_name"] != nil && equalQuals["major_engine_version"] != nil {
		input.MajorEngineVersion = aws.String(equalQuals["major_engine_version"].GetStringValue())
	}

	// List call
	err = svc.DescribeOptionGroupsPages(
		input,
		func(page *rds.DescribeOptionGroupsOutput, isLast bool) bool {
			for _, optionGroup := range page.OptionGroupsList {
				d.StreamListItem(ctx, optionGroup)

				// Check if context has been cancelled or if the limit has been reached (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
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

func getRDSDBOptionGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeOptionGroupsInput{
		OptionGroupName: aws.String(name),
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

	optionGroup := h.Item.(*rds.OptionGroup)

	// Create service
	svc, err := RDSService(ctx, d)
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

//// TRANSFORM FUNCTIONS

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
