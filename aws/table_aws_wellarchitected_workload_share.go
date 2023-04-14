package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedWorkloadShare(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_workload_share",
		Description: "AWS Well-Architected Workload Share",
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       listWellArchitectedWorkloadShares,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
				{Name: "shared_with", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadId"),
			},
			{
				Name:        "permission_type",
				Description: "Permission granted on a workload share.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadShare.PermissionType"),
			},
			{
				Name:        "share_id",
				Description: "The ID associated with the workload share.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadShare.ShareId"),
			},
			{
				Name:        "shared_with",
				Description: "The Amazon Web Services account ID, IAM role, organization ID, or organizational unit (OU) ID with which the workload is shared.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadShare.SharedWith"),
			},
			{
				Name:        "status",
				Description: "The status of a workload share.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadShare.Status"),
			},
			{
				Name:        "status_message",
				Description: "Optional message to compliment the Status field.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadShare.StatusMessage"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadShare.ShareId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorkloadSharesAkas,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type WorkloadShareInfo struct {
	WorkloadId    *string
	WorkloadShare types.WorkloadShareSummary
}

//// LIST FUNCTION

func listWellArchitectedWorkloadShares(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workload := h.Item.(types.WorkloadSummary)

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &wellarchitected.ListWorkloadSharesInput{
		MaxResults: maxLimit,
	}

	// Validate - User input must match the parent hydrate WorkloadId
	if d.EqualsQuals["workload_id"] != nil && d.EqualsQualString("workload_id") != *workload.WorkloadId {
		return nil, nil
	}
	input.WorkloadId = aws.String(*workload.WorkloadId)
	if d.EqualsQualString("status") != "" {
		input.Status = types.ShareStatus(d.EqualsQuals["status"].GetStringValue())
	}
	if d.EqualsQualString("shared_with") != "" {
		input.SharedWithPrefix = aws.String(d.EqualsQualString("shared_with"))
	}

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_workload_share.listWellArchitectedWorkloadShares", "client_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	paginator := wellarchitected.NewListWorkloadSharesPaginator(svc, input, func(o *wellarchitected.ListWorkloadSharesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotFoundException") || strings.Contains(err.Error(), "ValidationException") {
				plugin.Logger(ctx).Debug("aws_wellarchitected_workload_share.listWellArchitectedWorkloadShares", "checked_error", err)
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_wellarchitected_workload_share.listWellArchitectedWorkloadShares", "api_error", err)
			return nil, err
		}

		for _, item := range output.WorkloadShareSummaries {
			d.StreamListItem(ctx, WorkloadShareInfo{output.WorkloadId, item})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getWellArchitectedWorkloadSharesAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	workloadShare := h.Item.(WorkloadShareInfo)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":wellarchitected:" + region + ":" + commonColumnData.AccountId + ":workload/" + *workloadShare.WorkloadId + "/share/" + *workloadShare.WorkloadShare.ShareId}

	return akas, nil
}
