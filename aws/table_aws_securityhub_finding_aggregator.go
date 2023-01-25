package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"

	securityhubv1 "github.com/aws/aws-sdk-go/service/securityhub"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHubFindingAggregator(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_finding_aggregator",
		Description: "AWS Security Hub Finding Aggregator",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getSecurityHubFindingAggregator,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubFindingAggregators,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAccessException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(securityhubv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the finding aggregator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FindingAggregatorArn"),
			},
			{
				Name:        "finding_aggregation_region",
				Description: "The aggregation Region.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityHubFindingAggregator,
			},
			{
				Name:        "region_linking_mode",
				Description: "Indicates whether to link all Regions, all Regions except for a list of excluded Regions, or a list of included Regions.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityHubFindingAggregator,
			},
			{
				Name:        "regions",
				Description: "The list of excluded Regions or included Regions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityHubFindingAggregator,
			},

			/// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FindingAggregatorArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubFindingAggregators(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_finding_aggregator.listSecurityHubFindingAggregators", "client_error", err)
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

	input := &securityhub.ListFindingAggregatorsInput{
		MaxResults: maxLimit,
	}
	// List Call
	paginator := securityhub.NewListFindingAggregatorsPaginator(svc, input, func(o *securityhub.ListFindingAggregatorsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Handle error for accounts that are not subscribed to AWS Security Hub
			if strings.Contains(err.Error(), "not subscribed") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_securityhub_finding_aggregator.listSecurityHubFindingAggregators", "api_error", err)
			return nil, err
		}

		for _, aggregator := range output.FindingAggregators {
			d.StreamListItem(ctx, aggregator)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHubFindingAggregator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var aggregatorArn string
	if h.Item != nil {
		aggregatorArn = *h.Item.(types.FindingAggregator).FindingAggregatorArn
	} else {
		aggregatorArn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	// Empty check
	if aggregatorArn == "" {
		return nil, nil
	}

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_finding_aggregator.getSecurityHubFindingAggregator", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.GetFindingAggregatorInput{
		FindingAggregatorArn: &aggregatorArn,
	}

	// Get call
	op, err := svc.GetFindingAggregator(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_securityhub_finding_aggregator.getSecurityHubFindingAggregator", "api_error", err)
		return nil, err
	}
	return op, nil
}
