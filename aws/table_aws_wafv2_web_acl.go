package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafv2WebAcl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_web_acl",
		Description: "AWS WAFv2 Web ACL",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "name", "scope"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"WAFNonexistentItemException", "WAFInvalidParameterException"}),
			},
			Hydrate: getAwsWafv2WebAcl,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2WebAcls,
		},
		GetMatrixItemFunc: BuildWafRegionList,
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
				Name:        "associated_resources",
				Description: "The array of Amazon Resource Names (ARNs) of the associated resources.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAssociatedResources,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "default_action",
				Description: "The action to perform if none of the Rules contained in the Web ACL match.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsWafv2WebAcl,
			},
			{
				Name:        "logging_configuration",
				Description: "The logging configuration for the specified web ACL.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLoggingConfiguration,
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
				Hydrate:     listTagsForAwsWafv2WebAcl,
				Transform:   transform.FromField("TagInfoForResource.TagList").Transform(webAclTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(transform.EnsureStringArray),
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	scope := aws.String("REGIONAL")

	if region == "global" {
		region = "us-east-1"
		scope = aws.String("CLOUDFRONT")
	}
	plugin.Logger(ctx).Trace("listAwsWafv2WebAcls", "AWS_REGION", region)

	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	pagesLeft := true
	params := &wafv2.ListWebACLsInput{
		Scope: scope,
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

func getAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWafv2WebAcl")

	region := d.KeyColumnQualString(matrixKeyRegion)

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

	/*
	 * The region endpoint is same for both Global Web ACL and the Regional Web ACL created in us-east-1.
	 * The following checks are required to remove duplicate resource entries due to above mentioned condition, when performing GET operation.
	 * To work with CloudFront, you must specify the Region US East (N. Virginia) or us-east-1
	 * For the Regional Web ACL, region value should not be 'global', as 'global' region is only used to get Global Web ACLs.
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
		plugin.Logger(ctx).Debug("GetWebACL", "ERROR", err)
		return nil, err
	}

	return op.WebACL, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listTagsForAwsWafv2WebAcl")

	region := d.KeyColumnQualString(matrixKeyRegion)

	if region == "global" {
		region = "us-east-1"
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

	// Build param with maximum limit set
	param := &wafv2.ListTagsForResourceInput{
		ResourceARN: aws.String(data["Arn"]),
		Limit:       aws.Int64(100),
	}

	webAclTags, err := svc.ListTagsForResource(param)
	if err != nil {
		return nil, err
	}
	return webAclTags, nil
}

func getLoggingConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoggingConfiguration")

	region := d.KeyColumnQualString(matrixKeyRegion)

	if region == "global" {
		region = "us-east-1"
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
	param := &wafv2.GetLoggingConfigurationInput{
		ResourceArn: aws.String(data["Arn"]),
	}

	op, err := svc.GetLoggingConfiguration(param)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "WAFNonexistentItemException" {
				return nil, nil
			}
		}
		return nil, err
	}
	return op, nil
}

func listAssociatedResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAssociatedResources")

	region := d.KeyColumnQualString(matrixKeyRegion)

	if region == "global" {
		region = "us-east-1"
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
	param := &wafv2.ListResourcesForWebACLInput{
		WebACLArn: aws.String(data["Arn"]),
	}

	op, err := svc.ListResourcesForWebACL(param)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "WAFNonexistentItemException" {
				return nil, nil
			}
		}
		return nil, err
	}

	if len(op.ResourceArns) == 0 {
		return nil, nil
	}

	return op.ResourceArns, nil
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

	region := d.MatrixItem[matrixKeyRegion]

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
