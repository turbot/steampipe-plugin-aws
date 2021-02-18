package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsRDSDBClusterSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_cluster_snapshot",
		Description: "AWS RDS DB Cluster Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("db_cluster_snapshot_identifier"),
			ShouldIgnoreError: isNotFoundError([]string{"DBSnapshotNotFound", "DBClusterSnapshotNotFoundFault"}),
			ItemFromKey:       dbClusterSnapshotIdentifierFromKey,
			Hydrate:           getRDSDBClusterSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBClusterSnapshots,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_cluster_snapshot_identifier",
				Description: "The friendly name to identify the DB Cluster Snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB Cluster Snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotArn"),
			},
			{
				Name:        "type",
				Description: "The type of the DB Cluster Snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotType"),
			},
			{
				Name:        "status",
				Description: "Specifies the status of this DB Cluster Snapshot",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_cluster_identifier",
				Description: "The friendly name to identify the DB Cluster, that the snapshot snapshot was created from",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "create_time",
				Description: "The time when the snapshot was taken",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SnapshotCreateTime"),
			},
			{
				Name:        "allocated_storage",
				Description: "Specifies the allocated storage size in gibibytes (GiB)",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cluster_create_time",
				Description: "Specifies the time when the DB cluster was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "engine",
				Description: "Specifies the name of the database engine",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Specifies the version of the database engine for this DB cluster snapshot",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_database_authentication_enabled",
				Description: "Specifies whether mapping of AWS Identity and Access Management (IAM) accounts to database accounts is enabled, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IAMDatabaseAuthenticationEnabled"),
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS KMS key identifier for the AWS KMS customer master key (CMK)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_model",
				Description: "Provides the license model information for this DB cluster snapshot",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_user_name",
				Description: "Provides the master username for the DB cluster snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "percent_progress",
				Description: "Specifies the percentage of the estimated data that has been transferred",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "port",
				Description: "Specifies the port that the DB cluster was listening on at the time of the snapshot",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "source_db_cluster_snapshot_arn",
				Description: "The Amazon Resource Name (ARN) for the source DB cluster snapshot, if the DB cluster snapshot was copied from a source DB cluster snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceDBClusterSnapshotArn"),
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the DB cluster snapshot is encrypted, or not",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VPC ID associated with the DB cluster snapshot",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zones",
				Description: "A list of Availability Zones (AZs) where instances in the DB cluster snapshot can be restored",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_cluster_snapshot_attributes",
				Description: "A list of DB cluster snapshot attribute names and values for a manual DB cluster snapshot",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSDBClusterSnapshotAttributes,
				Transform:   transform.FromField("DBClusterSnapshotAttributesResult.DBClusterSnapshotAttributes"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB Cluster Snapshot",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getRDSDBClusterSnapshotTurbotTags),
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

//// BUILD HYDRATE INPUT

func dbClusterSnapshotIdentifierFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	snapshotIdentifier := quals["db_cluster_snapshot_identifier"].GetStringValue()
	item := &rds.DBClusterSnapshot{
		DBClusterSnapshotIdentifier: &snapshotIdentifier,
	}
	return item, nil
}

//// LIST FUNCTION

func listRDSDBClusterSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listRDSDBClusterSnapshots", "AWS_REGION", region)

	// Create Session
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeDBClusterSnapshotsPages(
		&rds.DescribeDBClusterSnapshotsInput{},
		func(page *rds.DescribeDBClusterSnapshotsOutput, isLast bool) bool {
			for _, dbClusterSnapshot := range page.DBClusterSnapshots {
				d.StreamListItem(ctx, dbClusterSnapshot)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBClusterSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	dbClusterSnapshot := h.Item.(*rds.DBClusterSnapshot)

	// Create service
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: aws.String(*dbClusterSnapshot.DBClusterSnapshotIdentifier),
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

func getAwsRDSDBClusterSnapshotAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRDSDBClusterSnapshotAttributes")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	dbClusterSnapshot := h.Item.(*rds.DBClusterSnapshot)

	// Create service
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: dbClusterSnapshot.DBClusterSnapshotIdentifier,
	}

	dbClusterSnapshotData, err := svc.DescribeDBClusterSnapshotAttributes(params)
	if err != nil {
		return nil, err
	}

	return dbClusterSnapshotData, nil
}

//// TRANSFORM FUNCTIONS ////

func getRDSDBClusterSnapshotTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbClusterSnapshot := d.HydrateItem.(*rds.DBClusterSnapshot)

	if dbClusterSnapshot.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbClusterSnapshot.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
