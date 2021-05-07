package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsRoute53Domains(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_domains",
		Description: "AWS Route53 Domains",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("domain_name"),
			Hydrate:    getAwsRoute53Domain,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRoute53Domains,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "abuse_contact_email",
				Description: "Email address to contact to report incorrect contact information for a domain,to report that the domain is being used to send spam, to report that someone is cyber squatting on a domain name, or report some other type of abuse.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "abuse_contact_phone",
				Description: "Phone number for reporting abuse.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "admin_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "auto_renew",
				Description: "Indicates whether the domain is automatically renewed upon expiration.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "dns_sec",
				Description: "Reserved for future use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "expiration_date",
				Description: "The date when the registration for the domain is set to expire. The date and time is in Unix time format and Coordinated Universal time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "expiry",
				Description: "Expiration date of the domain in Unix time format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpirationDate"),
			},
			{
				Name:        "nameservers",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "registrant_contact",
				Description: "Provides details about the domain registrant.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "registrant_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "registrar_name",
				Description: "Name of the registrar of the domain as identified in the registry. Domains with a .com, .net, or .org TLD are registered by Amazon Registrar.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "registrar_url",
				Description: "Web address of the registrar.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "registry_domain_id",
				Description: "Reserved for future use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "reseller",
				Description: "Reseller of the domain. Domains registered or transferred using Route 53 domains will have Amazon as the reseller.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "status_list",
				Description: "An array of domain name status codes, also known as Extensible Provisioning Protocol (EPP) status codes.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "tech-contact",
				Description: "Provides details about the domain technical contact.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "tech_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "transfer_lock",
				Description: "Indicates whether a domain is locked from unauthorized transfer to another party.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "updated_date",
				Description: "The last updated date of the domain as found in the response to a WHOIS query.The date and time is in Unix time format and Coordinated Universal time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "Who_is_server",
				Description: "The fully qualified name of the WHOIS server that can answer the WHOIS query for the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "admin_contact",
				Description: "Provides details about the domain administrative contact.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53Domain,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53DomainTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53DomainTags,
				Transform:   transform.FromField("TagList").Transform(route53DomainsTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRoute53DomainAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsRoute53Domains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsRoute53Domains")

	// Create session
	svc, err := Route53DomainsService(ctx, d, "us-east-1")
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListDomainsPages(
		&route53domains.ListDomainsInput{},
		func(page *route53domains.ListDomainsOutput, isLast bool) bool {
			for _, domain := range page.Domains {
				d.StreamListItem(ctx, domain)
			}
			return !isLast
		},
	)
	return nil, err
}

////  HYDRATE FUNCTIONS

func getAwsRoute53Domain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsRoute53Domain")

	var name string
	if h.Item != nil {
		name = *h.Item.(*route53domains.DomainSummary).DomainName
	} else {
		name = d.KeyColumnQuals["domain_name"].GetStringValue()
	}
	// Create session
	svc, err := Route53DomainsService(ctx, d, "us-east-1")
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53domains.GetDomainDetailInput{
		DomainName: &name,
	}

	// Get call
	data, err := svc.GetDomainDetail(params)
	if err != nil {
		logger.Debug("getAwsRoute53Domain", "ERROR", err)
		return nil, err
	}
	return data, nil
}

func getAwsRoute53DomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsRoute53DomainTags")

	name := domainName(h.Item)

	// Create Session
	svc, err := Route53DomainsService(ctx, d, "us-east-1")
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53domains.ListTagsForDomainInput{
		DomainName: &name,
	}

	// Get call
	op, err := svc.ListTagsForDomain(params)
	if err != nil {
		logger.Debug("getAwsRoute53DomainTags", "ERROR", err)
		return nil, err
	}
	return op, nil
}

func getAwsRoute53DomainAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRoute53DomainAkas")

	name := domainName(h.Item)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":route53domains:" + ":" + commonColumnData.AccountId + ":" + "name" + "/" + name
	return []string{aka}, nil
}

//// TRANSFORM FUNCTIONS

func route53DomainsTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("route53DomainsTagListToTurbotTags")

	tags := d.HydrateItem.(*route53domains.ListTagsForDomainOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func domainName(item interface{}) string {
	switch item.(type) {
	case *route53domains.DomainSummary:
		return *item.(*route53domains.DomainSummary).DomainName
	case *route53domains.GetDomainDetailOutput:
		return *item.(*route53domains.GetDomainDetailOutput).DomainName
	}
	return ""
}
