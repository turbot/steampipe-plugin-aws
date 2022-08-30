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

//// TABLE DEFINITION

func tableAwsWafWebAcl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_waf_web_acl",
		Description: "AWS WAF Web ACL",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"web_acl_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"WAFNonexistentItemException", "WAFInvalidParameterException"}),
			},
			Hydrate: getWafWebAcl,
		},
		List: &plugin.ListConfig{
			Hydrate: listWafWebAcls,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Web ACL. You cannot change the name of a Web ACL after you create it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafWebAcl,
				Transform:   transform.FromField("WebACLArn"),
			},
			{
				Name:        "web_acl_id",
				Description: "The unique identifier for the Web ACL.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebACLId"),
			},
			{
				Name:        "default_action_type",
				Description: "The action to perform if none of the Rules contained in the WebACL match.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafWebAcl,
				Transform:   transform.FromField("DefaultAction.Type"),
			},
			{
				Name:        "metric_name",
				Description: "A friendly name or description for the metrics for this WebACL.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafWebAcl,
			},
			{
				Name:        "logging_configuration",
				Description: "The logging configuration for the specified web ACL.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getClassicLoggingConfiguration,
			},
			{
				Name:        "rules",
				Description: "The Rule statements used to identify the web requests that you want to allow, block, or count.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAcl,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForWafWebAcl,
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
				Hydrate:     listTagsForWafWebAcl,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(classicWebAclTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAcl,
				Transform:   transform.FromField("WebACLArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listWafWebAcls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WAFService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_web_acl.listWafWebAcls", "service_creation_error", err)
		return nil, err
	}

	pagesLeft := true
	params := &waf.ListWebACLsInput{
		Limit: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.Limit {
			if *limit < 1 {
				params.Limit = aws.Int64(1)
			} else {
				params.Limit = limit
			}
		}
	}

	for pagesLeft {
		response, err := svc.ListWebACLs(params)
		if err != nil {
			return nil, err
		}

		for _, webAcl := range response.WebACLs {
			d.StreamListItem(ctx, webAcl)

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

func getWafWebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var id string
	if h.Item != nil {
		data := classicWebAclData(h.Item, ctx, d, h)
		id = data["ID"]

	} else {
		id = d.KeyColumnQuals["web_acl_id"].GetStringValue()
	}

	// Create Session
	svc, err := WAFService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_web_acl.getWafWebAcl", "service_creation_error", err)
		return nil, err
	}

	params := &waf.GetWebACLInput{
		WebACLId: aws.String(id),
	}

	op, err := svc.GetWebACL(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_web_acl.getWafWebAcl", "api_error", err)
		return nil, err
	}

	return op.WebACL, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForWafWebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := classicWebAclData(h.Item, ctx, d, h)

	// Create session
	svc, err := WAFService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_web_acl.listTagsForWafWebAcl", "service_creation_error", err)
		return nil, err
	}

	// Build param with maximum limit set
	param := &waf.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
		Limit:       aws.Int64(100),
	}

	webAclTags, err := svc.ListTagsForResource(param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_waf_web_acl.listTagsForWafWebAcl", "api_error", err)
		return nil, err
	}
	return webAclTags, nil
}

func getClassicLoggingConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := classicWebAclData(h.Item, ctx, d, h)

	// Create session
	svc, err := WAFService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &waf.GetLoggingConfigurationInput{
		ResourceArn: aws.String(data["Arn"]),
	}

	op, err := svc.GetLoggingConfiguration(param)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NonexistentItemException" {
				return nil, nil
			}
		}
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func classicWebAclTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*waf.ListTagsForResourceOutput)

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

func classicWebAclData(item interface{}, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *waf.WebACL:
		data["ID"] = *item.WebACLId
		data["Arn"] = *item.WebACLArn
		data["Name"] = *item.Name
	case *waf.WebACLSummary:
		data["ID"] = *item.WebACLId
		getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
		commonData, err := getCommonColumnsCached(ctx, d, h)
		if err != nil {
			return nil
		}
		commonColumnData := commonData.(*awsCommonColumnData)
		data["Arn"] = "arn:aws:waf::" + commonColumnData.AccountId + ":webacl/" + (*aws.String(*item.WebACLId))
		data["Name"] = *item.Name
	}
	return data
}
