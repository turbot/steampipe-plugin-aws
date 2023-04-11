package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedLensReviewReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_lens_review_report",
		Description: "AWS Well-Architected Lens Review Report",
		List: &plugin.ListConfig{
			Hydrate: getWellArchitectedLensReviewReports,
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
				Name:        "base64_string",
				Description: "The Base64-encoded string representation of a lens review report. This data can be used to create a PDF file.",
				Type:        proto.ColumnType_STRING,
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

type ReviewReportInfo struct {
	Base64String    *string
	LensAlias       *string
	LensArn         *string
	MilestoneNumber int32
	WorkloadId      *string
}

//// LIST FUNCTION

func getWellArchitectedLensReviewReports(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	workloadId := d.EqualsQualString("workload_id")
	lensAlias := d.EqualsQualString("lens_alias")

	input := &wellarchitected.GetLensReviewReportInput{
		WorkloadId: aws.String(workloadId),
		LensAlias:  aws.String(lensAlias),
	}

	op, err := svc.GetLensReviewReport(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_review_report.getWellArchitectedLensReviewReports", "api_error", err)
		return nil, err
	}

	if op != nil {
		return ReviewReportInfo{
			Base64String:    op.LensReviewReport.Base64String,
			LensAlias:       op.LensReviewReport.LensAlias,
			LensArn:         op.LensReviewReport.LensArn,
			MilestoneNumber: op.MilestoneNumber,
			WorkloadId:      op.WorkloadId,
		}, nil
	}

	return nil, nil
}
