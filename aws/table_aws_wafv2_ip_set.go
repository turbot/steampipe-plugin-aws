package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafv2IpSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_ip_set",
		Description: "AWS WAFv2 IP Set",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "name", "scope"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"WAFNonexistentItemException", "WAFInvalidParameterException", "InvalidParameter", "ValidationException"}),
			},
			Hydrate: getAwsWafv2IpSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2IpSets,
		},
		GetMatrixItemFunc: BuildWafRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the IP set.",
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
				Description: "A unique identifier for the IP set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "Specifies the scope of the IP Set. Possible values are: 'REGIONAL' and 'CLOUDFRONT'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ipSetLocation),
			},
			{
				Name:        "description",
				Description: "A description of the IP set that helps with identification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address_version",
				Description: "Specifies the IP address type. Possible values are: 'IPV4' and 'IPV6'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWafv2IpSet,
				Transform:   transform.FromField("IPAddressVersion"),
			},
			{
				Name:        "lock_token",
				Description: "A token used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "addresses",
				Description: "An array of strings that specify one or more IP addresses or blocks of IP addresses in Classless Inter-Domain Routing (CIDR) notation.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2IpSet,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAwsWafv2IpSet,
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
				Hydrate:     listTagsForAwsWafv2IpSet,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(ipSetTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},

			// AWS standard columns
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
				Transform:   transform.From(ipSetRegion),
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

func listAwsWafv2IpSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	scope := types.ScopeRegional

	if region == "global" {
		region = "us-east-1"
		scope = types.ScopeCloudfront
	}

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_ip_set.listAwsWafv2IpSets", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// List all IP sets
	pagesLeft := true
	maxLimit := int32(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := d.QueryContext.Limit
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	params := &wafv2.ListIPSetsInput{
		Scope: scope,
		Limit: aws.Int32(maxLimit),
	}

	// ListIPSets API doesn't support aws-sdk-go-v2 paginator yet
	for pagesLeft {
		response, err := svc.ListIPSets(ctx, params)
		if err != nil {
			// Service Client doesn't throw any error if region is not supported but the API throws no such host error for that region
			if strings.Contains(err.Error(), "no such host") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_wafv2_ip_set.listAwsWafv2IpSets", "api_error", err)
			return nil, err
		}

		for _, ipSets := range response.IPSets {
			d.StreamListItem(ctx, ipSets)

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

func getAwsWafv2IpSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	var id, name, scope string
	if h.Item != nil {
		data := ipSetData(h.Item)
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

	/*
	 * The region endpoint is same for both Global IP Set and the Regional IP Set created in us-east-1.
	 * The following checks are required to remove duplicate resource entries due to above mentioned condition, when performing GET operation.
	 * To work with CloudFront, you must specify the Region US East (N. Virginia) or us-east-1
	 * For the Regional IP Set, region value should not be 'global', as 'global' region is only used to get Global IP Sets.
	 * For any other region, region value will be same as working region.
	 */
	if scope == "REGIONAL" && region == "global" {
		return nil, nil
	}

	if strings.ToLower(scope) == "cloudfront" && region != "global" {
		return nil, nil
	}

	if region == "global" {
		region = "us-east-1"
	}

	// Create Session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_ip_set.getAwsWafv2IpSet", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	params := &wafv2.GetIPSetInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: types.Scope(scope),
	}

	op, err := svc.GetIPSet(ctx, params)
	if err != nil {
		// Service Client doesn't throw any error if region is not supported but the API throws no such host error for that region
		if strings.Contains(err.Error(), "no such host") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_wafv2_ip_set.getAwsWafv2IpSet", "api_error", err)
		return nil, err
	}

	return op.IPSet, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForAwsWafv2IpSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	if region == "global" {
		region = "us-east-1"
	}
	data := ipSetData(h.Item)
	locationType := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]
	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if locationType == "global" && region != "us-east-1" {
		return nil, nil
	}

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_ip_set.listTagsForAwsWafv2IpSet", "connection_error", err)
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

	ipSetTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_ip_set.listTagsForAwsWafv2IpSet", "api_error", err)
		return nil, err
	}
	return ipSetTags, nil
}

//// TRANSFORM FUNCTIONS

func ipSetLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := ipSetData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]
	if loc == "regional" {
		return "REGIONAL", nil
	}
	return "CLOUDFRONT", nil
}

func ipSetTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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

func ipSetRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := ipSetData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

	region := d.MatrixItem[matrixKeyRegion]

	if loc == "global" {
		return "global", nil
	}
	return region, nil
}

func ipSetData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *types.IPSet:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	case types.IPSetSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	}
	return data
}
