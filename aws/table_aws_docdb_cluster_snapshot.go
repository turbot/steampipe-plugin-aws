package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsDocDBClusterSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_docdb_cluster_snapshot",
		Description: "AWS Doc DB Cluster Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("db_cluster_snapshot_identifier"),
			ShouldIgnoreError: isNotFoundError([]string{"DBSnapshotNotFound", "DBClusterSnapshotNotFoundFault"}),
			Hydrate:           getDocDBClusterSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate: listDocDBClusterSnapshots,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_cluster_identifier", Require: plugin.Optional},
				{Name: "db_cluster_snapshot_identifier", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_cluster_snapshot_identifier",
				Description: "The friendly name to identify the cluster snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the cluster snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotArn"),
			},
			{
				Name:        "type",
				Description: "The type of the cluster snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotType"),
			},
			{
				Name:        "status",
				Description: "Specifies the status of this cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_cluster_identifier",
				Description: "The friendly name to identify the cluster, that the snapshot snapshot was created from.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "create_time",
				Description: "The time when the snapshot was taken.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SnapshotCreateTime"),
			},
			{
				Name:        "cluster_create_time",
				Description: "Specifies the time when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "engine",
				Description: "Specifies the name of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Specifies the version of the database engine for this cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS KMS key identifier for the AWS KMS customer master key (CMK).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_user_name",
				Description: "Provides the master username for the cluster snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "percent_progress",
				Description: "Specifies the percentage of the estimated data that has been transferred.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "port",
				Description: "Specifies the port that the cluster was listening on at the time of the snapshot.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "source_db_cluster_snapshot_arn",
				Description: "The Amazon Resource Name (ARN) for the source cluster snapshot, if the cluster snapshot was copied from a source cluster snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceDBClusterSnapshotArn"),
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the cluster snapshot is encrypted, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VPC ID associated with the cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zones",
				Description: "A list of Availability Zones (AZs) where instances in the cluster snapshot can be restored.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_cluster_snapshot_attributes",
				Description: "A list of DB cluster snapshot attribute names and values for a manual cluster snapshot.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsDocDBClusterSnapshotAttributes,
				Transform:   transform.FromField("DBClusterSnapshotAttributesResult.DBClusterSnapshotAttributes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the cluster snapshot.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDocDBClusterSnapshotTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDocDBClusterSnapshotTags,
				Transform:   transform.FromField("TagList").Transform(getDocDBClusterSnapshotTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotIdentifier"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterSnapshotArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listDocDBClusterSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDocDBClusterSnapshots")

	// Create Session
	svc, err := DocDBService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := docdb.DescribeDBClusterSnapshotsInput{
		MaxRecords: types.Int64(100),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxRecords {
			if *limit < 20 {
				input.MaxRecords = types.Int64(100)
			} else {
				input.MaxRecords = limit
			}
		}
	}
	filters := buildDocDbClusterSnapshotFilter(d.KeyColumnQuals)

	// List call
	err = svc.DescribeDBClusterSnapshotsPages(filters, func(page *docdb.DescribeDBClusterSnapshotsOutput, isLast bool) bool {
		for _, dbClusterSnapshot := range page.DBClusterSnapshots {
			d.StreamListItem(ctx, dbClusterSnapshot)
			// Check if context has been cancelled or if the limit has been hit (if specified)
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

func getDocDBClusterSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	snapshotIdentifier := d.KeyColumnQuals["db_cluster_snapshot_identifier"].GetStringValue()

	// Create service
	svc, err := DocDBService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &docdb.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: aws.String(snapshotIdentifier),
	}

	op, err := svc.DescribeDBClusterSnapshots(params)
	if err != nil {
		return nil, err
	}

	if op.DBClusterSnapshots != nil && len(op.DBClusterSnapshots) > 0 {
		return op.DBClusterSnapshots[0], nil
	}
	return nil, nil
}

func getAwsDocDBClusterSnapshotAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsDocDBClusterSnapshotAttributes")

	dbClusterSnapshot := h.Item.(*docdb.DBClusterSnapshot)

	// Create service
	svc, err := DocDBService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &docdb.DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: dbClusterSnapshot.DBClusterSnapshotIdentifier,
	}

	dbClusterSnapshotData, err := svc.DescribeDBClusterSnapshotAttributes(params)
	if err != nil {
		return nil, err
	}

	return dbClusterSnapshotData, nil
}

func getDocDBClusterSnapshotTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getDocDBClusterSnapshotTags")
	cluster := h.Item.(*docdb.DBClusterSnapshot)

	// Create Session
	svc, err := DocDBService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &docdb.ListTagsForResourceInput{
		ResourceName: cluster.DBClusterSnapshotArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Error("getDocDBClusterSnapshotTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getDocDBClusterSnapshotTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDocDBClusterSnapshotTurbotTags")
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

//// UTILITY FUNCTIONS

// build snapshots list call input filter
func buildDocDbClusterSnapshotFilter(quals plugin.KeyColumnEqualsQualMap) *docdb.DescribeDBClusterSnapshotsInput {
	var filters *docdb.DescribeDBClusterSnapshotsInput

	if quals["db_cluster_identifier"] != nil {
		*filters.DBClusterIdentifier = quals["db_cluster_identifier"].GetStringValue()
		*filters.DBClusterSnapshotIdentifier = quals["db_cluster_snapshot_identifier"].GetStringValue()
		*filters.SnapshotType = quals["type"].GetStringValue()
	}

	return filters
}
