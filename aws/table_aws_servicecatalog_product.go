package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"

	servicecatalogv1 "github.com/aws/aws-sdk-go/service/servicecatalog"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServicecatalogProduct(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicecatalog_product",
		Description: "AWS Service Catalog Product",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("product_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getServiceCatalogProduct,
			Tags:    map[string]string{"service": "servicecatalog", "action": "DescribeProduct"},
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceCatalogProducts,
			Tags:    map[string]string{"service": "servicecatalog", "action": "SearchProducts"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "accept_language",
					Require: plugin.Optional,
				},
				{
					Name:    "full_text_search",
					Require: plugin.Optional,
				},
				{
					Name:    "owner",
					Require: plugin.Optional,
				},
				{
					Name:    "type",
					Require: plugin.Optional,
				},
				{
					Name:    "source_product_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicecatalogv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.Name"),
			},
			{
				Name:        "id",
				Description: "The product view identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.Id"),
			},
			{
				Name:        "product_id",
				Description: "The product identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.ProductId"),
			},
			{
				Name:        "source_product_id",
				Description: "The source product identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("source_product_id"),
			},
			{
				Name:        "distributor",
				Description: "The distributor of the product. Contact the product administrator for the significance of this value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.Distributor"),
			},
			{
				Name:        "accept_language",
				Description: "The language code.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("accept_language"),
			},
			{
				Name:        "full_text_search",
				Description: "The full text for the product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("full_text_search"),
			},
			{
				Name:        "has_default_path",
				Description: "Indicates whether the product has a default path. If the product does not have a default path, call ListLaunchPaths to disambiguate between paths.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ProductViewSummary.HasDefaultPath"),
			},
			{
				Name:        "owner",
				Description: "The owner of the product. Contact the product administrator for the significance of this value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.Owner"),
			},
			{
				Name:        "short_description",
				Description: "Short description of the product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.ShortDescription"),
			},
			{
				Name:        "support_description",
				Description: "The description of the support for this product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.SupportDescription"),
			},
			{
				Name:        "support_email",
				Description: "The email contact information to obtain support for this product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.SupportEmail"),
			},
			{
				Name:        "support_url",
				Description: "The URL information to obtain support for this product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.SupportUrl"),
			},
			{
				Name:        "type",
				Description: "The product type. Contact the product administrator for the significance of this value. If this value is MARKETPLACE, the product was created by Amazon Web Services Marketplace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.Type"),
			},
			{
				Name:        "budgets",
				Description: "Information about the associated budgets.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceCatalogProduct,
				Transform:   transform.FromField("Budgets"),
			},
			{
				Name:        "launch_paths",
				Description: "Information about the associated launch paths.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceCatalogProduct,
				Transform:   transform.FromField("LaunchPaths"),
			},
			{
				Name:        "provisioning_artifacts",
				Description: "Information about the provisioning artifacts for the specified product.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceCatalogProduct,
				Transform:   transform.FromField("ProvisioningArtifacts"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductViewSummary.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProductArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceCatalogProducts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Client
	svc, err := ServiceCatalogClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_product.listServiceCatalogProducts", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &servicecatalog.SearchProductsInput{
		PageSize: maxLimit,
	}

	if d.EqualsQualString("accept_language") != "" {
		input.AcceptLanguage = aws.String(d.EqualsQualString("accept_language"))
	}

	filters := buildServiceCatalogProductFilter(ctx, d.Quals)
	input.Filters = filters

	paginator := servicecatalog.NewSearchProductsPaginator(svc, input, func(o *servicecatalog.SearchProductsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicecatalog_product.listServiceCatalogProducts", "api_error", err)
			return nil, err
		}

		for _, item := range output.ProductViewSummaries {
			d.StreamListItem(ctx, &servicecatalog.DescribeProductOutput{
				ProductViewSummary: &types.ProductViewSummary{
					Distributor:        item.Distributor,
					HasDefaultPath:     item.HasDefaultPath,
					Id:                 item.Id,
					Name:               item.Name,
					Owner:              item.Owner,
					ProductId:          item.ProductId,
					ShortDescription:   item.ShortDescription,
					SupportDescription: item.SupportDescription,
					SupportEmail:       item.SupportEmail,
					SupportUrl:         item.SupportUrl,
					Type:               item.Type,
				},
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

func getServiceCatalogProduct(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		data := h.Item.(*servicecatalog.DescribeProductOutput)
		id = *data.ProductViewSummary.ProductId
	} else {
		id = d.EqualsQuals["product_id"].GetStringValue()
	}

	if id == "" {
		return nil, nil
	}

	// Create client
	svc, err := ServiceCatalogClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_product.getServiceCatalogProduct", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &servicecatalog.DescribeProductInput{
		Id: aws.String(id),
	}

	// Get call
	op, err := svc.DescribeProduct(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_product.getServiceCatalogProduct", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getProductArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	product := h.Item.(*servicecatalog.DescribeProductOutput)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":catalog:" + region + ":" + commonColumnData.AccountId + ":product/" + *product.ProductViewSummary.ProductId
	return arn, nil
}

//// UTILITY FUNCTIONS

// Buid servicecatalog product list call filter param

func buildServiceCatalogProductFilter(ctx context.Context, quals plugin.KeyColumnQualMap) map[string][]string {
	filterQuals := map[string]string{
		"full_text_search":  "FullTextSearch",
		"owner":             "Owner",
		"type":              "ProductType",
		"source_product_id": "SourceProductId",
	}

	filter := make(map[string][]string)
	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter[filterName] = []string{val}
			}
		}
	}

	return filter
}
