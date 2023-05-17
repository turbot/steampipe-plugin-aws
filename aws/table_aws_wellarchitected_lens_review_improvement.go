package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"
	"github.com/aws/smithy-go"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedLensReviewImprovement(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_lens_review_improvement",
		Description: "AWS Well-Architected Lens Review Improvement",
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       listWellArchitectedLensReviewImprovements,
			// TODO: Uncomment and remove extra check in
			// listWellArchitectedLensReviewImprovements function once this works again
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			// },
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
				{Name: "lens_alias", Require: plugin.Optional},
				{Name: "milestone_number", Require: plugin.Optional},
				{Name: "pillar_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "lens_alias",
				Description: "The alias of the lens. For Amazon Web Services official lenses, this is either the lens alias, such as serverless, or the lens ARN, such as arn:aws:wellarchitected:us-east-1:123456789012:lens/my-lens. Each lens is identified by its LensSummary$LensAlias.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lens_arn",
				Description: "The ARN for the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "milestone_number",
				Description: "The milestone number. A workload can have a maximum of 100 milestones.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "improvement_plan_url",
				Description: "The improvement plan URL for a question. This value is only available if the question has been answered.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImprovementSummary.ImprovementPlanUrl"),
			},
			{
				Name:        "pillar_id",
				Description: "The ID used to identify a pillar, for example, security. A pillar is identified by its PillarReviewSummary$PillarId.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImprovementSummary.PillarId"),
			},
			{
				Name:        "question_id",
				Description: "The ID of the question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImprovementSummary.QuestionId"),
			},
			{
				Name:        "question_title",
				Description: "The title of the question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImprovementSummary.QuestionTitle"),
			},
			{
				Name:        "risk",
				Description: "The risk for a given workload, lens review, pillar, or question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImprovementSummary.Risk"),
			},
			{
				Name:        "improvement_plans",
				Description: "The improvement plan details.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ImprovementSummary.ImprovementPlans"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImprovementSummary.QuestionTitle"),
			},
		}),
	}
}

type ReviewImprovementInfo struct {
	LensAlias       *string
	LensArn         *string
	MilestoneNumber int32
	WorkloadId      *string
	types.ImprovementSummary
}

//// LIST FUNCTION

func listWellArchitectedLensReviewImprovements(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workload := h.Item.(types.WorkloadSummary)

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_review_improvement.listWellArchitectedLensReviewImprovements", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	if d.EqualsQualString("workload_id") != "" {
		if d.EqualsQualString("workload_id") != *workload.WorkloadId {
			return nil, nil
		}
	}

	for _, lensAlias := range workload.Lenses {
		// Check for reduce the numbers of API call if lens_alias or lens_arn is provided in query parameter
		if d.EqualsQualString("lens_alias") != "" {
			if d.EqualsQualString("lens_alias") != lensAlias {
				continue
			}
		}

		// Limiting the results
		maxLimit := int32(100)
		if d.QueryContext.Limit != nil {
			limit := int32(*d.QueryContext.Limit)
			if limit < maxLimit {
				maxLimit = limit
			}
		}

		input := &wellarchitected.ListLensReviewImprovementsInput{
			WorkloadId: workload.WorkloadId,
			LensAlias:  aws.String(lensAlias),
			MaxResults: maxLimit,
		}

		if d.EqualsQuals["pillar_id"] != nil {
			input.PillarId = aws.String(d.EqualsQuals["pillar_id"].GetStringValue())
		}
		if d.EqualsQuals["milestone_number"] != nil {
			milestoneNumber := int32(d.EqualsQuals["milestone_number"].GetInt64Value())
			if milestoneNumber < 1 || milestoneNumber > 100 {
				return nil, fmt.Errorf("MilestoneNumber must have minimum value of 1 and maximum value of 100")
			}
			input.MilestoneNumber = int32(d.EqualsQuals["milestone_number"].GetInt64Value())
		}

		_, err := stereamlistLensReviewImprovements(ctx, d, h, svc, input, maxLimit)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func stereamlistLensReviewImprovements(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, svc *wellarchitected.Client, input *wellarchitected.ListLensReviewImprovementsInput, maxLimit int32) (interface{}, error) {
	paginator := wellarchitected.NewListLensReviewImprovementsPaginator(svc, input, func(o *wellarchitected.ListLensReviewImprovementsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Adding the igone confog in the list config does not seems to work, so we have handles it here.
			var ae smithy.APIError
			if errors.As(err, &ae) {
				// In order to handle the validation exception, we should account for potential errors thrown by the API when querying this table with the pillar_id provided in the WHERE clause. If the specified pillar_id is not available within a workload, the API may generate an error that needs to be handled appropriately.
				// Error: operation error WellArchitected: ListLensReviewImprovements, https response error StatusCode: 400, RequestID: 8af3784e-b94e-4403-8821-a7700f23b341, ValidationException: [Validation] No pillar with ID operationalExcellence was found in workload 4fca39b680a31bb118be6bc0d177849d.
				if ae.ErrorCode() == "ResourceNotFoundException" || ae.ErrorCode() == "ValidationException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_wellarchitected_lens_review_improvement.stereamlistLensReviewImprovements", "api_error", err)
			return nil, err
		}

		for _, item := range output.ImprovementSummaries {
			// For Custom Lenses the lens alias is same as lens arn.
			// https://docs.aws.amazon.com/wellarchitected/latest/APIReference/API_LensSummary.html#wellarchitected-Type-LensSummary-LensAlias
			if output.LensAlias == nil {
				output.LensAlias = output.LensArn
			}

			d.StreamListItem(ctx, ReviewImprovementInfo{
				LensAlias:          output.LensAlias,
				LensArn:            output.LensArn,
				MilestoneNumber:    output.MilestoneNumber,
				WorkloadId:         output.WorkloadId,
				ImprovementSummary: item,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}
