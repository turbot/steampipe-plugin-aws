package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudsearch"
	"github.com/aws/aws-sdk-go-v2/service/cloudsearch/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudSearchDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudsearch_domain",
		Description: "AWS CloudSearch Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("domain_name"),
			Hydrate:    getCloudSearchDomain,
			Tags:       map[string]string{"service": "cloudsearch", "action": "DescribeDomains"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudSearchDomains,
			Tags:    map[string]string{"service": "cloudsearch", "action": "ListDomainNames"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudSearchDomain,
				Tags: map[string]string{"service": "cloudsearch", "action": "DescribeDomains"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CLOUDSEARCH_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "A string that represents the name of a domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_id",
				Description: "An internally generated unique identifier for a domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the search domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudSearchDomain,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "created",
				Description: "True if the search domain is created.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "deleted",
				Description: "True if the search domain has been deleted.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "processing",
				Description: "True if processing is being done to activate the current domain configuration.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "requires_index_documents",
				Description: "True if Index Documents need to be called to activate the current domain configuration.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "search_instance_count",
				Description: "The number of search instances that are available to process search requests.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "search_instance_type",
				Description: "The instance type that is being used to process search requests.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "search_partition_count",
				Description: "The number of partitions across which the search index is spread.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "doc_service",
				Description: "The service endpoint for updating documents in a search domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "limits",
				Description: "Limit details for a search domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudSearchDomain,
			},
			{
				Name:        "search_service",
				Description: "The service endpoint for requesting search results from a search domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudSearchDomain,
			},

			// steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudSearchDomain,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudSearchDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudSearchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudsearch_domain.listCloudsearchDomains", "client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &cloudsearch.ListDomainNamesInput{}

	resp, err := svc.ListDomainNames(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudsearch_domain.listCloudsearchDomains", "api_error", err)
		return nil, nil
	}

	// Doesn't support paginator for CloudSearch ListDomainNames API
	for domainName := range resp.DomainNames {
		d.StreamListItem(ctx, types.DomainStatus{DomainName: aws.String(domainName)})

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudSearchDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["domain_name"].GetStringValue()
	var domainName string

	if name != "" {
		domainName = name
	}
	if h.Item != nil {
		domainName = *h.Item.(types.DomainStatus).DomainName
	}

	if domainName == "" {
		return nil, nil
	}

	// Get client
	svc, err := CloudSearchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudsearch_domain.getCloudsearchDomain", "client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &cloudsearch.DescribeDomainsInput{
		DomainNames: []string{domainName},
	}

	// execute list call
	item, err := svc.DescribeDomains(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudsearch_domain.getCloudsearchDomain", "api_error", err)
		return nil, err
	}

	if len(item.DomainStatusList) > 0 {
		return item.DomainStatusList[0], nil
	}

	return nil, nil
}
