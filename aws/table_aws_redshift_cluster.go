package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsRedshiftCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_cluster",
		Description: "AWS Redshift Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("cluster_identifier"),
			ShouldIgnoreError: isNotFoundError([]string{"ClusterNotFound"}),
			ItemFromKey:       clusterIdentifierFromKey,
			Hydrate:           getRedshiftCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftClusters,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_identifier",
				Description: "The unique identifier of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The namespace Amazon Resource Name (ARN) of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterNamespaceArn"),
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
				Name:        "encrypted",
				Description: "A boolean value that, if true, indicates that data in the cluster is encrypted at rest.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "manual_snapshot_retention_period",
				Description: "The default number of days to retain a manual snapshot. If the value is -1, the snapshot is retained indefinitely. This setting doesn't change the retention period of existing snapshots. The value must be either -1 or an integer between 1 and 3,653.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "publicly_accessible",
				Description: "A boolean value that, if true, indicates that the cluster can be accessed from a public network.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "tags",
				Description: "The list of tags for the cluster.",
				Type:        proto.ColumnType_JSON,
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
			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getRedshiftClusterTurbotTags),
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

func clusterIdentifierFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	clusterIdentifier := quals["cluster_identifier"].GetStringValue()
	item := &redshift.Cluster{
		ClusterIdentifier: &clusterIdentifier,
	}
	return item, nil
}

//// LIST FUNCTION

func listRedshiftClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listRedshiftClusters", "AWS_REGION", defaultRegion)

	// Create Session
	svc, err := RedshiftService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	/*
		var params []*redshift.Filter
		if defaultRegion == "sa-east-1" {
			params = []*redshift.Filter{
				{
					Name: aws.String("engine"),
					Values: []*string{
						aws.String("mariadb"),
						aws.String("mysql"),
						aws.String("oracle-ee"),
						aws.String("oracle-se"),
						aws.String("oracle-se1"),
						aws.String("oracle-se2"),
						aws.String("postgres"),
						aws.String("sqlserver-ee"),
						aws.String("sqlserver-ex"),
						aws.String("sqlserver-se"),
						aws.String("sqlserver-web"),
					},
				},
			}
		} else if defaultRegion == "eu-north-1" {
			params = []*redshift.Filter{
				{
					Name: aws.String("engine"),
					Values: []*string{
						aws.String("mariadb"),
						aws.String("mysql"),
						aws.String("oracle-ee"),
						aws.String("oracle-se"),
						aws.String("oracle-se1"),
						aws.String("oracle-se2"),
						aws.String("postgres"),
						aws.String("sqlserver-ee"),
						aws.String("sqlserver-se"),
						aws.String("sqlserver-web"),
					},
				},
			}
		} else {
			params = []*redshift.Filter{
				{
					Name: aws.String("engine"),
					Values: []*string{
						aws.String("aurora"),
						aws.String("aurora-mysql"),
						aws.String("aurora-postgresql"),
						aws.String("mariadb"),
						aws.String("mysql"),
						aws.String("oracle-ee"),
						aws.String("oracle-se"),
						aws.String("oracle-se1"),
						aws.String("oracle-se2"),
						aws.String("postgres"),
						aws.String("sqlserver-ee"),
						aws.String("sqlserver-se"),
						aws.String("sqlserver-ex"),
						aws.String("sqlserver-web"),
					},
				},
			}
		}
	*/

	// List call
	err = svc.DescribeClustersPages(
		&redshift.DescribeClustersInput{},
		func(page *redshift.DescribeClustersOutput, isLast bool) bool {
			for _, cluster := range page.Clusters {
				d.StreamListItem(ctx, cluster)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRedshiftCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	cluster := h.Item.(*redshift.Cluster)

	// Create service
	svc, err := RedshiftService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &redshift.DescribeClustersInput{
		ClusterIdentifier: aws.String(*cluster.ClusterIdentifier),
	}

	op, err := svc.DescribeClusters(params)
	if err != nil {
		return nil, err
	}

	if op.Clusters != nil && len(op.Clusters) > 0 {
		return op.Clusters[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS ////

func getRedshiftClusterTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cluster := d.HydrateItem.(*redshift.Cluster)

	if cluster.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range cluster.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
