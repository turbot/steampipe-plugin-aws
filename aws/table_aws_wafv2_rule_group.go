package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafv2RuleGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_rule_group",
		Description: "AWS WAFv2 Rule Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"id", "name", "scope"}),
			Hydrate:           getAwsWafv2RuleGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2RuleGroups,
		},
		GetMatrixItem: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the rule group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "id",
				Description: "A unique identifier for the rule group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "Specifies the scope of the Regex Pattern Set. Possible values are: 'REGIONAL' and 'CLOUDFRONT'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(regexPatternSetLocation),
			},
			{
				Name:        "description",
				Description: "A description of the rule group that helps with identification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity",
				Description: "The web ACL capacity units (WCUs) required for this rule group.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsWafv2RuleGroup,
			},
			{
				Name:        "lock_token",
				Description: "A token used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rules",
				Description: "The Rule statements used to identify the web requests that you want to allow, block, or count.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2RuleGroup,
			},
			{
				Name:        "visibility_config",
				Description: "Defines and enables Amazon CloudWatch metrics and web request sample collection.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2RuleGroup,
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
				Transform:   transform.FromField("Name"),
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
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},

			// aws standard columns
			{
				Name:        "partition",
				Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ruleGroupRegion),
			},
			{
				Name:        "account_id",
				Description: "The AWS Account ID in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
			},
		},
	}
}

//// LIST FUNCTION

func listAwsWafv2RuleGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsWafv2RuleGroups", "AWS_REGION", region)

	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List all regional Regex Pattern Sets
	pagesLeft := true
	regionalRuleGroupParams := &wafv2.ListRuleGroupsInput{
		Scope: aws.String("REGIONAL"),
	}
	for pagesLeft {
		response, err := svc.ListRuleGroups(regionalRuleGroupParams)
		if err != nil {
			return nil, err
		}

		for _, regionalRuleGroups := range response.RuleGroups {
			d.StreamListItem(ctx, regionalRuleGroups)
		}

		if response.NextMarker != nil {
			pagesLeft = true
			regionalRuleGroupParams.NextMarker = response.NextMarker
		} else {
			pagesLeft = false
		}
	}

	// List all global Regex Pattern Sets
	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if region == "us-east-1" {
		pagesLeft = true
		globalRuleGroupParams := &wafv2.ListRuleGroupsInput{
			Scope: aws.String("CLOUDFRONT"),
		}
		for pagesLeft {
			response, err := svc.ListRuleGroups(globalRuleGroupParams)
			if err != nil {
				return nil, err
			}

			for _, globalRuleGroups := range response.RuleGroups {
				d.StreamListItem(ctx, globalRuleGroups)
			}

			if response.NextMarker != nil {
				pagesLeft = true
				globalRuleGroupParams.NextMarker = response.NextMarker
			} else {
				pagesLeft = false
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsWafv2RuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWafv2EuleGroup")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

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
		id = d.KeyColumnQuals["id"].GetStringValue()
		name = d.KeyColumnQuals["name"].GetStringValue()
		scope = d.KeyColumnQuals["scope"].GetStringValue()
	}

	// Create Session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if strings.ToLower(scope) == "cloudfront" && region != "us-east-1" {
		return nil, nil
	}

	params := &wafv2.GetRuleGroupInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: aws.String(scope),
	}

	op, err := svc.GetRuleGroup(params)
	if err != nil {
		plugin.Logger(ctx).Debug("GetRuleGroup", "ERROR", err)
		return nil, err
	}

	return op.RuleGroup, nil
}

func listTagsForAwsWafv2RuleGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listTagsForAwsWafv2RuleGroup")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	data := ruleGroupData(h.Item)
	locationType := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if locationType == "global" && region != "us-east-1" {
		return nil, nil
	}

	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &wafv2.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
	}

	ruleGroupTags, err := svc.ListTagsForResource(param)
	if err != nil {
		return nil, err
	}
	return ruleGroupTags, nil
}

//// TRANSFORM FUNCTIONS

func regexPatternSetLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := ruleGroupData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]
	if loc == "regional" {
		return "REGIONAL", nil
	}
	return "CLOUDFRONT", nil
}

func ruleGroupTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ruleGroupTagListToTurbotTags")
	data := d.HydrateItem.(*wafv2.ListTagsForResourceOutput)

	if data.TagInfoForResource.TagList == nil || len(data.TagInfoForResource.TagList) < 1 {
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

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	if loc == "global" {
		return "global", nil
	}
	return region, nil
}

func ruleGroupData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *wafv2.RuleGroup:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	case *wafv2.RuleGroupSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	}
	return data
}
