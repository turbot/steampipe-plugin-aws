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

func tableAwsRDSDBClusterParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_cluster_parameter_group",
		Description: "AWS RDS DB Cluster Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"DBParameterGroupNotFound"}),
			ItemFromKey:       clusterParameterGroupNameFromKey,
			Hydrate:           getRDSDBClusterParameterGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBClusterParameterGroups,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the DB cluster parameter group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterParameterGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB cluster parameter group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterParameterGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides the customer-specified description for this DB cluster parameter group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_parameter_group_family",
				Description: "The name of the DB parameter group family that this DB cluster parameter group is compatible with",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBParameterGroupFamily"),
			},
			{
				Name:        "parameters",
				Description: "A list of detailed parameter for a particular DB Cluster parameter group",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSClusterParameterGroupParameters,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tag_list",
				Description: "A list of tags attached to the DB Cluster parameter group",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSClusterParameterGroupTags,
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

//// BUILD HYDRATE INPUT

func clusterParameterGroupNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &rds.DBClusterParameterGroup{
		DBClusterParameterGroupName: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listRDSDBClusterParameterGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listRDSDBClusterParameterGroups", "AWS_REGION", defaultRegion)

	// Create Session
	svc, err := RDSService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeDBClusterParameterGroupsPages(
		&rds.DescribeDBClusterParameterGroupsInput{},
		func(page *rds.DescribeDBClusterParameterGroupsOutput, isLast bool) bool {
			for _, dbClusterParameterGroup := range page.DBClusterParameterGroups {
				d.StreamListItem(ctx, dbClusterParameterGroup)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBClusterParameterGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	dbClusterParameterGroup := h.Item.(*rds.DBClusterParameterGroup)

	// Create service
	svc, err := RDSService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: aws.String(*dbClusterParameterGroup.DBClusterParameterGroupName),
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
	defaultRegion := GetDefaultRegion()
	dbClusterParameterGroup := h.Item.(*rds.DBClusterParameterGroup)

	// Create service
	svc, err := RDSService(ctx, d, defaultRegion)
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

	return items, nil
}

func getAwsRDSClusterParameterGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRDSClusterParameterGroupTags")
	defaultRegion := GetDefaultRegion()
	dbClusterParameterGroup := h.Item.(*rds.DBClusterParameterGroup)

	// Create service
	svc, err := RDSService(ctx, d, defaultRegion)
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
