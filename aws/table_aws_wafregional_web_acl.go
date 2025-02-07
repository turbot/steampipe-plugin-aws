package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	"github.com/aws/aws-sdk-go-v2/service/wafregional/types"
	"github.com/aws/smithy-go"
	wafregionalEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafRegionalWebAcl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafregional_web_acl",
		Description: "AWS WAF Regional Web ACL",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"web_acl_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"WAFNonexistentItemException", "WAFInvalidParameterException"}),
			},
			Hydrate: getWafRegionalWebAcl,
			Tags:    map[string]string{"service": "waf-regional", "action": "GetWebACL"},
		},
		List: &plugin.ListConfig{
			Hydrate: listWafRegionalWebAcls,
			Tags:    map[string]string{"service": "waf-regional", "action": "ListWebACLs"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wafregionalEndpoint.AWS_WAF_REGIONAL_SERVICE_ID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getWafRegionalWebAcl,
				Tags: map[string]string{"service": "waf-regional", "action": "GetWebACL"},
			},
			{
				Func: getWafRegionalLoggingConfiguration,
				Tags: map[string]string{"service": "waf-regional", "action": "GetLoggingConfiguration"},
			},
			{
				Func: getWafRegionalResources,
				Tags: map[string]string{"service": "waf-regional", "action": "ListResourcesForWebACL"},
			},
			{
				Func: listTagsForWafRegionalWebAcl,
				Tags: map[string]string{"service": "waf-regional", "action": "ListTagsForResource"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Web ACL. You cannot change the name of a Web ACL after you create it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafRegionalWebAcl,
				Transform:   transform.FromField("WebACLArn"),
			},
			{
				Name:        "web_acl_id",
				Description: "The unique identifier for the Web ACL.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebACLId"),
			},
			{
				Name:        "default_action",
				Description: "The action to perform if none of the Rules contained in the WebACL match.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafRegionalWebAcl,
				Transform:   transform.FromField("DefaultAction.Type"),
			},
			{
				Name:        "metric_name",
				Description: "A friendly name or description for the metrics for this WebACL.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafRegionalWebAcl,
			},
			{
				Name:        "logging_configuration",
				Description: "The logging configuration for the web ACL.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRegionalLoggingConfiguration,
			},
			{
				Name:        "resources",
				Description: "An array of ARNs (Amazon Resource Names) of the resources associated with the web ACL.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRegionalResources,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "rules",
				Description: "The Rule statements used to identify the web requests that you want to allow, block, or count.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRegionalWebAcl,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForWafRegionalWebAcl,
				Transform:   transform.FromField("TagInfoForResource.TagList"),
			},

			// Steampipe standard columns
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
				Hydrate:     listTagsForWafRegionalWebAcl,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(wafRegionalWebAclTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafRegionalWebAcl,
				Transform:   transform.FromField("WebACLArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listWafRegionalWebAcls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.listWafRegionalWebAcls", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	params := &wafregional.ListWebACLsInput{
		Limit: maxItems,
	}

	// API doesn't support aws-sdk-go-v2 paginator as of date
	pagesLeft := true
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.ListWebACLs(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafregional_web_acl.listWafRegionalWebAcls", "api_error", err)
			return nil, err
		}

		for _, webAcl := range response.WebACLs {
			d.StreamListItem(ctx, webAcl)

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

func getWafRegionalWebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var id string
	if h.Item != nil {
		dataMap, err := getWafRegionalWebAclData(ctx, d, h)
		if err != nil {
			return nil, nil
		}
		data := dataMap.(map[string]string)
		id = data["ID"]
	} else {
		id = d.EqualsQuals["web_acl_id"].GetStringValue()
		if id == "" {
			return nil, nil
		}
	}

	// Create Session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.getWafRegionalWebAcl", "client_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	params := &wafregional.GetWebACLInput{
		WebACLId: aws.String(id),
	}

	op, err := svc.GetWebACL(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.getWafRegionalWebAcl", "api_error", err)
		return nil, err
	}

	return op.WebACL, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForWafRegionalWebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dataMap, err := getWafRegionalWebAclData(ctx, d, h)
	if err != nil {
		return nil, nil
	}
	data := dataMap.(map[string]string)

	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.listTagsForWafRegionalWebAcl", "client_error", err)
		return nil, err
	}

	// Build param with maximum limit set
	param := &wafregional.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
		Limit:       int32(100),
	}

	webAclTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.listTagsForWafRegionalWebAcl", "api_error", err)
		return nil, err
	}
	return webAclTags, nil
}

func getWafRegionalLoggingConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dataMap, err := getWafRegionalWebAclData(ctx, d, h)
	if err != nil {
		return nil, nil
	}
	data := dataMap.(map[string]string)

	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.getWafRegionalLoggingConfiguration", "client_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	// Build param
	param := &wafregional.GetLoggingConfigurationInput{
		ResourceArn: aws.String(data["Arn"]),
	}

	op, err := svc.GetLoggingConfiguration(ctx, param)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "WAFNonexistentItemException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.getWafRegionalLoggingConfiguration", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getWafRegionalResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dataMap, err := getWafRegionalWebAclData(ctx, d, h)
	if err != nil {
		return nil, err
	}
	data := dataMap.(map[string]string)

	// Create session
	svc, err := WAFRegionalClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.getWafRegionalResources", "client_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	// Build param
	param := &wafregional.ListResourcesForWebACLInput{
		WebACLId: aws.String(data["ID"]),
	}

	op, err := svc.ListResourcesForWebACL(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafregional_web_acl.getWafRegionalResources", "api_error", err)
		return nil, err
	}
	return op.ResourceArns, nil
}

func getWafRegionalWebAclData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := map[string]string{}
	switch item := h.Item.(type) {
	case *types.WebACL:
		data["ID"] = *item.WebACLId
		data["Arn"] = *item.WebACLArn
		data["Name"] = *item.Name
	case types.WebACLSummary:
		data["ID"] = *item.WebACLId
		commonData, err := getCommonColumns(ctx, d, h)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafregional_web_acl.getWafRegionalWebAclData", "cache_error", err)
			return nil, err
		}
		region := d.EqualsQualString(matrixKeyRegion)
		commonColumnData := commonData.(*awsCommonColumnData)
		data["Arn"] = fmt.Sprintf("arn:%s:waf-regional:%s:%s:webacl/%s", commonColumnData.Partition, region, commonColumnData.AccountId, *item.WebACLId)
		data["Name"] = *item.Name
	}
	return data, nil
}

//// TRANSFORM FUNCTIONS

func wafRegionalWebAclTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*wafregional.ListTagsForResourceOutput)

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
