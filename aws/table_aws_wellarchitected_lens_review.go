package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedLensReview(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_lens_review",
		Description: "AWS Well-Architected Lens Review",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Required},
				{Name: "lens_alias", Require: plugin.Required},
				{Name: "milestone_number", Require: plugin.Optional, CacheMatch: "exact"},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getWellArchitectedLensReview,
			Tags:    map[string]string{"service": "wellarchitected", "action": "GetLensReview"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       listWellArchitectedLensReviews,
			Tags:          map[string]string{"service": "wellarchitected", "action": "ListLensReviews"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
				{Name: "milestone_number", Require: plugin.Optional, CacheMatch: "exact"},
			},
			// TODO: Uncomment and remove extra check in
			// listWellArchitectedLensReviews function once this works again
			//IgnoreConfig: &plugin.IgnoreConfig{
			//	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			//},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getWellArchitectedLensReview,
				Tags: map[string]string{"service": "wellarchitected", "action": "GetLensReview"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_WELLARCHITECTED_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "lens_name",
				Description: "The full name of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensReview.LensName"),
			},
			{
				Name:        "lens_arn",
				Description: "The ARN for the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensReview.LensArn"),
			},
			{
				Name:        "lens_alias",
				Description: "The alias of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensReview.LensAlias"),
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lens_status",
				Description: "The status of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensReview.LensStatus"),
			},
			{
				Name:        "lens_version",
				Description: "The version of the lens.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensReview.LensVersion"),
			},
			{
				Name:        "milestone_number",
				Description: "The milestone number. A workload can have a maximum of 100 milestones.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "notes",
				Description: "The notes associated with the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedLensReview,
				Transform:   transform.FromField("LensReview.Notes"),
			},
			{
				Name:        "updated_at",
				Description: "The date and time of the last update.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LensReview.UpdatedAt"),
			},
			{
				Name:        "pillar_review_summaries",
				Description: "A map from risk names to the count of how questions have that rating.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedLensReview,
				Transform:   transform.FromField("LensReview.PillarReviewSummaries"),
			},
			{
				Name:        "risk_counts",
				Description: "A map from risk names to the count of how questions have that rating.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LensReview.RiskCounts"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LensName"),
			},
		}),
	}
}

type LensReviewInfo struct {
	MilestoneNumber int32
	WorkloadId      *string
	*types.LensReview
}

//// LIST FUNCTION

func listWellArchitectedLensReviews(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workloadId := h.Item.(types.WorkloadSummary).WorkloadId

	// Reduce number of API call if the workload id has been provided in query parameter.
	if d.EqualsQualString("workload_id") != "" {
		if d.EqualsQualString("workload_id") != *workloadId {
			return nil, nil
		}
	}

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_review.listWellArchitectedLensReviews", "connection_error", err)
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

	input := &wellarchitected.ListLensReviewsInput{
		WorkloadId: workloadId,
		MaxResults: aws.Int32(maxLimit),
	}

	if d.EqualsQuals["milestone_number"] != nil {
		input.MilestoneNumber = aws.Int32(int32(d.EqualsQuals["milestone_number"].GetInt64Value()))
	}

	paginator := wellarchitected.NewListLensReviewsPaginator(svc, input, func(o *wellarchitected.ListLensReviewsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// TODO: Shouldn't be needed, but the List IgnoreConfig doesn't seem to
			// be working (maybe due to ParentHydrate use?)
			if strings.Contains(err.Error(), "ResourceNotFoundException") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_wellarchitected_lens_review.listWellArchitectedLensReviews", "api_error", err)
			return nil, err
		}

		for _, item := range output.LensReviewSummaries {
			reviewSummary := LensReviewInfo{
				WorkloadId:      workloadId,
				LensReview: &types.LensReview{
					LensAlias:   item.LensAlias,
					LensArn:     item.LensArn,
					LensName:    item.LensName,
					LensStatus:  item.LensStatus,
					LensVersion: item.LensVersion,
					RiskCounts:  item.RiskCounts,
					UpdatedAt:   item.UpdatedAt,
				},
			}
			if output.MilestoneNumber != nil {
				reviewSummary.MilestoneNumber = *output.MilestoneNumber
			}
			d.StreamListItem(ctx, reviewSummary)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWellArchitectedLensReview(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var workloadId, lensAlias string
	if h.Item != nil {
		workloadId = *h.Item.(LensReviewInfo).WorkloadId
		lensAlias = *h.Item.(LensReviewInfo).LensAlias
	} else {
		quals := d.EqualsQuals
		workloadId = quals["workload_id"].GetStringValue()
		lensAlias = quals["lens_alias"].GetStringValue()
	}

	// Empty Check
	if workloadId == "" || lensAlias == "" {
		return nil, nil
	}

	// Create Session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_review.getWellArchitectedLensReview", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &wellarchitected.GetLensReviewInput{
		WorkloadId: aws.String(workloadId),
		LensAlias:  aws.String(lensAlias),
	}

	if d.EqualsQuals["milestone_number"] != nil {
		params.MilestoneNumber = aws.Int32(int32(d.EqualsQuals["milestone_number"].GetInt64Value()))
	}

	op, err := svc.GetLensReview(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_review.getWellArchitectedLensReview", "api_error", err)
		return nil, err
	}

	review := LensReviewInfo{
		WorkloadId:      &workloadId,
		LensReview:      op.LensReview,
	}
	if op.MilestoneNumber != nil {
		review.MilestoneNumber = *op.MilestoneNumber
	}
	return review, nil
}
