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

func tableAwsWellArchitectedLensReview(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_lens_review",
		Description: "AWS Well-Architected Lens Review",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"workload_id", "lens_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getWellArchitectedLensReview,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       listWellArchitectedLensReviews,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
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
				Name:        "notes",
				Description: "The notes associated with the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedLensReview,
				Transform:   transform.FromField("LensReview.Notes"),
			},
			{
				Name:        "updated_at",
				Description: "The date and time recorded.",
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
	WorkloadId *string
	*types.LensReview
}

//// LIST FUNCTION

func listWellArchitectedLensReviews(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workloadId := h.Item.(types.WorkloadSummary).WorkloadId

	// Reduce number of API call if the workload id has been provided in query parameter.
	equalQuals := d.EqualsQuals
	if equalQuals["workload_id"] != nil {
		if equalQuals["workload_id"].GetStringValue() != *workloadId {
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
		MaxResults: maxLimit,
	}

	paginator := wellarchitected.NewListLensReviewsPaginator(svc, input, func(o *wellarchitected.ListLensReviewsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_lens_review.listWellArchitectedLensReviews", "api_error", err)
			return nil, err
		}

		for _, item := range output.LensReviewSummaries {
			d.StreamListItem(ctx, LensReviewInfo{
				WorkloadId: workloadId,
				LensReview: &types.LensReview{
					LensAlias:   item.LensAlias,
					LensArn:     item.LensArn,
					LensName:    item.LensName,
					LensStatus:  item.LensStatus,
					LensVersion: item.LensVersion,
					RiskCounts:  item.RiskCounts,
					UpdatedAt:   item.UpdatedAt,
				},
			})

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
	var workloadId, lensArn string
	if h.Item != nil {
		workloadId = *h.Item.(LensReviewInfo).WorkloadId
		lensArn = *h.Item.(LensReviewInfo).LensArn
	} else {
		quals := d.EqualsQuals
		workloadId = quals["workload_id"].GetStringValue()
		lensArn = quals["lens_arn"].GetStringValue()
	}

	// Empty Check
	if workloadId == "" || lensArn == "" {
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
		LensAlias:  aws.String(lensArn),
	}

	op, err := svc.GetLensReview(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_review.getWellArchitectedLensReview", "api_error", err)
		return nil, err
	}

	return LensReviewInfo{WorkloadId: &workloadId, LensReview: op.LensReview}, nil
}
