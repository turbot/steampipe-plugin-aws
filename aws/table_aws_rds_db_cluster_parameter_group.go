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

func tableAwsRDSDBClusterParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_cluster_parameter_group",
		Description: "AWS RDS DB Cluster Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"DBParameterGroupNotFound"}),
			},
			Hydrate: getRDSDBClusterParameterGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBClusterParameterGroups,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the DB cluster parameter group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterParameterGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB cluster parameter group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterParameterGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides the customer-specified description for this DB cluster parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_parameter_group_family",
				Description: "The name of the DB parameter group family that this DB cluster parameter group is compatible with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBParameterGroupFamily"),
			},
			{
				Name:        "parameters",
				Description: "A list of detailed parameter for a particular DB Cluster parameter group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSClusterParameterGroupParameters,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB Cluster parameter group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSClusterParameterGroupTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSClusterParameterGroupTags,
				Transform:   transform.From(getRDSDBClusterParameterGroupTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterParameterGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterParameterGroupArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBClusterParameterGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRDSDBClusterParameterGroups")

	// Create Session
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &rds.DescribeDBClusterParameterGroupsInput{
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
	err = svc.DescribeDBClusterParameterGroupsPages(
		input,
		func(page *rds.DescribeDBClusterParameterGroupsOutput, isLast bool) bool {
			for _, dbClusterParameterGroup := range page.DBClusterParameterGroups {
				d.StreamListItem(ctx, dbClusterParameterGroup)

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

func getRDSDBClusterParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: aws.String(name),
	}

	op, err := svc.DescribeDBClusterParameterGroups(params)
	if err != nil {
		return nil, err
	}

	if op.DBClusterParameterGroups != nil && len(op.DBClusterParameterGroups) > 0 {
		return op.DBClusterParameterGroups[0], nil
	}
	return nil, nil
}

func getAwsRDSClusterParameterGroupParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRDSClusterParameterGroupParameters")

	dbClusterParameterGroup := h.Item.(*rds.DBClusterParameterGroup)

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	var items []*rds.Parameter
	err = svc.DescribeDBClusterParametersPages(
		&rds.DescribeDBClusterParametersInput{
			DBClusterParameterGroupName: dbClusterParameterGroup.DBClusterParameterGroupName,
		},
		func(page *rds.DescribeDBClusterParametersOutput, isLast bool) bool {
			items = append(items, page.Parameters...)
			return !isLast
		},
	)

	return items, err
}

func getAwsRDSClusterParameterGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRDSClusterParameterGroupTags")

	dbClusterParameterGroup := h.Item.(*rds.DBClusterParameterGroup)

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: dbClusterParameterGroup.DBClusterParameterGroupArn,
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS ////

func getRDSDBClusterParameterGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbClusterParameterGroup := d.HydrateItem.(*rds.ListTagsForResourceOutput)

	if dbClusterParameterGroup.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbClusterParameterGroup.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
