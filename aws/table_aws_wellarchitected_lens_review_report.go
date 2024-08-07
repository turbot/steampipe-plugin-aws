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

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedLensReviewReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_lens_review_report",
		Description: "AWS Well-Architected Lens Review Report",
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedWorkloads,
			Hydrate:       getWellArchitectedLensReviewReports,
			Tags:          map[string]string{"service": "wellarchitected", "action": "GetLensReviewReport"},
			// TODO: Uncomment and remove extra check in
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			// },
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_id", Require: plugin.Optional},
				{Name: "lens_alias", Require: plugin.Optional},
				{Name: "milestone_number", Require: plugin.Optional},
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
				Name:        "base64_string",
				Description: "The Base64-encoded string representation of a lens review report. This data can be used to create a PDF file.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

type ReviewReportInfo struct {
	Base64String    *string
	LensAlias       *string
	LensArn         *string
	MilestoneNumber int32
	WorkloadId      *string
}

//// LIST FUNCTION

func getWellArchitectedLensReviewReports(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workload := h.Item.(types.WorkloadSummary)

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_review_report.getWellArchitectedLensReviewReports", "connection_error", err)
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
		input := &wellarchitected.GetLensReviewReportInput{
			WorkloadId: workload.WorkloadId,
			LensAlias:  aws.String(lensAlias),
		}
		equalQuals := d.EqualsQuals
		if equalQuals != nil {
			if equalQuals["milestone_number"] != nil {
				milestoneNumber := int32(equalQuals["milestone_number"].GetInt64Value())
				if milestoneNumber < 1 || milestoneNumber > 100 {
					return nil, fmt.Errorf("MilestoneNumber must have minimum value of 1 and maximum value of 100")
				}
				input.MilestoneNumber = aws.Int32(milestoneNumber)
			}
		}

		op, err := svc.GetLensReviewReport(ctx, input)
		if err != nil {
			// If user provided milestone number does not exist then the API will throw ResourceNotFoundException error.
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if helpers.StringSliceContains([]string{"ResourceNotFoundException"}, ae.ErrorCode()) {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_wellarchitected_lens_review_report.getWellArchitectedLensReviewReports", "api_error", err)
			return nil, err
		}

		d.StreamListItem(ctx, ReviewReportInfo{
			Base64String:    op.LensReviewReport.Base64String,
			LensAlias:       op.LensReviewReport.LensAlias,
			LensArn:         op.LensReviewReport.LensArn,
			MilestoneNumber: *op.MilestoneNumber,
			WorkloadId:      op.WorkloadId,
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
