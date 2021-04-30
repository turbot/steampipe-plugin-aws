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
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getRoute53Domain,
		},
		List: &plugin.ListConfig{
			Hydrate: listRoute53Domains,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
			{
				Name:        "auto_renew",
				Description: "Indicates whether the domain is automatically renewed upon expiration.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "expiry",
				Description: "Expiration date of the domain in Unix time format and Coordinated Universal Time (UTC).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "transfer_lock",
				Description: "Indicates whether a domain is locked from unauthorized transfer to another party.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the parameter.",
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

func listRoute53Domains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

func getRoute53Domain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getRoute53Domain")

	// Create session
	svc, err := Route53DomainsService(ctx, d, "us-east-1")
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53domains.ListDomainsInput{}

	// Get call
	data, err := svc.ListDomains(params)
	if err != nil {
		logger.Debug("getRoute53Domain", "ERROR", err)
		return nil, err
	}
	return data.Domains[0], nil
}

func getAwsRoute53DomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsRoute53DomainTags")

	name := h.Item.(*route53domains.DomainSummary).DomainName

	// Create Session
	svc, err := Route53DomainsService(ctx, d, "us-east-1")
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &route53domains.ListTagsForDomainInput{
		DomainName: name,
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
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":route53domains:" + ":" + commonColumnData.AccountId
	return []string{aka}, nil
}

//// TRANSFORM FUNCTION

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
