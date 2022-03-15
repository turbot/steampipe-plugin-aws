package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableAwsWAFRegionalRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_waf_regional_rule",
		Description: "AWS WAF Regional Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("rule_id"),
			ShouldIgnoreError: isNotFoundError([]string{"WAFNonexistentItemException"}),
			Hydrate:           getAwsWAFRegionalRule,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"WAFNonexistentItemException"}),
			Hydrate:           listAwsWAFRegionalRules,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name for the rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rule_id",
				Description: "The ID of the Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metric_name",
				Description: "The name of the metric for the Rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWAFRegionalRule,
			},
			{
				Name:        "predicates",
				Description: "The Predicates object contains one Predicate element for each ByteMatchSet,IPSet, or SqlInjectionMatchSet object that you want to include in a Rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRegionalRule,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRegionalRuleTags,
				Transform:   transform.FromField("TagInfoForResource.TagList"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRegionalRuleTags,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(wafRegionalRuleTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWAFRegionalRuleAkas,
				Transform:   transform.FromValue(),
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

	// List call
	params := &waf.ListRulesInput{Limit: aws.Int64(100)}

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
		plugin.Logger(ctx).Trace("rules len : ", len(response.Rules))
		for _, rule := range response.Rules {
			d.StreamListItem(ctx, rule)
			plugin.Logger(ctx).Trace("response : ", rule)
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
		return nil, err
	}

	return data.Rule, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func getAwsWAFRegionalRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWAFRuleTags")

	region := d.KeyColumnQualString(matrixKeyRegion)
	id := regionalRuleData(h.Item)

	commonAwsColumns, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonAwsColumns.(*awsCommonColumnData)

	// Create Session
	svc, err := WAFRegionalService(ctx, d)
	if err != nil {
		return nil, err
	}

	aka := "arn:" + commonColumnData.Partition + ":waf-regional:" + region + ":" + commonColumnData.AccountId + ":rule" + "/" + id

	// Build param with maximum limit set
	params := &waf.ListTagsForResourceInput{
		ResourceARN: &aka,
		Limit:       aws.Int64(100),
	}

	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getAwsWAFRegionalRuleAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWAFRuleAkas")

	region := d.KeyColumnQualString(matrixKeyRegion)
	id := regionalRuleData(h.Item)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":waf-regional:" + region + ":" + commonColumnData.AccountId + ":rule" + "/" + id

	return []string{aka}, nil
}

//// TRANSFORM FUNCTION

func wafRegionalRuleTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tagListToTurbotTags")
	tagList := d.HydrateItem.(*waf.ListTagsForResourceOutput)

	if tagList.TagInfoForResource.TagList == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.TagInfoForResource.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
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
