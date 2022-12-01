package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securitylake"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityLakeSubscriber(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securitylake_subscriber",
		Description: "AWS Security Lake Subscriber",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.SingleColumn("product_arn"),
		// 	IgnoreConfig: &plugin.IgnoreConfig{
		// 		ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
		// 	},
		// 	Hydrate: getSecurityHubProduct,
		// },
		List: &plugin.ListConfig{
			Hydrate: listSecurityLakeSubscribers,
		},
		GetMatrixItemFunc: BuildRegionList,
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

func listSecurityLakeSubscribers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Client
	svc, err := SecurityLakeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securitylake_subscriber.listSecurityLakeSubscribers", "client_error", err)
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

	input := &securitylake.ListSubscribersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := securitylake.NewListSubscribersPaginator(svc, input, func(o *securitylake.ListSubscribersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_securitylake_subscriber.listSecurityLakeSubscribers", "api_error", err)
			return nil, err
		}

		for _, item := range output.Subscribers {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityLakeSubscriber(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create session
	svc, err := SecurityLakeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securitylake_subscriber.getSecurityLakeSubscriber", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &securitylake.GetSubscriberInput{
		Id: aws.String(id),
	}

	// Get call
	op, err := svc.GetSubscriber(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securitylake_subscriber.getSecurityLakeSubscriber", "api_error", err)
		return nil, err
	}

	return op.Subscriber, nil
}
