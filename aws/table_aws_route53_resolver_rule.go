package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRoute53ResolverRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_resolver_rule",
		Description: "AWS Route53 Resolver Rule",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getAwsRoute53ResolverRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRoute53ResolverRules,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name for the Resolver rule, which you specified when you created the Resolver rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID that Resolver assigned to the Resolver rule when you created it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) for the Resolver rule specified by Id.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "A code that specifies the current status of the Resolver rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creator_request_id",
				Description: "A unique string that you specified when you created the Resolver rule. CreatorRequestId identifies the request and allows failed requests to be retried without the risk of executing the operation twice.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_name",
				Description: "DNS queries for this domain name are forwarded to the IP addresses that are specified in TargetIps.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "When a rule is shared with another AWS account, the account ID of the account that the rule is shared with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resolver_endpoint_id",
				Description: "The ID of the endpoint that the rule is associated with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rule_type",
				Description: "When you want to forward DNS queries for specified domain name to resolvers on your network, specify FORWARD.When you have a forwarding rule to forward DNS queries for a domain to your network and you want Resolver to process queries for a subdomain of that domain, specify SYSTEM.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "share_status",
				Description: "Indicates whether the rules is shared and, if so, whether the current account is sharing the rule with another account, or another account is sharing the rule with the current account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_message",
				Description: "A detailed description of the status of a Resolver rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time that the Resolver rule was created, in Unix time format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "modification_time",
				Description: "The date and time that the Resolver rule was last updated, in Unix time format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resolver_rule_associations",
				Description: "The associations that were created between Resolver rules and VPCs using the current AWS account, and that match the specified filters, if any.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listResolverRuleAssociation,
			},
			{
				Name:        "target_ips",
				Description: "An array that contains the IP addresses and ports that an outbound endpoint forwards DNS queries to. Typically, these are the IP addresses of DNS resolvers on your network. Specify IPv4 addresses. IPv6 is not supported.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Resolver Rule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53ResolverRuleTags,
				Transform:   transform.FromField("Tags"),
			},
			// Standard columns for all tables
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
				Hydrate:     getAwsRoute53ResolverRuleTags,
				Transform:   transform.FromField("Tags").Transform(route53resolverRuleTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsRoute53ResolverRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsRoute53ResolverRules", "AWS_REGION", region)

	// Create session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListResolverRulesPages(
		&route53resolver.ListResolverRulesInput{},
		func(page *route53resolver.ListResolverRulesOutput, isLast bool) bool {
			for _, parameter := range page.ResolverRules {
				d.StreamListItem(ctx, parameter)

			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsRoute53ResolverRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsRoute53ResolverRule")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53resolver.GetResolverRuleInput{
		ResolverRuleId: &id,
	}

	// Get call
	data, err := svc.GetResolverRule(params)
	if err != nil {
		logger.Debug("getAwsRoute53ResolverRule", "ERROR", err)
		return nil, err
	}

	return data.ResolverRule, nil
}

func listResolverRuleAssociation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listResolverRuleAssociation")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	resolverRuleData := h.Item.(*route53resolver.ResolverRule)
	id := resolverRuleData.Id
	// Create Session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53resolver.ListResolverRuleAssociationsInput{
		Filters: []*route53resolver.Filter{
			{
				Name: types.String("ResolverRuleId"),
				Values: []*string{
					id,
				},
			},
		},
	}

	op, err := svc.ListResolverRuleAssociations(params)
	if err != nil {
		logger.Debug("listResolverRuleAssociation", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsRoute53ResolverRuleTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsRoute53ResolverRuleTags")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	resolverRuleData := h.Item.(*route53resolver.ResolverRule)

	// For default resolver rule i.e not supported tag
	defaultID := "rslvr-autodefined-rr-internet-resolver"
	if *resolverRuleData.Id == defaultID {
		return nil, nil
	}

	// Create Session
	svc, err := Route53ResolverService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53resolver.ListTagsForResourceInput{
		ResourceArn: resolverRuleData.Arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getAwsRoute53ResolverRuleTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func route53resolverRuleTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("route53resolverTagListToTurbotTags")

	if d.Value == nil {
		return nil, nil
	}
	tagList := d.Value.([]*route53resolver.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	} else {
		return nil, nil
	}

	return turbotTagsMap, nil
}
