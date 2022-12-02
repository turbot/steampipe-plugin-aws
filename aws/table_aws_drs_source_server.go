package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/drs"
	"github.com/aws/aws-sdk-go-v2/service/drs/types"
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
				{Name: "staging_account_id", Require: plugin.Optional},
				{Name: "hardware_id", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				// UninitializedAccountException - This error comes up when default replication settings are not set for a particular region.
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UninitializedAccountException", "BadRequestException"}),
			},
			Hydrate: listAwsDRSSourceServers,
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
			{
				Name:        "staging_account_id",
				Description: "The staging account ID that extended source servers belong to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StagingArea.StagingAccountID"),
			},
			{
				Name:        "hardware_id",
				Description: "An ID that describes the hardware of the Source Server. This is either an EC2 instance id, a VMware uuid or a mac address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("hardware_id"),
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

func listAwsDRSSourceServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := DRSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_drs_source_server.listAwsDRSSourceServers", "connection_error", err)
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
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = int32(maxItems)
	sourceServerId := d.KeyColumnQualString("source_server_id")
	stagingAccountId := d.KeyColumnQualString("staging_account_id")
	hardwareId := d.KeyColumnQualString("hardware_id")

	filter := &types.DescribeSourceServersRequestFilters{}

	if sourceServerId != "" {
		filter.SourceServerIDs = []string{sourceServerId}
	}

	if stagingAccountId != "" {
		filter.StagingAccountIDs = []string{stagingAccountId}
	}

	if hardwareId != "" {
		filter.HardwareId = aws.String(hardwareId)
	}

	input.Filters = filter

	paginator := drs.NewDescribeSourceServersPaginator(svc, &input, func(o *drs.DescribeSourceServersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_drs_source_server.listAwsDRSSourceServers", "api_error", err)
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
