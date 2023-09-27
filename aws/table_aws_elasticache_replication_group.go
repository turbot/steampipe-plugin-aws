package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	elasticachev1 "github.com/aws/aws-sdk-go/service/elasticache"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElastiCacheReplicationGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_replication_group",
		Description: "AWS ElastiCache Replication Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("replication_group_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ReplicationGroupNotFoundFault", "InvalidParameterValue"}),
			},
			Hydrate: getElastiCacheReplicationGroup,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeReplicationGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheReplicationGroups,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeReplicationGroups"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elasticachev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "replication_group_id",
				Description: "The identifier for the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the replication group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "description",
				Description: "The user supplied description of the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "at_rest_encryption_enabled",
				Description: "A flag that enables encryption at-rest when set to true.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the KMS key used to encrypt the disk in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auth_token_enabled",
				Description: "A flag that enables using an AuthToken (password) when issuing Redis commands.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "auth_token_last_modified_date",
				Description: "The date when the auth token was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "automatic_failover",
				Description: "Indicates the status of automatic failover for this Redis replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_node_type",
				Description: "The name of the compute and memory capacity node type for each node in the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_enabled",
				Description: "A flag indicating whether or not this replication group is cluster enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "multi_az",
				Description: "A flag indicating if you have Multi-AZ enabled to enhance fault tolerance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MultiAZ"),
			},
			{
				Name:        "snapshot_retention_limit",
				Description: "The number of days for which ElastiCache retains automatic cluster snapshots before deleting them.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "snapshot_window",
				Description: "The daily time range (in UTC) during which ElastiCache begins taking a daily snapshot of your node group (shard).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshotting_cluster_id",
				Description: "The cluster ID that is used as the daily snapshot source for the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current state of this replication group - creating, available, modifying, deleting, create-failed, snapshotting.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_encryption_enabled",
				Description: "A flag that enables in-transit encryption when set to true.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "configuration_endpoint",
				Description: "The configuration endpoint for this replication group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "global_replication_group_info",
				Description: "The name of the Global Datastore and role of this replication group in the Global Datastore.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "member_clusters",
				Description: "The names of all the cache clusters that are part of this replication group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "member_clusters_outpost_arns",
				Description: "The outpost ARNs of the replication group's member clusters.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "node_groups",
				Description: "A list of node groups in this replication group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pending_modified_values",
				Description: "A group of settings to be applied to the replication group, either immediately or during the next maintenance window.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_group_ids",
				Description: "The list of user group IDs that have access to the replication group.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReplicationGroupId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listElastiCacheReplicationGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_replication_group.listElastiCacheReplicationGroups", "get_client_error", err)
		return nil, err
	}

	input := &elasticache.DescribeReplicationGroupsInput{
		MaxRecords: aws.Int32(100),
	}

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

	paginator := elasticache.NewDescribeReplicationGroupsPaginator(svc, input, func(o *elasticache.DescribeReplicationGroupsPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_replication_group.listElastiCacheParameterGroup", "api_error", err)
			return nil, err
		}

		for _, replicationGroup := range output.ReplicationGroups {
			d.StreamListItem(ctx, replicationGroup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getElastiCacheReplicationGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_replication_group.getElastiCacheReplicationGroup", "get_client_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	replicationGroupId := quals["replication_group_id"].GetStringValue()

	params := &elasticache.DescribeReplicationGroupsInput{
		ReplicationGroupId: aws.String(replicationGroupId),
	}

	op, err := svc.DescribeReplicationGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_replication_group.getElastiCacheReplicationGroup", "api_error", err)
		return nil, err
	}

	if op.ReplicationGroups != nil && len(op.ReplicationGroups) > 0 {
		return op.ReplicationGroups[0], nil
	}
	return nil, nil
}
