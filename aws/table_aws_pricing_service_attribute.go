package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pricing"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsPricingServiceAttribute(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pricing_service_attribute",
		Description: "AWS Pricing Service Attribute",
		List: &plugin.ListConfig{
			Hydrate: listPricingServiceAttributes,
			Tags:    map[string]string{"service": "pricing", "action": "DescribeServices"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listAttributeValues,
				Tags: map[string]string{"service": "pricing", "action": "GetAttributeValues"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "service_code",
				Description: "The service code of the AWS service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attribute_name",
				Description: "The supported attribute names for the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attribute_values",
				Description: "The supported attribute values for the service and attribute name.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAttributeValues,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type ServiceDetail struct {
	AttributeName string
	ServiceCode   *string
}

//// LIST FUNCTION

func listPricingServiceAttributes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := PricingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pricing_service_attribute.listPricingServiceAttributes", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &pricing.DescribeServicesInput{
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int32(maxLimit),
	}

	equalQual := d.EqualsQuals
	if equalQual["service_code"] != nil {
		input.ServiceCode = aws.String(equalQual["service_code"].GetStringValue())
	}

	paginator := pricing.NewDescribeServicesPaginator(svc, input, func(o *pricing.DescribeServicesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_pricing_service_attribute.listPricingServiceAttributes", "api_error", err)
			return nil, err
		}

		for _, items := range output.Services {
			for _, attributeName := range items.AttributeNames {
				d.StreamListItem(ctx, ServiceDetail{attributeName, items.ServiceCode})
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func listAttributeValues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	serviceCode := h.Item.(ServiceDetail).ServiceCode
	attributeName := h.Item.(ServiceDetail).AttributeName

	// Create Session
	svc, err := PricingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pricing_service_attribute.listAttributeValues", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &pricing.GetAttributeValuesInput{
		AttributeName: &attributeName,
		ServiceCode:   serviceCode,
		MaxResults:    aws.Int32(maxLimit),
	}

	attributeValues := []string{}

	paginator := pricing.NewGetAttributeValuesPaginator(svc, input, func(o *pricing.GetAttributeValuesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_pricing_service_attribute.listAttributeValues", "api_error", err)
			return nil, err
		}

		for _, items := range output.AttributeValues {
			attributeValues = append(attributeValues, *items.Value)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return attributeValues, nil
}
