package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice/types"

	elasticsearchservicev1 "github.com/aws/aws-sdk-go/service/elasticsearchservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsElasticsearchDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticsearch_domain",
		Description: "AWS Elasticsearch Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("domain_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsElasticsearchDomain,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsElasticsearchDomains,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:    listAwsElasticsearchDomainTags,
				Depends: []plugin.HydrateFunc{getAwsElasticsearchDomain},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elasticsearchservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_id",
				Description: "The id of the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsElasticsearchDomain,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "elasticsearch_version",
				Description: "The version for the Elasticsearch domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "endpoint",
				Description: "The Elasticsearch domain endpoint that use to submit index and search requests.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "access_policies",
				Description: "IAM access policy as a JSON-formatted string.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "created",
				Description: "The domain creation status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "deleted",
				Description: "The domain deletion status.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "processing",
				Description: "The status of the Elasticsearch domain configuration.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "upgrade_processing",
				Description: "The status of an Elasticsearch domain version upgrade.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "enabled",
				Description: "Specifies the status of the NodeToNodeEncryptionOptions.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsElasticsearchDomain,
				Transform:   transform.FromField("NodeToNodeEncryptionOptions.Enabled"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
				Transform:   transform.FromField("AccessPolicies").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "ebs_options",
				Description: "Specifies whether EBS-based storage is enabled.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
				Transform:   transform.FromField("EBSOptions"),
			},
			{
				Name:        "advanced_options",
				Description: "Specifies the status of the AdvancedOptions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "advanced_security_options",
				Description: "Specifies The current status of the Elasticsearch domain's advanced security options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "auto_tune_options",
				Description: "The current status of the Elasticsearch domain's Auto-Tune options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "cognito_options",
				Description: "The CognitoOptions for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "domain_endpoint_options",
				Description: "The current status of the Elasticsearch domain's endpoint options.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "elasticsearch_cluster_config",
				Description: "The type and number of instances in the domain cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "encryption_at_rest_options",
				Description: "Specifies the status of the EncryptionAtRestOptions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "log_publishing_options",
				Description: "Log publishing options for the given domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "service_software_options",
				Description: "The current status of the Elasticsearch domain's service software.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "snapshot_options",
				Description: "Specifies the status of the SnapshotOptions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
			},
			{
				Name:        "vpc_options",
				Description: "The VPCOptions for the specified domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
				Transform:   transform.FromField("VPCOptions"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsElasticsearchDomainTags,
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
				Hydrate:     listAwsElasticsearchDomainTags,
				Transform:   transform.FromField("TagList").Transform(getAwsElasticsearchDomaintagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticsearchDomain,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsElasticsearchDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ElasticsearchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticsearch_domain.listAwsElasticsearchDomain", "connection_error", err)
		return nil, err
	}

	// List call
	params := &elasticsearchservice.ListDomainNamesInput{}

	// API doesn't support pagination as of date
	op, err := svc.ListDomainNames(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticsearch_domain.listAwsElasticsearchDomain", "api_error", err)
		return nil, err
	}

	for _, domainname := range op.DomainNames {
		d.StreamListItem(ctx, types.ElasticsearchDomainStatus{
			DomainName: domainname.DomainName,
		})

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsElasticsearchDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var domainname string
	if h.Item != nil {
		domainname = *h.Item.(types.ElasticsearchDomainStatus).DomainName
	} else {
		domainname = d.EqualsQuals["domain_name"].GetStringValue()
	}

	// Create Session
	svc, err := ElasticsearchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticsearch_domain.getAwsElasticsearchDomain", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &elasticsearchservice.DescribeElasticsearchDomainInput{
		DomainName: &domainname,
	}

	// Get call
	data, err := svc.DescribeElasticsearchDomain(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticsearch_domain.getAwsElasticsearchDomain", "api_error", err)
		return nil, err
	}

	return data.DomainStatus, nil
}

func listAwsElasticsearchDomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Domain will be nil if getAwsElasticsearchDomain returned an error but
	// was ignored through ignore_error_codes config arg
	if h.HydrateResults["getAwsElasticsearchDomain"] == nil {
		return nil, nil
	}

	arn := h.HydrateResults["getAwsElasticsearchDomain"].(*types.ElasticsearchDomainStatus).ARN

	// Create Session
	svc, err := ElasticsearchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticsearch_domain.listAwsElasticsearchDomainTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &elasticsearchservice.ListTagsInput{
		ARN: arn,
	}

	// Get call
	op, err := svc.ListTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticsearch_domain.listAwsElasticsearchDomainTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func getAwsElasticsearchDomaintagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(*elasticsearchservice.ListTagsOutput)

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
