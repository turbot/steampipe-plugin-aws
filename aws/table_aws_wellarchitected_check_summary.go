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

func tableAwsWellArchitectedCheckSummary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_check_summary",
		Description: "AWS Well-Architected Summary",
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       listWellArchitectedCheckSummaries,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
				{Name: "choice_id", Require: plugin.Optional},
				{Name: "lens_arn", Require: plugin.Optional},
				{Name: "pillar_id", Require: plugin.Optional},
				{Name: "question_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Trusted Advisor check ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.Id"),
			},
			{
				Name:        "choice_id",
				Description: "The ID of a choice.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.ChoiceId"),
			},
			{
				Name:        "description",
				Description: "Trusted Advisor check description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.Description"),
			},
			{
				Name:        "lens_arn",
				Description: "Well-Architected Lens ARN associated to the check.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.LensArn"),
			},
			{
				Name:        "name",
				Description: "Trusted Advisor check name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.Name"),
			},
			{
				Name:        "pillar_id",
				Description: "The ID used to identify a pillar, for example, security. A pillar is identified by its PillarReviewSummary$PillarId.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.PillarId"),
			},
			{
				Name:        "provider",
				Description: "Provider of the check related to the best practice.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.Provider"),
			},
			{
				Name:        "question_id",
				Description: "The ID of the question.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.QuestionId"),
			},
			{
				Name:        "status",
				Description: "Status associated to the check.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CheckSummary.Status"),
			},
			{
				Name:        "updated_at",
				Description: "The date and time recorded.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CheckSummary.UpdatedAt"),
			},
			{
				Name:        "workload_id",
				Description: "The ID of the question.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_summary",
				Description: "Account summary associated to the check.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CheckSummary.AccountSummary"),
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

type CheckSummaryInfo struct {
	types.CheckSummary
	WorkloadId *string
}

//// LIST FUNCTION

func listWellArchitectedCheckSummaries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	answerList, err := getAnswerDetailsForWorkload(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_check_summary.listWellArchitectedCheckSummaries", "error", err)
		return nil, err
	}

	if len(answerList) == 0 {
		return nil, nil
	}

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_check_summary.listWellArchitectedCheckSummaries", "client_error", err)
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
			continue
		}
		if d.EqualsQualString("question_id") != "" && d.EqualsQualString("question_id") != *answer.QuestionId {
			continue
		}
		_, err := fetchWellArchitectedCheckSummaries(ctx, d, svc, answer)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_check_summary.listWellArchitectedCheckSummaries", "api_error", err)
			return nil, err
		}
	}

	return nil, nil
}

func fetchWellArchitectedCheckSummaries(ctx context.Context, d *plugin.QueryData, svc *wellarchitected.Client, answer AnswerInfo) (interface{}, error) {

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	for _, choice := range answer.Choices {
		if d.EqualsQualString("choice_id") != "" && d.EqualsQualString("choice_id") != *choice.ChoiceId {
			continue
		}

		input := &wellarchitected.ListCheckSummariesInput{
			MaxResults: maxLimit,
			LensArn:    aws.String(*answer.LensArn),
			PillarId:   aws.String(*answer.PillarId),
			QuestionId: aws.String(*answer.QuestionId),
			WorkloadId: aws.String(*answer.WorkloadId),
			ChoiceId:   aws.String(*choice.ChoiceId),
		}

		paginator := wellarchitected.NewListCheckSummariesPaginator(svc, input, func(o *wellarchitected.ListCheckSummariesPaginatorOptions) {
			o.Limit = maxLimit
			o.StopOnDuplicateToken = true
		})

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "ResourceNotFoundException") || strings.Contains(err.Error(), "ValidationException") {
					return nil, nil
				}
				return nil, err
			}

			for _, item := range output.CheckSummaries {
				d.StreamListItem(ctx, CheckSummaryInfo{item, answer.WorkloadId})

				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}
	return nil, nil
}
