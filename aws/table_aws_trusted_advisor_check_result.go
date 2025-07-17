package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/aws/aws-sdk-go-v2/service/support/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// TrustedAdvisorCheckResultData represents a flattened structure combining check result data with individual flagged resource details
type TrustedAdvisorCheckResultData struct {
	CheckId                 *string                                      `json:"checkId"`
	Timestamp               *string                                      `json:"timestamp"`
	Status                  *string                                      `json:"status"`
	ResourcesSummary        *types.TrustedAdvisorResourcesSummary        `json:"resourcesSummary"`
	CategorySpecificSummary *types.TrustedAdvisorCategorySpecificSummary `json:"categorySpecificSummary"`

	// Flagged resource properties (flattened from flaggedResources array)
	FlaggedResourceId           *string   `json:"flaggedResourceId"`
	FlaggedResourceStatus       *string   `json:"flaggedResourceStatus"`
	FlaggedResourceRegion       *string   `json:"flaggedResourceRegion"`
	FlaggedResourceIsSuppressed bool      `json:"flaggedResourceIsSuppressed"`
	FlaggedResourceMetadata     []*string `json:"flaggedResourceMetadata"`
}

//// TABLE DEFINITION

func tableAwsTrustedAdvisorCheckResult(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_trusted_advisor_check_result",
		Description: "AWS Trusted Advisor Check Result",
		List: &plugin.ListConfig{
			Hydrate: listTrustedAdvisorCheckResults,
			Tags:    map[string]string{"service": "support", "action": "DescribeTrustedAdvisorCheckResult"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "language",
					Require: plugin.Required,
				},
				{
					Name:    "check_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "check_id",
				Description: "The unique identifier for the Trusted Advisor check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "language",
				Description: "The ISO 639-1 code for the language that you want your checks to appear in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("language"),
			},
			{
				Name:        "status",
				Description: "The overall status of the Trusted Advisor check: 'ok' (green), 'warning' (yellow), 'error' (red), or 'not_available'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timestamp",
				Description: "The time when the Trusted Advisor check was last refreshed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "flagged_resource_id",
				Description: "The unique identifier for the flagged resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "flagged_resource_status",
				Description: "The status of the flagged resource: 'ok' (green), 'warning' (yellow), 'error' (red), or 'not_available'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "flagged_resource_region",
				Description: "The AWS region of the flagged resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "flagged_resource_is_suppressed",
				Description: "Specifies whether the flagged AWS resource was ignored by Trusted Advisor because it was marked as suppressed by the user.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "flagged_resource_metadata",
				Description: "Additional information about the flagged resource. The exact metadata and its order can be obtained by calling DescribeTrustedAdvisorChecks.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resources_summary",
				Description: "Summary information about the resources analyzed by the Trusted Advisor check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "category_specific_summary",
				Description: "Summary information that relates to the category of the check. Cost Optimizing is the only category that is currently supported.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FlaggedResourceId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listTrustedAdvisorCheckResults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SupportClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_trusted_advisor_check_result.listTrustedAdvisorCheckResults", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	language := d.EqualsQualString("language")
	checkId := d.EqualsQualString("check_id")

	// Empty check
	if language == "" || checkId == "" {
		return nil, nil
	}

	input := &support.DescribeTrustedAdvisorCheckResultInput{
		Language: aws.String(language),
		CheckId:  aws.String(checkId),
	}

	// Get call
	result, err := svc.DescribeTrustedAdvisorCheckResult(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_trusted_advisor_check_result.listTrustedAdvisorCheckResults", "api_error", err)
		return nil, err
	}

	if result != nil && result.Result != nil {
		checkResult := result.Result

		// Stream each flagged resource as a separate row
		for _, resource := range checkResult.FlaggedResources {
			resourceData := &TrustedAdvisorCheckResultData{
				CheckId:                     checkResult.CheckId,
				Timestamp:                   checkResult.Timestamp,
				Status:                      checkResult.Status,
				ResourcesSummary:            checkResult.ResourcesSummary,
				CategorySpecificSummary:     checkResult.CategorySpecificSummary,
				FlaggedResourceId:           resource.ResourceId,
				FlaggedResourceStatus:       resource.Status,
				FlaggedResourceRegion:       resource.Region,
				FlaggedResourceIsSuppressed: resource.IsSuppressed,
				FlaggedResourceMetadata:     resource.Metadata,
			}

			d.StreamListItem(ctx, resourceData)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}