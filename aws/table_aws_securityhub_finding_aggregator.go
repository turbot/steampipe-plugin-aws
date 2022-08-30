package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
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
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listSecurityHubFindingAggregators")

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &securityhub.ListFindingAggregatorsInput{
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
	err = svc.ListFindingAggregatorsPages(
		input,
		func(page *securityhub.ListFindingAggregatorsOutput, isLast bool) bool {
			for _, aggregator := range page.FindingAggregators {
				d.StreamListItem(ctx, aggregator)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		// Handeled error for not subscribed region
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSecurityHubFindingAggregator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityHubFindingAggregator")

	var aggregatorArn string
	if h.Item != nil {
		aggregatorArn = *h.Item.(*securityhub.FindingAggregator).FindingAggregatorArn
	} else {
		aggregatorArn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	// Empty check
	if aggregatorArn == "" {
		return nil, nil
	}

	// get service
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &securityhub.GetFindingAggregatorInput{
		FindingAggregatorArn: &aggregatorArn,
	}

	// Get call
	op, err := svc.GetFindingAggregator(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getSecurityHubFindingAggregator", "ERROR", err)
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		return nil, err
	}
	return op, nil
}
