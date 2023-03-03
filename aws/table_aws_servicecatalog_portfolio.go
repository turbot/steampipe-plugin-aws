package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"

	servicecatalogv1 "github.com/aws/aws-sdk-go/service/servicecatalog"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServicecatalogPortfolio(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicecatalog_portfolio",
		Description: "AWS Service Catalog Portfolio",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getServiceCatalogPortfolio,
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceCatalogPortfolios,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicecatalogv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "The name to use for display purposes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PortfolioDetail.DisplayName"),
			},
			{
				Name:        "id",
				Description: "The portfolio identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PortfolioDetail.Id"),
			},
			{
				Name:        "arn",
				Description: "The ARN assigned to the portfolio.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PortfolioDetail.ARN"),
			},
			{
				Name:        "created_time",
				Description: "The UTC timestamp of the creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("PortfolioDetail.CreatedTime"),
			},
			{
				Name:        "description",
				Description: "The description of the portfolio.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PortfolioDetail.Description"),
			},
			{
				Name:        "provider_name",
				Description: "The name of the portfolio provider.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PortfolioDetail.ProviderName"),
			},
			{
				Name:        "budgets",
				Description: "Information about the associated budgets.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceCatalogPortfolio,
				Transform:   transform.FromField("Budgets"),
			},
			{
				Name:        "tag_options",
				Description: "Information about the tag options associated with the portfolio.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceCatalogPortfolio,
				Transform:   transform.FromField("TagOptions"),
			},
			{
				Name:        "tags_src",
				Description: "Information about the tags associated with the portfolio.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceCatalogPortfolio,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PortfolioDetail.DisplayName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceCatalogPortfolio,
				Transform:   transform.From(portfolioTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PortfolioDetail.ARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceCatalogPortfolios(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Client
	svc, err := ServiceCatalogClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_portfolio.listServiceCatalogPortfolios", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limiting the results
	// Document(https://docs.aws.amazon.com/servicecatalog/latest/dg/API_ListPortfolios.html) says the PageSize value can be between 1-100.
	// But the API throws error: api error ValidationException: 1 validation error detected: Value '100' at 'pageSize' failed to satisfy constraint: Member must have value less than or equal to 20 (SQLSTATE HV000)
	maxLimit := int32(20)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &servicecatalog.ListPortfoliosInput{
		PageSize: maxLimit,
	}

	paginator := servicecatalog.NewListPortfoliosPaginator(svc, input, func(o *servicecatalog.ListPortfoliosPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicecatalog_portfolio.listServiceCatalogPortfolios", "api_error", err)
			return nil, err
		}

		for _, item := range output.PortfolioDetails {
			d.StreamListItem(ctx, &servicecatalog.DescribePortfolioOutput{
				PortfolioDetail: &item,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceCatalogPortfolio(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		data := h.Item.(*servicecatalog.DescribePortfolioOutput)
		id = *data.PortfolioDetail.Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}

	if id == "" {
		return nil, nil
	}

	// Create client
	svc, err := ServiceCatalogClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_portfolio.getServiceCatalogPortfolio", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &servicecatalog.DescribePortfolioInput{
		Id: aws.String(id),
	}

	// Get call
	op, err := svc.DescribePortfolio(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_portfolio.getServiceCatalogPortfolio", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func portfolioTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*servicecatalog.DescribePortfolioOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
