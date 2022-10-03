package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsPricingServiceAttribute(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pricing_service_attribute",
		Description: "AWS Pricing Service Attribute",
		List: &plugin.ListConfig{
			Hydrate: listPricingServiceAttributes,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_code", Require: plugin.Optional},
			},
		},
		Columns: awsDefaultColumns([]*plugin.Column{
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
	AttributeName *string
	ServiceCode   *string
}

//// LIST FUNCTION

func listPricingServiceAttributes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := PricingService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &pricing.DescribeServicesInput{
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(100),
	}

	equalQual := d.KeyColumnQuals
	if equalQual["service_code"] != nil {
		input.ServiceCode = aws.String(equalQual["service_code"].GetStringValue())
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
	err = svc.DescribeServicesPages(
		input,
		func(page *pricing.DescribeServicesOutput, isLast bool) bool {
			for _, service := range page.Services {
				for _, attributeName := range service.AttributeNames {
					d.StreamListItem(ctx, ServiceDetail{attributeName, service.ServiceCode})
				}
				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listPricingServiceAttributes", "err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func listAttributeValues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	serviceCode := h.Item.(ServiceDetail).ServiceCode
	attributeName := h.Item.(ServiceDetail).AttributeName

	// Create Session
	svc, err := PricingService(ctx, d)
	if err != nil {
		return nil, err
	}
	input := &pricing.GetAttributeValuesInput{
		AttributeName: attributeName,
		ServiceCode:   serviceCode,
		MaxResults:    aws.Int64(100),
	}

	attributeValues := []string{}

	// List call
	err = svc.GetAttributeValuesPages(
		input,
		func(page *pricing.GetAttributeValuesOutput, isLast bool) bool {
			for _, attributeValue := range page.AttributeValues {
				attributeValues = append(attributeValues, *attributeValue.Value)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listAttributeValues", "err", err)
		return nil, err
	}

	return attributeValues, nil
}
