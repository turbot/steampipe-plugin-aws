package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	"github.com/aws/aws-sdk-go-v2/service/wafregional/types"

	wafregionalv1 "github.com/aws/aws-sdk-go/service/wafregional"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsWAFRegionalRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafregional_rule",
		Description: "AWS WAF Regional Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("rule_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"WAFNonexistentItemException"}),
			},
			Hydrate: getAwsWAFRegionalRule,
			Tags:    map[string]string{"service": "waf-regional", "action": "GetRule"},
		},
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"WAFNonexistentItemException"}),
			},
			Hydrate: listAwsWAFRegionalRules,
			Tags:    map[string]string{"service": "waf-regional", "action": "ListRules"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wafregionalv1.EndpointsID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsWAFRegionalRule,
				Tags: map[string]string{"service": "waf-regional", "action": "GetRule"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name or description for the Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "Amazon Resource Name (ARN) of the Rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWAFRegionalRuleArn,
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
				Hydrate:     getAwsWAFRegionalRuleArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWAFRegionalRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule.listAwsWAFRegionalRules", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// List call
	maxItems := int32(100)
	params := &wafregional.ListRulesInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	// Minimunm limit is 0
	// https://docs.aws.amazon.com/wafregional/latest/APIReference/API_waf_ListRules.html
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.Limit = limit
		}
	}

	// API doesn't support aws-sdk-go-v2 paginator as of date
	pagesLeft := true
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.ListRules(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafregional_rule.listAwsWAFRegionalRules", "api_error", err)
			return nil, err
		}
		for _, rule := range response.Rules {
			d.StreamListItem(ctx, rule)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
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
	// Create Session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule.getAwsWAFRegionalRule", "api_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var id string
	if h.Item != nil {
		id = regionalRuleData(h.Item)
	} else {
		id = d.EqualsQuals["rule_id"].GetStringValue()
	}

	// Build the params
	param := &wafregional.GetRuleInput{
		RuleId: aws.String(id),
	}

	// Get call
	data, err := svc.GetRule(ctx, param)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "WAFNonexistentItemException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_wafregional_rule.getAwsWAFRegionalRule", "api_error", err)
		return nil, err
	}

	return data.Rule, nil
}

func getAwsWAFRegionalRuleArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	id := regionalRuleData(h.Item)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_rule.getAwsWAFRegionalRuleArn", "api_error", err)
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	aka := fmt.Sprintf("arn:%s:waf-regional:%s:%s:rule/%s", commonColumnData.Partition, region, commonColumnData.AccountId, id)

	return aka, nil
}

func regionalRuleData(item interface{}) string {
	switch item := item.(type) {
	case types.RuleSummary:
		return *item.RuleId
	case *types.Rule:
		return *item.RuleId
	}
	return ""
}
