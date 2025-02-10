package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	gluev1 "github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueCrawler(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_crawler",
		Description: "AWS Glue Crawler",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException", "InvalidParameter"}),
			},
			Hydrate: getGlueCrawler,
			Tags:    map[string]string{"service": "glue", "action": "GetCrawler"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueCrawlers,
			Tags:    map[string]string{"service": "glue", "action": "GetCrawlers"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getGlueCrawler,
				Tags: map[string]string{"service": "glue", "action": "GetCrawler"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(gluev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the crawler.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the crawler.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueCrawlerArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "database_name",
				Description: "The name of the database in which the crawler's output is stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Indicates whether the crawler is running or pending.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The Amazon Resource Name (ARN) of an IAM role that's used to access customer resources, such as Amazon Simple Storage Service (Amazon S3) data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time that the crawler was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the crawler.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "crawl_elapsed_time",
				Description: "If the crawler is running, contains the total time elapsed since the last crawl began.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "crawler_lineage_settings",
				Description: "Specifies whether data lineage is enabled for the crawler.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LineageConfiguration.CrawlerLineageSettings"),
			},
			{
				Name:        "crawler_security_configuration",
				Description: "The name of the SecurityConfiguration structure to be used by this crawler.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated",
				Description: "The time that the crawler was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "recrawl_behavior",
				Description: "Specifies whether to crawl the entire dataset again or to crawl only folders that were added since the last crawler run. A value of CRAWL_EVERYTHING specifies crawling the entire dataset again. A value of CRAWL_NEW_FOLDERS_ONLY specifies crawling only folders that were added since the last crawler run. A value of CRAWL_EVENT_MODE specifies crawling only the changes identified by Amazon S3 events.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecrawlPolicy.RecrawlBehavior"),
			},
			{
				Name:        "table_prefix",
				Description: "The prefix added to the names of tables that are created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The version of the crawler.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "classifiers",
				Description: "A list of UTF-8 strings that specify the custom classifiers that are associated with the crawler.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lake_formation_configuration",
				Description: "Specifies whether the crawler should use Lake Formation credentials for the crawler instead of the IAM role credentials.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "configuration",
				Description: "Crawler configuration information.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_crawl",
				Description: "The status of the last crawl, and potentially error information if an error occurred.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "schedule",
				Description: "For scheduled crawlers, the schedule when the crawler runs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "schema_change_policy",
				Description: "The policy that specifies update and delete behaviors for the crawler.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "targets",
				Description: "A collection of targets to crawl.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
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
				Hydrate:     getTagsForGlueCrawler,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueCrawlerArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueCrawlers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_crawler.listGlueCrawlers", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &glue.GetCrawlersInput{}

	// List call

	paginator := glue.NewGetCrawlersPaginator(svc, input, func(o *glue.GetCrawlersPaginatorOptions) {
		o.StopOnDuplicateToken = true
		o.Limit = maxLimit
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_crawler.listGlueCrawlers", "api_error", err)
			return nil, err
		}
		for _, crawler := range output.Crawlers {
			d.StreamListItem(ctx, crawler)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueCrawler(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_crawler.getGlueCrawler", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &glue.GetCrawlerInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetCrawler(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_crawler.getGlueCrawler", "api_error", err)
		return nil, err
	}

	return *data.Crawler, nil
}

func getTagsForGlueCrawler(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn, _ := getGlueCrawlerArn(ctx, d, h)
	return getTagsForGlueResource(ctx, d, arn.(string))
}

func getGlueCrawlerArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.Crawler)

	// Get common columns

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_crawler.getGlueCrawlerArn", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn format - https://docs.aws.amazon.com/glue/latest/dg/glue-specifying-resource-arns.html
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":crawler/" + *data.Name

	return arn, nil
}
