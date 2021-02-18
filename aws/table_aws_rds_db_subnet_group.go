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

func tableAwsRDSDBSubnetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_subnet_group",
		Description: "AWS RDS DB Subnet Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"DBSubnetGroupNotFoundFault"}),
			ItemFromKey:       subnetGroupNameFromKey,
			Hydrate:           getRDSDBSubnetGroup,
		},
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listRDSDBSubnetGroups,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides the description of the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroupDescription"),
			},
			{
				Name:        "status",
				Description: "Provides the status of the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubnetGroupStatus"),
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VpcId of the DB subnet group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnets",
				Description: "A list of Subnet elements",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB subnet group",
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

//// BUILD HYDRATE INPUT

func subnetGroupNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &rds.DBSubnetGroup{
		DBSubnetGroupName: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listRDSDBSubnetGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listRDSDBSubnetGroups", "AWS_REGION", region)

	// Create Session
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeDBSubnetGroupsPages(
		&rds.DescribeDBSubnetGroupsInput{},
		func(page *rds.DescribeDBSubnetGroupsOutput, isLast bool) bool {
			for _, dbSubnetGroup := range page.DBSubnetGroups {
				d.StreamListItem(ctx, dbSubnetGroup)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBSubnetGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	dbSubnetGroup := h.Item.(*rds.DBSubnetGroup)

	// Create service
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: aws.String(*dbSubnetGroup.DBSubnetGroupName),
	}

	op, err := svc.DescribeDBSubnetGroups(params)
	if err != nil {
		return nil, err
	}

	if op.DBSubnetGroups != nil && len(op.DBSubnetGroups) > 0 {
		return op.DBSubnetGroups[0], nil
	}
	return nil, nil
}

func getRDSDBSubnetGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRDSDBSubnetGroupTags")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	dbSubnetGroup := h.Item.(*rds.DBSubnetGroup)

	// Create service
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: dbSubnetGroup.DBSubnetGroupArn,
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS ////

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
