package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_parameter_group",
		Description: "AWS RDS DB Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBParameterGroupNotFound"}),
			},
			Hydrate: getRDSDBParameterGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBParameterGroups,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
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

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_parameter_group.listRDSDBParameterGroups", "connection_error", err)
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

	input := &rds.DescribeDBParameterGroupsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	paginator := rds.NewDescribeDBParameterGroupsPaginator(svc, input, func(o *rds.DescribeDBParameterGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_parameter_group.listRDSDBParameterGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBParameterGroups {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_parameter_group.getRDSDBParameterGroup", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBParameterGroupsInput{
		DBParameterGroupName: aws.String(name),
	}

	op, err := svc.DescribeDBParameterGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_parameter_group.getRDSDBParameterGroup", "api_error", err)
		return nil, err
	}

	if op.DBParameterGroups != nil && len(op.DBParameterGroups) > 0 {
		return op.DBParameterGroups[0], nil
	}
	return nil, nil
}

func getRDSParameterGroupParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	dbParameterGroup := h.Item.(types.DBParameterGroup)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_parameter_group.getRDSParameterGroupParameters", "connection_error", err)
		return nil, err
	}

	input := &rds.DescribeDBParametersInput{
		DBParameterGroupName: dbParameterGroup.DBParameterGroupName,
	}

	var items []types.Parameter

	paginator := rds.NewDescribeDBParametersPaginator(svc, input, func(o *rds.DescribeDBParametersPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_parameter_group.getRDSParameterGroupParameters", "api_error", err)
			return nil, err
		}
		items = append(items, output.Parameters...)
	}

	return items, err
}

func getRDSParameterGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dbParameterGroup := h.Item.(types.DBParameterGroup)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_parameter_group.getRDSParameterGroupTags", "connection_error", err)
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: dbParameterGroup.DBParameterGroupArn,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_parameter_group.getRDSParameterGroupTags", "api_error", err)
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
