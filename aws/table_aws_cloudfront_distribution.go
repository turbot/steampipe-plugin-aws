package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFrontDistribution(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudfront_distribution",
		Description: "AWS CloudFront Distribution",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"NoSuchDistribution"}),
			},
			Hydrate: getCloudFrontDistribution,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsCloudFrontDistributions,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The identifier for the Distribution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id", "Distribution.Id"),
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) for the distribution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN", "Distribution.ARN"),
			},
			{
				Name:        "status",
				Description: "The current status of the Distribution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status", "Distribution.Status"),
			},
			{
				Name:        "caller_reference",
				Description: "A unique value that ensures that the request can't be replayed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFrontDistribution,
				Transform:   transform.FromField("Distribution.DistributionConfig.CallerReference"),
			},
			{
				Name:        "comment",
				Description: "The comment originally specified when this Distribution was created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Comment", "Distribution.DistributionConfig.Comment"),
			},
			{
				Name:        "default_root_object",
				Description: "The object that you want CloudFront to request from your origin.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFrontDistribution,
				Transform:   transform.FromField("Distribution.DistributionConfig.DefaultRootObject"),
			},
			{
				Name:        "domain_name",
				Description: "The domain name that corresponds to the Distribution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName", "Distribution.DomainName"),
			},
			{
				Name:        "enabled",
				Description: "Whether the Distribution is enabled to accept user requests for content.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Enabled", "Distribution.DistributionConfig.Enabled"),
			},
			{
				Name:        "e_tag",
				Description: "The current version of the configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFrontDistribution,
			},
			{
				Name:        "http_version",
				Description: "Specify the maximum HTTP version that you want viewers to use to communicate with CloudFront. The default value for new web Distributions is http2. Viewers that don't support HTTP/2 will automatically use an earlier version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("HttpVersion", "Distribution.DistributionConfig.HttpVersion"),
			},
			{
				Name:        "is_ipv6_enabled",
				Description: "Whether CloudFront responds to IPv6 DNS requests with an IPv6 address for your Distribution.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsIPV6Enabled", "Distribution.DistributionConfig.IsIPV6Enabled"),
			},
			{
				Name:        "in_progress_invalidation_batches",
				Description: "The number of invalidation batches currently in progress.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudFrontDistribution,
				Transform:   transform.FromField("Distribution.InProgressInvalidationBatches"),
			},
			{
				Name:        "last_modified_time",
				Description: "The date and time the Distribution was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastModifiedTime", "Distribution.DistributionConfig.LastModifiedTime"),
			},
			{
				Name:        "price_class",
				Description: "A complex type that contains information about price class for this streaming Distribution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PriceClass", "Distribution.DistributionConfig.PriceClass"),
			},
			{
				Name:        "web_acl_id",
				Description: "The Web ACL Id (if any) associated with the distribution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebACLId", "Distribution.DistributionConfig.WebACLId"),
			},
			{
				Name:        "active_trusted_key_groups",
				Description: "CloudFront automatically adds this field to the response if you’ve configured a cache behavior in this distribution to serve private content using key groups.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFrontDistribution,
				Transform:   transform.FromField("Distribution.ActiveTrustedKeyGroups"),
			},
			{
				Name:        "active_trusted_signers",
				Description: "A list of AWS accounts and the identifiers of active CloudFront key pairs in each account that CloudFront can use to verify the signatures of signed URLs and signed cookies.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFrontDistribution,
				Transform:   transform.FromField("Distribution.ActiveTrustedSigners"),
			},
			{
				Name:        "aliases",
				Description: "A complex type that contains information about CNAMEs (alternate domain names),if any, for this distribution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Aliases", "Distribution.DistributionConfig.Aliases"),
			},
			{
				Name:        "alias_icp_recordals",
				Description: "AWS services in China customers must file for an Internet Content Provider (ICP) recordal if they want to serve content publicly on an alternate domain name, also known as a CNAME, that they've added to CloudFront. AliasICPRecordal provides the ICP recordal status for CNAMEs associated with distributions.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AliasICPRecordals", "Distribution.AliasICPRecordals"),
			},
			{
				Name:        "cache_behaviors",
				Description: "The number of cache behaviors for this Distribution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CacheBehaviors", "Distribution.DistributionConfig.CacheBehaviors"),
			},
			{
				Name:        "custom_error_responses",
				Description: "A complex type that contains zero or more CustomErrorResponses elements.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CustomErrorResponses", "Distribution.DistributionConfig.CustomErrorResponses"),
			},
			{
				Name:        "default_cache_behavior",
				Description: "A complex type that describes the default cache behavior if you don't specify a CacheBehavior element or if files don't match any of the values of PathPattern in CacheBehavior elements. You must create exactly one default cache behavior.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefaultCacheBehavior", "Distribution.DistributionConfig.DefaultCacheBehavior"),
			},
			{
				Name:        "logging",
				Description: "A complex type that controls whether access logs are written for the distribution.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFrontDistribution,
				Transform:   transform.FromField("Distribution.DistributionConfig.Logging"),
			},
			{
				Name:        "origins",
				Description: "A complex type that contains information about origins for this distribution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Origins.Items", "Distribution.DistributionConfig.Origins.Items"),
			},
			{
				Name:        "origin_groups",
				Description: "A complex type that contains information about origin groups for this distribution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OriginGroups", "Distribution.DistributionConfig.OriginGroups"),
			},
			{
				Name:        "restrictions",
				Description: "A complex type that identifies ways in which you want to restrict distribution of your content.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Restrictions", "Distribution.DistributionConfig.Restrictions"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Maintenance Window",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFrontDistributionTags,
				Transform:   transform.FromField("Tags.Items").Transform(handleEmptySliceAndMap),
			},
			{
				Name:        "viewer_certificate",
				Description: "A complex type that determines the distribution’s SSL/TLS configuration for communicating with viewers.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ViewerCertificate", "Distribution.DistributionConfig.ViewerCertificate"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFrontDistributionTags,
				Transform:   transform.FromField("Tags.Items").Transform(cloudFrontDistributionTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN", "Distribution.ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsCloudFrontDistributions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_distribution.listAwsCloudFrontDistributions", "client_error", err)
		return nil, err
	}

	maxItems := int32(1000)

	// The maximum number for MaxItems parameter is not defined by the API
	// We have set the MaxItems to 1000 based on our test
	input := &cloudfront.ListDistributionsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxItems = &maxItems
	paginator := cloudfront.NewListDistributionsPaginator(svc, input, func(o *cloudfront.ListDistributionsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudfront_distribution.listAwsCloudFrontDistributions", "api_error", err)
			return nil, err
		}

		for _, distribution := range output.DistributionList.Items {
			d.StreamListItem(ctx, distribution)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudFrontDistribution(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_distribution.getCloudFrontDistribution", "client_error", err)
		return nil, err
	}

	var distributionID string
	if h.Item != nil {
		distributionID = *h.Item.(types.DistributionSummary).Id
	} else {
		distributionID = d.KeyColumnQuals["id"].GetStringValue()
	}

	if strings.TrimSpace(distributionID) == "" {
		return nil, nil
	}

	params := &cloudfront.GetDistributionInput{
		Id: &distributionID,
	}

	op, err := svc.GetDistribution(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_distribution.getCloudFrontDistribution", "api_error", err)
		return nil, err
	}

	return *op, nil
}

func getCloudFrontDistributionTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudFrontClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_distribution.getCloudFrontDistributionTags", "client_error", err)
		return nil, err
	}

	distributionAka := cloudFrontDistributionAka(h.Item)

	// Build the params
	params := &cloudfront.ListTagsForResourceInput{
		Resource: distributionAka,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudfront_distribution.getCloudFrontDistributionTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func cloudFrontDistributionTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(tagList) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	} else {
		return nil, nil
	}

	return turbotTagsMap, nil
}

func cloudFrontDistributionAka(item interface{}) *string {
	switch item := item.(type) {
	case cloudfront.GetDistributionOutput:
		return item.Distribution.ARN
	case types.DistributionSummary:
		return item.ARN
	}
	return nil
}
