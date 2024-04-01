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

//// TABLE DEFINITION

func tableAwsTrustedAdvisorCheckSummary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_trusted_advisor_check_summary",
		Description: "AWS Trusted Advisor Check Summary",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("check_id"),
			Hydrate:    getTrustedAdvisorCheckSummary,
			Tags:       map[string]string{"service": "support", "action": "DescribeTrustedAdvisorCheckSummaries"},
		},
		List: &plugin.ListConfig{
			Hydrate: listTrustedAdvisorCheckSummaries,
			Tags:    map[string]string{"service": "support", "action": "DescribeTrustedAdvisorChecks"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "language",
					Require: plugin.Required,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getTrustedAdvisorCheckSummary,
				Tags: map[string]string{"service": "support", "action": "DescribeTrustedAdvisorCheckSummaries"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The display name for the Trusted Advisor check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "check_id",
				Description: "The unique identifier for the Trusted Advisor check.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "CheckId"),
			},
			{
				Name:        "category",
				Description: "The category of the Trusted Advisor check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "language",
				Description: "The ISO 639-1 code for the language that you want your checks to appear in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("language"),
			},
			{
				Name:        "description",
				Description: "The description of the Trusted Advisor check, which includes the alert criteria and recommended operations (contains HTML markup).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The alert status of the check: 'ok' (green), 'warning' (yellow), 'error' (red), or 'not_available'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTrustedAdvisorCheckSummary,
			},
			{
				Name:        "timestamp",
				Description: "The time of the last refresh of the check.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getTrustedAdvisorCheckSummary,
			},
			{
				Name:        "resources_flagged",
				Description: "The number of Amazon Web Services resources that were flagged (listed) by the Trusted Advisor check.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTrustedAdvisorCheckSummary,
				Transform:   transform.FromField("ResourcesSummary.ResourcesFlagged"),
			},
			{
				Name:        "resources_ignored",
				Description: "The number of Amazon Web Services resources ignored by Trusted Advisor because information was unavailable.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTrustedAdvisorCheckSummary,
				Transform:   transform.FromField("ResourcesSummary.ResourcesIgnored"),
			},
			{
				Name:        "resources_processed",
				Description: "The number of Amazon Web Services resources that were analyzed by the Trusted Advisor check.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTrustedAdvisorCheckSummary,
				Transform:   transform.FromField("ResourcesSummary.ResourcesProcessed"),
			},
			{
				Name:        "resources_suppressed",
				Description: "The number of Amazon Web Services resources ignored by Trusted Advisor because they were marked as suppressed by the user.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getTrustedAdvisorCheckSummary,
				Transform:   transform.FromField("ResourcesSummary.ResourcesSuppressed"),
			},
			{
				Name:        "category_specific_summary",
				Description: "Summary information that relates to the category of the check. Cost Optimizing is the only category that is currently supported.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTrustedAdvisorCheckSummary,
			},
			{
				Name:        "metadata",
				Description: "The column headings for the data returned by the Trusted Advisor check. The order of the headings corresponds to the order of the data in the Metadata element of the TrustedAdvisorResourceDetail for the check. Metadata contains all the data that is shown in the Excel download, even in those cases where the UI shows just summary data.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listTrustedAdvisorCheckSummaries(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SupportClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_trusted_advisor_check_summary.listTrustedAdvisorCheckSummaries", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	language := d.EqualsQualString("language")

	// Empty check
	if language == "" {
		return nil, nil
	}

	input := &support.DescribeTrustedAdvisorChecksInput{
		Language: aws.String(language),
	}

	// List call
	result, err := svc.DescribeTrustedAdvisorChecks(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_trusted_advisor_check_summary.listTrustedAdvisorCheckSummaries", "api_error", err)
		return nil, err
	}

	for _, check := range result.Checks {
		d.StreamListItem(ctx, check)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getTrustedAdvisorCheckSummary(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	checkId := ""

	if h.Item != nil {
		checkSummary := h.Item.(types.TrustedAdvisorCheckDescription)
		checkId = *checkSummary.Id
	} else {
		checkId = d.EqualsQualString("check_id")
	}

	// Empty check
	if checkId == "" {
		return nil, nil
	}

	// Create session
	svc, err := SupportClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_trusted_advisor_check_summary.getTrustedAdvisorCheckSummary", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	input := &support.DescribeTrustedAdvisorCheckSummariesInput{
		CheckIds: []*string{aws.String(checkId)},
	}

	// Get call
	result, err := svc.DescribeTrustedAdvisorCheckSummaries(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_trusted_advisor_check_summary.getTrustedAdvisorCheckSummary", "api_error", err)
		return nil, err
	}

	if result != nil && len(result.Summaries) > 0 {
		return result.Summaries[0], nil
	}
	return nil, nil
}
