package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/drs"
	"github.com/aws/aws-sdk-go-v2/service/drs/types"

	drsEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDRSRecoverySnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_drs_recovery_snapshot",
		Description: "AWS DRS Recovery Snapshot",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "source_server_id", Require: plugin.Optional},
				{Name: "timestamp", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				// UninitializedAccountException - This error comes up when default replication settings are not set for a particular region.
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UninitializedAccountException", "BadRequestException"}),
			},
			ParentHydrate: listAwsDRSSourceServers,
			Hydrate:       listAwsDRSRecoverySnapshots,
			Tags:          map[string]string{"service": "drs", "action": "DescribeRecoverySnapshots"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(drsEndpoint.AWS_DRS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "snapshot_id",
				Description: "The ID of the snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotID"),
			},
			{
				Name:        "source_server_id",
				Description: "The ID of the source server that the snapshot was taken for.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceServerID"),
			},
			{
				Name:        "expected_timestamp",
				Description: "The timestamp of when we expect the snapshot to be taken.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "timestamp",
				Description: "The actual timestamp when the snapshot was taken.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "ebs_snapshots",
				Description: "A list of EBS snapshots.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotID"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsDRSRecoverySnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var sourceServerID string
	if h.Item != nil {
		sourceServerID = *h.Item.(types.SourceServer).SourceServerID
	}

	if d.EqualsQualString("source_server_id") != "" {
		if sourceServerID != d.EqualsQualString("source_server_id") {
			return nil, nil
		}
	}
	// Create service
	svc, err := DRSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_drs_recovery_snapshot.listAwsDRSRecoverySnapshots", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// The API has no limit on MaxResults, so use 1000 as a sensible default
	maxItems := int32(1000)
	input := drs.DescribeRecoverySnapshotsInput{
		SourceServerID: aws.String(sourceServerID),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)

	filter := &types.DescribeRecoverySnapshotsRequestFilters{}

	quals := d.Quals
	if quals["timestamp"] != nil {
		for _, q := range quals["timestamp"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime().Format(time.RFC3339)
			switch q.Operator {
			case ">=", ">":
				filter.FromDateTime = aws.String(timestamp)
			case "<", "<=":
				filter.ToDateTime = aws.String(timestamp)
			case "=":
				filter.FromDateTime = aws.String(timestamp)
				filter.ToDateTime = aws.String(timestamp)
			}
		}
	}

	input.Filters = filter

	paginator := drs.NewDescribeRecoverySnapshotsPaginator(svc, &input, func(o *drs.DescribeRecoverySnapshotsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_drs_recovery_snapshot.listAwsDRSRecoverySnapshots", "api_error", err)
			return nil, err
		}

		for _, recoverySnapshot := range output.Items {
			d.StreamListItem(ctx, recoverySnapshot)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
