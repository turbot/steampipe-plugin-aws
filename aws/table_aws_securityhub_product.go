package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsSecurityhubProduct(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_product",
		Description: "AWS Securityhub Product",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("product_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidAccessException"}),
			Hydrate:           getSecurityHubProduct,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubProducts,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductName"),
			},
			{
				Name:        "product_arn",
				Description: "The ARN assigned to the product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activation_url",
				Description: "The URL used to activate the product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "company_name",
				Description: "The name of the company that provides the product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "marketplace_url",
				Description: "The URL for the page that contains more information about the product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "categories",
				Description: "The categories assigned to the product.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "integration_types",
				Description: "The types of integration that the product supports.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product_subscription_resource_policy",
				Description: "The resource policy associated with the product.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProductName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProductArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubProducts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listSecurityyHubProducts", "AWS_REGION", region)

	// Create Session
	svc, err := SecurityHubService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeProductsPages(
		&securityhub.DescribeProductsInput{},
		func(page *securityhub.DescribeProductsOutput, isLast bool) bool {
			for _, product := range page.Products {
				d.StreamListItem(ctx, product)
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listSecurityHubProducts", "query_error", err)
		return nil, nil
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getSecurityHubProduct(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	productArn := d.KeyColumnQuals["product_arn"].GetStringValue()

	// Create service
	svc, err := SecurityHubService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &securityhub.DescribeProductsInput{
		ProductArn: &productArn,
	}

	op, err := svc.DescribeProducts(params)
	if err != nil {
		return nil, err
	}

	if op.Products != nil && len(op.Products) > 0 {
		return op.Products[0], nil
	}

	return nil, nil
}
