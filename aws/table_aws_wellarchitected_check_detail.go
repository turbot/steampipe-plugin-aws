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

func tableAwsWellArchitectedCheckDetail(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_check_detail",
		Description: "AWS Well-Architected Check Detail",
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedLenses,
			Hydrate:       listWellArchitectedCheckDetails,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Required},
				{Name: "choice_id", Require: plugin.Optional},
				{Name: "lens_arn", Require: plugin.Optional},
				{Name: "pillar_id", Require: plugin.Optional},
				{Name: "question_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// {
			// 	Name:        "data",
			// 	Description: "The ID of the question.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromValue(),
			// },
			{
				Name:        "id",
				Description: "Trusted Advisor check ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "choice_id",
				Description: "The ID of a choice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Trusted Advisor check description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "flagged_resources",
				Description: "Count of flagged resources associated to the check.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "lens_arn",
				Description: "Well-Architected Lens ARN associated to the check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Trusted Advisor check name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pillar_id",
				Description: "The ID used to identify a pillar, for example, security. A pillar is identified by its PillarReviewSummary$PillarId.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider",
				Description: "Provider of the check related to the best practice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "question_id",
				Description: "The ID of the question.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reason",
				Description: "Reason associated to the check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Status associated to the check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The date and time recorded.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "workload_id",
				Description: "The ID of the question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("workload_id"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

type CheckDetailInfo struct {
	types.CheckDetail
	LensAlias  *string
	WorkloadId *string
}

//// LIST FUNCTION

func listWellArchitectedCheckDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// if d.EqualsQualString("choice_id") == "" {
	// 	return nil, nil
	// }
	// choiceId := d.EqualsQualString("choice_id")

	answerList, err := listAnswerDetailsForWorkload(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_check_detail.listWellArchitectedCheckDetails", "client_error", err)
		return nil, err
	}

	if len(answerList) == 0 {
		return nil, nil
	}

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_check_detail.listWellArchitectedCheckDetails", "client_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	for _, answer := range answerList {
		// List call
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.EqualsQualString("pillar_id") != "" && d.EqualsQualString("pillar_id") != *answer.PillarId {
			return nil, nil
		}
		if d.EqualsQualString("question_id") != "" && d.EqualsQualString("question_id") != *answer.QuestionId {
			return nil, nil
		}
		_, err := fetchWellArchitectedCheckDetails(ctx, d, svc, answer)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_check_detail.listWellArchitectedCheckDetails", "api_error", err)
			return nil, err
		}
	}

	return nil, nil
}

func fetchWellArchitectedCheckDetails(ctx context.Context, d *plugin.QueryData, svc *wellarchitected.Client, answer AnswerInfo) (interface{}, error) {

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

	for _, choice := range answer.Choices {

		if d.EqualsQualString("choice_id") != "" && d.EqualsQualString("choice_id") != *choice.ChoiceId {
			return nil, nil
		}

		input := &wellarchitected.ListCheckDetailsInput{
			MaxResults: maxLimit,
			LensArn:    aws.String(*answer.LensArn),
			PillarId:   aws.String(*answer.PillarId),
			QuestionId: aws.String(*answer.QuestionId),
			WorkloadId: aws.String(*answer.WorkloadId),
			ChoiceId:   aws.String(*choice.ChoiceId),
		}

		paginator := wellarchitected.NewListCheckDetailsPaginator(svc, input, func(o *wellarchitected.ListCheckDetailsPaginatorOptions) {
			o.Limit = maxLimit
			o.StopOnDuplicateToken = true
		})

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "ResourceNotFoundException") || strings.Contains(err.Error(), "ValidationException") {
					return nil, nil
				}
			}

			for _, item := range output.CheckDetails {
				d.StreamListItem(ctx, item)

				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}
	return nil, nil
}

func listAnswerDetailsForWorkload(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) ([]AnswerInfo, error) {
	lens := h.Item.(types.LensSummary)

	// Validate - User inputs must not be blank
	if d.EqualsQualString("workload_id") == "" {
		return nil, nil
	}

	if d.EqualsQualString("lens_arn") != "" && d.EqualsQualString("lens_arn") != *lens.LensArn {
		return nil, nil
	}

	input := &wellarchitected.ListAnswersInput{
		MaxResults: int32(50),
		LensAlias:  aws.String(*lens.LensAlias),
		WorkloadId: aws.String(d.EqualsQualString("workload_id")),
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

	var answerList []AnswerInfo

	paginator := wellarchitected.NewListAnswersPaginator(svc, input, func(o *wellarchitected.ListAnswersPaginatorOptions) {
		o.Limit = int32(50)
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotFoundException") || strings.Contains(err.Error(), "ValidationException") {
				return nil, nil
			}
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

			if output.LensAlias == nil {
				output.LensAlias = output.LensArn
			}
			answerList = append(answerList, AnswerInfo{answer, output.LensAlias, output.LensArn, &output.MilestoneNumber, output.WorkloadId})
		}
	}
	return answerList, nil
}
