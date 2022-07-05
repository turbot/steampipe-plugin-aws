package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableAwsPricingProduct(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pricing_product",
		Description: "AWS Pricing Product",
		List: &plugin.ListConfig{
			Hydrate: listPricingProducts,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "service_code",
					Require: plugin.Required,
				},
				{
					Name:    "field",
					Require: plugin.Required,
				},
				{
					Name:    "value",
					Require: plugin.Required,
				},
			},
		},
		Columns: awsS3Columns([]*plugin.Column{
			{
				Name:        "service_code",
				Description: "The service code of the AWS service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "field",
				Description: "The product metadata field that you want to filter on. You can filter by just the service code to see all products for a specific service, filter by just the attribute name to see a specific attribute for multiple services, or use both a service code and an attribute name to retrieve only products that match both fields.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The service code or attribute value that you want to filter by.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "publication_date",
				Description: "The publication date.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version",
				Description: "The version information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product",
				Description: "The product details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "terms",
				Description: "The terms details for the product.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

type PricingInfo struct {
	ServiceCode     string
	Field           string
	Value           string
	Terms           interface{}
	Product         interface{}
	Version         interface{}
	PublicationDate interface{}
}

//// LIST FUNCTION

func listPricingProducts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	serviceCode := d.KeyColumnQuals["service_code"].GetStringValue()
	field := d.KeyColumnQuals["field"].GetStringValue()
	value := d.KeyColumnQuals["value"].GetStringValue()

	// check if serviceCode or field or value is empty
	if serviceCode == "" || field == "" || value == "" {
		return nil, nil
	}

	// Create Session
	svc, err := PricingService(ctx, d)
	if err != nil {
		return nil, err
	}

	filter := []*pricing.Filter{
		{
			Type:  aws.String("TERM_MATCH"),
			Field: &field,
			Value: &value,
		},
	}

	input := &pricing.GetProductsInput{
		Filters:       filter,
		FormatVersion: aws.String("aws_v1"),
		ServiceCode:   &serviceCode,
		MaxResults:    aws.Int64(100),
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = types.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.GetProductsPages(
		input,
		func(pages *pricing.GetProductsOutput, isLast bool) bool {
			for _, priceList := range pages.PriceList {
				d.StreamListItem(ctx, PricingInfo{serviceCode, field, value, priceList["terms"], priceList["product"], priceList["version"], priceList["publicationDate"]})

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listPricingProducts", "list", err)
	}

	return nil, err
}
