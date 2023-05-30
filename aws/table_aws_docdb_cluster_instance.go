package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"

	docdbv1 "github.com/aws/aws-sdk-go/service/docdb"

	"github.com/turbot/go-kit/helpers"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDocDBClusterInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_docdb_cluster_instance",
		Description: "AWS DocumentDB Cluster Instance",
		List: &plugin.ListConfig{
			Hydrate: listDocDBClusterInstances,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue", "InvalidParameterCombination", "DBInstanceNotFound"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_cluster_identifier", Require: plugin.Optional},
				{Name: "db_instance_identifier", Require: plugin.Optional},
				{Name: "db_instance_arn", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(docdbv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_instance_identifier",
				Description: "Contains a user-provided database identifier. This identifier is the unique key that identifies an instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceIdentifier"),
			},
			{
				Name:        "db_instance_arn",
				Description: "The Amazon Resource Name (ARN) for the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceArn"),
			},
			{
				Name:        "db_cluster_identifier",
				Description: "Contains the name of the cluster that the instance is a member of if the instance is a member of a cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "db_instance_status",
				Description: "Specifies the current state of this database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceStatus"),
			},
			{
				Name:        "db_instance_class",
				Description: "Contains the name of the compute and memory capacity class of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceClass"),
			},
			{
				Name:        "dbi_resource_id",
				Description: "The Amazon Web Services Region-unique, immutable identifier for the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DbiResourceId"),
			},
			{
				Name:        "availability_zone",
				Description: "Specifies the name of the availability zone the instance is located in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_retention_period",
				Description: "Specifies the number of days for which automatic snapshots are retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ca_certificate_identifier",
				Description: "The identifier of the CA certificate for this DB instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CACertificateIdentifier"),
			},
			{
				Name:        "copy_tags_to_snapshot",
				Description: "Specifies whether tags are copied from the DB instance to snapshots of the DB instance, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "db_subnet_group_arn",
				Description: "The Amazon Resource Name (ARN) for the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.DBSubnetGroupArn"),
			},
			{
				Name:        "db_subnet_group_description",
				Description: "Provides the description of the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.DBSubnetGroupDescription"),
			},
			{
				Name:        "db_subnet_group_name",
				Description: "The name of the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.DBSubnetGroupName"),
			},
			{
				Name:        "db_subnet_group_status",
				Description: "Provides the status of the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.SubnetGroupStatus"),
			},
			{
				Name:        "endpoint_address",
				Description: "Specifies the DNS address of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Endpoint.Address"),
			},
			{
				Name:        "endpoint_hosted_zone_id",
				Description: "Specifies the ID that Amazon Route 53 assigns when you create a hosted zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Endpoint.HostedZoneId"),
			},
			{
				Name:        "endpoint_port",
				Description: "Specifies the port that the database engine is listening on.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Endpoint.Port"),
			},
			{
				Name:        "engine",
				Description: "The name of the database engine to be used for this instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Indicates the database engine version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_create_time",
				Description: "Provides the date and time the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kms_key_id",
				Description: "If StorageEncrypted is true, the KMS key identifier for the encrypted instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_restorable_time",
				Description: "Specifies the latest time to which a database can be restored with point-in-time restore.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "preferred_backup_window",
				Description: "Specifies the daily time range during which automated backups are created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "Specifies the weekly time range during which system maintenance can occur.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "promotion_tier",
				Description: "A value that specifies the order in which an Amazon DocumentDB replica is promoted to the primary instance after a failure of the existing primary instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "publicly_accessible",
				Description: "Specifies the accessibility options for the DB instance.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether or not the instance is encrypted.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VpcId of the DB subnet group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBSubnetGroup.VpcId"),
			},
			{
				Name:        "enabled_cloudwatch_logs_exports",
				Description: "A list of log types that this instance is configured to export to CloudWatch Logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pending_modified_values",
				Description: "Specifies that changes to the instance are pending.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status_infos",
				Description: "The status of a read replica.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnets",
				Description: "A list of subnet elements.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBSubnetGroup.Subnets"),
			},
			{
				Name:        "vpc_security_groups",
				Description: "A list of VPC security group elements that the instance belongs to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the Instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDocDBClusterInstanceTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDocDBClusterInstanceTags,
				Transform:   transform.From(DocDBClusterInstanceTagListToTurbotTags),
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

//// LIST FUNCTION

func listDocDBClusterInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := DocDBClient(ctx, d)
	if err != nil {
		logger.Error("aws_docdb_cluster_instance.listDocDBClusterInstances", "client_error", err)
		return nil, err
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	input := docdb.DescribeDBInstancesInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	filters := buildDocDBInstanceFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := docdb.NewDescribeDBInstancesPaginator(svc, &input, func(o *docdb.DescribeDBInstancesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_docdb_cluster_instance.listDocDBClusterInstances", "api_error", err)
			return nil, err
		}

		for _, instance := range output.DBInstances {
			// The DescribeDBInstances API returns non-DocDB clusters as well, but we only want DocDB clusters here.
			if helpers.StringSliceContains([]string{"docdb"}, *instance.Engine) {
				d.StreamListItem(ctx, instance)
			}

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDocDBClusterInstanceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	cluster := h.Item.(types.DBInstance)

	// Create Session
	svc, err := DocDBClient(ctx, d)
	if err != nil {
		logger.Error("aws_docdb_cluster_instance.getDocDBClusterInstanceTags", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &docdb.ListTagsForResourceInput{
		ResourceName: cluster.DBInstanceArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		logger.Error("aws_docdb_cluster_instance.getDocDBClusterInstanceTags", "api_error", err)
		return nil, err
	}

	return op.TagList, nil
}

//// TRANSFORM FUNCTIONS

func DocDBClusterInstanceTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

// build instance list call input filter
func buildDocDBInstanceFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)
	filterQuals := map[string]string{
		"db_cluster_identifier":  "db-cluster-id",
		"db_instance_identifier": "db-instance-id",
		"db_instance_arn":        "db-instance-id",
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
