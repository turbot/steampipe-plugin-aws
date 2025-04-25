package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBClusterParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_cluster_parameter_group",
		Description: "AWS RDS DB Cluster Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBParameterGroupNotFound"}),
			},
			Hydrate: getRDSDBClusterParameterGroup,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBClusterParameterGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBClusterParameterGroups,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBClusterParameterGroups"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsRDSClusterParameterGroupParameters,
				Tags: map[string]string{"service": "rds", "action": "DescribeDBClusterParameters"},
			},
			{
				Func: getAwsRDSClusterParameterGroupTags,
				Tags: map[string]string{"service": "rds", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
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

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.listRDSDBClusterParameterGroups", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	input := &rds.DescribeDBClusterParameterGroupsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	paginator := rds.NewDescribeDBClusterParameterGroupsPaginator(svc, input, func(o *rds.DescribeDBClusterParameterGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.listRDSDBClusterParameterGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBClusterParameterGroups {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBClusterParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.getRDSDBClusterParameterGroup", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: aws.String(name),
	}

	op, err := svc.DescribeDBClusterParameterGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.getRDSDBClusterParameterGroup", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.DBClusterParameterGroups) > 0 {
		return op.DBClusterParameterGroups[0], nil
	}
	return nil, nil
}

func getAwsRDSClusterParameterGroupParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	dbClusterParameterGroup := h.Item.(types.DBClusterParameterGroup)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.getAwsRDSClusterParameterGroupParameters", "connection_error", err)
		return nil, err
	}

	input := &rds.DescribeDBClusterParametersInput{
		DBClusterParameterGroupName: dbClusterParameterGroup.DBClusterParameterGroupName,
	}

	paginator := rds.NewDescribeDBClusterParametersPaginator(svc, input, func(o *rds.DescribeDBClusterParametersPaginatorOptions) {
		o.Limit = int32(100)
		o.StopOnDuplicateToken = true
	})

	var parameters []types.Parameter

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.getAwsRDSClusterParameterGroupParameters", "api_error", err)
			return nil, err
		}
		parameters = append(parameters, output.Parameters...)
	}

	return parameters, err
}

func getAwsRDSClusterParameterGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	dbClusterParameterGroup := h.Item.(types.DBClusterParameterGroup)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.getAwsRDSClusterParameterGroupTags", "connection_error", err)
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: dbClusterParameterGroup.DBClusterParameterGroupArn,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_parameter_group.getAwsRDSClusterParameterGroupTags", "api_error", err)
		return nil, err
	}
	if len(op.TagList) > 0 {
		return op, nil
	}

	return nil, nil
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
