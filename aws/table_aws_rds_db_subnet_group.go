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

func tableAwsRDSDBSubnetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_subnet_group",
		Description: "AWS RDS DB Subnet Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBSubnetGroupNotFoundFault"}),
			},
			Hydrate: getRDSDBSubnetGroup,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBSubnetGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBSubnetGroups,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBSubnetGroups"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getRDSDBSubnetGroupTags,
				Tags: map[string]string{"service": "rds", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides the description of the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroupDescription"),
			},
			{
				Name:        "status",
				Description: "Provides the status of the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubnetGroupStatus"),
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VpcId of the DB subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnets",
				Description: "A list of Subnet elements.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB subnet group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRDSDBSubnetGroupTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRDSDBSubnetGroupTags,
				Transform:   transform.From(getRDSDBSubnetGroupTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBSubnetGroupArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBSubnetGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_subnet_group.listRDSDBSubnetGroups", "connection_error", err)
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

	input := &rds.DescribeDBSubnetGroupsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	paginator := rds.NewDescribeDBSubnetGroupsPaginator(svc, input, func(o *rds.DescribeDBSubnetGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_subnet_group.listRDSDBSubnetGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBSubnetGroups {
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

func getRDSDBSubnetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_subnet_group.getRDSDBSubnetGroup", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: aws.String(name),
	}

	op, err := svc.DescribeDBSubnetGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_subnet_group.getRDSDBSubnetGroup", "api_error", err)
		return nil, err
	}

	if op.DBSubnetGroups != nil && len(op.DBSubnetGroups) > 0 {
		return op.DBSubnetGroups[0], nil
	}
	return nil, nil
}

func getRDSDBSubnetGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	dbSubnetGroup := h.Item.(types.DBSubnetGroup)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_subnet_group.getRDSDBSubnetGroupTags", "connection_error", err)
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: dbSubnetGroup.DBSubnetGroupArn,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_subnet_group.getRDSDBSubnetGroupTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getRDSDBSubnetGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbSubnetGroup := d.HydrateItem.(*rds.ListTagsForResourceOutput)

	if dbSubnetGroup.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbSubnetGroup.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
