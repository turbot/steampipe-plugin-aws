package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafv2RuleGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_rule_group",
		Description: "AWS WAFv2 Rule Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "name", "scope"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"WAFInvalidParameterException", "WAFNonexistentItemException", "ValidationException", "InvalidParameter"}),
			},
			Hydrate: getAwsWafv2RuleGroup,
			Tags:    map[string]string{"service": "wafv2", "action": "GetRuleGroup"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2RuleGroups,
			Tags:    map[string]string{"service": "wafv2", "action": "ListRuleGroups"},
		},
		GetMatrixItemFunc: WAFRegionMatrix,
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsWafv2RuleGroup,
				Tags: map[string]string{"service": "wafv2", "action": "GetRuleGroup"},
			},
			{
				Func: listTagsForAwsWafv2RuleGroup,
				Tags: map[string]string{"service": "wafv2", "action": "ListTagsForResource"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the rule group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "RuleGroup.Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN", "RuleGroup.ARN"),
			},
			{
				Name:        "id",
				Description: "A unique identifier for the rule group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "RuleGroup.Id"),
			},
			{
				Name:        "scope",
				Description: "Specifies the scope of the rule group. Possible values are: 'REGIONAL' and 'CLOUDFRONT'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ruleGroupLocation),
			},
			{
				Name:        "capacity",
				Description: "The web ACL capacity units (WCUs) required for this rule group.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsWafv2RuleGroup,
				Transform:   transform.FromField("RuleGroup.Capacity"),
			},
			{
				Name:        "description",
				Description: "A description of the rule group that helps with identification.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description", "RuleGroup.Description"),
			},
			{
				Name:        "lock_token",
				Description: "A token used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LockToken", "RuleGroupSummary.LockToken"),
			},
			{
				Name:        "rules",
				Description: "The Rule statements used to identify the web requests that you want to allow, block, or count.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2RuleGroup,
				Transform:   transform.FromField("RuleGroup.Rules"),
			},
			{
				Name:        "visibility_config",
				Description: "Defines and enables Amazon CloudWatch metrics and web request sample collection.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2RuleGroup,
				Transform:   transform.FromField("RuleGroup.VisibilityConfig"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAwsWafv2RuleGroup,
				Transform:   transform.FromField("TagInfoForResource.TagList"),
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "RuleGroup.Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAwsWafv2RuleGroup,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(ruleGroupTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN", "RuleGroup.ARN").Transform(arnToAkas),
			},

			// AWS standard columns
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ruleGroupRegion),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWafv2RuleGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	scope := types.ScopeRegional

	if region == "global" {
		scope = types.ScopeCloudfront
	}

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_rule_group.listAwsWafv2RuleGroups", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// List all rule groups
	pagesLeft := true
	maxLimit := int32(100)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	params := &wafv2.ListRuleGroupsInput{
		Scope: scope,
		Limit: aws.Int32(maxLimit),
	}

	// ListRuleGroups API doesn't support aws-sdk-go-v2 paginator yet
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.ListRuleGroups(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafv2_rule_group.listAwsWafv2RuleGroups", "api_error", err)
			return nil, err
		}

		for _, ruleGroups := range response.RuleGroups {
			d.StreamListItem(ctx, ruleGroups)

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

func getAwsWafv2RuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region := d.EqualsQualString(matrixKeyRegion)

	var id, name, scope string
	if h.Item != nil {
		data := ruleGroupData(h.Item)
		id = data["ID"]
		name = data["Name"]
		locationType := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

		if locationType == "regional" {
			scope = "REGIONAL"
		} else {
			scope = "CLOUDFRONT"
		}
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
		name = d.EqualsQuals["name"].GetStringValue()
		scope = d.EqualsQuals["scope"].GetStringValue()
	}

	/*
	 * The region endpoint is same for both Global Rule Group and the Regional Rule Group created in us-east-1.
	 * The following checks are required to remove duplicate resource entries due to above mentioned condition, when performing GET operation.
	 * To work with CloudFront, you must specify the Region US East (N. Virginia) or us-east-1
	 * For the Regional Rule Group, region value should not be 'global', as 'global' region is only used to get Global Rule Groups.
	 * For any other region, region value will be same as working region.
	 */
	if scope == "REGIONAL" && region == "global" {
		return nil, nil
	}

	if strings.ToLower(scope) == "cloudfront" && region != "global" {
		return nil, nil
	}

	// Create Session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_rule_group.getAwsWafv2RuleGroup", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	params := &wafv2.GetRuleGroupInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: types.Scope(scope),
	}

	op, err := svc.GetRuleGroup(ctx, params)
	if err != nil {

		plugin.Logger(ctx).Error("aws_wafv2_rule_group.getAwsWafv2RuleGroup", "api_error", err)
		return nil, err
	}

	return op, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go-v2/issues/3513
func listTagsForAwsWafv2RuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	data := ruleGroupData(h.Item)

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_rule_group.listTagsForAwsWafv2RuleGroup", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// Build param with maximum limit set
	param := &wafv2.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
		Limit:       aws.Int32(100),
	}

	ruleGroupTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_rule_group.listTagsForAwsWafv2RuleGroup", "api_error", err)
		return nil, err
	}
	return ruleGroupTags, nil
}

//// TRANSFORM FUNCTIONS

func ruleGroupLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := ruleGroupData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]
	if loc == "regional" {
		return "REGIONAL", nil
	}
	return "CLOUDFRONT", nil
}

func ruleGroupTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {

	data := d.HydrateItem.(*wafv2.ListTagsForResourceOutput)

	if len(data.TagInfoForResource.TagList) < 1 {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if data.TagInfoForResource.TagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range data.TagInfoForResource.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func ruleGroupRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := ruleGroupData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

	region := d.MatrixItem[matrixKeyRegion]

	if loc == "global" {
		return "global", nil
	}
	return region, nil
}

func ruleGroupData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *wafv2.GetRuleGroupOutput:
		data["ID"] = *item.RuleGroup.Id
		data["Arn"] = *item.RuleGroup.ARN
		data["Name"] = *item.RuleGroup.Name
		data["Description"] = *item.RuleGroup.Description
	case types.RuleGroupSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	}
	return data
}
