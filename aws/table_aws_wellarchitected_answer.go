package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"
	"github.com/aws/smithy-go"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedAnswer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_answer",
		Description: "AWS Well-Architected Answer",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "lens_alias", Require: plugin.Required},
				{Name: "question_id", Require: plugin.Required},
				{Name: "workload_id", Require: plugin.Required},
				{Name: "milestone_number", Require: plugin.Optional},
			},
			Hydrate: getWellArchitectedAnswer,
			Tags:    map[string]string{"service": "wellarchitected", "action": "GetAnswer"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       listWellArchitectedAnswers,
			Tags:          map[string]string{"service": "wellarchitected", "action": "ListAnswers"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "lens_alias", Require: plugin.Optional},
				{Name: "pillar_id", Require: plugin.Optional},
				{Name: "workload_id", Require: plugin.Optional},
				{Name: "milestone_number", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getWellArchitectedAnswer,
				Tags: map[string]string{"service": "wellarchitected", "action": "GetAnswer"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "question_id",
				Description: "The ID of the question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Answer.QuestionId"),
			},
			{
				Name:        "lens_alias",
				Description: "The alias of the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "helpful_resource_display_text",
				Description: "The helpful resource text to be displayed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedAnswer,
				Transform:   transform.FromField("Answer.HelpfulResourceDisplayText"),
			},
			{
				Name:        "helpful_resource_url",
				Description: "The helpful resource URL for a question.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedAnswer,
				Transform:   transform.FromField("Answer.HelpfulResourceUrl"),
			},
			{
				Name:        "improvement_plan_url",
				Description: "The improvement plan URL for a question. This value is only available if the question has been answered.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedAnswer,
				Transform:   transform.FromField("Answer.ImprovementPlanUrl"),
			},
			{
				Name:        "is_applicable",
				Description: "Defines whether this question is applicable to a lens review.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Answer.IsApplicable"),
			},
			{
				Name:        "lens_arn",
				Description: "The Amazon Resource Name (ARN) of the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "milestone_number",
				Description: "The milestone number.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "notes",
				Description: "The notes associated with the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedAnswer,
				Transform:   transform.FromField("Answer.Notes"),
			},
			{
				Name:        "pillar_id",
				Description: "The ID used to identify a pillar, for example, security. A pillar is identified by its PillarReviewSummary$PillarId.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Answer.PillarId"),
			},
			{
				Name:        "question_description",
				Description: "The description of the question.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedAnswer,
				Transform:   transform.FromField("Answer.QuestionDescription"),
			},
			{
				Name:        "question_title",
				Description: "The title of the question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Answer.QuestionTitle"),
			},
			{
				Name:        "reason",
				Description: "The reason why the question is not applicable to your workload.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Answer.Reason"),
			},
			{
				Name:        "risk",
				Description: "The risk for a given workload, lens review, pillar, or question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Answer.Risk"),
			},
			{
				Name:        "choice_answers",
				Description: "A list of selected choices to a question in your workload.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedAnswer,
				Transform:   transform.FromField("Answer.ChoiceAnswers"),
			},
			{
				Name:        "choices",
				Description: "List of choices available for a question.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Answer.Choices"),
			},
			{
				Name:        "selected_choices",
				Description: "List of selected choice IDs in a question answer. The values entered replace the previously selected choices.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Answer.SelectedChoices"),
			},

			// Seampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Answer.QuestionTitle"),
			},
		}),
	}
}

type AnswerInfo struct {
	types.Answer
	LensAlias       *string
	LensArn         *string
	MilestoneNumber *int32
	WorkloadId      *string
}

//// LIST FUNCTION

func listWellArchitectedAnswers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workload := h.Item.(types.WorkloadSummary)

	// Validate - User inputs must not be blank and return nil if doesn't match the hydrated workload ID
	if d.EqualsQualString("workload_id") != "" && d.EqualsQualString("workload_id") != *workload.WorkloadId {
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

	input := &wellarchitected.ListAnswersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_answer.listWellArchitectedAnswers", "client_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	for _, lensAlias := range workload.Lenses {
		if d.EqualsQualString("lens_alias") != "" && d.EqualsQualString("lens_alias") != lensAlias {
			continue
		}
		if d.EqualsQualString("pillar_id") != "" {
			input.PillarId = aws.String(d.EqualsQualString("pillar_id"))
		}
		if d.EqualsQuals["milestone_number"] != nil {
			input.MilestoneNumber = aws.Int32(int32(d.EqualsQuals["milestone_number"].GetInt64Value()))
		}
		input.LensAlias = aws.String(lensAlias)
		input.WorkloadId = aws.String(*workload.WorkloadId)

		paginator := wellarchitected.NewListAnswersPaginator(svc, input, func(o *wellarchitected.ListAnswersPaginatorOptions) {
			o.Limit = maxLimit
			o.StopOnDuplicateToken = true
		})

		// List call
		for paginator.HasMorePages() {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			output, err := paginator.NextPage(ctx)
			if err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) {
					if helpers.StringSliceContains([]string{"ResourceNotFoundException"}, ae.ErrorCode()) {
						return nil, nil
					}
				}

				plugin.Logger(ctx).Error("aws_wellarchitected_answer.listWellArchitectedAnswers", "api_error", err)
				return nil, err
			}

			for _, item := range output.AnswerSummaries {

				answer := types.Answer{
					Choices:       item.Choices,
					IsApplicable:  item.IsApplicable,
					PillarId:      item.PillarId,
					QuestionId:    item.QuestionId,
					QuestionTitle: item.QuestionTitle,
					Reason:        item.Reason,
					Risk:          item.Risk,
				}

				d.StreamListItem(ctx, &AnswerInfo{answer, output.LensAlias, output.LensArn, output.MilestoneNumber, output.WorkloadId})

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWellArchitectedAnswer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var questionId, lensAlias, workloadId string
	var milestoneNumber int32
	if h.Item != nil {
		answer := h.Item.(*AnswerInfo)
		questionId = *answer.QuestionId
		lensAlias = *answer.LensAlias
		workloadId = *answer.WorkloadId
		milestoneNumber = *answer.MilestoneNumber
	} else {
		questionId = d.EqualsQualString("question_id")
		lensAlias = d.EqualsQualString("lens_alias")
		workloadId = d.EqualsQualString("workload_id")
		if d.EqualsQuals["milestone_number"] != nil {
			milestoneNumber = int32(d.EqualsQuals["milestone_number"].GetInt64Value())
		}
	}

	// Validate - User inputs must not be blank
	if questionId == "" || lensAlias == "" || workloadId == "" {
		return nil, nil
	}

	params := &wellarchitected.GetAnswerInput{
		QuestionId: aws.String(questionId),
		LensAlias:  aws.String(lensAlias),
		WorkloadId: aws.String(workloadId),
	}
	if milestoneNumber != 0 {
		*params.MilestoneNumber = milestoneNumber
	}

	// Create Session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_answer.getWellArchitectedAnswer", "client_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	op, err := svc.GetAnswer(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_answer.getWellArchitectedAnswer", "api_error", err)
		return nil, err
	}

	return &AnswerInfo{*op.Answer, op.LensAlias, op.LensArn, op.MilestoneNumber, op.WorkloadId}, nil
}
