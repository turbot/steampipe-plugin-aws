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

func tableAwsRDSDBCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_cluster",
		Description: "AWS RDS DB Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_cluster_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBClusterNotFoundFault"}),
			},
			Hydrate: getRDSDBCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBClusters,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "clone_group_id", Require: plugin.Optional},
				{Name: "engine", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_cluster_identifier",
				Description: "The friendly name to identify the DB Cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB Cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterArn"),
			},
			{
				Name:        "status",
				Description: "Specifies the status of this DB Cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_id",
				Description: "The AWS Region-unique, immutable identifier for the DB cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DbClusterResourceId"),
			},
			{
				Name:        "create_time",
				Description: "Specifies the time when the DB cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ClusterCreateTime"),
			},
			{
				Name:        "activity_stream_kinesis_stream_name",
				Description: "The name of the Amazon Kinesis data stream used for the database activity stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activity_stream_kms_key_id",
				Description: "The AWS KMS key identifier used for encrypting messages in the database activity stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activity_stream_mode",
				Description: "The mode of the database activity stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activity_stream_status",
				Description: "The status of the database activity stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allocated_storage",
				Description: "Specifies the allocated storage size in gibibytes (GiB).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "backtrack_consumed_change_records",
				Description: "The number of change records stored for Backtrack.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "backtrack_window",
				Description: "The target backtrack window, in seconds.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "backup_retention_period",
				Description: "Specifies the number of days for which automatic DB snapshots are retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "capacity",
				Description: "The current capacity of an Aurora Serverless DB cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "character_set_name",
				Description: "Specifies the name of the character set that this cluster is associated with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "clone_group_id",
				Description: "Identifies the clone group to which the DB cluster is associated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "copy_tags_to_snapshot",
				Description: "Specifies whether tags are copied from the DB cluster to snapshots of the DB cluster, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cross_account_clone",
				Description: "Specifies whether the DB cluster is a clone of a DB cluster owned by a different AWS account, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "database_name",
				Description: "Contains the name of the initial database of this DB cluster that was provided at create time.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_cluster_parameter_group",
				Description: "Specifies the name of the DB cluster parameter group for the DB cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterParameterGroup"),
			},
			{
				Name:        "db_subnet_group",
				Description: "Specifies information on the subnet group associated with the DB cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup"),
			},
			{
				Name:        "deletion_protection",
				Description: "Specifies whether the DB cluster has deletion protection enabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "earliest_backtrack_time",
				Description: "The earliest time to which a DB cluster can be backtracked.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "earliest_restorable_time",
				Description: "The earliest time to which a database can be restored with point-in-time restore.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "endpoint",
				Description: "Specifies the connection endpoint for the primary instance of the DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine",
				Description: "The name of the database engine to be used for this DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_mode",
				Description: "The DB engine mode of the DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Indicates the database engine version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "global_write_forwarding_requested",
				Description: "Specifies whether you have requested to enable write forwarding for a secondary cluster in an Aurora global database, or not.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "global_write_forwarding_status",
				Description: "Specifies whether a secondary cluster in an Aurora global database has write forwarding enabled, or not.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hosted_zone_id",
				Description: "Specifies the ID that Amazon Route 53 assigns when you create a hosted zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "http_endpoint_enabled",
				Description: "Specifies whether the HTTP endpoint for an Aurora Serverless DB cluster is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "iam_database_authentication_enabled",
				Description: "Specifies whether the the mapping of AWS IAM accounts to database accounts is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IAMDatabaseAuthenticationEnabled"),
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS KMS key identifier for the encrypted DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_restorable_time",
				Description: "Specifies the latest time to which a database can be restored with point-in-time restore.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "master_user_name",
				Description: "Contains the master username for the DB cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "multi_az",
				Description: "Specifies whether the DB cluster has instances in multiple Availability Zones, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("MultiAZ"),
			},
			{
				Name:        "percent_progress",
				Description: "Specifies the progress of the operation as a percentage.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "Specifies the port that the database engine is listening on.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "preferred_backup_window",
				Description: "Specifies the daily time range during which automated backups are created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "Specifies the weekly time range during which system maintenance can occur",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reader_endpoint",
				Description: "The reader endpoint for the DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the DB cluster is encrypted, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "associated_roles",
				Description: "A list of AWS IAM roles that are associated with the DB cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "availability_zones",
				Description: "A list of Availability Zones (AZs) where instances in the DB cluster can be created.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "custom_endpoints",
				Description: "A list of all custom endpoints associated with the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "members",
				Description: "A list of instances that make up the DB cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "option_group_memberships",
				Description: "A list of option group memberships for this DB cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterOptionGroupMemberships"),
			},
			{
				Name:        "domain_memberships",
				Description: "A list of Active Directory Domain membership records associated with the DB cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "enabled_cloudwatch_logs_exports",
				Description: "A list of log types that this DB cluster is configured to export to CloudWatch Logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pending_maintenance_actions",
				Description: "A list that provides details about the pending maintenance actions for the resource.",
				Hydrate:     getRDSDBClusterPendingMaintenanceAction,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "read_replica_identifiers",
				Description: "A list of identifiers of the read replicas associated with this DB cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_security_groups",
				Description: "A list of VPC security groups that the DB cluster belongs to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB Cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getRDSDBClusterTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster.listRDSDBClusters", "connection_error", err)
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

	input := &rds.DescribeDBClustersInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	filters := buildRdsDbClusterFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := rds.NewDescribeDBClustersPaginator(svc, input, func(o *rds.DescribeDBClustersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_cluster.listRDSDBClusters", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBClusters {
			// The DescribeDBClusters API returns non-RDS DB clusters as well,
			// but we only want RDS clusters here, even if the 'engine' qual
			// isn't passed in.
			if isSuppportedRDSEngine(*items.Engine) {
				d.StreamListItem(ctx, items)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	dbClusterIdentifier := d.EqualsQuals["db_cluster_identifier"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster.getRDSDBCluster", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(dbClusterIdentifier),
	}

	op, err := svc.DescribeDBClusters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster.getRDSDBCluster", "api_error", err)
		return nil, err
	}

	if op.DBClusters != nil && len(op.DBClusters) > 0 {

		cluster := op.DBClusters[0]
		if isSuppportedRDSEngine(*cluster.Engine) {
			return cluster, nil
		}
	}
	return nil, nil
}

func getRDSDBClusterPendingMaintenanceAction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dbClusterIdentifier := *h.Item.(types.DBCluster).DBClusterIdentifier

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster.getRDSDBClusterPendingMaintenanceAction", "connection_error", err)
		return nil, err
	}

	filter := &types.Filter{
		Name:   aws.String("db-cluster-id"),
		Values: []string{dbClusterIdentifier},
	}
	params := &rds.DescribePendingMaintenanceActionsInput{
		Filters: []types.Filter{*filter},
	}

	op, err := svc.DescribePendingMaintenanceActions(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster.getRDSDBClusterPendingMaintenanceAction", "api_error", err)
		return nil, err
	}

	if len(op.PendingMaintenanceActions) > 0 {
		return op.PendingMaintenanceActions, nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getRDSDBClusterTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbCluster := d.HydrateItem.(types.DBCluster)

	if dbCluster.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbCluster.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

//// UTILITY FUNCTIONS

// build rdds db cluster list call input filter
func buildRdsDbClusterFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)
	filterQuals := map[string]string{
		"clone_group_id": "clone-group-id",
		"engine":         "engine",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
