package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsRedshiftSubnetGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_subnet_group",
		Description: "AWS Redshift Subnet Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_subnet_group_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ClusterSubnetGroupNotFoundFault"}),
			},
			Hydrate: getRedshiftSubnetGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftSubnetGroups,
		},
		GetMatrixItemFunc: BuildRegionList,
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

func listRedshiftSubnetGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_subnet_group.listRedshiftSubnetGroups", "connection_error", err)
		return nil, err
	}

	input := &redshift.DescribeClusterSubnetGroupsInput{
		MaxRecords: aws.Int32(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxRecords {
			if limit < 20 {
				input.MaxRecords = aws.Int32(20)
			} else {
				input.MaxRecords = aws.Int32(limit)
			}
		}
	}

	// List call
	paginator := redshift.NewDescribeClusterSubnetGroupsPaginator(svc, input, func(o *redshift.DescribeClusterSubnetGroupsPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_redshift_subnet_group.listRedshiftSubnetGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.ClusterSubnetGroups {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedshiftSubnetGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	clusterSubnetGroupName := d.KeyColumnQuals["cluster_subnet_group_name"].GetStringValue()

	// Return nil, if no input provided
	if clusterSubnetGroupName == "" {
		return nil, nil
	}

	// Create service
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_subnet_group.getRedshiftSubnetGroup", "connection_error", err)
		return nil, err
	}

	params := &redshift.DescribeClusterSubnetGroupsInput{
		ClusterSubnetGroupName: aws.String(clusterSubnetGroupName),
	}

	op, err := svc.DescribeClusterSubnetGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_subnet_group.getRedshiftSubnetGroup", "api_error", err)
		return nil, err
	}

	if len(op.ClusterSubnetGroups) > 0 {
		return op.ClusterSubnetGroups[0], nil
	}

	return nil, nil
}

func getRedshiftSubnetGroupAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(types.ClusterSubnetGroup)

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_subnet_group.getRedshiftSubnetGroupAkas", "getCommonColumnsCached_error", err)
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
	clusterSubnetGroup := d.HydrateItem.(types.ClusterSubnetGroup)

	// Get the resource tags
	var turbotTagsMap map[string]string
	if len(clusterSubnetGroup.Tags) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range clusterSubnetGroup.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
