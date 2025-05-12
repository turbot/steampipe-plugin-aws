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

func tableAwsWafv2RegexPatternSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_regex_pattern_set",
		Description: "AWS WAFv2 Regex Pattern Set",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "name", "scope"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"WAFInvalidParameterException", "WAFNonexistentItemException", "ValidationException", "InvalidParameter"}),
			},
			Hydrate: getAwsWafv2RegexPatternSet,
			Tags:    map[string]string{"service": "wafv2", "action": "GetRegexPatternSet"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2RegexPatternSets,
			Tags:    map[string]string{"service": "wafv2", "action": "ListRegexPatternSets"},
		},
		GetMatrixItemFunc: WAFRegionMatrix,
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsWafv2RegexPatternSet,
				Tags: map[string]string{"service": "wafv2", "action": "GetRegexPatternSet"},
			},
			{
				Func: listTagsForAwsWafv2RegexPatternSet,
				Tags: map[string]string{"service": "wafv2", "action": "ListTagsForResource"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Regex Pattern set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "RegexPatternSet.Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN", "RegexPatternSet.ARN"),
			},
			{
				Name:        "id",
				Description: "A unique identifier for the Regex Pattern set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "RegexPatternSet.Id"),
			},
			{
				Name:        "scope",
				Description: "Specifies the scope of the Regex Pattern Set. Possible values are: 'REGIONAL' and 'CLOUDFRONT'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(regexPatternSetLocation),
			},
			{
				Name:        "description",
				Description: "A description of the Regex Pattern set that helps with identification.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description", "RegexPatternSet.Description"),
			},
			{
				Name:        "lock_token",
				Description: "A token used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LockToken", "RegexPatternSetSummary.LockToken"),
			},
			{
				Name:        "regular_expressions",
				Description: "The list of regular expression patterns in the set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2RegexPatternSet,
				Transform:   transform.FromField("RegexPatternSet.RegularExpressionList").Transform(regularExpressionObjectListToRegularExpressionList),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAwsWafv2RegexPatternSet,
				Transform:   transform.FromField("TagInfoForResource.TagList"),
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "RegexPatternSet.Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAwsWafv2RegexPatternSet,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(regexPatternSetTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN", "RegexPatternSet.ARN").Transform(arnToAkas),
			},

			// AWS standard columns
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region", "RegexPatternSet.Region").Transform(regexPatternSetRegion),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWafv2RegexPatternSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	scope := types.ScopeRegional

	if region == "global" {
		scope = types.ScopeCloudfront
	}
	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_regex_pattern_set.listAwsWafv2RegexPatternSets", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// List all Regex Pattern Sets
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
	params := &wafv2.ListRegexPatternSetsInput{
		Scope: scope,
		Limit: aws.Int32(maxLimit),
	}

	// ListRegexPatternSets API doesn't support aws-sdk-go-v2 paginator yet
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.ListRegexPatternSets(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafv2_regex_pattern_set.listAwsWafv2RegexPatternSets", "api_error", err)
			return nil, err
		}

		for _, regexPatternSets := range response.RegexPatternSets {
			d.StreamListItem(ctx, regexPatternSets)

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

func getAwsWafv2RegexPatternSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	var id, name, scope string
	if h.Item != nil {
		data := regexPatternSetData(h.Item)
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
	 * The region endpoint is same for both Global Regex Pattern Set and the Regional Regex Pattern Set created in us-east-1.
	 * The following checks are required to remove duplicate resource entries due to above mentioned condition, when performing GET operation.
	 * To work with CloudFront, you must specify the Region US East (N. Virginia) or us-east-1
	 * For the Regional Regex Pattern Set, region value should not be 'global', as 'global' region is only used to get Global Regex Pattern Sets.
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
		plugin.Logger(ctx).Error("aws_wafv2_regex_pattern_set.getAwsWafv2RegexPatternSet", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	params := &wafv2.GetRegexPatternSetInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: types.Scope(scope),
	}

	op, err := svc.GetRegexPatternSet(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_regex_pattern_set.getAwsWafv2RegexPatternSet", "api_error", err)
		return nil, err
	}

	return op, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForAwsWafv2RegexPatternSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region := d.EqualsQualString(matrixKeyRegion)

	data := regexPatternSetData(h.Item)

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_regex_pattern_set.listTagsForAwsWafv2RegexPatternSet", "connection_error", err)
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

	regexPatternSetTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_regex_pattern_set.listTagsForAwsWafv2RegexPatternSet", "api_error", err)
		return nil, err
	}
	return regexPatternSetTags, nil
}

//// TRANSFORM FUNCTIONS

func regexPatternSetLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := regexPatternSetData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]
	if loc == "regional" {
		return "REGIONAL", nil
	}
	return "CLOUDFRONT", nil
}

func regexPatternSetTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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

func regularExpressionObjectListToRegularExpressionList(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*wafv2.GetRegexPatternSetOutput)

	if len(data.RegexPatternSet.RegularExpressionList) < 1 {
		return nil, nil
	}

	// Fetching the regex patterns from the array of object & storing as a list
	var regexList []string
	if data.RegexPatternSet.RegularExpressionList != nil {
		for _, i := range data.RegexPatternSet.RegularExpressionList {
			regexList = append(regexList, *i.RegexString)
		}
	}

	return regexList, nil
}

func regexPatternSetRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := regexPatternSetData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

	region := d.MatrixItem[matrixKeyRegion]

	if loc == "global" {
		return "global", nil
	}
	return region, nil
}

func regexPatternSetData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *wafv2.GetRegexPatternSetOutput:
		data["ID"] = *item.RegexPatternSet.Id
		data["Arn"] = *item.RegexPatternSet.ARN
		data["Name"] = *item.RegexPatternSet.Name
		data["Description"] = *item.RegexPatternSet.Description
	case types.RegexPatternSetSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	}
	return data
}
