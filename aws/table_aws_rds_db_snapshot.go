package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsRDSDBSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_snapshot",
		Description: "AWS RDS DB Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_snapshot_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"DBSnapshotNotFound"}),
			},
			Hydrate: getRDSDBSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBSnapshots,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_instance_identifier", Require: plugin.Optional},
				{Name: "dbi_resource_id", Require: plugin.Optional},
				{Name: "engine", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listRDSDBSnapshots")

	// Create Session
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &rds.DescribeDBSnapshotsInput{
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

	filters := buildRdsDbSnapshotFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// List call
	err = svc.DescribeDBSnapshotsPages(
		input,
		func(page *rds.DescribeDBSnapshotsOutput, isLast bool) bool {
			for _, dbSnapshot := range page.DBSnapshots {
				d.StreamListItem(ctx, dbSnapshot)

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

func getRDSDBSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	dbSnapshotIdentifier := d.KeyColumnQuals["db_snapshot_identifier"].GetStringValue()

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: aws.String(dbSnapshotIdentifier),
	}

	op, err := svc.DescribeDBSnapshots(params)
	if err != nil {
		return nil, err
	}

	if op.DBSnapshots != nil && len(op.DBSnapshots) > 0 {
		return op.DBSnapshots[0], nil
	}
	return nil, nil
}

func getAwsRDSDBSnapshotAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRDSDBSnapshotAttributes")

	dbSnapshot := h.Item.(*rds.DBSnapshot)

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBSnapshotAttributesInput{
		DBSnapshotIdentifier: aws.String(*dbSnapshot.DBSnapshotIdentifier),
	}

	op, err := svc.DescribeDBSnapshotAttributes(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getRDSDBSnapshotTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbSnapshot := d.HydrateItem.(*rds.DBSnapshot)

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
func buildRdsDbSnapshotFilter(quals plugin.KeyColumnQualMap) []*rds.Filter {
	filters := make([]*rds.Filter, 0)
	filterQuals := map[string]string{
		"db_instance_identifier": "db-instance-id",
		"dbi_resource_id":        "dbi-resource-id",
		"engine":                 "engine",
		"type":                   "snapshot-type",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := rds.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []*string{aws.String(val)}
			} else {
				v := value.([]*string)
				filter.Values = v
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
