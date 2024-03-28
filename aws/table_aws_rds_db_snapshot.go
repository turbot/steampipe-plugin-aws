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

func tableAwsRDSDBSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_snapshot",
		Description: "AWS RDS DB Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_snapshot_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBSnapshotNotFound"}),
			},
			Hydrate: getRDSDBSnapshot,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBSnapshots"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBSnapshots,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBSnapshots"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_instance_identifier", Require: plugin.Optional},
				{Name: "dbi_resource_id", Require: plugin.Optional},
				{Name: "engine", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsRDSDBSnapshotAttributes,
				Tags: map[string]string{"service": "rds", "action": "DescribeDBSnapshotAttributes"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_snapshot_identifier",
				Description: "The friendly name to identify the DB snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSnapshotIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSnapshotArn"),
			},
			{
				Name:        "type",
				Description: "Provides the type of the DB snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotType"),
			},
			{
				Name:        "status",
				Description: "Specifies the status of this DB snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "Specifies when the snapshot was taken.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SnapshotCreateTime"),
			},
			{
				Name:        "allocated_storage",
				Description: "Specifies the allocated storage size in gibibytes(GiB).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "Specifies the name of the Availability Zone the DB instance was located in, at the time of the DB snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_instance_identifier",
				Description: "Specifies the DB instance identifier of the DB instance this DB snapshot was created from.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceIdentifier"),
			},
			{
				Name:        "dbi_resource_id",
				Description: "The identifier for the source DB instance, which can't be changed and which is unique to an AWS Region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encrypted",
				Description: "Specifies whether the DB snapshot is encrypted, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "engine",
				Description: "Specifies the name of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Specifies the version of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_database_authentication_enabled",
				Description: "Specifies whether the mapping of AWS IAM accounts to database accounts is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IAMDatabaseAuthenticationEnabled"),
			},
			{
				Name:        "instance_create_time",
				Description: "Specifies the time when the DB instance, from which the snapshot was taken, was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "iops",
				Description: "Specifies the Provisioned IOPS (I/O operations per second) value of the DB instance at the time of the snapshot.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kms_key_id",
				Description: "Specifies the AWS KMS key identifier for the encrypted DB snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_model",
				Description: "Specifies the License model information for the restored DB instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_user_name",
				Description: "Provides the master username for the DB snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "option_group_name",
				Description: "Provides the option group name for the DB snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "percent_progress",
				Description: "The percentage of the estimated data that has been transferred.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "port",
				Description: "Specifies the port that the database engine was listening on at the time of the snapshot.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "source_db_snapshot_identifier",
				Description: "The DB snapshot ARN that the DB snapshot was copied from.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceDBSnapshotIdentifier"),
			},
			{
				Name:        "db_system_id",
				Description: "The Oracle system identifier (SID), which is the name of the Oracle database instance that manages your database files. The Oracle SID is also the name of your CDB.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSystemId"),
			},
			{
				Name:        "dedicated_log_volume",
				Description: "Indicates whether the DB instance has a dedicated log volume (DLV) enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "multi_tenant",
				Description: "Indicates whether the snapshot is of a DB instance using the multi-tenant configuration (TRUE) or the single-tenant configuration (FALSE).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "snapshot_database_time",
				Description: "The timestamp of the most recent transaction applied to the database that you're backing up.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "snapshot_target",
				Description: "Specifies where manual snapshots are stored: Amazon Web Services Outposts or the Amazon Web Services Region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_region",
				Description: "The AWS Region that the DB snapshot was created in or copied from.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_type",
				Description: "Specifies the storage type associated with DB snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tde_credential_arn",
				Description: "The ARN from the key store with which to associate the instance for TDE encryption.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timezone",
				Description: "The time zone of the DB snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VPC ID associated with the DB snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_snapshot_attributes",
				Description: "A list of DB snapshot attribute names and values for a manual DB snapshot.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSDBSnapshotAttributes,
				Transform:   transform.FromField("DBSnapshotAttributesResult.DBSnapshotAttributes"),
			},
			{
				Name:        "processor_features",
				Description: "The number of CPU cores and the number of threads per core for the DB instance class of the DB instance when the DB snapshot was created.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB snapshot.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getRDSDBSnapshotTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSnapshotIdentifier"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBSnapshotArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_snapshot.listRDSDBSnapshots", "connection_error", err)
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

	input := &rds.DescribeDBSnapshotsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	filters := buildRdsDbSnapshotFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := rds.NewDescribeDBSnapshotsPaginator(svc, input, func(o *rds.DescribeDBSnapshotsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_snapshot.listRDSDBSnapshots", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBSnapshots {
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

func getRDSDBSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	dbSnapshotIdentifier := d.EqualsQuals["db_snapshot_identifier"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_snapshot.getRDSDBSnapshot", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: aws.String(dbSnapshotIdentifier),
	}

	op, err := svc.DescribeDBSnapshots(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_snapshot.getRDSDBSnapshot", "api_error", err)
		return nil, err
	}

	if op.DBSnapshots != nil && len(op.DBSnapshots) > 0 {
		snapshot := op.DBSnapshots[0]
		if isSuppportedRDSEngine(*snapshot.Engine) {
			return snapshot, nil
		}
	}
	return nil, nil
}

func getAwsRDSDBSnapshotAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	dbSnapshot := h.Item.(types.DBSnapshot)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_snapshot.getAwsRDSDBSnapshotAttributes", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBSnapshotAttributesInput{
		DBSnapshotIdentifier: aws.String(*dbSnapshot.DBSnapshotIdentifier),
	}

	op, err := svc.DescribeDBSnapshotAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_snapshot.getAwsRDSDBSnapshotAttributes", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getRDSDBSnapshotTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbSnapshot := d.HydrateItem.(types.DBSnapshot)

	if dbSnapshot.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbSnapshot.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

//// UTILITY FUNCTIONS

// build snapshots list call input filter
func buildRdsDbSnapshotFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)
	filterQuals := map[string]string{
		"db_instance_identifier": "db-instance-id",
		"dbi_resource_id":        "dbi-resource-id",
		"engine":                 "engine",
		"type":                   "snapshot-type",
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
