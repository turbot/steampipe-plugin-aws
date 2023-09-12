package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"

	neptunev1 "github.com/aws/aws-sdk-go/service/neptune"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsNeptuneDBCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_neptune_db_cluster",
		Description: "AWS Neptune DB Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_cluster_identifier"),
			Hydrate:    getNeptuneDBCluster,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBClusterNotFoundFault"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listNeptuneDBClusters,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(neptunev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_cluster_identifier",
				Description: "Contains a user-supplied DB cluster identifier. This identifier is the unique key that identifies a DB cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterArn"),
			},
			{
				Name:        "status",
				Description: "Specifies the current state of this DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_create_time",
				Description: "Specifies the time when the DB cluster was created, in Universal Coordinated Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "allocated_storage",
				Description: "AllocatedStorage always returns 1, because Neptune DB cluster storage size is not fixed, but instead automatically adjusts as needed.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "automatic_restart_time",
				Description: "Time at which the DB cluster will be automatically restarted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "backup_retention_period",
				Description: "Specifies the number of days for which automatic DB snapshots are retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "clone_group_id",
				Description: "Identifies the clone group to which the DB cluster is associated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "copy_tags_to_snapshot",
				Description: "If set to true, tags are copied to any snapshot of the DB cluster that is created.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cross_account_clone",
				Description: "If set to true, the DB cluster can be cloned across accounts.",
				Type:        proto.ColumnType_BOOL,
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
				Name:        "database_name",
				Description: "Contains the name of the initial database of this DB cluster that was provided.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_cluster_resource_id",
				Description: "The Amazon Region-unique, immutable identifier for the DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_protection",
				Description: "Indicates whether or not the DB cluster has deletion protection enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "earliest_restorable_time",
				Description: "Specifies the earliest time to which a database can be restored with point-in-time restore.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "endpoint",
				Description: "Specifies the connection endpoint for the primary instance of the DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine",
				Description: "Provides the name of the database engine to be used for this DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Indicates the database engine version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hosted_zone_id",
				Description: "Specifies the ID that Amazon Route 53 assigns when you create a hosted zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_database_authentication_enabled",
				Description: "True if mapping of Amazon Identity and Access Management (IAM) accounts to database accounts is enabled, and otherwise false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IAMDatabaseAuthenticationEnabled"),
			},
			{
				Name:        "kms_key_id",
				Description: "If StorageEncrypted is true, the Amazon KMS key identifier for the encrypted DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_restorable_time",
				Description: "Specifies the latest time to which a database can be restored with point-in-time restore.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "multi_az",
				Description: "Specifies whether the DB cluster has instances in multiple Availability Zones.",
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
				Description: "Specifies the daily time range during which automated backups are created if automated backups are enabled, as determined by the BackupRetentionPeriod.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "Specifies the weekly time range during which system maintenance can occur, in Universal Coordinated Time (UTC).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reader_endpoint",
				Description: "The reader endpoint for the DB cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the DB cluster is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "associated_roles",
				Description: "Provides a list of the Amazon Identity and Access Management (IAM) roles.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "availability_zones",
				Description: "Provides the list of EC2 Availability Zones that instances in the DB cluster can be created in.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_cluster_members",
				Description: "Provides the list of instances that make up the DB cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterMembers"),
			},
			{
				Name:        "enabled_cloudwatch_logs_exports",
				Description: "A list of log types that this DB cluster is configured to export to CloudWatch Logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "read_replica_identifiers",
				Description: "Contains one or more identifiers of the Read Replicas associated with this DB cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_security_groups",
				Description: "Provides a list of VPC security groups that the DB cluster belongs to.",
				Type:        proto.ColumnType_JSON,
			},

			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNeptuneDBClusterTags,
				Transform:   transform.FromField("TagList"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNeptuneDBClusterTags,
				Transform:   transform.From(neptuneDBClusterTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listNeptuneDBClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := NeptuneClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster.listNeptuneDBClusters", "get_client_error", err)
		return nil, err
	}

	// Filter parameter is not supported yet in this SDK version so optional quals can not be implemented
	input := &neptune.DescribeDBClustersInput{
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

	paginator := neptune.NewDescribeDBClustersPaginator(svc, input, func(o *neptune.DescribeDBClustersPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_neptune_db_cluster.listNeptuneDBClusters", "api_error", err)
			return nil, err
		}

		for _, cluster := range output.DBClusters {
			// The DescribeDBClusters API returns non-Neptune DB clusters as well,
			// but we only want Neptune clusters here. The input has a Filter param
			// which can help filter out non-Neptune clusters, but as of 2022/08/15,
			// the SDK says the Filter param is not currently supported.
			// Related issue: https://github.com/aws/aws-sdk-go/issues/4515
			if *cluster.Engine == "neptune" {
				d.StreamListItem(ctx, cluster)
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

func getNeptuneDBCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identifier := d.EqualsQuals["db_cluster_identifier"].GetStringValue()

	// Create session
	svc, err := NeptuneClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster.getNeptuneDBCluster", "get_client_error", err)
		return nil, err
	}

	// Build the params
	params := &neptune.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(identifier),
	}

	// Get call
	data, err := svc.DescribeDBClusters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster.getNeptuneDBCluster", "api_error", err)
		return nil, err
	}
	if len(data.DBClusters) > 0 {
		cluster := data.DBClusters[0]
		if *cluster.Engine == "neptune" {
			return cluster, nil
		}
	}
	return nil, nil
}

func getNeptuneDBClusterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	clusterArn := h.Item.(types.DBCluster).DBClusterArn

	// Create session
	svc, err := NeptuneClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster.getNeptuneDBClusterTags", "get_client_error", err)
		return nil, err
	}

	input := &neptune.ListTagsForResourceInput{
		ResourceName: clusterArn,
	}

	tags, err := svc.ListTagsForResource(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster.getNeptuneDBClusterTags", "api_error", err)
		return nil, err
	}

	if tags == nil {
		return nil, nil
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func neptuneDBClusterTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagsDetails := d.HydrateItem.(*neptune.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagsDetails != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagsDetails.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
