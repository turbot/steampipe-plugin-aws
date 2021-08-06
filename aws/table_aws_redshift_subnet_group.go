package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsRedshiftSubnetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_subnet_group",
		Description: "AWS Redshift Subnet Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("cluster_subnet_group_name"),
			ShouldIgnoreError: isNotFoundError([]string{"ClusterSubnetGroupNotFoundFault"}),
			Hydrate:           getRedshiftSubnetGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftSubnetGroup,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_subnet_group_name",
				Description: "The name of the cluster subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_group_status",
				Description: "The status of the cluster subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the cluster subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The VPC ID of the cluster subnet group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnets",
				Description: "A list of the VPC Subnet elements.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the subnet group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRedshiftSubnetGroup,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterSubnetGroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(redshiftSubnetGroupTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRedshiftSubnetGroupAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listRedshiftSubnetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listRedshiftSubnetGroup")

	// Create Session
	svc, err := RedshiftService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeClusterSubnetGroupsPages(
		&redshift.DescribeClusterSubnetGroupsInput{},
		func(page *redshift.DescribeClusterSubnetGroupsOutput, isLast bool) bool {
			for _, subnetGroup := range page.ClusterSubnetGroups {
				d.StreamListItem(ctx, subnetGroup)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRedshiftSubnetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	clusterSubnetGroupName := d.KeyColumnQuals["cluster_subnet_group_name"].GetStringValue()

	// Create service
	svc, err := RedshiftService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &redshift.DescribeClusterSubnetGroupsInput{
		ClusterSubnetGroupName: aws.String(clusterSubnetGroupName),
	}

	op, err := svc.DescribeClusterSubnetGroups(params)
	if err != nil {
		return nil, err
	}

	if op.ClusterSubnetGroups != nil && len(op.ClusterSubnetGroups) > 0 {
		return op.ClusterSubnetGroups[0], nil
	}
	return nil, nil
}

func getRedshiftSubnetGroupAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRedshiftSubnetGroupAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*redshift.ClusterSubnetGroup)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":redshift:" + region + ":" + commonColumnData.AccountId + ":subnetgroup:" + *data.ClusterSubnetGroupName

	// Get data for turbot defined properties
	akas := []string{arn}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func redshiftSubnetGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*redshift.ClusterSubnetGroup)
	if data.Tags == nil {
		return nil, nil
	}

	// Get the resource tags
	var turbotTagsMap map[string]string
	if data.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}