package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAwsPricingProductJson(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pricing_product_json",
		Description: "AWS Pricing Product",
		List: &plugin.ListConfig{
			Hydrate: listPricingProductJson,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "service_code",
					Require:    plugin.Required,
					CacheMatch: "exact",
				},
				{
					Name:       "attributes",
					Require:    plugin.Required,
					CacheMatch: "exact",
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
				Name:        "attributes",
				Description: "The attributes contains attribute name and value details.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("attributes"),
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

type PricingJsonInfo struct {
	ServiceCode     string
	Terms           interface{}
	Product         interface{}
	Version         interface{}
	PublicationDate interface{}
}

//// LIST FUNCTION

func listPricingProductJson(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	serviceCode := d.KeyColumnQuals["service_code"].GetStringValue()
	attributesString := d.KeyColumnQuals["attributes"].GetJsonbValue()
	attributes := make(map[string]string)

	// check if serviceCode is empty
	if serviceCode == "" {
		return nil, nil
	}

	err := json.Unmarshal([]byte(attributesString), &attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributesString: %v", err)
	}

	filter := setFilters(ctx, attributes)

	// Create Session
	svc, err := PricingService(ctx, d)
	if err != nil {
		return nil, err
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
				d.StreamListItem(ctx, PricingJsonInfo{serviceCode, priceList["terms"], priceList["product"], priceList["version"], priceList["publicationDate"]})

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

//// UTILITY FUNCTION

func setFilters(ctx context.Context, attributes map[string]string) []*pricing.Filter {
	filters := []*pricing.Filter{}

	for k, v := range attributes {
		filters = append(filters, &pricing.Filter{
			Type:  aws.String("TERM_MATCH"),
			Field: aws.String(k),
			Value: aws.String(v),
		})
	}

	return filters
}
