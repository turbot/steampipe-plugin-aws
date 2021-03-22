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

func tableAwsWafv2WebAcl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_web_acl",
		Description: "AWS WAFv2 Web ACL",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"id", "name", "scope"}),
			ShouldIgnoreError: isNotFoundError([]string{"WAFNonexistentItemException"}),
			Hydrate:           getAwsWafv2WebAcl,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2WebAcl,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Web ACL. You cannot change the name of a Web ACL after you create it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "id",
				Description: "The unique identifier for the Web ACL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the Web ACL that helps with identification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "A description of the Web ACL that helps with identification.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(webAclLocation),
			},
			{
				Name:        "capacity",
				Description: "The web ACL capacity units (WCUs) currently being used by this web ACL.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsWafv2WebAcl,
			},
			{
				Name:        "lock_token",
				Description: "A token used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_by_firewall_manager",
				Description: "Indicates whether this web ACL is managed by AWS Firewall Manager.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "default_action",
				Description: "DefaultAction is a required field.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2WebAcl,
			},
			{
				Name:        "pre_process_firewall_manager_rule_groups",
				Description: "The first set of rules for AWS WAF to process in the web ACL.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2WebAcl,
			},
			{
				Name:        "post_process_firewall_manager_rule_groups",
				Description: "The last set of rules for AWS WAF to process in the web ACL.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2WebAcl,
			},
			{
				Name:        "rules",
				Description: "The Rule statements used to identify the web requests that you want to allow, block, or count.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2WebAcl,
			},
			{
				Name:        "visibility_config",
				Description: "Defines and enables Amazon CloudWatch metrics and web request sample collection.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2WebAcl,
			},
			// {
			// 	Name:        "tags_src",
			// 	Description: "A list of tags associated with the vault.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     listTagsForAwsWafv2WebAcl,
			// 	Transform:   transform.FromField("TagList"),
			// },

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			// {
			// 	Name:        "tags",
			// 	Description: resourceInterfaceDescription("tags"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("TagList").Transform(webAclTagListToTurbotTags),
			// },
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsWafv2WebAcl", "AWS_REGION", region)

	if region != "us-east-1"{
		return nil, nil
	}
	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	regionalParam := "REGIONAL"
	param := &wafv2.ListWebACLsInput{
		Scope: &regionalParam,
	}
	for {
		response, err := svc.ListWebACLs(param)
		if err != nil {
			return nil, err
		}
		for _, reginalWebACLs := range response.WebACLs {
			d.StreamListItem(ctx, reginalWebACLs)
		}
		if response.NextMarker == nil {
			break
		}
	}
	// globalParam := "CLOUDFRONT"
	param2 := &wafv2.ListWebACLsInput{
		Scope: aws.String("CLOUDFRONT"),
	}
	for {
		response, err := svc.ListWebACLs(param2)
		if err != nil {
			return nil, err
		}
		for _, globalWebACLs := range response.WebACLs {
			d.StreamListItem(ctx, globalWebACLs)
		}
		if response.NextMarker == nil {
			break
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsWafv2WebAcl")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var id string
	if h.Item != nil {
		id = webAclId(h.Item)
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}
	var name string
	if h.Item != nil {
		name = webAclName(h.Item)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}
	var scope string
	if h.Item != nil {
		data := h.Item.(*wafv2.WebACLSummary)
		scope = strings.Split(string(*data.ARN), ":")[5]
	} else {
		scope = d.KeyColumnQuals["scope"].GetStringValue()
	}

	if scope == "regional" {
		scope = "REGIONAL"
	} else {
		scope = "CLOUDFRONT"
	}

	// Create Session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &wafv2.GetWebACLInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: aws.String(scope),
	}

	op, err := svc.GetWebACL(params)
	if err != nil {
		logger.Debug("GetWebACL", "ERROR", err)
		return nil, err
	}

	return op.WebACL, nil
}

func listTagsForAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listTagsForAwsWafv2WebAcl")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	data := h.Item.(*wafv2.WebACLSummary)

	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &wafv2.ListTagsForResourceInput{
		ResourceARN: data.ARN,
	}

	webAclTags, err := svc.ListTagsForResource(param)
	if err != nil {
		return nil, err
	}
	return webAclTags, nil
}

//// TRANSFORM FUNCTIONS

func webAclLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*wafv2.WebACLSummary)
	loc := strings.Split(string(*data.ARN), ":")[5]
	if loc == "regional" {
		return "REGIONAL", nil
	}
	return "GLOBAL", nil
}

func webAclTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("webAclTagListToTurbotTags")
	webAclTag := d.HydrateItem.(*wafv2.TagInfoForResource)

	if webAclTag.TagList == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if webAclTag.TagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range webAclTag.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func webAclId(item interface{}) string {
	switch item.(type) {
	case *wafv2.WebACL:
		return *item.(*wafv2.WebACL).Id
	case *wafv2.WebACLSummary:
		return *item.(*wafv2.WebACLSummary).Id
	}
	return ""
}

func webAclName(item interface{}) string {
	switch item.(type) {
	case *wafv2.WebACL:
		return *item.(*wafv2.WebACL).Name
	case *wafv2.WebACLSummary:
		return *item.(*wafv2.WebACLSummary).Name
	}
	return ""
}
