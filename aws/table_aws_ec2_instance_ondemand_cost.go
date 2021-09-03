package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsEc2InstanceCostHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_instance_on_demand_pricing",
		Description: "AWS EC2 Instance On-Demand Pricing",
		List: &plugin.ListConfig{
			Hydrate: ListEc2InstanceCostHourly,
		},
		// GetMatrixItem: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "service_code",
				Description: "The service code of the AWS service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.servicename"),
			},
			{
				Name:        "publication_date",
				Description: "The publication day.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "availabilityzone",
				Description: "The availability zone for the product is available.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.availabilityzone"),
			},
			{
				Name:        "capacity_status",
				Description: "The product capacity status.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.capacitystatus"),
			},
			{
				Name:        "clock_speed",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.clockSpeed"),
			},
			{
				Name:        "current_generation",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.currentGeneration"),
			},
			{
				Name:        "dedicated_ebs_throughput",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.dedicatedEbsThroughput"),
			},
			{
				Name:        "ecu",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.ecu"),
			},
			{
				Name:        "enhanced_networking_supported",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.enhancedNetworkingSupported"),
			},
			{
				Name:        "instance_family",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.instanceFamily"),
			},
			{
				Name:        "instance_type",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.instanceType"),
			},
			{
				Name:        "instance_sku",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.instancesku"),
			},
			{
				Name:        "intel_avx2_available",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.intelAvx2Available"),
			},
			{
				Name:        "intel_avx_available",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.intelAvxAvailable"),
			},
			{
				Name:        "intel_turbo_available",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.intelTurboAvailable"),
			},
			{
				Name:        "license_model",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.licenseModel"),
			},
			{
				Name:        "location",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.location"),
			},
			{
				Name:        "location_type",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.locationType"),
			},
			{
				Name:        "marketoption",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.marketoption"),
			},
			{
				Name:        "memory",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.memory"),
			},
			{
				Name:        "network_performance",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.networkPerformance"),
			},
			{
				Name:        "normalization_size_factor",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.normalizationSizeFactor"),
			},
			{
				Name:        "operating_system",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.operatingSystem"),
			},
			{
				Name:        "operation",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.operation"),
			},
			{
				Name:        "physical_processor",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.physicalProcessor"),
			},
			{
				Name:        "pre_installed_sw",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.preInstalledSw"),
			},
			{
				Name:        "processor_architecture",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.processorArchitecture"),
			},
			{
				Name:        "processor_features",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.processorFeatures"),
			},

			{
				Name:        "storage",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.storage"),
			},
			{
				Name:        "tenancy",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.tenancy"),
			},
			{
				Name:        "usage_type",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.attributes.usagetype"),
			},
			{
				Name:        "vcpu",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Product.attributes.vcpu"),
			},
			{
				Name:        "vpcnetworking_support",
				Description: "The name of the AWS service",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Product.attributes.vpcnetworkingsupport"),
			},
			{
				Name:        "product_family",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.productFamily"),
			},
			{
				Name:        "sku",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Product.sku"),
			},
			{
				Name:        "effective_date",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("OnDemandInfo.effectiveDate"),
			},
			{
				Name:        "offer_term_code",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OnDemandInfo.offerTermCode"),
			},
			{
				Name:        "description",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PriceDimensions.description"),
			},
			{
				Name:        "begin_range",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PriceDimensions.beginRange"),
			},
			{
				Name:        "end_range",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PriceDimensions.endRange"),
			},
			{
				Name:        "rate_code",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PriceDimensions.rateCode"),
			},
			{
				Name:        "unit",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PriceDimensions.unit"),
			},
			{
				Name:        "price_currency",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "price",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "applies_to",
				Description: "The name of the AWS service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PriceDimensions.appliesTo"),
			},
		}),
	}
}

type AttributeNameWithServiceCode struct {
	AttributeNames []*string
	ServiceCode    *string
}

type PricingInfo struct {
	ServiceCode     interface{}
	PriceCurrency   interface{}
	Price           interface{}
	PublicationDate interface{}
	Product         interface{}
	OnDemandInfo    interface{}
	PriceDimensions interface{}
}

//// LIST FUNCTION

func ListEc2InstanceCostHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// region := d.KeyColumnQualString(matrixKeyRegion)
	// logger := plugin.Logger(ctx)

	// Create Session
	svc, err := PricingService(ctx, d)
	if err != nil {
		return nil, err
	}

	var pricInfo PricingInfo
	err = svc.GetProductsPages(
		&pricing.GetProductsInput{
			FormatVersion: aws.String("aws_v1"),
			ServiceCode:   aws.String("AmazonEC2"),
		}, func(pages *pricing.GetProductsOutput, isLast bool) bool {
			for _, price := range pages.PriceList {
				terms := price["terms"]
				if terms != nil {
					if terms.(map[string]interface{})["OnDemand"] != nil {
						data := terms.(map[string]interface{})["OnDemand"].(map[string]interface{})
						for _, v := range data {
							pricInfo.OnDemandInfo = v
							break
						}
						if pricInfo.OnDemandInfo.(map[string]interface{})["priceDimensions"] != nil {
							pricingData := pricInfo.OnDemandInfo.(map[string]interface{})["priceDimensions"].(map[string]interface{})
							for _, v1 := range pricingData {
								pricInfo.PriceDimensions = v1
								break
							}
						}
						if pricInfo.PriceDimensions.(map[string]interface{})["pricePerUnit"] != nil {
							priceCurrency := pricInfo.PriceDimensions.(map[string]interface{})["pricePerUnit"].(map[string]interface{})
							for k2, v2 := range priceCurrency {
								pricInfo.PriceCurrency = k2
								pricInfo.Price = v2
								break
							}
						}
					}
				}

				pricInfo.Product = price["product"]
				pricInfo.ServiceCode = price["serviceCode"]
				pricInfo.PublicationDate = price["publicationDate"]
				d.StreamListItem(ctx, pricInfo)
			}
			return !isLast
		},
	)

	return nil, err
}
