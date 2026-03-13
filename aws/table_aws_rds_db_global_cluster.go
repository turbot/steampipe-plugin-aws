package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSGlobalCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_global_cluster",
		Description: "AWS RDS Global Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("global_cluster_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"GlobalClusterNotFoundFault"}),
			},
			Hydrate: getRDSGlobalCluster,
			Tags:    map[string]string{"service": "rds", "action": "DescribeGlobalClusters"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSGlobalClusters,
			Tags:    map[string]string{"service": "rds", "action": "DescribeGlobalClusters"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getRDSGlobalClusterTags,
				Tags: map[string]string{"service": "rds", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_RDS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "global_cluster_identifier",
				Description: "Contains a user-supplied global database cluster identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "global_cluster_arn",
				Description: "The Amazon Resource Name (ARN) for the global database cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the current state of this global database cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint",
				Description: "The writer endpoint for the primary DB cluster in this global database cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Endpoint"),
			},
			{
				Name:        "global_cluster_resource_id",
				Description: "The AWS Region-unique, immutable identifier for the global database cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_name",
				Description: "The default database name within the global database cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_protection",
				Description: "The deletion protection setting for the global database cluster.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "engine",
				Description: "The Aurora database engine used by the global database cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Indicates the database engine version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_encrypted",
				Description: "The storage encryption setting for the global database cluster.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "failover_state",
				Description: "Properties for the current state of an in-process or pending switchover/failover.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "global_cluster_members",
				Description: "The list of primary and secondary clusters within the global database cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the global database cluster.",
				Hydrate:     getRDSGlobalClusterTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GlobalClusterIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRDSGlobalClusterTags,
				Transform:   transform.FromValue().Transform(rdsGlobalClusterTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("GlobalClusterArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSGlobalClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_global_cluster.listRDSGlobalClusters", "connection_error", err)
		return nil, err
	}

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

	input := &rds.DescribeGlobalClustersInput{MaxRecords: aws.Int32(maxLimit)}
	paginator := rds.NewDescribeGlobalClustersPaginator(svc, input, func(o *rds.DescribeGlobalClustersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_global_cluster.listRDSGlobalClusters", "api_error", err)
			return nil, err
		}

		for _, globalCluster := range output.GlobalClusters {
			d.StreamListItem(ctx, globalCluster)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRDSGlobalCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	globalClusterIdentifier := d.EqualsQuals["global_cluster_identifier"].GetStringValue()
	if globalClusterIdentifier == "" {
		return nil, nil
	}

	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_global_cluster.getRDSGlobalCluster", "connection_error", err)
		return nil, err
	}

	input := &rds.DescribeGlobalClustersInput{
		GlobalClusterIdentifier: aws.String(globalClusterIdentifier),
	}

	op, err := svc.DescribeGlobalClusters(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_global_cluster.getRDSGlobalCluster", "api_error", err)
		return nil, err
	}

	if len(op.GlobalClusters) > 0 {
		return op.GlobalClusters[0], nil
	}

	return nil, nil
}

func getRDSGlobalClusterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_global_cluster.getRDSGlobalClusterTags", "connection_error", err)
		return nil, err
	}

	globalCluster := h.Item.(types.GlobalCluster)
	if globalCluster.GlobalClusterArn == nil {
		return nil, nil
	}

	input := &rds.ListTagsForResourceInput{ResourceName: globalCluster.GlobalClusterArn}
	op, err := svc.ListTagsForResource(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_global_cluster.getRDSGlobalClusterTags", "api_error", err)
		return nil, err
	}

	return op.TagList, nil
}

//// TRANSFORM FUNCTION

func rdsGlobalClusterTagListToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tagList, ok := d.Value.([]types.Tag)
	if !ok {
		return nil, nil
	}

	if len(tagList) == 0 {
		return nil, nil
	}

	tagsMap := map[string]string{}
	for _, tag := range tagList {
		if tag.Key != nil && tag.Value != nil {
			tagsMap[*tag.Key] = *tag.Value
		}
	}

	return tagsMap, nil
}
