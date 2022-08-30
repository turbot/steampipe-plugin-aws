package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsSecurityHubInsight(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_insight",
		Description: "AWS Securityhub Insight",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getSecurityHubInsight,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "InvalidInputException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubInsights,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of a Security Hub insight.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of a Security Hub insight.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InsightArn"),
			},
			{
				Name:        "group_by_attribute",
				Description: "The grouping attribute for the insight's findings. Indicates how to group the matching findings,and identifies the type of item that the insight applies to. For example, if an insight is grouped by resource identifier, then the insight produces a list of resource identifiers.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filters",
				Description: "One or more attributes used to filter the findings included in the insight. The insight only includes findings that match the criteria defined in the filters.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InsightArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubInsights(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSecurityHubInsights")

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &securityhub.GetInsightsInput{
		MaxResults: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.GetInsightsPages(
		input,
		func(page *securityhub.GetInsightsOutput, isLast bool) bool {
			for _, insight := range page.Insights {
				d.StreamListItem(ctx, insight)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		// Handle error for accounts that are not subscribed to AWS Security Hub
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listSecurityHubInsights", "list", err)
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getSecurityHubInsight(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSecurityHubInsight")

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// Entry check
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &securityhub.GetInsightsInput{
		InsightArns: []*string{aws.String(arn)},
	}

	// Get call
	data, err := svc.GetInsights(params)
	if err != nil {
		// Handle error for accounts that are not subscribed to AWS Security Hub
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("getSecurityHubInsight", "get", err)
	}
	if len(data.Insights) > 0 {
		return data.Insights[0], nil
	}
	return nil, nil
}
