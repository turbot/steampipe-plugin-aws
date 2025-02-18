package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafv2WebAcl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_web_acl",
		Description: "AWS WAFv2 Web ACL",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "name", "scope"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"WAFNonexistentItemException", "WAFInvalidParameterException"}),
			},
			Hydrate: getAwsWafv2WebAcl,
			Tags:    map[string]string{"service": "wafv2", "action": "GetWebACL"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2WebAcls,
			Tags:    map[string]string{"service": "wafv2", "action": "ListWebACLs"},
		},
		GetMatrixItemFunc: WAFRegionMatrix,
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsWafv2WebAcl,
				Tags: map[string]string{"service": "wafv2", "action": "GetWebACL"},
			},
			{
				Func: listAssociatedResources,
				Tags: map[string]string{"service": "wafv2", "action": "ListResourcesForWebACL"},
			},
			{
				Func: getLoggingConfiguration,
				Tags: map[string]string{"service": "wafv2", "action": "GetLoggingConfiguration"},
			},
			{
				Func: listTagsForAwsWafv2WebAcl,
				Tags: map[string]string{"service": "wafv2", "action": "ListTagsForResource"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
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
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(webAclRegion),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsWafv2WebAcls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	scope := types.ScopeRegional

	if region == "global" {
		scope = types.ScopeCloudfront
	}

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.listAwsWafv2WebAcls", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}
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
	params := &wafv2.ListWebACLsInput{
		Scope: scope,
		Limit: aws.Int32(maxLimit),
	}

	// ListWebACLs API doesn't support aws-sdk-go-v2 paginator yet
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.ListWebACLs(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafv2_web_acl.listAwsWafv2WebAcls", "api_error", err)
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

func getAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region := d.EqualsQualString(matrixKeyRegion)

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
		id = d.EqualsQuals["id"].GetStringValue()
		name = d.EqualsQuals["name"].GetStringValue()
		scope = d.EqualsQuals["scope"].GetStringValue()
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

	// Create Session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.getAwsWafv2WebAcl", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}
	params := &wafv2.GetWebACLInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: types.Scope(scope),
	}

	op, err := svc.GetWebACL(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.getAwsWafv2WebAcl", "api_error", err)
		return nil, err
	}

	return op.WebACL, nil
}

// ListTagsForResource.NextMarker return empty string in API call
// due to which pagination will not work properly
// https://github.com/aws/aws-sdk-go/issues/3513
func listTagsForAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	data := webAclData(h.Item)

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.listTagsForAwsWafv2WebAcl", "connection_error", err)
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

	webAclTags, err := svc.ListTagsForResource(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.listTagsForAwsWafv2WebAcl", "api_error", err)
		return nil, err
	}
	return webAclTags, nil
}

func getLoggingConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	data := webAclData(h.Item)

	// Create session
	svc, err := WAFV2Client(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.getLoggingConfiguration", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}
	// Build param
	param := &wafv2.GetLoggingConfigurationInput{
		ResourceArn: aws.String(data["Arn"]),
	}

	op, err := svc.GetLoggingConfiguration(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.getLoggingConfiguration", "api_error", err)
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "WAFNonexistentItemException" {
				return nil, nil
			}
		}
		return nil, err
	}
	return op, nil
}

func listAssociatedResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region := d.EqualsQualString(matrixKeyRegion)

	data := webAclData(h.Item)
	locationType := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

	// Create session
	if locationType == "global" {

		svc, err := CloudFrontClient(ctx, d)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafv2_web_acl.listAssociatedResources", "connection_error", err)
			return nil, err
		}
		if svc == nil {
			// unsupported region check
			return nil, nil
		}

		// Doc(https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListDistributionsByWebACLId.html) says
		// We need to pass the Web ACL ID to get the associated distrubutions with it but it doesn't.
		// By passing the Web ACL ARN we are getting the associated disctributions with it.
		// The AWS CLI behaves the same way as the API is behaving.
		// Build param
		param := &cloudfront.ListDistributionsByWebACLIdInput{
			WebACLId: aws.String(data["Arn"]),
		}

		op, err := svc.ListDistributionsByWebACLId(ctx, param)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafv2_web_acl.listAssociatedResources", "api_error", err)
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "WAFNonexistentItemException" {
					return nil, nil
				}
			}
			return nil, err
		}

		var ARNs []string
		if op.DistributionList != nil {
			if len(op.DistributionList.Items) > 0 {
				for _, item := range op.DistributionList.Items {
					ARNs = append(ARNs, *item.ARN)
				}
			}
		}
		return ARNs, nil
	} else {
		svc, err := WAFV2Client(ctx, d, region)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wafv2_web_acl.listAssociatedResources", "connection_error", err)
			return nil, err
		}
		if svc == nil {
			// unsupported region check
			return nil, nil
		}

		// Build param
		param := &wafv2.ListResourcesForWebACLInput{
			WebACLArn: aws.String(data["Arn"]),
		}

		var resourceArns []string

		resourceTypes := []types.ResourceType{types.ResourceTypeApplicationLoadBalancer, types.ResourceTypeApiGateway, types.ResourceTypeAppsync, types.ResourceTypeCognitioUserPool}

		for _, resourceType := range resourceTypes {
			param.ResourceType = resourceType
			res, err := listAssociatedResourcesByResourceType(ctx, svc, param)

			if err != nil {
				return nil, err
			}

			resourceArns = append(resourceArns, res...)
		}
		return resourceArns, nil
	}
}

func listAssociatedResourcesByResourceType(ctx context.Context, svc *wafv2.Client, input *wafv2.ListResourcesForWebACLInput) ([]string, error) {
	op, err := svc.ListResourcesForWebACL(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wafv2_web_acl.listAssociatedResourcesByResourceType", "api_error", err)
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "WAFNonexistentItemException" {
				return nil, nil
			}
		}
		return nil, err
	}
	if len(op.ResourceArns) == 0 {
		return []string{}, nil
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
	case *types.WebACL:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
	case types.WebACLSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
	}
	return data
}
