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

func tableAwsWafv2RegexPatternSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_regex_pattern_set",
		Description: "AWS WAFv2 Regex Pattern Set",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"id", "name", "scope"}),
			Hydrate:           getAwsWafv2RegexPatternSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2RegexPatternSets,
		},
		GetMatrixItem: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Regex Pattern set.",
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
				Description: "A unique identifier for the Regex Pattern set.",
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
				Description: "A description of the Regex Pattern set that helps with identification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lock_token",
				Description: "A token used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "regular_expressions",
				Description: "The list of regular expression patterns in the set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2RegexPatternSet,
				Transform:   transform.FromField("RegularExpressionList").Transform(RegularExpressionObjectListToRegularExpressionList),
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
				Transform:   transform.FromField("Name"),
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
				Transform:   transform.From(regexPatternSetRegion),
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

func listAwsWafv2RegexPatternSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsWafv2RegexPatternSets", "AWS_REGION", region)

	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List all regional Regex Pattern Sets
	pagesLeft := true
	regionalRegexPatternSetParams := &wafv2.ListRegexPatternSetsInput{
		Scope: aws.String("REGIONAL"),
	}
	for pagesLeft {
		response, err := svc.ListRegexPatternSets(regionalRegexPatternSetParams)
		if err != nil {
			return nil, err
		}

		for _, regionalRegexPatternSets := range response.RegexPatternSets {
			d.StreamListItem(ctx, regionalRegexPatternSets)
		}

		if response.NextMarker != nil {
			pagesLeft = true
			regionalRegexPatternSetParams.NextMarker = response.NextMarker
		} else {
			pagesLeft = false
		}
	}

	// List all global Regex Pattern Sets
	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if region == "us-east-1" {
		pagesLeft = true
		globalRegexPatternSetParams := &wafv2.ListRegexPatternSetsInput{
			Scope: aws.String("CLOUDFRONT"),
		}
		for pagesLeft {
			response, err := svc.ListRegexPatternSets(globalRegexPatternSetParams)
			if err != nil {
				return nil, err
			}

			for _, globalRegexPatternSets := range response.RegexPatternSets {
				d.StreamListItem(ctx, globalRegexPatternSets)
			}

			if response.NextMarker != nil {
				pagesLeft = true
				globalRegexPatternSetParams.NextMarker = response.NextMarker
			} else {
				pagesLeft = false
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsWafv2RegexPatternSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWafv2RegexPatternSet")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

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

	params := &wafv2.GetRegexPatternSetInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: aws.String(scope),
	}

	op, err := svc.GetRegexPatternSet(params)
	if err != nil {
		plugin.Logger(ctx).Debug("GetRegexPatternSet", "ERROR", err)
		return nil, err
	}

	return op.RegexPatternSet, nil
}

func listTagsForAwsWafv2RegexPatternSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listTagsForAwsWafv2RegexPatternSet")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	data := regexPatternSetData(h.Item)
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

	regexPatternSetTags, err := svc.ListTagsForResource(param)
	if err != nil {
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
	plugin.Logger(ctx).Trace("regexPatternSetTagListToTurbotTags")
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

func RegularExpressionObjectListToRegularExpressionList(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("RegularExpressionObjectListToRegularExpressionList")
	data := d.HydrateItem.(*wafv2.RegexPatternSet)

	if data.RegularExpressionList == nil || len(data.RegularExpressionList) < 1 {
		return nil, nil
	}

	// Fetching the regex patterns from the array of object & storing as a list
	var regexList []string
	if data.RegularExpressionList != nil {
		for _, i := range data.RegularExpressionList {
			regexList = append(regexList, *i.RegexString)
		}
	}

	return regexList, nil
}

func regexPatternSetRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := regexPatternSetData(d.HydrateItem)
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

func regexPatternSetData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *wafv2.RegexPatternSet:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	case *wafv2.RegexPatternSetSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	}
	return data
}
