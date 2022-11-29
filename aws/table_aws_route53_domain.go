package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/aws/aws-sdk-go-v2/service/route53domains/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRoute53Domain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_domain",
		Description: "AWS Route53 Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("domain_name"),
			Hydrate:    getRoute53Domain,
		},
		List: &plugin.ListConfig{
			Hydrate: listRoute53Domains,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53DomainARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "abuse_contact_email",
				Description: "Email address to contact to report incorrect contact information for a domain,to report that the domain is being used to send spam, to report that someone is cyber squatting on a domain name, or report some other type of abuse.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "abuse_contact_phone",
				Description: "Phone number for reporting abuse.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "admin_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "auto_renew",
				Description: "Indicates whether the domain is automatically renewed upon expiration.",
				Type:        proto.ColumnType_BOOL,
			},
			// As of May 25, 2021, API doesn't return the DNSSEC configuration in response.
			// {
			// 	Name:        "dns_sec",
			// 	Description: "Reserved for future use.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Hydrate:     getRoute53Domain,
			// },
			{
				Name:        "creation_date",
				Description: "The date when the domain was created as found in the response to a WHOIS query.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "expiration_date",
				Description: "The date when the registration for the domain is set to expire. The date and time is in Unix time format and Coordinated Universal time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "registrant_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "registrar_name",
				Description: "Name of the registrar of the domain as identified in the registry. Domains with a .com, .net, or .org TLD are registered by Amazon Registrar.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "registrar_url",
				Description: "Web address of the registrar.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "registry_domain_id",
				Description: "Reserved for future use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "reseller",
				Description: "Reseller of the domain. Domains registered or transferred using Route 53 domains will have Amazon as the reseller.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "tech_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getRoute53Domain,
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
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "who_is_server",
				Description: "The fully qualified name of the WHOIS server that can answer the WHOIS query for the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "nameservers",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "registrant_contact",
				Description: "Provides details about the domain registrant.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "status_list",
				Description: "An array of domain name status codes, also known as Extensible Provisioning Protocol (EPP) status codes.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "tech_contact",
				Description: "Provides details about the domain technical contact.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "admin_contact",
				Description: "Provides details about the domain administrative contact.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53Domain,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53DomainTags,
				Transform:   transform.FromField("TagList"),
			},

			// Steampipe standard columns
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
				Hydrate:     getRoute53DomainTags,
				Transform:   transform.FromField("TagList").Transform(route53DomainTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53DomainARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listRoute53Domains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := Route53DomainsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_domain.listRoute53Domains", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)
	input := route53domains.ListDomainsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 20 {
				maxItems = int32(20)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxItems = aws.Int32(maxItems)
	paginator := route53domains.NewListDomainsPaginator(svc, &input, func(o *route53domains.ListDomainsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_domain.listRoute53Domains", "api_error", err)
			return nil, err
		}

		for _, domain := range output.Domains {
			d.StreamListItem(ctx, domain)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRoute53Domain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		name = *h.Item.(types.DomainSummary).DomainName
	} else {
		name = d.KeyColumnQuals["domain_name"].GetStringValue()
	}

	// Create session
	svc, err := Route53DomainsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_domain.getRoute53Domain", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &route53domains.GetDomainDetailInput{
		DomainName: &name,
	}

	// Get call
	data, err := svc.GetDomainDetail(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_domain.getRoute53Domain", "api_error", err)
		return nil, err
	}
	return data, nil
}

func getRoute53DomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := domainName(h.Item)

	// Create session
	svc, err := Route53DomainsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_domain.getRoute53DomainTags", "client_error", err)
		return nil, err
	}

	// Build the params
	params := route53domains.ListTagsForDomainInput{
		DomainName: &name,
	}

	// Get call
	op, err := svc.ListTagsForDomain(ctx, &params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_domain.getRoute53DomainTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getRoute53DomainARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := domainName(h.Item)

	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":route53domains:::domain/" + name
	return arn, nil
}

//// TRANSFORM FUNCTIONS

func route53DomainTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {

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
	switch item := item.(type) {
	case types.DomainSummary:
		return *item.DomainName
	case route53domains.GetDomainDetailOutput:
		return *item.DomainName
	}
	return ""
}
