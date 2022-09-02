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

func tableAwsRDSDBParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_parameter_group",
		Description: "AWS RDS DB Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"DBParameterGroupNotFound"}),
			},
			Hydrate: getRDSDBParameterGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBParameterGroups,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the DB parameter group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBParameterGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB parameter group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBParameterGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides the customer-specified description for this DB parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_parameter_group_family",
				Description: "The name of the DB parameter group family that this DB parameter group is compatible with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBParameterGroupFamily"),
			},
			{
				Name:        "parameters",
				Description: "A list of detailed parameter for a particular DB parameter group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRDSParameterGroupParameters,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB parameter group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRDSParameterGroupTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRDSParameterGroupTags,
				Transform:   transform.From(getRDSDBParameterGroupTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBParameterGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBParameterGroupArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBParameterGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRDSDBParameterGroups")

	// Create Session
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &rds.DescribeDBParameterGroupsInput{
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

	// List call
	err = svc.DescribeDBParameterGroupsPages(
		input,
		func(page *rds.DescribeDBParameterGroupsOutput, isLast bool) bool {
			for _, dbParameterGroup := range page.DBParameterGroups {
				d.StreamListItem(ctx, dbParameterGroup)

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

func getRDSDBParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBParameterGroupsInput{
		DBParameterGroupName: aws.String(name),
	}

	op, err := svc.DescribeDBParameterGroups(params)
	if err != nil {
		return nil, err
	}

	if op.DBParameterGroups != nil && len(op.DBParameterGroups) > 0 {
		return op.DBParameterGroups[0], nil
	}
	return nil, nil
}

func getRDSParameterGroupParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRDSParameterGroupParameters")

	dbParameterGroup := h.Item.(*rds.DBParameterGroup)

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	var items []*rds.Parameter
	err = svc.DescribeDBParametersPages(
		&rds.DescribeDBParametersInput{
			DBParameterGroupName: dbParameterGroup.DBParameterGroupName,
		},
		func(page *rds.DescribeDBParametersOutput, isLast bool) bool {
			items = append(items, page.Parameters...)
			return !isLast
		},
	)

	return items, err
}

func getRDSParameterGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRDSParameterGroupTags")

	dbParameterGroup := h.Item.(*rds.DBParameterGroup)

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: dbParameterGroup.DBParameterGroupArn,
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getRDSDBParameterGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbParameterGroup := d.HydrateItem.(*rds.ListTagsForResourceOutput)

	if dbParameterGroup.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbParameterGroup.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
