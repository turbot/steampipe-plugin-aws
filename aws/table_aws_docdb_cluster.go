package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsDocDBCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_docdb_cluster",
		Description: "AWS Doc DB Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("db_cluster_identifier"),
			ShouldIgnoreError: isNotFoundError([]string{"DBClusterNotFoundFault"}),
			Hydrate:           getDocDBCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listDocDBClusters,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_cluster_identifier",
				Description: "Contains a user-supplied cluster identifier. This identifier is the unique key that identifies a cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the Cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterArn"),
			},
			{
				Name:        "status",
				Description: "Specifies the current state of this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "Specifies the time when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ClusterCreateTime"),
			},
			{
				Name:        "backup_retention_period",
				Description: "Specifies the number of days for which automatic snapshots are retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "db_cluster_parameter_group",
				Description: "Specifies the name of the cluster parameter group for the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterParameterGroup"),
			},
			{
				Name:        "db_cluster_resource_id",
				Description: "The Region-unique, immutable identifier for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_subnet_group",
				Description: "Specifies information on the subnet group associated with the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup"),
			},
			{
				Name:        "deletion_protection",
				Description: "Specifies whether the cluster has deletion protection enabled, or not.",
				Type:        proto.ColumnType_BOOL,
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
				Name:        "kms_key_id",
				Description: "The AWS KMS key identifier for the encrypted cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_restorable_time",
				Description: "Specifies the latest time to which a database can be restored with point-in-time restore.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "master_user_name",
				Description: "Contains the master username for the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "multi_az",
				Description: "Specifies whether the cluster has instances in multiple Availability Zones, or not.",
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
				Name:        "replication_source_identifier",
				Description: "Contains the identifier of the source cluster if this cluster is a secondary cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the cluster is encrypted, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "associated_roles",
				Description: "A list of AWS IAM roles that are associated with the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "availability_zones",
				Description: "A list of Availability Zones (AZs) where instances in the cluster can be created.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "members",
				Description: "A list of instances that make up the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterMembers"),
			},
			{
				Name:        "read_replica_identifiers",
				Description: "A list of identifiers of the read replicas associated with this cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_security_groups",
				Description: "A list of VPC security groups that the DB cluster belongs to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the Cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDocDBClusterTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDocDBClusterTags,
				Transform:   transform.FromField("TagList").Transform(docDBClusterTagListToTurbotTags),
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
				Transform:   transform.FromField("DBClusterArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDocDBClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDocDBClusters")

	// Create Session
	svc, err := DocDBService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &docdb.DescribeDBClustersInput{
		MaxRecords: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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
	err = svc.DescribeDBClustersPages(
		input,
		func(page *docdb.DescribeDBClustersOutput, isLast bool) bool {
			for _, dbCluster := range page.DBClusters {
				d.StreamListItem(ctx, dbCluster)

				// Check if context has been cancelled or if the limit has been reached (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDocDBCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	dbClusterIdentifier := d.KeyColumnQuals["db_cluster_identifier"].GetStringValue()
	if len(dbClusterIdentifier) < 1 {
		return nil, nil
	}

	// Create service
	svc, err := DocDBService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &docdb.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(dbClusterIdentifier),
	}

	op, err := svc.DescribeDBClusters(params)
	if err != nil {
		return nil, err
	}

	if op.DBClusters != nil && len(op.DBClusters) > 0 {
		return op.DBClusters[0], nil
	}
	return nil, nil
}

func getDocDBClusterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDocDBClusterTags")
	cluster := h.Item.(*docdb.DBCluster)

	// Create Session
	svc, err := DocDBService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &docdb.ListTagsForResourceInput{
		ResourceName: cluster.DBClusterArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Error("getDocDBClusterTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func docDBClusterTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("docDBClusterTagListToTurbotTags")
	tagList := d.Value.([]*docdb.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
