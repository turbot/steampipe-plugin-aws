package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"

	rdsEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsNeptuneDBClusterSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_neptune_db_cluster_snapshot",
		Description: "AWS Neptune DB Cluster Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_cluster_snapshot_identifier"),
			Hydrate:    getNeptuneDBClusterSnapshot,
			Tags:       map[string]string{"service": "neptune", "action": "DescribeDBClusterSnapshots"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBClusterSnapshotNotFoundFault"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listNeptuneDBClusterSnapshots,
			Tags:    map[string]string{"service": "neptune", "action": "DescribeDBClusterSnapshots"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_cluster_identifier", Require: plugin.Optional},
				{Name: "snapshot_type", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getNeptuneDBClusterSnapshotAttributes,
				Tags: map[string]string{"service": "neptune", "action": "DescribeDBClusterSnapshotAttributes"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsEndpoint.AWS_RDS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_cluster_snapshot_identifier",
				Description: "Specifies the identifier for a DB cluster snapshot. Must match the identifier of an existing snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotIdentifier"),
			},
			{
				Name:        "db_cluster_snapshot_arn",
				Description: "The Amazon Resource Name (ARN) for the DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotArn"),
			},
			{
				Name:        "snapshot_type",
				Description: "Provides the type of the DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the status of this DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_cluster_identifier",
				Description: "Specifies the DB cluster identifier of the DB cluster that this DB cluster snapshot was created from.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},

			{
				Name:        "cluster_create_time",
				Description: "Specifies the time when the DB cluster was created, in Universal Coordinated Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "allocated_storage",
				Description: "Specifies the allocated storage size in gibibytes (GiB).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "snapshot_create_time",
				Description: "Provides the time when the snapshot was taken, in Universal Coordinated Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "percent_progress",
				Description: "Specifies the percentage of the estimated data that has been transferred.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "source_db_cluster_snapshot_arn",
				Description: "If the DB cluster snapshot was copied from a source DB cluster snapshot, the Amazon Resource Name (ARN) for the source DB cluster snapshot, otherwise, a null value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceDBClusterSnapshotArn"),
			},
			{
				Name:        "engine",
				Description: "Specifies the name of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Provides the version of the database engine for this DB cluster snapshot.",
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
				Description: "If StorageEncrypted is true, the Amazon KMS key identifier for the encrypted DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_model",
				Description: "Provides the license model information for this DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_username",
				Description: "Not supported by Neptune.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "Specifies the port that the DB cluster was listening on at the time of the snapshot.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the DB cluster snapshot is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VPC ID associated with the DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_type",
				Description: "The storage type associated with the DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zones",
				Description: "Provides the list of EC2 Availability Zones that instances in the DB cluster snapshot can be restored in.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_cluster_snapshot_attributes",
				Description: "A list of DB cluster snapshot attribute names and values for a manual DB cluster snapshot.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNeptuneDBClusterSnapshotAttributes,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
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
				Transform:   transform.FromField("DBClusterSnapshotArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listNeptuneDBClusterSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := NeptuneClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster_snapshot.listNeptuneDBClusterSnapshots", "get_client_error", err)
		return nil, err
	}

	input := &neptune.DescribeDBClusterSnapshotsInput{
		MaxRecords: aws.Int32(100),
	}

	if d.EqualsQualString("db_cluster_identifier") != "" {
		input.DBClusterIdentifier = aws.String(d.EqualsQualString("db_cluster_identifier"))
	}

	if d.EqualsQualString("snapshot_type") != "" {
		input.SnapshotType = aws.String(d.EqualsQualString("snapshot_type"))
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

	paginator := neptune.NewDescribeDBClusterSnapshotsPaginator(svc, input, func(o *neptune.DescribeDBClusterSnapshotsPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_neptune_db_cluster_snapshot.listNeptuneDBClusterSnapshots", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBClusterSnapshots {

			if *items.Engine == "neptune" {
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

func getNeptuneDBClusterSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	snapshotIdentifier := d.EqualsQuals["db_cluster_snapshot_identifier"].GetStringValue()

	// Create service
	svc, err := NeptuneClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster_snapshot.getNeptuneDBClusterSnapshot", "connection_error", err)
		return nil, err
	}

	params := &neptune.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: aws.String(snapshotIdentifier),
	}

	op, err := svc.DescribeDBClusterSnapshots(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster_snapshot.getNeptuneDBClusterSnapshot", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.DBClusterSnapshots) > 0 {
		snapshot := op.DBClusterSnapshots[0]
		if *snapshot.Engine == "neptune" {
			return snapshot, nil
		}
	}
	return nil, nil
}

func getNeptuneDBClusterSnapshotAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	dbClusterSnapshot := h.Item.(types.DBClusterSnapshot)

	// Create service
	svc, err := NeptuneClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster_snapshot.getNeptuneDBClusterSnapshotAttributes", "connection_error", err)
		return nil, err
	}

	params := &neptune.DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: dbClusterSnapshot.DBClusterSnapshotIdentifier,
	}

	dbClusterSnapshotData, err := svc.DescribeDBClusterSnapshotAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_neptune_db_cluster_snapshot.getNeptuneDBClusterSnapshotAttributes", "api_error", err)
		return nil, err
	}

	var attributes = make([]map[string]interface{}, 0)

	if dbClusterSnapshotData.DBClusterSnapshotAttributesResult != nil {

		for _, attribute := range dbClusterSnapshotData.DBClusterSnapshotAttributesResult.DBClusterSnapshotAttributes {
			var result = make(map[string]interface{})

			result["AttributeName"] = attribute.AttributeName
			if len(attribute.AttributeValues) == 0 {
				result["AttributeValues"] = nil
			} else {
				result["AttributeValues"] = attribute.AttributeValues
			}

			attributes = append(attributes, result)

		}
	}

	return attributes, nil
}
