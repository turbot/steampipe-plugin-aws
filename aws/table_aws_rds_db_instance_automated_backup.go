package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBInstanceAutomatedBackup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_instance_automated_backup",
		Description: "AWS RDS DB Instance Automated Backup",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBInstanceAutomatedBackupNotFound"}),
			},
			Hydrate: getRDSDBInstanceAutomatedBackup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBInstanceAutomatedBackups,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_instance_identifier", Require: plugin.Optional},
				{Name: "dbi_resource_id", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_instance_identifier",
				Description: "The friendly name to identify the DB Instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the replicated automated backups.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceAutomatedBackupsArn"),
			},
			{
				Name:        "db_instance_arn",
				Description: "The Amazon Resource Name (ARN) for the automated backups.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceArn"),
			},
			{
				Name:        "status",
				Description: "Specifies the current state of this database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allocated_storage",
				Description: "Specifies the allocated storage size in gibibytes (GiB).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone that the automated backup was created in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_retention_period",
				Description: "The retention period for the automated backups.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "backup_target",
				Description: "Specifies where automated backups are stored: Amazon Web Services Outposts or the Amazon Web Services Region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dbi_resource_id",
				Description: "The identifier for the source DB instance, which can't be changed and which is unique to an Amazon Web Services Region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encrypted",
				Description: "Specifies whether the automated backup is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "engine",
				Description: "The name of the database engine for this automated backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "The version of the database engine for the automated backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_database_authentication_enabled",
				Description: "True if mapping of Amazon Web Services Identity and Access Management (IAM) accounts to database accounts is enabled, and otherwise false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IAMDatabaseAuthenticationEnabled"),
			},
			{
				Name:        "instance_create_time",
				Description: "True if mapping of Amazon Web Services Identity and Access Management (IAM) accounts to database accounts is enabled, and otherwise false.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "iops",
				Description: "True if mapping of Amazon Web Services Identity and Access Management (IAM) accounts to database accounts is enabled, and otherwise false.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kms_key_id",
				Description: "The Amazon Web Services KMS key ID for an automated backup. The Amazon Web Services KMS key identifier is the key ARN, key ID, alias ARN, or alias name for the KMS key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_model",
				Description: "The Amazon Web Services KMS key ID for an automated backup. The Amazon Web Services KMS key identifier is the key ARN, key ID, alias ARN, or alias name for the KMS key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_username",
				Description: "The license model of an automated backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "option_group_name",
				Description: "The option group the automated backup is associated with. If omitted, the default option group for the engine specified is used.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The port number that the automated backup used for connections. Default: Inherits from the source DB instance Valid Values: 1150-65535.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "storage_throughput",
				Description: "Specifies the storage throughput for the automated backup.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "storage_type",
				Description: "Specifies the storage type associated with the automated backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tde_credential_arn",
				Description: "The ARN from the key store with which the automated backup is associated for TDE encryption.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timezone",
				Description: "The time zone of the automated backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VPC ID associated with the DB instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_instance_automated_backups_replications",
				Description: "The list of replications to different Amazon Web Services Regions associated with the automated backup.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBInstanceAutomatedBackupsReplications"),
			},
			{
				Name:        "restore_window",
				Description: "Earliest and latest time an instance can be restored to.",
				Type:        proto.ColumnType_JSON,
			},

			// Stteampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceIdentifier"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBInstanceAutomatedBackupsArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBInstanceAutomatedBackups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_instance_automated_backup.listRDSDBInstanceAutomatedBackups", "connection_error", err)
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

	input := &rds.DescribeDBInstanceAutomatedBackupsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	filters := buildRdsDbInstanceAutomatedBackupFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := rds.NewDescribeDBInstanceAutomatedBackupsPaginator(svc, input, func(o *rds.DescribeDBInstanceAutomatedBackupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_instance_automated_backup.listRDSDBInstanceAutomatedBackups", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBInstanceAutomatedBackups {
			if helpers.StringSliceContains(
				[]string{
					"aurora",
					"aurora-mysql",
					"aurora-postgresql",
					"mysql",
					"postgres",
				},
				*items.Engine) {
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

func getRDSDBInstanceAutomatedBackup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQualString("arn")

	// empty arn check
	if arn == "" {
		return nil, nil
	}

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_instance_automated_backup.getRDSDBInstanceAutomatedBackup", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBInstanceAutomatedBackupsInput{
		DBInstanceAutomatedBackupsArn: aws.String(arn),
	}

	op, err := svc.DescribeDBInstanceAutomatedBackups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_instance_automated_backup.getRDSDBInstanceAutomatedBackup", "api_error", err)
		return nil, err
	}

	if op.DBInstanceAutomatedBackups != nil && len(op.DBInstanceAutomatedBackups) > 0 {
		backup := op.DBInstanceAutomatedBackups[0]
		if helpers.StringSliceContains(
			[]string{
				"aurora",
				"aurora-mysql",
				"aurora-postgresql",
				"mysql",
				"postgres",
			},
			*backup.Engine) {
			return backup, nil
		}
	}
	return nil, nil
}

//// UTILITY FUNCTIONS

// build rds db instance automated backup list call input filter
func buildRdsDbInstanceAutomatedBackupFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)
	filterQuals := map[string]string{
		"db_instance_identifier": "db-instance-id",
		"dbi_resource_id":        "dbi-resource-id",
		"status":                 "status",
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
