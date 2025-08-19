package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServicecatalogPortfolioShare(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicecatalog_portfolio_share",
		Description: "AWS Service Catalog Portfolio Share",
		List: &plugin.ListConfig{
			ParentHydrate: listServiceCatalogPortfolios,
			Hydrate:       listServiceCatalogPortfolioShares,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "type", Require: plugin.Required},
				{Name: "portfolio_id", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "servicecatalog", "action": "DescribePortfolioShares"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SERVICECATALOG_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "portfolio_id",
				Description: "The unique identifier of the portfolio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "portfolio_arn",
				Description: "The ARN of the portfolio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "portfolio_display_name",
				Description: "The display name of the portfolio.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the portfolio share.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PortfolioShareDetail.Type"),
			},
			{
				Name:        "principal_id",
				Description: "The identifier of the recipient entity that received the portfolio share.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "accepted",
				Description: "Indicates whether the shared portfolio is imported by the recipient account.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PortfolioShareDetail.Accepted"),
			},
			{
				Name:        "share_principals",
				Description: "Indicates if Principal sharing is enabled or disabled for the portfolio share.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PortfolioShareDetail.SharePrincipals"),
			},
			{
				Name:        "share_tag_options",
				Description: "Indicates whether TagOptions sharing is enabled or disabled for the portfolio share.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PortfolioShareDetail.ShareTagOptions"),
			},

			// Steampipe standard Columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getServiceCatalogPortfolioShareTitle),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceCatalogPortfolioShares(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get portfolio details from parent
	portfolio := h.Item.(*servicecatalog.DescribePortfolioOutput)
	portfolioId := *portfolio.PortfolioDetail.Id

	if d.EqualsQualString("portfolio_id") != "" && d.EqualsQualString("portfolio_id") != portfolioId {
		return nil, nil
	}

	// Create service
	svc, err := ServiceCatalogClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_portfolio_share.listServiceCatalogPortfolioShares", "connection_error", err)
		return nil, err
	}

	// Portfolio share types
	shareType := d.EqualsQualString("type")

	params := &servicecatalog.DescribePortfolioSharesInput{
		PortfolioId: aws.String(portfolioId),
		Type:        types.DescribePortfolioShareType(shareType),
	}

	paginator := servicecatalog.NewDescribePortfolioSharesPaginator(svc, params, func(o *servicecatalog.DescribePortfolioSharesPaginatorOptions) {
		o.Limit = 100
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicecatalog_portfolio_share.listServiceCatalogPortfolioShares", "api_error", err, "portfolio_id", portfolioId, "type", shareType)
			continue
		}

		for _, share := range output.PortfolioShareDetails {
			// Create a custom struct to include portfolio details
			portfolioShareData := &PortfolioShareData{
				PortfolioShareDetail: share,
				PortfolioId:          portfolioId,
				PortfolioArn:         *portfolio.PortfolioDetail.ARN,
				PortfolioDisplayName: *portfolio.PortfolioDetail.DisplayName,
			}

			d.StreamListItem(ctx, portfolioShareData)

			// Check if context is cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getServiceCatalogPortfolioShareTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*PortfolioShareData)

	title := data.PortfolioDisplayName + " - " + string(data.PortfolioShareDetail.Type) + " - " + *data.PortfolioShareDetail.PrincipalId
	return title, nil
}

//// STRUCTS

type PortfolioShareData struct {
	PortfolioShareDetail types.PortfolioShareDetail
	PortfolioId          string
	PortfolioArn         string
	PortfolioDisplayName string
}
