package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/opensearchservice"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsOpenSearchDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_opensearch_domain",
		Description: "AWS OpenSearch Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("domain_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getOpenSearchDomain,
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenSearchDomains,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:    listOpenSearchDomainTags,
				Depends: []plugin.HydrateFunc{getOpenSearchDomain},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_id",
				Description: "The unique identifier for the specified domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOpenSearchDomain,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "access_policies",
				Description: "The IAM access policies of the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "created",
				Description: "The domain creation status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "deleted",
				Description: "The domain deletion status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "endpoint",
				Description: "The domain endpoint that is used to submit index and search requests.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "engine_version",
				Description: "The domain's OpenSearch version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "processing",
				Description: "The status of the domain configuration.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "upgrade_processing",
				Description: "The status of the domain version upgrade.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "node_to_node_encryption_options_enabled",
				Description: "Specifies the status of the node to node encryption status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getOpenSearchDomain,
				Transform:   transform.FromField("NodeToNodeEncryptionOptions.Enabled"),
			},
			{
				Name:        "advanced_options",
				Description: "Specifies the status of the advanced options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "advanced_security_options",
				Description: "Specifies The current status of the OpenSearch domain's advanced security options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "auto_tune_options",
				Description: "The current status of the domain's auto-tune options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "cluster_config",
				Description: "The type and number of instances in the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "cognito_options",
				Description: "The cognito options for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "domain_endpoint_options",
				Description: "The current status of the domain's endpoint options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "ebs_options",
				Description: "The EBSOptions for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
				Transform:   transform.FromField("EBSOptions"),
			},
			{
				Name:        "encryption_at_rest_options",
				Description: "The status of the encryption at rest options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "endpoints",
				Description: "Map containing the domain endpoints used to submit index and search requests.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "log_publishing_options",
				Description: "Log publishing options for the given domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "service_software_options",
				Description: "The current status of the domain's service software.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "snapshot_options",
				Description: "Specifies the status of the snapshot options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
			},
			{
				Name:        "vpc_options",
				Description: "The vpc options for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
				Transform:   transform.FromField("VPCOptions"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listOpenSearchDomainTags,
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
				Hydrate:     listOpenSearchDomainTags,
				Transform:   transform.FromField("TagList").Transform(openSearchDomaintagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchDomain,
				Transform:   transform.FromField("ARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listOpenSearchDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := OpenSearchService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	params := &opensearchservice.ListDomainNamesInput{}

	op, err := svc.ListDomainNames(params)
	if err != nil {
		return nil, err
	}

	for _, domainname := range op.DomainNames {
		d.StreamListItem(ctx, &opensearchservice.DomainStatus{
			DomainName: domainname.DomainName,
		})

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenSearchDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getOpenSearchDomain")

	var domainName string
	if h.Item != nil {
		domainName = *h.Item.(*opensearchservice.DomainStatus).DomainName
	} else {
		domainName = d.KeyColumnQuals["domain_name"].GetStringValue()

		// Validate user input
		if len(domainName) < 1 {
			return nil, nil
		}
	}

	// Create Session
	svc, err := OpenSearchService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &opensearchservice.DescribeDomainInput{
		DomainName: &domainName,
	}

	// Get call
	data, err := svc.DescribeDomain(params)
	if err != nil {
		logger.Error("getOpenSearchDomain", "ERROR", err)
		return nil, err
	}

	return data.DomainStatus, nil
}

func listOpenSearchDomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listOpenSearchDomainTags")

	// Domain will be nil if getOpenSearchDomain returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getOpenSearchDomain"] == nil {
		return nil, nil
	}

	arn := h.HydrateResults["getOpenSearchDomain"].(*opensearchservice.DomainStatus).ARN

	// Create Session
	svc, err := OpenSearchService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &opensearchservice.ListTagsInput{
		ARN: arn,
	}

	// Get call
	op, err := svc.ListTags(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func openSearchDomaintagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("openSearchDomaintagListToTurbotTags")
	tagList := d.HydrateItem.(*opensearchservice.ListTagsOutput)

	if tagList.TagList == nil {
		return nil, nil
	}
	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
