package aws

import (
	"context"

	"github.com/turbot/go-kit/helpers"
	go_kit "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"

	rdsEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAwsDocDBClusterSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_docdb_cluster_snapshot",
		Description: "AWS DocumentDB Cluster Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_cluster_snapshot_identifier"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBSnapshotNotFound", "DBClusterSnapshotNotFoundFault"}),
			},
			Hydrate: getDocDBClusterSnapshot,
			Tags:    map[string]string{"service": "docdb-elastic", "action": "DescribeDBClusterSnapshots"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDocDBClusterSnapshots,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "db_cluster_identifier", Require: plugin.Optional},
				{Name: "snapshot_type", Require: plugin.Optional},
				{Name: "include_public", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "include_shared", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Tags: map[string]string{"service": "docdb-elastic", "action": "DescribeDBClusterSnapshots"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getDocDBClusterTags,
				Tags: map[string]string{"service": "docdb-elastic", "action": "ListTagsForResource"},
			},
			{
				Func: getAwsDocDBClusterSnapshotAttributes,
				Tags: map[string]string{"service": "docdb-elastic", "action": "DescribeDBClusterSnapshotAttributes"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsEndpoint.RDSServiceID),
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
				Name:        "snapshot_type",
				Description: "The type of the cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the status of this cluster snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "include_public",
				Description: "Set to true to include manual cluster snapshots that are public and can be copied or restored by any Amazon Web Services account, and otherwise false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("include_public"),
				Default:     false,
			},
			{
				Name:        "include_shared",
				Description: "Set to true to include shared manual cluster snapshots from other Amazon Web Services accounts that this Amazon Web Services account has been given permission to copy or restore, and otherwise false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("include_shared"),
				Default:     false,
			},
			{
				Name:        "db_cluster_identifier",
				Description: "The friendly name to identify the cluster, that the snapshot snapshot was created from.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "snapshot_create_time",
				Description: "The time when the snapshot was taken.",
				Type:        proto.ColumnType_TIMESTAMP,
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
	// Create Session
	svc, err := DocDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.listDocDBClusterSnapshots", "connection", err)
		return nil, err
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
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

	input := docdb.DescribeDBClusterSnapshotsInput{
		MaxRecords: go_kit.Int32(maxLimit),
	}

	if d.EqualsQualString("db_cluster_identifier") != "" {
		input.DBClusterIdentifier = aws.String(d.EqualsQualString("db_cluster_identifier"))
	}
	if d.EqualsQualString("snapshot_type") != "" {
		input.SnapshotType = aws.String(d.EqualsQualString("snapshot_type"))
	}
	if d.EqualsQuals["include_public"] != nil {
		input.IncludePublic = aws.Bool(d.EqualsQuals["include_public"].GetBoolValue())
	}
	if d.EqualsQuals["include_shared"] != nil {
		input.IncludePublic = aws.Bool(d.EqualsQuals["include_shared"].GetBoolValue())
	}

	// List call
	paginator := docdb.NewDescribeDBClusterSnapshotsPaginator(svc, &input, func(o *docdb.DescribeDBClusterSnapshotsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.listDocDBClusterSnapshots", "api_error", err)
			return nil, err
		}

		for _, cluster := range output.DBClusterSnapshots {
			// The DescribeDBClusters API returns non-DocDB clusters as well, but we only want DocDB clusters here.
			if helpers.StringSliceContains([]string{"docdb"}, *cluster.Engine) {
				d.StreamListItem(ctx, cluster)
			}

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDocDBClusterSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	snapshotIdentifier := d.EqualsQualString("db_cluster_snapshot_identifier")

	// Create service
	svc, err := DocDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.getDocDBClusterSnapshot", "connection_error", err)
		return nil, err
	}

	params := &docdb.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: aws.String(snapshotIdentifier),
	}

	op, err := svc.DescribeDBClusterSnapshots(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.getDocDBClusterSnapshot", "api_error", err)
		return nil, err
	}

	if len(op.DBClusterSnapshots) > 0 {
		return op.DBClusterSnapshots[0], nil
	}
	return nil, nil
}

func getAwsDocDBClusterSnapshotAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dbClusterSnapshot := h.Item.(types.DBClusterSnapshot)

	// Create service
	svc, err := DocDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.getAwsDocDBClusterSnapshotAttributes", "connection_error", err)
		return nil, err
	}

	params := &docdb.DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: dbClusterSnapshot.DBClusterSnapshotIdentifier,
	}

	dbClusterSnapshotData, err := svc.DescribeDBClusterSnapshotAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.getAwsDocDBClusterSnapshotAttributes", "api_error", err)
		return nil, err
	}

	return dbClusterSnapshotData, nil
}

func getDocDBClusterSnapshotTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(types.DBClusterSnapshot)

	// Create Session
	svc, err := DocDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.getDocDBClusterSnapshotTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &docdb.ListTagsForResourceInput{
		ResourceName: cluster.DBClusterSnapshotArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.getDocDBClusterSnapshotTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getDocDBClusterSnapshotTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

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
