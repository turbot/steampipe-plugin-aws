package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/memorydb"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINATION

func tableAwsMemoryDBCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_memorydb_cluster",
		Description: "AWS MemoryDB Cluster",
		List: &plugin.ListConfig{
			Hydrate: listAwsMemoryDBClusters,
			Tags:    map[string]string{"service": "memorydb", "action": "DescribeClusters"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ClusterNotFoundFault"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_MEMORY_DB_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The user-supplied name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "status",
				Description: "The status of the cluster (e.g., Available, Updating, Creating).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tls_enabled",
				Description: "Indicates if In-transit encryption is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("TLSEnabled"),
			},
			{
				Name:        "acl_name",
				Description: "The name of the Access Control List associated with this cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ACLName"),
			},
			{
				Name:        "auto_minor_version_upgrade",
				Description: "When set to true, the cluster will automatically receive minor engine version upgrades after launch.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "availability_mode",
				Description: "Indicates if the cluster has a Multi-AZ configuration (multiaz) or not (singleaz).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AvailabilityMode"),
			},
			{
				Name:        "data_tiering",
				Description: "Enables data tiering, supported only for clusters using the r6gd node type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_patch_version",
				Description: "The Redis engine patch version used by the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "The Redis engine version used by the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the KMS key used to encrypt the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "maintenance_window",
				Description: "Specifies the weekly time range during which maintenance on the cluster is performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type",
				Description: "The cluster's node type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_shards",
				Description: "The number of shards in the cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "parameter_group_name",
				Description: "The name of the parameter group used by the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameter_group_status",
				Description: "The status of the parameter group used by the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_retention_limit",
				Description: "The number of days for which MemoryDB retains automatic snapshots before deleting them.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "snapshot_window",
				Description: "The daily time range (in UTC) during which MemoryDB begins taking a daily snapshot of your shard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sns_topic_arn",
				Description: "The Amazon Resource Name (ARN) of the SNS notification topic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sns_topic_status",
				Description: "The SNS topic must be in Active status to receive notifications.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_group_name",
				Description: "The name of the subnet group used by the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pending_updates",
				Description: "A group of settings that are currently being applied.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PendingUpdates"),
			},
			{
				Name:        "cluster_endpoint",
				Description: "The cluster's configuration endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_groups",
				Description: "A list of security groups used by the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroups"),
			},
			{
				Name:        "shards",
				Description: "A list of shards that are members of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Shards"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

// // LIST FUNCTION
func listAwsMemoryDBClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	svc, err := MemoryDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_memorydb_cluster.listAwsMemoryDBClusters", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = int32(20)
		} else {
			maxLimit = limit
		}
	}

	// Page size must be greater than 0 and less than or equal to 1000
	input := &memorydb.DescribeClustersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := memorydb.NewDescribeClustersPaginator(svc, input, func(o *memorydb.DescribeClustersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_memorydb_cluster.listAwsMemoryDBClusters", "connection_error", err)
			return nil, err
		}

		for _, item := range output.Clusters {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
