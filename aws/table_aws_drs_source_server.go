package aws

import (
	"context"
	"math"

	"github.com/aws/aws-sdk-go-v2/service/drs"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsDRSSourceServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_drs_source_server",
		Description: "AWS DRS Source Server",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "source_server_id", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UninitializedAccountException"}),
			},
			Hydrate: listAwsDRSSourceServer,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "source_server_id",
				Description: "The ID of the Source Server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceServerID"),
			},
			{
				Name:        "arn",
				Description: "The ARN of the Source Server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recovery_instance_id",
				Description: "The ID of the Recovery Instance associated with this Source Server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecoveryInstanceId"),
			},
			{
				Name:        "source_properties",
				Description: "The source properties of the Source Server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "data_replication_info",
				Description: "The Data Replication Info of the Source Server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_launch_result",
				Description: "The status of the last recovery launch of this Source Server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "life_cycle",
				Description: "The lifecycle information of this Source Server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replication_direction",
				Description: "Replication direction of the Source Server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "reversed_direction_source_server_arn",
				Description: "For EC2-originated Source Servers which have been failed over and then failed back, this value will mean the ARN of the Source Server on the opposite replication direction.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReversedDirectionSourceServerArn"),
			},
			{
				Name:        "source_cloud_properties",
				Description: "Source cloud properties of the Source Server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "staging_area",
				Description: "The staging area of the source server.",
				Type:        proto.ColumnType_JSON,
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceServerID"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsDRSSourceServer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := DRSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_drs_source_server.listAwsDRSSourceServer", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(10000)
	input := drs.DescribeSourceServersInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		maxItems = int32(math.Min(float64(maxItems), float64(limit)))
		maxItems = int32(math.Max(1, float64(maxItems)))
	}
	input.MaxResults = int32(maxItems)

	sourceServerId := d.KeyColumnQualString("source_server_id")

	if sourceServerId != "" {
		input.Filters.SourceServerIDs = []string{sourceServerId}
	}

	paginator := drs.NewDescribeSourceServersPaginator(svc, &input, func(o *drs.DescribeSourceServersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_drs_source_server.listAwsDRSSourceServer", "api_error", err)
			return nil, err
		}

		for _, sourceServer := range output.Items {
			d.StreamListItem(ctx, sourceServer)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
