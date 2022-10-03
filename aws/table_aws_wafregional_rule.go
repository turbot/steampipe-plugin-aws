package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/waf"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsWAFRegionalRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafregional_rule",
		Description: "AWS WAF Regional Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("rule_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"WAFNonexistentItemException"}),
			},
			Hydrate: getAwsWAFRegionalRule,
		},
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"WAFNonexistentItemException"}),
			},
			Hydrate: listAwsWAFRegionalRules,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name or description for the Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "Amazon Resource Name (ARN) of the Rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWAFRegionalRuleAkas,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "rule_id",
				Description: "A unique identifier for a Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_name",
				Description: "A friendly name or description for the metrics for this Rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWAFRegionalRule,
			},
			{
				Name:        "predicates",
				Description: "The Predicates object contains one Predicate element for each ByteMatchSet,IPSet, or SqlInjectionMatchSet object that you want to include in a Rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRegionalRule,
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
				Hydrate:     getAwsWAFRegionalRuleAkas,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWAFRegionalRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsWAFRegionalRules")
	// Create session
	svc, err := WAFRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	// List call
	params := &waf.ListRulesInput{
		Limit: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	// Minimunm limit is 0
	// https://docs.aws.amazon.com/waf/latest/APIReference/API_waf_ListRules.html
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.Limit {
			params.Limit = limit
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListRules(params)
		if err != nil {
			return nil, err
		}
		for _, rule := range response.Rules {
			d.StreamListItem(ctx, rule)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.NextMarker != nil {
			pagesLeft = true
			params.NextMarker = response.NextMarker
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsWAFRegionalRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsWAFRegionalRule")

	// Create Session
	svc, err := WAFRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	var id string
	if h.Item != nil {
		id = regionalRuleData(h.Item)
	} else {
		id = d.KeyColumnQuals["rule_id"].GetStringValue()
	}

	// Build the params
	param := &waf.GetRuleInput{
		RuleId: aws.String(id),
	}

	// Get call
	data, err := svc.GetRule(param)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "WAFNonexistentItemException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return data.Rule, nil
}

func getAwsWAFRegionalRuleAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWAFRegionalRuleAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	id := regionalRuleData(h.Item)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":waf-regional:" + region + ":" + commonColumnData.AccountId + ":rule" + "/" + id

	return aka, nil
}

func regionalRuleData(item interface{}) string {
	switch item := item.(type) {
	case *waf.RuleSummary:
		return *item.RuleId
	case *waf.Rule:
		return *item.RuleId
	}
	return ""
}
