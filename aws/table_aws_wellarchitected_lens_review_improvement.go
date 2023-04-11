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

func tableAwsWellArchitectedLensReviewImprovement(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_lens_review_improvement",
		Description: "AWS Well-Architected Lens Review Improvement",
		List: &plugin.ListConfig{
			Hydrate: listWellArchitectedLensReviewImprovements,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Required},
				{Name: "lens_alias", Require: plugin.Required},
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

func listWellArchitectedLensReviewImprovements(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	workloadId := d.EqualsQualString("workload_id")
	lensAlias := d.EqualsQualString("lens_alias")

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &wellarchitected.ListLensReviewImprovementsInput{
		WorkloadId: aws.String(workloadId),
		LensAlias:  aws.String(lensAlias),
		MaxResults: maxLimit,
	}

	paginator := wellarchitected.NewListLensReviewImprovementsPaginator(svc, input, func(o *wellarchitected.ListLensReviewImprovementsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_lens_review_improvement.listWellArchitectedLensReviewImprovements", "api_error", err)
			return nil, err
		}

		for _, item := range output.ImprovementSummaries {
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
