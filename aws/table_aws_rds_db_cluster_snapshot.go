package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBClusterSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_cluster_snapshot",
		Description: "AWS RDS DB Cluster Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_cluster_snapshot_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBSnapshotNotFound", "DBClusterSnapshotNotFoundFault"}),
			},
			Hydrate: getRDSDBClusterSnapshot,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBClusterSnapshots"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBClusterSnapshots,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBClusterSnapshots"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_cluster_identifier", Require: plugin.Optional},
				{Name: "db_cluster_snapshot_identifier", Require: plugin.Optional},
				{Name: "engine", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsRDSDBClusterSnapshotAttributes,
				Tags: map[string]string{"service": "rds", "action": "DescribeDBClusterSnapshotAttributes"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsEndpoint.RDSServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_cluster_snapshot_identifier",
				Description: "The friendly name to identify the DB Cluster Snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotIdentifier"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the DB Cluster Snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterSnapshotArn"),
			},
			{
				Name:        "type",
				Description: "The type of the DB Cluster Snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotType"),
			},
			{
				Name:        "db_cluster_resource_id",
				Description: "The resource ID of the DB cluster that this DB cluster snapshot was created from.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_mode",
				Description: "The engine mode of the database engine for this DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "snapshot_type",
				Description: "The type of the DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the status of this DB Cluster Snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_cluster_identifier",
				Description: "The friendly name to identify the DB Cluster, that the snapshot snapshot was created from.",
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
				Name:        "allocated_storage",
				Description: "Specifies the allocated storage size in gibibytes (GiB).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cluster_create_time",
				Description: "Specifies the time when the DB cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "engine",
				Description: "Specifies the name of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "Specifies the version of the database engine for this DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_database_authentication_enabled",
				Description: "Specifies whether mapping of AWS Identity and Access Management (IAM) accounts to database accounts is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IAMDatabaseAuthenticationEnabled"),
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS KMS key identifier for the AWS KMS customer master key (CMK).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_model",
				Description: "Provides the license model information for this DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_user_name",
				Description: "Provides the master username for the DB cluster snapshot.",
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
				Description: "Specifies the port that the DB cluster was listening on at the time of the snapshot.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "source_db_cluster_snapshot_arn",
				Description: "The Amazon Resource Name (ARN) for the source DB cluster snapshot, if the DB cluster snapshot was copied from a source DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceDBClusterSnapshotArn"),
			},
			{
				Name:        "storage_encrypted",
				Description: "Specifies whether the DB cluster snapshot is encrypted, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VPC ID associated with the DB cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zones",
				Description: "A list of Availability Zones (AZs) where instances in the DB cluster snapshot can be restored.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_cluster_snapshot_attributes",
				Description: "A list of DB cluster snapshot attribute names and values for a manual DB cluster snapshot.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSDBClusterSnapshotAttributes,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the DB Cluster Snapshot.",
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

//// LIST FUNCTION

func listRDSDBClusterSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_snapshot.listRDSDBClusterSnapshots", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	// select * from aws_rds_db_cluster_snapshot limit 3
	// Error: InvalidParameterValue: Invalid value 3 for MaxRecords. Must be between 20 and 100
	// 	status code: 400, request id: c39eead1-96e0-49c8-a927-aa9a3131836d
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

	input := &rds.DescribeDBClusterSnapshotsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	filters := buildRdsDbClusterSnapshotFilter(d.Quals)

	if len(filters) != 0 {
		input.Filters = filters
	}

	paginator := rds.NewDescribeDBClusterSnapshotsPaginator(svc, input, func(o *rds.DescribeDBClusterSnapshotsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_cluster_snapshot.listRDSDBClusterSnapshots", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBClusterSnapshots {
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

func getRDSDBClusterSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	snapshotIdentifier := d.EqualsQuals["db_cluster_snapshot_identifier"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_snapshot.getRDSDBClusterSnapshot", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: aws.String(snapshotIdentifier),
	}

	op, err := svc.DescribeDBClusterSnapshots(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_snapshot.getRDSDBClusterSnapshot", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.DBClusterSnapshots) > 0 {
		snapshot := op.DBClusterSnapshots[0]
		if isSuppportedRDSEngine(*snapshot.Engine) {
			return snapshot, nil
		}
	}
	return nil, nil
}

func getAwsRDSDBClusterSnapshotAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	dbClusterSnapshot := h.Item.(types.DBClusterSnapshot)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_snapshot.getAwsRDSDBClusterSnapshotAttributes", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: dbClusterSnapshot.DBClusterSnapshotIdentifier,
	}

	dbClusterSnapshotData, err := svc.DescribeDBClusterSnapshotAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_cluster_snapshot.getAwsRDSDBClusterSnapshotAttributes", "api_error", err)
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

//// TRANSFORM FUNCTIONS

func getRDSDBClusterSnapshotTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dbClusterSnapshot := d.HydrateItem.(types.DBClusterSnapshot)

	if dbClusterSnapshot.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range dbClusterSnapshot.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

//// UTILITY FUNCTIONS

// build snapshots list call input filter
func buildRdsDbClusterSnapshotFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)
	filterQuals := map[string]string{
		"db_cluster_identifier":          "db-cluster-id",
		"db_cluster_snapshot_identifier": "db-cluster-snapshot-id",
		"engine":                         "engine",
		"type":                           "snapshot-type",
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
