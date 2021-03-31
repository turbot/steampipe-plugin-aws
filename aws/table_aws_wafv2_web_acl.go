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
			Hydrate: listAwsWafv2WebAcls,
		},
		GetMatrixItem: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Web ACL. You cannot change the name of a Web ACL after you create it.",
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
				Description: "The unique identifier for the Web ACL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "Specifies the scope of the Web ACL. Possibles values are: 'REGIONAL' and 'CLOUDFRONT'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(webAclLocation),
			},
			{
				Name:        "description",
				Description: "A description of the Web ACL that helps with identification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity",
				Description: "The Web ACL capacity units(WCUs) currently being used by this resource.",
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
				Hydrate:     getAwsWafv2WebAcl,
			},
			{
				Name:        "default_action",
				Description: "The action to perform if none of the Rules contained in the Web ACL match.",
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
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForAwsWafv2WebAcl,
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
				Hydrate:     listTagsForAwsWafv2WebAcl,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(webAclTagListToTurbotTags),
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
				Transform:   transform.From(webAclRegion),
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

func listAwsWafv2WebAcls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsWafv2WebAcls", "AWS_REGION", region)

	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List all regional web acls
	pagesLeft := true
	regionalWebAclParams := &wafv2.ListWebACLsInput{
		Scope: aws.String("REGIONAL"),
	}
	for pagesLeft {
		response, err := svc.ListWebACLs(regionalWebAclParams)
		if err != nil {
			return nil, err
		}

		for _, reginalWebACLs := range response.WebACLs {
			d.StreamListItem(ctx, reginalWebACLs)
		}

		if response.NextMarker != nil {
			pagesLeft = true
			regionalWebAclParams.NextMarker = response.NextMarker
		} else {
			pagesLeft = false
		}
	}

	// List all global web acls
	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if region == "us-east-1" {
		pagesLeft = true
		globalWebAclParams := &wafv2.ListWebACLsInput{
			Scope: aws.String("CLOUDFRONT"),
		}
		for pagesLeft {
			response, err := svc.ListWebACLs(globalWebAclParams)
			if err != nil {
				return nil, err
			}

			for _, globalWebACLs := range response.WebACLs {
				d.StreamListItem(ctx, globalWebACLs)
			}

			if response.NextMarker != nil {
				pagesLeft = true
				globalWebAclParams.NextMarker = response.NextMarker
			} else {
				pagesLeft = false
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWafv2WebAcl")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var id, name, scope string
	if h.Item != nil {
		data := webAclData(h.Item)
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

	params := &wafv2.GetWebACLInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: aws.String(scope),
	}

	op, err := svc.GetWebACL(params)
	if err != nil {
		plugin.Logger(ctx).Debug("GetWebACL", "ERROR", err)
		return nil, err
	}

	return op.WebACL, nil
}

func listTagsForAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listTagsForAwsWafv2WebAcl")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	data := webAclData(h.Item)
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

	webAclTags, err := svc.ListTagsForResource(param)
	if err != nil {
		return nil, err
	}
	return webAclTags, nil
}

//// TRANSFORM FUNCTIONS

func webAclLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := webAclData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]
	if loc == "regional" {
		return "REGIONAL", nil
	}
	return "CLOUDFRONT", nil
}

func webAclTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("webAclTagListToTurbotTags")
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

func webAclRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := webAclData(d.HydrateItem)
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

func webAclData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *wafv2.WebACL:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
	case *wafv2.WebACLSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
	}
	return data
}
