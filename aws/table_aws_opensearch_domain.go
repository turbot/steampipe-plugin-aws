package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/opensearchservice"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsOpenSearchDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_opensearch_domain",
		Description: "AWS OpenSearch Domain",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("domain_name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getAwsOpenSearchDomain,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsOpenSearchDomains,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
			{
				Func:    listAwsOpenSearchDomainTags,
				Depends: []plugin.HydrateFunc{getAwsOpenSearchDomain},
			},
		},
		GetMatrixItem: BuildRegionList,
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
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsOpenSearchDomain,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "access_policies",
				Description: "The IAM access policies of the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "created",
				Description: "The domain creation status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "deleted",
				Description: "The domain deletion status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "endpoint",
				Description: "The domain endpoint that is used to submit index and search requests.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "engine_version",
				Description: "The domain's OpenSearch version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "processing",
				Description: "The status of the domain configuration.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "upgrade_processing",
				Description: "The status of the domain version upgrade.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "node_to_node_encryption_enabled",
				Description: "Specifies the status of the node to node encryption status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsOpenSearchDomain,
				Transform:   transform.FromField("NodeToNodeEncryptionOptions.Enabled"),
			},
			{
				Name:        "advanced_options",
				Description: "Specifies the status of the advanced options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "advanced_security_options",
				Description: "Specifies The current status of the OpenSearch domain's advanced security options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "auto_tune_options",
				Description: "The current status of the domain's auto-tune options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "cluster_config",
				Description: "The type and number of instances in the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "cognito_options",
				Description: "The cognito options for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "domain_endpoint_options",
				Description: "The current status of the domain's endpoint options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "ebs_options",
				Description: "The EBSOptions for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
				Transform:   transform.FromField("EBSOptions"),
			},
			{
				Name:        "encryption_at_rest_options",
				Description: "The status of the encryption at rest options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "endpoints",
				Description: "Map containing the domain endpoints used to submit index and search requests.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "log_publishing_options",
				Description: "Log publishing options for the given domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "service_software_options",
				Description: "The current status of the domain's service software.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "snapshot_options",
				Description: "Specifies the status of the snapshot options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
			},
			{
				Name:        "vpc_options",
				Description: "The vpc options for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
				Transform:   transform.FromField("VPCOptions"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsOpenSearchDomainTags,
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
				Hydrate:     listAwsOpenSearchDomainTags,
				Transform:   transform.FromField("TagList").Transform(getAwsOpenSearchDomaintagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOpenSearchDomain,
				Transform:   transform.FromField("ARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsOpenSearchDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

func getAwsOpenSearchDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsOpenSearchDomain")

	var domainname string
	if h.Item != nil {
		domainname = *h.Item.(*opensearchservice.DomainStatus).DomainName
	} else {
		domainname = d.KeyColumnQuals["domain_name"].GetStringValue()
	}

	// Create Session
	svc, err := OpenSearchService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &opensearchservice.DescribeDomainInput{
		DomainName: &domainname,
	}

	// Get call
	data, err := svc.DescribeDomain(params)
	if err != nil {
		logger.Error("getAwsOpenSearchDomain", "ERROR", err)
		return nil, err
	}

	return data.DomainStatus, nil
}

func listAwsOpenSearchDomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsOpenSearchDomainTags")

	arn := h.HydrateResults["getAwsOpenSearchDomain"].(*opensearchservice.DomainStatus).ARN

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

func getAwsOpenSearchDomaintagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsOpenSearchDomaintagListToTurbotTags")
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
