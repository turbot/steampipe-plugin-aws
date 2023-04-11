package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedMilestone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_milestone",
		Description: "AWS Well-Architected Milestone",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "milestone_number", Require: plugin.Required},
				{Name: "workload_id", Require: plugin.Required},
			},
			Hydrate: getWellArchitectedMilestone,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       listWellArchitectedMilestones,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "milestone_name",
				Description: "The name of the milestone in a workload. Milestone names must be unique within a workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "milestone_number",
				Description: "The milestone number. A workload can have a maximum of 100 milestones.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "recorded_at",
				Description: "The date and time recorded.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadSummary.WorkloadId", "Workload.WorkloadId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listWellArchitectedMilestones(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workloadId := h.Item.(types.WorkloadSummary).WorkloadId

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

	input := &wellarchitected.ListMilestonesInput{
		MaxResults: maxLimit,
	}

	// Validate - User input must match the hydrated WorkloadId
	if d.EqualsQualString("workload_id") != "" && *workloadId != d.EqualsQualString("workload_id") {
		return nil, nil
	}
	input.WorkloadId = aws.String(*workloadId)

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_milestone.listWellArchitectedMilestones", "client_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	paginator := wellarchitected.NewListMilestonesPaginator(svc, input, func(o *wellarchitected.ListMilestonesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_milestone.listWellArchitectedMilestones", "api_error", err)
			return nil, err
		}

		for _, item := range output.MilestoneSummaries {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWellArchitectedMilestone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("workload_id")
	number := int32(d.EqualsQuals["milestone_number"].GetInt64Value())

	// Create Session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_milestone.getWellArchitectedMilestone", "client_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	params := &wellarchitected.GetMilestoneInput{
		WorkloadId:      aws.String(id),
		MilestoneNumber: *aws.Int32(number),
	}

	op, err := svc.GetMilestone(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_milestone.getWellArchitectedMilestone", "api_error", err)
		return nil, err
	}
	return op.Milestone, nil
}
