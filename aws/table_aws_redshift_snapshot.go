package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsRedshiftSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_snapshot",
		Description: "AWS Redshift Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("snapshot_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ClusterSnapshotNotFound"}),
			},
			Hydrate: getAwsRedshiftSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRedshiftSnapshots,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_identifier", Require: plugin.Optional},
				{Name: "owner_account", Require: plugin.Optional},
				{Name: "snapshot_type", Require: plugin.Optional},
				{Name: "snapshot_create_time", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "snapshot_identifier",
				Description: "The unique identifier of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_identifier",
				Description: "The identifier of the cluster for which the snapshot was taken.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_type",
				Description: "The snapshot type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actual_incremental_backup_size_in_mega_bytes",
				Description: "The size of the incremental backup.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone in which the cluster was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_progress_in_mega-bytes",
				Description: "The number of megabytes that have been transferred to the snapshot backup.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "cluster_create_time",
				Description: "The time (UTC) when the cluster was originally created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "cluster_version",
				Description: "The version ID of the Amazon Redshift engine that is running on the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "current_backup_rate_in_mega_bytes_per_second",
				Description: "The number of megabytes per second being transferred to the snapshot backup.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "db_name",
				Description: "The name of the database that was created when the cluster was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "elapsed_time_in_seconds",
				Description: "The amount of time an in-progress snapshot backup has been running, or the amount of time it took a completed backup to finish.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encrypted",
				Description: "If true, the data in the snapshot is encrypted at rest.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "encrypted_with_hsm",
				Description: "A boolean that indicates whether the snapshot data is encrypted using the HSM keys of the source cluster.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "engine_full_version",
				Description: "The cluster version of the cluster used to create the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enhanced_vpc_routing",
				Description: "An option that specifies whether to create the cluster with enhanced VPC routing enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "estimated_seconds_to_completion",
				Description: "The estimate of the time remaining before the snapshot backup will complete.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS KMS key ID of the encryption key that was used to encrypt data in the cluster from which the snapshot was taken.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "maintenance_track_name",
				Description: "The name of the maintenance track for the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "manual_snapshot_remaining_days",
				Description: "The number of days until a manual snapshot will pass its retention period.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "manual_snapshot_retention_period",
				Description: "The number of days that a manual snapshot is retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "master_username",
				Description: "The master user name for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type",
				Description: "The node type of the nodes in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_nodes",
				Description: "The number of nodes in the cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "owner_account",
				Description: "The AWS customer account used to create or copy the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The port that the cluster is listening on.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "snapshot_create_time",
				Description: "The time (in UTC format) when Amazon Redshift began the snapshot.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "snapshot_retention_start_time",
				Description: "A timestamp representing the start of the retention period for the snapshot.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "source_region",
				Description: "The source region from which the snapshot was copied.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The snapshot status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "total_backup_size_in_mega_bytes",
				Description: "The size of the complete set of backup data that would be used to restore the cluster.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "vpc_id",
				Description: "The VPC identifier of the cluster if the snapshot is from a cluster in a VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "accounts_with_restore_access",
				Description: "A list of the AWS customer accounts authorized to restore the snapshot.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "restorable_node_types",
				Description: "The list of node types that this cluster snapshot is able to restore into.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(redshiftSnapshotTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRedshiftSnapshotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

func listAwsRedshiftSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsRedshiftSnapshots")

	// Create Session
	svc, err := RedshiftService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &redshift.DescribeClusterSnapshotsInput{
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

	equalQuals := d.KeyColumnQuals
	if equalQuals["cluster_identifier"] != nil {
		if equalQuals["cluster_identifier"].GetStringValue() != "" {
			input.ClusterIdentifier = aws.String(equalQuals["cluster_identifier"].GetStringValue())
		}
	}
	if equalQuals["owner_account"] != nil {
		if equalQuals["owner_account"].GetStringValue() != "" {
			input.OwnerAccount = aws.String(equalQuals["owner_account"].GetStringValue())
		}
	}
	if equalQuals["snapshot_type"] != nil {
		if equalQuals["snapshot_type"].GetStringValue() != "" {
			input.SnapshotType = aws.String(equalQuals["snapshot_type"].GetStringValue())
		}
	}
	if equalQuals["snapshot_create_time"] != nil {
		input.StartTime = aws.Time(equalQuals["snapshot_create_time"].GetTimestampValue().AsTime())
	}

	// List call
	err = svc.DescribeClusterSnapshotsPages(
		input,
		func(page *redshift.DescribeClusterSnapshotsOutput, isLast bool) bool {
			for _, snapshot := range page.Snapshots {
				d.StreamListItem(ctx, snapshot)

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

//// HYDRATE FUNCTIONS

func getAwsRedshiftSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(*redshift.Snapshot).SnapshotIdentifier
	} else {
		name = d.KeyColumnQuals["snapshot_identifier"].GetStringValue()
	}

	// Create service
	svc, err := RedshiftService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &redshift.DescribeClusterSnapshotsInput{
		SnapshotIdentifier: aws.String(name),
	}

	op, err := svc.DescribeClusterSnapshots(params)
	if err != nil {
		return nil, err
	}

	if op != nil && len(op.Snapshots) > 0 {
		return op.Snapshots[0], nil
	}
	return nil, nil
}

func getAwsRedshiftSnapshotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRedshiftSnapshotAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	snapshot := h.Item.(*redshift.Snapshot)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":redshift:" + region + ":" + commonColumnData.AccountId + ":snapshot:" + *snapshot.ClusterIdentifier + "/" + *snapshot.SnapshotIdentifier

	// Get data for turbot defined properties
	akas := []string{arn}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func redshiftSnapshotTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	snapshot := d.HydrateItem.(*redshift.Snapshot)

	// Get the resource tags
	if snapshot.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range snapshot.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
