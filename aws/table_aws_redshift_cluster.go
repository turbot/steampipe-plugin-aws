package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/redshift"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRedshiftCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_cluster",
		Description: "AWS Redshift Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ClusterNotFound", "InvalidParameterValue", "NoSuchEntity", "ValidationError"}),
			ItemFromKey:       redshiftClusterFromKey,
			Hydrate:           getAwsRedshiftCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRedshiftClusters,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The unique identifier of the cluster",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The namespace Amazon Resource Name (ARN) of the cluster",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterNamespaceArn"),
			},
			{
				Name:        "master_user_name",
				Description: "The master user name for the cluster. This name is used to connect to the database that is specified in the DBName parameter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "status",
				Description: "The current state of the cluster",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterStatus"),
			},
			{
				Name:        "encrypted",
				Description: "A boolean value that, if true, indicates that data in the cluster is encrypted at rest",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name: "allow_version_upgrade",
				Description: "A boolean value that, if true, indicates that major version upgrades will	be applied automatically to the cluster during the maintenance window.",
				Type: proto.ColumnType_BOOL,
			},
			{
				Name:        "automated_snapshot_retention_period",
				Description: "The number of days that automatic cluster snapshots are retained",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "manual_snapshot_retention_period",
				Description: "The value must be either -1 or an integer between 1 and 3,653",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "The name of the Availability Zone in which the cluster is located",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone_relocation_status",
				Description: "Describes the status of the Availability Zone relocation operation",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_availability_status",
				Description: "The availability status of the cluster for queries",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_create_time",
				Description: "The date and time that the cluster was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "cluster_public_key",
				Description: "The public key for the cluster",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_revision_number",
				Description: "The specific revision number of the database in the cluster",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_subnet_group_name",
				Description: "The name of the subnet group that is associated with the cluster. This parameter is valid only when the cluster is in a VPC",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_version",
				Description: "The version ID of the Amazon Redshift engine that is running on the cluster",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_name",
				Description: "The name of the initial database that was created when the cluster was created",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBName"),
			},
			{
				Name:        "elastic_resize_number_of_node_options",
				Description: "The number of nodes that you can resize the cluster to with the elastic resize method",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enhanced_vpc_routing",
				Description: "If this option is true, enhanced VPC routing is enabled",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS Key Management Service (AWS KMS) key ID of the encryption key used to encrypt data in the cluster",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The identifier of the VPC the cluster is in, if the cluster is in a VPC",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "modify_status",
				Description: "The status of a modify operation, if any, initiated for the cluster",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type",
				Description: "The node type for the nodes in the cluster",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_nodes",
				Description: "The number of compute nodes in the cluster",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "publicly_accessible",
				Description: "A boolean value that, if true, indicates that the cluster can be accessed from a public network",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "expected_next_snapshot_schedule_time",
				Description: "The date and time when the next snapshot is expected to be taken for clusters with a valid snapshot schedule and backups enabled",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "expected_next_snapshot_schedule_time_status",
				Description: "The status of next expected snapshot for clusters having a valid snapshot",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "maintenance_track_name",
				Description: "The name of the maintenance track for the cluster",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_maintenance_window_start_time",
				Description: "The date and time in UTC when system maintenance can begin",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "The weekly time range, in Universal Coordinated Time (UTC), during which system maintenance can occur",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_schedule_identifier",
				Description: "A unique identifier for the cluster snapshot schedule",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_schedule_state",
				Description: "The current state of the cluster snapshot schedule",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pending_actions",
				Description: "Cluster operations that are waiting to be started",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PendingActions"),
			},
			{
				Name:        "pending_modified_values",
				Description: "A value that, if present, indicates that changes to the cluster are pending. Specific pending changes are identified by subelements",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PendingModifiedValues"),
			},
			{
				Name:        "vpc_security_groups",
				Description: "A list of Amazon Virtual Private Cloud (Amazon VPC) security groups that are associated with the cluster. This parameter is returned only if the cluster is in a VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VpcSecurityGroups"),
			},
			{
				Name:        "cluster_nodes",
				Description: "The nodes in the cluster",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterNodes"),
			},
			{
				Name:        "cluster_parameter_groups",
				Description: "The list of cluster parameter groups that are associated with this cluster",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterParameterGroups"),
			},
			{
				Name:        "cluster_security_groups",
				Description: "A list of cluster security group that are associated with the cluster",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterSecurityGroups"),
			},
			{
				Name:        "cluster_snapshot_copy_status",
				Description: "A value that returns the destination region and retention period that are configured for cross-region snapshot copy",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterSnapshotCopyStatus"),
			},
			{
				Name:        "data_transfer_progress",
				Description: "Describes the status of a cluster while it is in the process of resizing with an incremental resize.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DataTransferProgress"),
			},
			{
				Name:        "deferred_maintenance_windows",
				Description: "Describes a group of DeferredMaintenanceWindow objects",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DeferredMaintenanceWindows"),
			},
			{
				Name:        "elastic_ip_status",
				Description: "The status of the elastic IP (EIP) address",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ElasticIpStatus"),
			},
			{
				Name:        "endpoint",
				Description: "The connection endpoint",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Endpoint"),
			},
			{
				Name:        "iam_roles",
				Description: "A list of AWS Identity and Access Management (IAM) roles that can be used by the cluster to access other AWS services",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IamRoles"),
			},
			{
				Name:        "hsm_status",
				Description: "A value that reports whether the Amazon Redshift cluster has finished applying any hardware security module (HSM) settings changes specified in a modify cluster command",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HsmStatus"),
			},
			{
				Name:        "resize_info",
				Description: "A boolean value indicating if the resize operation can be cancelled",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResizeInfo"),
			},
			{
				Name:        "restore_status",
				Description: "A value that describes the status of a cluster restore action",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RestoreStatus"),
			},
			{
				Name:        "tags_src",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(redshiftClusterTurbotTags, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterIdentifier"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterNamespaceArn").Transform(arnToAkas),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func redshiftClusterFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	clusterName := quals["name"].GetStringValue()
	item := &redshift.Cluster{
		ClusterIdentifier: &clusterName,
	}
	return item, nil
}

//// LIST FUNCTION

func listAwsRedshiftClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := RedshiftService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.DescribeClustersPages(
		&redshift.DescribeClustersInput{},
		func(clusterOutput *redshift.DescribeClustersOutput, isLast bool) bool {
			for _, cluster := range clusterOutput.Clusters {
				d.StreamListItem(ctx, cluster)
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsRedshiftCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	clusterIdentifier := d.KeyColumnQuals["name"].GetStringValue()

	// create service
	svc, err := RedshiftService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &redshift.DescribeClustersInput{
		ClusterIdentifier: &clusterIdentifier,
	}

	// execute list call
	op, err := svc.DescribeClusters(params)
	if err != nil {
		return nil, err
	}

	if len(op.Clusters) > 0 {
		return op.Clusters[0], nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func redshiftClusterTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("redshiftClusterTurbotTags")
	cluster := d.HydrateItem.(*redshift.Cluster)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if cluster.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range cluster.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
