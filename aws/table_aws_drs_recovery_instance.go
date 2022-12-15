package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/drs"
	"github.com/aws/aws-sdk-go-v2/service/drs/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsDRSRecoveryInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_drs_recovery_instance",
		Description: "AWS DRS recovery instance",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "source_server_id", Require: plugin.Optional},
				{Name: "recovery_instance_id", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				// UninitializedAccountException - This error comes up when default replication settings are not set for a particular region.
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UninitializedAccountException", "BadRequestException"}),
			},
			Hydrate: listAwsDRSRecoveryInstances,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "recovery_instance_id",
				Description: "The ID of the recovery instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecoveryInstanceID"),
			},
			{
				Name:        "arn",
				Description: "The ARN of the recovery instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_server_id",
				Description: "The Source Server ID that this recovery instance is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceServerID"),
			},
			{
				Name:        "ec2_instance_id",
				Description: "The EC2 instance ID of the recovery instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Ec2InstanceID"),
			},
			{
				Name:        "ec2_instance_state",
				Description: "The state of the EC2 instance for this recovery instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_drill",
				Description: "Whether this recovery instance was created for a drill or for an actual recovery event.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "job_id",
				Description: "The ID of the Job that created the recovery instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobID"),
			},
			{
				Name:        "origin_environment",
				Description: "Environment (On Premises / AWS) of the instance that the recovery instance originated from.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "point_in_time_snapshot_date_time",
				Description: "The date and time of the Point in Time (PIT) snapshot that this recovery instance was launched from.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_replication_info",
				Description: "The Data Replication Info of the recovery instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "failback",
				Description: "An object representing failback related information of the recovery instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "recovery_instance_properties",
				Description: "Properties of the recovery instance machine.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecoveryInstanceID"),
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
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsDRSRecoveryInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := DRSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_drs_recovery_instance.listAwsDRSRecoveryInstances", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// The API has no limit on MaxResults and the default max number of recovery instance per region is 3000, so use 1000 as a sensible default
	maxItems := int32(1000)
	input := drs.DescribeRecoveryInstancesInput{}

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
	sourceServerID := d.KeyColumnQualString("source_server_id")
	recoveryInstanceId := d.KeyColumnQualString("recovery_instance_id")

	filter := &types.DescribeRecoveryInstancesRequestFilters{}

	if sourceServerID != "" {
		filter.SourceServerIDs = []string{sourceServerID}
	}

	if recoveryInstanceId != "" {
		filter.RecoveryInstanceIDs = []string{recoveryInstanceId}
	}

	input.Filters = filter

	paginator := drs.NewDescribeRecoveryInstancesPaginator(svc, &input, func(o *drs.DescribeRecoveryInstancesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_drs_recovery_instance.listAwsDRSRecoveryInstances", "api_error", err)
			return nil, err
		}

		for _, resoveryInstance := range output.Items {
			d.StreamListItem(ctx, resoveryInstance)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
