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

func tableAwsRDSDBInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_instance",
		Description: "AWS RDS DB Instance",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("db_instance_identifier"),
			ShouldIgnoreError: isNotFoundError([]string{"DBInstanceNotFound"}),
			ItemFromKey:       dbInstanceIdentifierFromKey,
			Hydrate:           getRDSDBInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBInstances,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_instance_identifier",
				Description: "The friendly name to identify the DB Instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB Instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceArn"),
			},
			{
				Name:        "db_cluster_identifier",
				Description: "The friendly name to identify the DB cluster, that the DB instance is a member of",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "status",
				Description: "Specifies the current state of this database",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceStatus"),
			},
			{
				Name:        "class",
				Description: "Contains the name of the compute and memory capacity class of the DB instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceClass"),
			},
			{
				Name:        "resource_id",
				Description: "The AWS Region-unique, immutable identifier for the DB instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DbiResourceId"),
			},
			{
				Name:        "allocated_storage",
				Description: "Specifies the allocated storage size specified in gibibytes(GiB)",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "auto_minor_version_upgrade",
				Description: "Specifies whether minor version patches are applied automatically, or not",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "availability_zone",
				Description: "Specifies the name of the Availability Zone the DB instance is located in",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_retention_period",
				Description: "Specifies the number of days for which automatic DB snapshots are retained",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ca_certificate_identifier",
				Description: "The identifier of the CA certificate for this DB instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CACertificateIdentifier"),
			},
			{
				Name:        "character_set_name",
				Description: "Specifies the name of the character set that this instance is associated with",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "copy_tags_to_snapshot",
				Description: "Specifies whether tags are copied from the DB instance to snapshots of the DB instance, or not",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "customer_owned_ip_enabled",
				Description: "Specifies whether a customer-owned IP address (CoIP) is enabled for an RDS on Outposts DB instance, or not",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "port",
				Description: "Specifies the port that the DB instance listens on",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DbInstancePort"),
			},
			{
				Name:        "db_name",
				Description: "Contains the name of the initial database of this instance that was provided at create time",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBName"),
			},
			{
				Name:        "db_subnet_group_arn",
				Description: "The Amazon Resource Name (ARN) for the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.DBSubnetGroupArn"),
			},
			{
				Name:        "db_subnet_group_description",
				Description: "Provides the description of the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.DBSubnetGroupDescription"),
			},
			{
				Name:        "db_subnet_group_name",
				Description: "The name of the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.DBSubnetGroupName"),
			},
			{
				Name:        "db_subnet_group_status",
				Description: "Provides the status of the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.SubnetGroupStatus"),
			},
			{
				Name:        "deletion_protection",
				Description: "Specifies whether the DB instance has deletion protection enabled, or not",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "endpoint_address",
				Description: "Specifies the DNS address of the DB instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Endpoint.Address"),
			},
			{
				Name:        "endpoint_hosted_zone_id",
				Description: "Specifies the ID that Amazon Route 53 assigns when you create a hosted zone",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Endpoint.HostedZoneId"),
			},
			{
				Name:        "endpoint_port",
				Description: "Specifies the port that the database engine is listening on",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Endpoint.Port"),
			},
			{
				Name:        "engine",
				Description: "The name of the database engine to be used for this DB instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Indicates the database engine version",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enhanced_monitoring_resource_arn",
				Description: "The ARN of the Amazon CloudWatch Logs log stream that receives the Enhanced Monitoring metrics data for the DB instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_database_authentication_enabled",
				Description: "Specifies whether the the mapping of AWS IAM accounts to database accounts is enabled, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IAMDatabaseAuthenticationEnabled"),
			},
			{
				Name:        "create_time",
				Description: "Provides the date and time the DB instance was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("InstanceCreateTime"),
			},
			{
				Name:        "iops",
				Description: "Specifies the Provisioned IOPS (I/O operations per second) value",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS KMS key identifier for the encrypted DB instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_restorable_time",
				Description: "Specifies the latest time to which a database can be restored with point-in-time restore",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "license_model",
				Description: "License model information for this DB instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_user_name",
				Description: "Contains the master username for the DB instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "max_allocated_storage",
				Description: "The upper limit to which Amazon RDS can automatically scale the storage of the DB instance",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "monitoring_interval",
				Description: "The interval, in seconds, between points when Enhanced Monitoring metrics are collected for the DB instance",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "monitoring_role_arn",
				Description: "The ARN for the IAM role that permits RDS to send Enhanced Monitoring metrics to Amazon CloudWatch Logs",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "multi_az",
				Description: "Specifies if the DB instance is a Multi-AZ deployment",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("MultiAZ"),
			},
			{
				Name:        "nchar_character_set_name",
				Description: "The name of the NCHAR character set for the Oracle DB instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "performance_insights_enabled",
				Description: "Specifies whether Performance Insights is enabled for the DB instance, or not",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "performance_insights_kms_key_id",
				Description: "The AWS KMS key identifier for encryption of Performance Insights data",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PerformanceInsightsKMSKeyId"),
			},
			{
				Name:        "performance_insights_retention_period",
				Description: "The amount of time, in days, to retain Performance Insights data",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "preferred_backup_window",
				Description: "Specifies the daily time range during which automated backups are created",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "Specifies the weekly time range during which system maintenance can occur",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "promotion_tier",
				Description: "Specifies the order in which an Aurora Replica is promoted to the primary instance after a failure of the existing primary instance",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "publicly_accessible",
				Description: "Specifies the accessibility options for the DB instance",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "read_replica_source_db_instance_identifier",
				Description: "Contains the identifier of the source DB instance if this DB instance is a read replica",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReadReplicaSourceDBInstanceIdentifier"),
			},
			{
				Name:        "replica_mode",
				Description: "The mode of an Oracle read replica",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secondary_availability_zone",
				Description: "Specifies the name of the secondary Availability Zone for a DB instance with multi-AZ support",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the DB instance is encrypted, or not",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "storage_type",
				Description: "Specifies the storage type associated with DB instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tde_credential_arn",
				Description: " The ARN from the key store with which the instance is associated for TDE encryption",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timezone",
				Description: "The time zone of the DB instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VpcId of the DB subnet group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.VpcId"),
			},
			{
				Name:        "associated_roles",
				Description: "A list of AWS IAM roles that are associated with the DB instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_parameter_groups",
				Description: "A list of DB parameter groups applied to this DB instance",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBParameterGroups"),
			},
			{
				Name:        "db_security_groups",
				Description: "A list of DB security group associated with the DB instance",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBSecurityGroups"),
			},
			{
				Name:        "domain_memberships",
				Description: "A list of Active Directory Domain membership records associated with the DB instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "enabled_cloudwatch_logs_exports",
				Description: "A list of log types that this DB instance is configured to export to CloudWatch Logs",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "option_group_memberships",
				Description: "A list of option group memberships for this DB instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "processor_features",
				Description: "The number of CPU cores and the number of threads per core for the DB instance class of the DB instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "read_replica_db_cluster_identifiers",
				Description: "A list of identifiers of Aurora DB clusters to which the RDS DB instance is replicated as a read replica",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ReadReplicaDBClusterIdentifiers"),
			},
			{
				Name:        "read_replica_db_instance_identifiers",
				Description: "A list of identifiers of the read replicas associated with this DB instance",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ReadReplicaDBInstanceIdentifiers"),
			},
			{
				Name:        "status_infos",
				Description: "The status of a read replica",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnets",
				Description: "A list of subnet elements",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBSubnetGroup.Subnets"),
			},
			{
				Name:        "vpc_security_groups",
				Description: "A list of VPC security group elements that the DB instance belongs to",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB Instance",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getRDSDBInstanceTurbotTags),
			},
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
				Transform:   transform.FromField("DBInstanceArn").Transform(arnToAkas),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func dbInstanceIdentifierFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	dbInstanceIdentifier := quals["db_instance_identifier"].GetStringValue()
	item := &rds.DBInstance{
		DBInstanceIdentifier: &dbInstanceIdentifier,
	}
	return item, nil
}

//// LIST FUNCTION

func listRDSDBInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listRDSDBInstances", "AWS_REGION", region)

	// Create Session
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	var params []*rds.Filter
	if region == "sa-east-1" {
		params = []*rds.Filter{
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
	} else if region == "eu-north-1" {
		params = []*rds.Filter{
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
		params = []*rds.Filter{
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

	// List call
	err = svc.DescribeDBInstancesPages(
		&rds.DescribeDBInstancesInput{
			Filters: params,
		},
		func(page *rds.DescribeDBInstancesOutput, isLast bool) bool {
			for _, dbInstance := range page.DBInstances {
				d.StreamListItem(ctx, dbInstance)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	dbInstance := h.Item.(*rds.DBInstance)

	// Create service
	svc, err := RDSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(*dbInstance.DBInstanceIdentifier),
	}

	op, err := svc.DescribeDBInstances(params)
	if err != nil {
		return nil, err
	}

	if op.DBInstances != nil && len(op.DBInstances) > 0 {
		return op.DBInstances[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS ////

func getRDSDBInstanceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbInstance := d.HydrateItem.(*rds.DBInstance)

	if dbInstance.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbInstance.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
