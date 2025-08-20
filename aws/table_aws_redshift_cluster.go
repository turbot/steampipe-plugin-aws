package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRedshiftCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_cluster",
		Description: "AWS Redshift Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ClusterNotFound"}),
			},
			Hydrate: getRedshiftCluster,
			Tags:    map[string]string{"service": "redshift", "action": "DescribeClusters"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftClusters,
			Tags:    map[string]string{"service": "redshift", "action": "DescribeClusters"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getRedshiftLoggingDetails,
				Tags: map[string]string{"service": "redshift", "action": "DescribeLoggingStatus"},
			},
			{
				Func: getClusterScheduledActions,
				Tags: map[string]string{"service": "redshift", "action": "DescribeScheduledActions"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_REDSHIFT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_identifier",
				Description: "The unique identifier of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRedshiftClusterARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "cluster_namespace_arn",
				Description: "The namespace Amazon Resource Name (ARN) of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allow_version_upgrade",
				Description: "A boolean value that, if true, indicates that major version upgrades will be applied automatically to the cluster during the maintenance window.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "automated_snapshot_retention_period",
				Description: "The number of days that automatic cluster snapshots are retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "The name of the Availability Zone in which the cluster is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone_relocation_status",
				Description: "Describes the status of the Availability Zone relocation operation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_availability_status",
				Description: "The availability status of the cluster for queries.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_create_time",
				Description: "The date and time that the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "cluster_nodes",
				Description: "The nodes in the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cluster_parameter_groups",
				Description: "The list of cluster parameter groups that are associated with this cluster. Each parameter group in the list is returned with its status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cluster_public_key",
				Description: "The public key for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_revision_number",
				Description: "The specific revision number of the database in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_security_groups",
				Description: "A list of cluster security group that are associated with the cluster. Each security group is represented by an element that contains ClusterSecurityGroup.Name and ClusterSecurityGroup.Status subelements. Cluster security groups are used when the cluster is not created in an Amazon Virtual Private Cloud (VPC). Clusters that are created in a VPC use VPC security groups, which are listed by the VpcSecurityGroups parameter.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterSecurityGroups"),
			},
			{
				Name:        "cluster_snapshot_copy_status",
				Description: "A value that returns the destination region and retention period that are configured for cross-region snapshot copy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cluster_status",
				Description: "The current state of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_subnet_group_name",
				Description: "The name of the subnet group that is associated with the cluster. This parameter is valid only when the cluster is in a VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_version",
				Description: "The version ID of the Amazon Redshift engine that is running on the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "custom_domain_certificate_arn",
				Description: "The certificate Amazon Resource Name (ARN) for the custom domain name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "total_storage_capacity_in_mega_bytes",
				Description: "The total storage capacity of the cluster in megabytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "custom_domain_name",
				Description: "The custom domain name associated with the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_iam_role_arn",
				Description: "The Amazon Resource Name (ARN) for the IAM role set as default for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_password_secret_arn",
				Description: "The Amazon Resource Name (ARN) for the cluster's admin user credentials secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_password_secret_kms_key_id",
				Description: "The ID of the Key Management Service (KMS) key used to encrypt and store the cluster's admin credentials secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "custom_domain_certificate_expiry_date",
				Description: "The expiration date for the certificate associated with the custom domain name.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CustomDomainCertificateExpiryDate").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "aqua_configuration",
				Description: "This field is retired. Amazon Redshift automatically determines whether to use AQUA (Advanced Query Accelerator).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "reserved_node_exchange_status",
				Description: "The status of the reserved-node exchange request. Statuses include in-progress and requested.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "data_transfer_progress",
				Description: "Describes the status of a cluster while it is in the process of resizing with an incremental resize.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_name",
				Description: "The name of the initial database that was created when the cluster was created. This same name is returned for the life of the cluster. If an initial database was not specified, a database named devdev was created by default.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBName"),
			},
			{
				Name:        "deferred_maintenance_windows",
				Description: "Describes a group of DeferredMaintenanceWindow objects.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DeferredMaintenanceWindows"),
			},
			{
				Name:        "elastic_ip_status",
				Description: "The status of the elastic IP (EIP) address.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "elastic_resize_number_of_node_options",
				Description: "The number of nodes that you can resize the cluster to with the elastic resize method.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encrypted",
				Description: "A boolean value that, if true, indicates that data in the cluster is encrypted at rest.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "endpoint",
				Description: "The connection endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "enhanced_vpc_routing",
				Description: "An option that specifies whether to create the cluster with enhanced VPC routing enabled. To create a cluster that uses enhanced VPC routing, the cluster must be in a VPC. If this option is true, enhanced VPC routing is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "expected_next_snapshot_schedule_time",
				Description: "The date and time when the next snapshot is expected to be taken for clusters with a valid snapshot schedule and backups enabled.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "expected_next_snapshot_schedule_time_status",
				Description: "The status of next expected snapshot for clusters having a valid snapshot schedule and backups enabled.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hsm_status",
				Description: "A value that reports whether the Amazon Redshift cluster has finished applying any hardware security module (HSM) settings changes specified in a modify cluster command.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_roles",
				Description: "A list of AWS Identity and Access Management (IAM) roles that can be used by the cluster to access other AWS services.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IamRoles"),
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS Key Management Service (AWS KMS) key ID of the encryption key used to encrypt data in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "maintenance_track_name",
				Description: "The name of the maintenance track for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "manual_snapshot_retention_period",
				Description: "The default number of days to retain a manual snapshot. If the value is -1, the snapshot is retained indefinitely. This setting doesn't change the retention period of existing snapshots. The value must be either -1 or an integer between 1 and 3,653.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "master_username",
				Description: "The master user name for the cluster. This name is used to connect to the database that is specified in the DBName parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "modify_status",
				Description: "The status of a modify operation, if any, initiated for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "multi_az",
				Description: "A boolean value that, if true, indicates that the cluster is deployed in two Availability Zones.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("MultiAZ"),
			},
			{
				Name:        "next_maintenance_window_start_time",
				Description: "The date and time in UTC when system maintenance can begin.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "node_type",
				Description: "The node type for the nodes in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_nodes",
				Description: "The number of compute nodes in the cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "pending_actions",
				Description: "Cluster operations that are waiting to be started.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pending_modified_values",
				Description: "A value that, if present, indicates that changes to the cluster are pending. Specific pending changes are identified by subelements.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "The weekly time range, in Universal Coordinated Time (UTC), during which system maintenance can occur.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "publicly_accessible",
				Description: "A boolean value that, if true, indicates that the cluster can be accessed from a public network.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "resize_info",
				Description: "Describes a resize operation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "restore_status",
				Description: "A value that describes the status of a cluster restore action. This parameter returns null if the cluster was not created by restoring a snapshot.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "snapshot_schedule_identifier",
				Description: "A unique identifier for the cluster snapshot schedule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_schedule_state",
				Description: "The current state of the cluster snapshot schedule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotScheduleState"),
			},
			{
				Name:        "vpc_id",
				Description: "The identifier of the VPC the cluster is in, if the cluster is in a VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_security_groups",
				Description: "A list of Amazon Virtual Private Cloud (Amazon VPC) security groups that are associated with the cluster. This parameter is returned only if the cluster is in a VPC.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "logging_status",
				Description: "Describes the status of logging for a cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRedshiftLoggingDetails,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "scheduled_actions",
				Description: "A list of scheduled actions for specified cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getClusterScheduledActions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getRedshiftClusterTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRedshiftClusterARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listRedshiftClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_cluster.listRedshiftClusters", "connection_error", err)
		return nil, err
	}

	input := &redshift.DescribeClustersInput{
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
	paginator := redshift.NewDescribeClustersPaginator(svc, input, func(o *redshift.DescribeClustersPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_redshift_cluster.listRedshiftClusters", "api_error", err)
			return nil, err
		}

		for _, items := range output.Clusters {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedshiftCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_cluster.getRedshiftCluster", "connection_error", err)
		return nil, err
	}

	name := d.EqualsQuals["cluster_identifier"].GetStringValue()

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Build params
	params := &redshift.DescribeClustersInput{
		ClusterIdentifier: aws.String(name),
	}

	op, err := svc.DescribeClusters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_cluster.getRedshiftCluster", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.Clusters) > 0 {
		return op.Clusters[0], nil
	}
	return nil, nil
}

func getRedshiftLoggingDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := h.Item.(types.Cluster).ClusterIdentifier

	// Create client
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_cluster.getRedshiftLoggingDetails", "connection_error", err)
		return nil, err
	}

	params := &redshift.DescribeLoggingStatusInput{
		ClusterIdentifier: name,
	}

	op, err := svc.DescribeLoggingStatus(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_cluster.getRedshiftLoggingDetails", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getClusterScheduledActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_cluster.getClusterScheduledActions", "connection_error", err)
		return nil, err
	}
	name := *h.Item.(types.Cluster).ClusterIdentifier

	params := &redshift.DescribeScheduledActionsInput{
		Filters: []types.ScheduledActionFilter{
			{
				Name:   "cluster-identifier",
				Values: []string{name},
			},
		},
	}

	var scheduledActions []types.ScheduledAction
	// List call
	paginator := redshift.NewDescribeScheduledActionsPaginator(svc, params, func(o *redshift.DescribeScheduledActionsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_redshift_cluster.getClusterScheduledActions", "api_error", err)
			return nil, err
		}

		scheduledActions = append(scheduledActions, output.ScheduledActions...)
	}

	return scheduledActions, nil
}

func getRedshiftClusterARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(types.Cluster)
	region := d.EqualsQualString(matrixKeyRegion)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_cluster.getRedshiftClusterARN", "getCommonColumns_error", err)
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":redshift:" + region + ":" + commonColumnData.AccountId + ":cluster:" + *cluster.ClusterIdentifier

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getRedshiftClusterTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cluster := d.HydrateItem.(types.Cluster)

	if len(cluster.Tags) > 0 {
		turbotTagsMap := map[string]string{}
		for _, i := range cluster.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
