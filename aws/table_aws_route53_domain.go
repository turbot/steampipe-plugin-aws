package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsRoute53Domain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_domain",
		Description: "AWS Route53 Domain",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("domain_name"),
			ItemFromKey: domains,
			Hydrate:     getRoute53Domains,
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
				Name:        "creation_date",
				Description: "A comment for the zone",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "expiration_date",
				Description: "Expiration date of the domain.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "auto_renew",
				Description: "Indicates whether the domain is automatically renewed upon expiration.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "transferlock",
				Description: "Indicates whether a domain is locked from unauthorized transfer to another party",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "admin_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries. If the value is true, WHOIS (who is) queries return contact information either for Amazon Registrar (for .com, .net, and .org domains) or for our registrar associate, Gandi (for all other TLDs). If the value is false, WHOIS queries return the information that you entered for the admin contact.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "registrant_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries. If the value is true, WHOIS (who is) queries return contact information either for Amazon Registrar (for .com, .net, and .org domains) or for our registrar associate, Gandi (for all other TLDs). If the value is false, WHOIS queries return the information that you entered for the registrant contact (domain owner).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "tech_privacy",
				Description: "Specifies whether contact information is concealed from WHOIS queries. If the value is true, WHOIS (who is) queries return contact information either for Amazon Registrar (for .com, .net, and .org domains) or for our registrar associate, Gandi (for all other TLDs). If the value is false, WHOIS queries return the information that you entered for the technical contact.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "registrar_name",
				Description: "Name of the registrar of the domain as identified in the registry",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "who_is_server",
				Description: "The fully qualified name of the WHOIS server that can answer the WHOIS query for the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "registrar_url",
				Description: "Web address of the registrar.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "abuse_contact_email",
				Description: "Email address to contact to report incorrect contact information for a domain, to report that the domain is being used to send spam, to report that someone is cybersquatting on a domain name, or report some other type of abuse.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "abuse_contact_phone",
				Description: "Phone number for reporting abuse.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_list",
				Description: "Specifies the status of a variety of operations on a domain name",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "nameservers",
				Description: "Specifies the fully qualified host name of the name server and Glue IP address of a name server entry",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "admin_contact",
				Description: "Provides details about the domain administrative contact.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "registrant_contact",
				Description: "Provides details about the domain registrant.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tech_contact",
				Description: "Provides details about the domain technical contact.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
		}),
	}
}

// BUILD HYDRATE INPUT

func domains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["domain_name"].GetStringValue()
	item := &route53domains.DomainSummary{
		DomainName: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listRoute53Domains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := "us-east-1"
	plugin.Logger(ctx).Trace("listRoute53Domains", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Route53DomainsService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.ListDomainsPages(
		&route53domains.ListDomainsInput{},
		func(page *route53domains.ListDomainsOutput, isLast bool) bool {
			for _, domains := range page.Domains {
				d.StreamListItem(ctx, domains)
			}
			return true
		},
	)

	return nil, err
}

// HYDRATE FUNCTIONS

func getRoute53Domains(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRoute53Domains")

	defaultRegion := "us-east-1"
	domains := h.Item.(*route53domains.DomainSummary)

	// Create session
	svc, err := Route53DomainsService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	if len(*domains.DomainName) < 1 {
		return nil, nil
	}

	params := &route53domains.GetDomainDetailInput{
		DomainName: domains.DomainName,
	}

	item, err := svc.GetDomainDetail(params)
	if err != nil {
		return nil, err
	}

	return item.DomainName, nil
}
