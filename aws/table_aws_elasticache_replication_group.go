package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsElastiCacheReplicationGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_replication_group",
		Description: "AWS ElastiCache Replication Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("replication_group_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ReplicationGroupNotFoundFault", "InvalidParameterValue"}),
			},
			Hydrate: getElastiCacheReplicationGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listElastiCacheReplicationGroups,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &elasticache.DescribeReplicationGroupsInput{
		MaxRecords: aws.Int64(100),
	}

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
	err = svc.DescribeReplicationGroupsPages(
		input,
		func(page *elasticache.DescribeReplicationGroupsOutput, isLast bool) bool {
			for _, replicationGroup := range page.ReplicationGroups {
				d.StreamListItem(ctx, replicationGroup)

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

//// HYDRATE FUNCTION

func getElastiCacheReplicationGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getElastiCacheReplicationGroup")

	// Create service
	svc, err := ElastiCacheService(ctx, d)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	replicationGroupId := quals["replication_group_id"].GetStringValue()

	params := &elasticache.DescribeReplicationGroupsInput{
		ReplicationGroupId: aws.String(replicationGroupId),
	}

	op, err := svc.DescribeReplicationGroups(params)
	if err != nil {
		return nil, err
	}

	if op.ReplicationGroups != nil && len(op.ReplicationGroups) > 0 {
		return op.ReplicationGroups[0], nil
	}
	return nil, nil
}
