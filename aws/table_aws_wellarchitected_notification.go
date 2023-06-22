package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedNotification(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_notification",
		Description: "AWS Well-Architected Notification",
		List: &plugin.ListConfig{
			Hydrate: listWellArchitectedNotifications,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "current_lens_version",
				Description: "The current version of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensUpgradeSummary.CurrentLensVersion"),
			},
			{
				Name:        "latest_lens_version",
				Description: "The latest version of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensUpgradeSummary.LatestLensVersion"),
			},
			{
				Name:        "lens_alias",
				Description: "The alias of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensUpgradeSummary.LensAlias"),
			},
			{
				Name:        "lens_arn",
				Description: "The ARN of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensUpgradeSummary.LensArn"),
			},
			{
				Name:        "type",
				Description: "The type of notification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensUpgradeSummary.WorkloadId"),
			},
			{
				Name:        "workload_name",
				Description: "The name of the workload.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensUpgradeSummary.WorkloadName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listWellArchitectedNotifications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_notification.listWellArchitectedNotifications", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &wellarchitected.ListNotificationsInput{
		MaxResults: maxLimit,
	}

	if d.EqualsQualString("workload_id") != "" {
		input.WorkloadId = aws.String(d.EqualsQualString("workload_id"))
	}

	paginator := wellarchitected.NewListNotificationsPaginator(svc, input, func(o *wellarchitected.ListNotificationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_notification.listWellArchitectedNotifications", "api_error", err)
			return nil, err
		}

		for _, item := range output.NotificationSummaries {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
