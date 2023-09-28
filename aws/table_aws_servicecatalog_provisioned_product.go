package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"

	servicecatalogv1 "github.com/aws/aws-sdk-go/service/servicecatalog"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServicecatalogProvisionedProduct(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicecatalog_provisioned_product",
		Description: "AWS Service Catalog Provisioned Product",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id", "name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getServiceCatalogProvisionedProduct,
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceCatalogProvisionedProducts,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "accept_language",
					Require: plugin.Optional,
				},
				{
					Name:    "search_query",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(servicecatalogv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The ARN of the provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.Arn"),
			},
			{
				Name:        "id",
				Description: "The identifier of the provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.Id"),
			},
			{
				Name:        "created_time",
				Description: "The UTC time stamp of the creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ProvisionedProductDetail.CreatedTime"),
			},
			{
				Name:        "idempotency_token",
				Description: "A unique identifier that you provide to ensure idempotency. If multiple requests differ only by the idempotency token, the same response is returned for each repeated request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.IdempotencyToken"),
			},
			{
				Name:        "last_provisioning_record_id",
				Description: "The record identifier of the last request performed on this provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.LastProvisioningRecordId"),
			},
			{
				Name:        "accept_language",
				Description: "The language code.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("accept_language"),
			},
			{
				Name:        "search_query",
				Description: "The search query for the provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("search_query"),
			},
			{
				Name:        "last_record_id",
				Description: "The record identifier of the last request performed on this provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.LastRecordId"),
			},
			{
				Name:        "last_successful_provisioning_record_id",
				Description: "The record identifier of the last successful request performed on this provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.LastSuccessfulProvisioningRecordId"),
			},
			{
				Name:        "launch_role_arn",
				Description: "The ARN of the launch role associated with the provisioned product.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getServiceCatalogProvisionedProduct,
				Transform:   transform.FromField("ProvisionedProductDetail.LaunchRoleArn"),
			},
			{
				Name:        "name",
				Description: "The user-friendly name of the provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.Name"),
			},
			{
				Name:        "product_id",
				Description: "The product identifier. For example, prod-abcdzk7xy33qa.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.ProductId"),
			},
			{
				Name:        "provisioning_artifact_id",
				Description: "The identifier of the provisioning artifact. For example, pa-4abcdjnxjj6ne.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.ProvisioningArtifactId"),
			},
			{
				Name:        "status",
				Description: "The current status of the provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.Status"),
			},
			{
				Name:        "status_message",
				Description: "The current status message of the provisioned product.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.StatusMessage"),
			},
			{
				Name:        "type",
				Description: "The type of provisioned product. The supported values are CFN_STACK and CFN_STACKSET.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.Type"),
			},
			{
				Name:        "cloud_watch_dashboards",
				Description: "Any CloudWatch dashboards that were created when provisioning the product.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvisionedProductDetail.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProvisionedProductDetail.Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listServiceCatalogProvisionedProducts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Client
	svc, err := ServiceCatalogClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_provisioned_product.listServiceCatalogProvisionedProducts", "client_error", err)
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

	input := &servicecatalog.SearchProvisionedProductsInput{
		PageSize: maxLimit,
	}

	if d.EqualsQualString("accept_language") != "" {
		input.AcceptLanguage = aws.String(d.EqualsQualString("accept_language"))
	}

	filters := buildServiceCatalogProvisionedProductFilter(ctx, d.Quals)
	input.Filters = filters

	paginator := servicecatalog.NewSearchProvisionedProductsPaginator(svc, input, func(o *servicecatalog.SearchProvisionedProductsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicecatalog_provisioned_product.listServiceCatalogProvisionedProducts", "api_error", err)
			return nil, err
		}

		for _, item := range output.ProvisionedProducts {
			d.StreamListItem(ctx, &servicecatalog.DescribeProvisionedProductOutput{
				ProvisionedProductDetail: &types.ProvisionedProductDetail{
					Arn:                                item.Arn,
					CreatedTime:                        item.CreatedTime,
					Name:                               item.Name,
					Id:                                 item.Id,
					IdempotencyToken:                   item.IdempotencyToken,
					LastProvisioningRecordId:           item.LastProvisioningRecordId,
					LastRecordId:                       item.LastRecordId,
					LastSuccessfulProvisioningRecordId: item.LastSuccessfulProvisioningRecordId,
					ProductId:                          item.ProductId,
					ProvisioningArtifactId:             item.ProvisioningArtifactId,
					Status:                             item.Status,
					StatusMessage:                      item.StatusMessage,
					Type:                               item.Type,
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

func getServiceCatalogProvisionedProduct(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id, name string
	id = d.EqualsQualString("id")
	name = d.EqualsQualString("name")

	if id == "" && name == "" {
		return nil, nil
	}

	if id != "" && name != "" {
		return nil, fmt.Errorf("Both ProvisionedProductName and ProvisionedProductId cannot be passed in the where clause.")
	}

	// Create client
	svc, err := ServiceCatalogClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicecatalog_provisioned_product.getServiceCatalogProvisionedProduct", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Get call by passing id
	if id != "" {
		params := &servicecatalog.DescribeProvisionedProductInput{
			Id: aws.String(id),
		}
		op, err := svc.DescribeProvisionedProduct(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicecatalog_provisioned_product.getServiceCatalogProvisionedProduct", "api_error", err)
			return nil, err
		}
		return op, nil
	}

	// Get call by passing name
	if name != "" {
		params := &servicecatalog.DescribeProvisionedProductInput{
			Name: aws.String(name),
		}
		op, err := svc.DescribeProvisionedProduct(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_servicecatalog_provisioned_product.getServiceCatalogProvisionedProduct", "api_error", err)
			return nil, err
		}
		return op, nil
	}
	return nil, nil
}

//// UTILITY FUNCTIONS

// Buid servicecatalog product list call filter param

func buildServiceCatalogProvisionedProductFilter(ctx context.Context, quals plugin.KeyColumnQualMap) map[string][]string {
	filterQuals := map[string]string{
		"search_query": "SearchQuery",
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
