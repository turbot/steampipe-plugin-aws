package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"

	securityhubEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityhubEnabledProductSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_enabled_product_subscription",
		Description: "AWS Securityhub Enabled Product Subscription",
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubEnabledProductSubscriptions,
			Tags:    map[string]string{"service": "securityhub", "action": "ListEnabledProductsForImport"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(securityhubEndpoint.AWS_SECURITYHUB_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the product subscription.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(geSecurityHubProductSubscriptionTurbotTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type ProductSubscription struct {
	Arn string
}

//// LIST FUNCTION

func listSecurityHubEnabledProductSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_enabled_product_subscription.listSecurityHubEnabledProductSubscriptions", "client_error", err)
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

	input := &securityhub.ListEnabledProductsForImportInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := securityhub.NewListEnabledProductsForImportPaginator(svc, input, func(o *securityhub.ListEnabledProductsForImportPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Handle error for accounts that are not subscribed to AWS Security Hub
			if strings.Contains(err.Error(), "is not subscribed to AWS Security Hub") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_securityhub_enabled_product_subscription.listSecurityHubEnabledProductSubscriptions", "api_error", err)
			return nil, err
		}

		for _, product := range output.ProductSubscriptions {
			p := ProductSubscription{
				Arn: product,
			}

			d.StreamListItem(ctx, p)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS ////

func geSecurityHubProductSubscriptionTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(ProductSubscription)
	splitID := strings.Split(string(data.Arn), "/")
	title := splitID[len(splitID)-1]
	return title, nil
}
