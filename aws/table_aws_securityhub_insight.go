package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"

	securityhubEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHubInsight(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_insight",
		Description: "AWS Securityhub Insight",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getSecurityHubInsight,
			Tags:       map[string]string{"service": "securityhub", "action": "GetInsights"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidInputException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubInsights,
			Tags:    map[string]string{"service": "securityhub", "action": "GetInsights"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(securityhubEndpoint.AWS_SECURITYHUB_SERVICE_ID),
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

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_insight.listSecurityHubInsights", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
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

	input := &securityhub.GetInsightsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := securityhub.NewGetInsightsPaginator(svc, input, func(o *securityhub.GetInsightsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Handle error for accounts that are not subscribed to AWS Security Hub
			if strings.Contains(err.Error(), "not subscribed") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_securityhub_insight.listSecurityHubInsights", "api_error", err)
			return nil, err
		}

		for _, insight := range output.Insights {
			d.StreamListItem(ctx, insight)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHubInsight(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	arn := d.EqualsQuals["arn"].GetStringValue()

	// Entry check
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_insight.getSecurityHubInsight", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.GetInsightsInput{
		InsightArns: []string{arn},
	}

	// Get call
	data, err := svc.GetInsights(ctx, params)
	if err != nil {
		// Handle error for accounts that are not subscribed to AWS Security Hub
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_securityhub_insight.getSecurityHubInsight", "api_error", err)
		return nil, err
	}
	if len(data.Insights) > 0 {
		return data.Insights[0], nil
	}
	return nil, nil
}
