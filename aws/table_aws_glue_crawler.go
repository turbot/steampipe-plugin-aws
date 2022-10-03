package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueCrawler(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_crawler",
		Description: "AWS Glue Crawler",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException", "InvalidParameter"}),
			},
			Hydrate: getGlueCrawler,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueCrawlers,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := GlueService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &glue.GetCrawlersInput{
		MaxResults: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.GetCrawlersPages(
		input,
		func(page *glue.GetCrawlersOutput, isLast bool) bool {
			for _, crawler := range page.Crawlers {
				d.StreamListItem(ctx, crawler)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listGlueCrawlers", "list", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueCrawler(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &glue.GetCrawlerInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetCrawler(params)
	if err != nil {
		plugin.Logger(ctx).Error("getGlueCrawler", "get", err)
		return nil, err
	}

	return data.Crawler, nil
}

func getGlueCrawlerArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGlueCrawlerArn")
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*glue.Crawler)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn format - https://docs.aws.amazon.com/glue/latest/dg/glue-specifying-resource-arns.html
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":crawler/" + *data.Name

	return arn, nil
}
