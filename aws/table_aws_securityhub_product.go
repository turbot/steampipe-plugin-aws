package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"

	securityhubv1 "github.com/aws/aws-sdk-go/service/securityhub"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityhubProduct(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_product",
		Description: "AWS Securityhub Product",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("product_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
			Hydrate: getSecurityHubProduct,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubProducts,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(securityhubv1.EndpointsID),
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

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_product.listSecurityHubProducts", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
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

	input := &securityhub.DescribeProductsInput{
		MaxResults: maxLimit,
	}

	paginator := securityhub.NewDescribeProductsPaginator(svc, input, func(o *securityhub.DescribeProductsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Handle error for accounts that are not subscribed to AWS Security Hub
			if strings.Contains(err.Error(), "is not subscribed to AWS Security Hub") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_securityhub_product.listSecurityHubProducts", "api_error", err)
			return nil, err
		}

		for _, product := range output.Products {
			d.StreamListItem(ctx, product)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHubProduct(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	productArn := d.KeyColumnQuals["product_arn"].GetStringValue()

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_product.getSecurityHubProduct", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.DescribeProductsInput{
		ProductArn: &productArn,
	}

	// Get call
	op, err := svc.DescribeProducts(ctx, params)
	if err != nil {
		// Handle error for accounts that are not subscribed to AWS Security Hub
		if strings.Contains(err.Error(), "is not subscribed to AWS Security Hub") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_securityhub_product.getSecurityHubProduct", "api_error", err)
		return nil, err
	}

	if op.Products != nil && len(op.Products) > 0 {
		return op.Products[0], nil
	}

	return nil, nil
}
