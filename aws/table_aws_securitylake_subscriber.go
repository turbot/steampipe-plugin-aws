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
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("subscription_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
			Hydrate: getSecurityLakeSubscriber,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityLakeSubscribers,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "subscriber_name",
				Description: "The name of your Amazon Security Lake subscriber account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_id",
				Description: "The subscription ID of the Amazon Security Lake subscriber account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The date and time when the subscription was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "subscription_status",
				Description: "Subscription status of the Amazon Security Lake subscriber account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The date and time when the subscription was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "external_id",
				Description: "The external ID of the subscriber.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) specifying the role of the subscriber.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "s3_bucket_arn",
				Description: "The Amazon Resource Name (ARN) for the Amazon S3 bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sns_arn",
				Description: "The Amazon Resource Name (ARN) for the Amazon Simple Notification Service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscriber_description",
				Description: "The subscriber descriptions for a subscriber account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_endpoint",
				Description: "The subscription endpoint to which exception messages are posted.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_protocol",
				Description: "The subscription protocol to which exception messages are posted.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "access_types",
				Description: "You can choose to notify subscribers of new objects with an Amazon Simple Queue Service (Amazon SQS) queue or through messaging to an HTTPS endpoint provided by the subscriber. Subscribers can consume data by directly querying Lake Formation tables in your S3 bucket via services like Amazon Athena. This subscription type is defined as LAKEFORMATION.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_types",
				Description: "Amazon Security Lake supports logs and events collection for the natively-supported Amazon Web Services services.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubscriberName"),
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
