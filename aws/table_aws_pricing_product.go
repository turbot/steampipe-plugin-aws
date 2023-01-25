package aws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsPricingProduct(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pricing_product",
		Description: "AWS Pricing Product",
		List: &plugin.ListConfig{
			Hydrate: listPricingProduct,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Required},
				{Name: "filters", Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{Name: "rate_code", Description: "A unique code for a product/ offer/ pricing-tier combination.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.PriceDimension.RateCode")},
			{Name: "service_code", Description: "This identifies the specific AWS service to the customer as a unique short abbreviation.", Type: proto.ColumnType_STRING, Transform: transform.FromField("ServiceCode")},
			{Name: "term", Description: "Whether your AWS usage is Reserved or On-Demand.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.Term")},
			{Name: "purchase_option", Description: "How you chose to pay for this line item (All Upfront, Partial Upfront, No Upfront).", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.TermAttributes.PurchaseOption")},
			{Name: "lease_contract_length", Description: "The length of time that your RI is reserved for.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.TermAttributes.LeaseContractLength")},
			{Name: "description", Description: "Description for a product / offer / pricing-tier combination.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.PriceDimension.Description")},
			{Name: "begin_range", Description: "Start of billing range, by unit", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.PriceDimension.BeginRange")},
			{Name: "end_range", Description: "Enf of billing range, by unit", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.PriceDimension.EndRange")},
			{Name: "unit", Description: "The pricing unit that AWS used for calculating your usage cost (ex: hours)", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.PriceDimension.Unit")},
			{Name: "price_per_unit", Description: "Price by unit", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.PriceDimension.PricePerUnit").Transform(extractPricePerUnit)},
			{Name: "currency", Description: "Currency used for the price", Type: proto.ColumnType_STRING, Transform: transform.FromField("Offer.PriceDimension.PricePerUnit").Transform(extractCurrency)},
			{Name: "publication_date", Description: "The publication date of the offer.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("PublicationDate")},
			{Name: "effective_date", Description: "The effective date of the pricing details.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Offer.EffectiveDate")},
			{Name: "version", Description: "The publication version of the offer.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Version")},
			// product attributes
			{Name: "attributes", Description: "Product attributes.", Type: proto.ColumnType_JSON, Transform: transform.FromField("Product.Attributes")},
			// product attributes filters
			{Name: "filters", Description: "Product filtering by attribute.", Type: proto.ColumnType_JSON, Transform: transform.FromQual("filters")},
		}),
	}
}

//// LIST FUNCTION

type Product struct {
	_             struct{} `type:"structure"`
	ProductFamily *string
	Attributes    map[string]*string
}

type Term map[string]*Offer

type PriceDimension struct {
	Unit         *string
	BeginRange   *string
	EndRange     *string
	Description  *string
	RateCode     *string
	PricePerUnit map[string]*string
}

type TermAttributes struct {
	_                   struct{} `type:"structure"`
	PurchaseOption      *string
	LeaseContractLength *string
}

type Offer struct {
	_               struct{} `type:"structure"`
	PriceDimensions map[string]*PriceDimension
	EffectiveDate   *time.Time
	OfferTermCode   *string
	TermAttributes  *TermAttributes
}

type PriceList struct {
	_               struct{} `type:"structure"`
	Product         *Product
	ServiceCode     *string
	Terms           map[string]*Term
	Version         *string
	PublicationDate *time.Time
}

func listPricingProduct(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := PricingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pricing_product.listPricingProduct", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxItems := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 20 {
				maxItems = 20
			} else {
				maxItems = limit
			}
		}
	}

	equalQual := d.KeyColumnQuals
	input := pricing.GetProductsInput{
		ServiceCode:   aws.String(equalQual["service_code"].GetStringValue()),
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    maxItems,
	}

	type OfferOutput struct {
		_              struct{} `type:"structure"`
		PriceDimension *PriceDimension
		Term           string
		EffectiveDate  *time.Time
		OfferTermCode  *string
		TermAttributes *TermAttributes
	}

	type PriceOutput struct {
		_               struct{} `type:"structure"`
		Product         *Product
		ServiceCode     *string
		Offer           *OfferOutput
		Version         *string
		PublicationDate *time.Time
	}

	filters, err := buildPricingFilter(d.Quals["filters"])
	if err != nil {
		plugin.Logger(ctx).Error("aws_pricing_product.listPricingProduct", "filters_building_error", err)
		return nil, err
	}

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := pricing.NewGetProductsPaginator(svc, &input, func(o *pricing.GetProductsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_pricing_product.listPricingProduct", "api_error", err)
			return nil, err
		}

		for _, priceList := range output.PriceList {
			var priceListObject PriceList

			err = json.Unmarshal([]byte(priceList), &priceListObject)
			if err != nil {
				plugin.Logger(ctx).Error("aws_pricing_product.listPricingProduct", "unmarshal_err", err)
				return nil, err
			}

			for term, termDetail := range priceListObject.Terms {
				for _, offer := range *termDetail {
					for _, priceDimension := range offer.PriceDimensions {
						d.StreamListItem(ctx, PriceOutput{
							Product:         priceListObject.Product,
							ServiceCode:     priceListObject.ServiceCode,
							Version:         priceListObject.Version,
							PublicationDate: priceListObject.PublicationDate,
							Offer:           &OfferOutput{PriceDimension: priceDimension, Term: term, EffectiveDate: offer.EffectiveDate, OfferTermCode: offer.OfferTermCode, TermAttributes: offer.TermAttributes},
						})
					}
				}
			}

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func extractPricePerUnit(_ context.Context, d *transform.TransformData) (interface{}, error) {
	priceMap := d.Value.(map[string]*string)

	for _, units := range priceMap {
		return units, nil
	}
	return nil, nil
}

func extractCurrency(_ context.Context, d *transform.TransformData) (interface{}, error) {
	priceMap := d.Value.(map[string]*string)

	for currency := range priceMap {
		return currency, nil
	}
	return nil, nil
}

// build pricing list call input filter
func buildPricingFilter(qual *plugin.KeyColumnQuals) ([]types.Filter, error) {
	if qual == nil {
		return nil, nil
	}

	filters := make([]types.Filter, 0)
	for _, qual := range qual.Quals {
		qualFilter := make(map[string]*string)
		err := json.Unmarshal([]byte(qual.Value.GetJsonbValue()), &qualFilter)
		if err != nil {
			return nil, err
		}

		for attributeName, attributeValue := range qualFilter {
			filter := types.Filter{
				Field: aws.String(attributeName),
				Type:  types.FilterTypeTermMatch,
			}
			filter.Value = attributeValue
			filters = append(filters, filter)
		}
	}
	return filters, nil
}
