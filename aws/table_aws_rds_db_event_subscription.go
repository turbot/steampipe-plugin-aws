package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBEventSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_event_subscription",
		Description: "AWS RDS DB Event Subscription",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cust_subscription_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"SubscriptionNotFound"}),
			},
			Hydrate: getRDSDBEventSubscription,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBEventSubscriptions,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cust_subscription_id",
				Description: "The RDS event notification subscription Id.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_aws_id",
				Description: "The AWS customer account associated with the RDS event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the event subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EventSubscriptionArn"),
			},
			{
				Name:        "status",
				Description: "The status of the RDS event notification subscription, it can be one of the following: creating | modifying | deleting | active | no-permission | topic-not-exist.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "A Boolean value indicating if the subscription is enabled. True indicates the subscription is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "sns_topic_arn",
				Description: "The topic ARN of the RDS event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_type",
				Description: "The source type for the RDS event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_creation_time",
				Description: "The time the RDS event notification subscription was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.From(convertStringToRFC3339Timestamp),
			},
			{
				Name:        "event_categories_list",
				Description: "A list of event categories for the RDS event notification subscription.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_ids_list",
				Description: "A list of source IDs for the RDS event notification subscription.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CustSubscriptionId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EventSubscriptionArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBEventSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_event_subscription.listRDSDBEventSubscriptions", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	input := &rds.DescribeEventSubscriptionsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	paginator := rds.NewDescribeEventSubscriptionsPaginator(svc, input, func(o *rds.DescribeEventSubscriptionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_event_subscription.listRDSDBEventSubscriptions", "api_error", err)
			return nil, err
		}

		for _, items := range output.EventSubscriptionsList {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBEventSubscription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	subscriptionId := d.KeyColumnQuals["cust_subscription_id"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_event_subscription.getRDSDBEventSubscription", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeEventSubscriptionsInput{
		SubscriptionName: aws.String(subscriptionId),
	}

	op, err := svc.DescribeEventSubscriptions(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_event_subscription.getRDSDBEventSubscription", "api_error", err)
		return nil, err
	}

	if op.EventSubscriptionsList != nil && len(op.EventSubscriptionsList) > 0 {
		return op.EventSubscriptionsList[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTION

func convertStringToRFC3339Timestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.EventSubscription)
	parsedTime := strings.Replace(*data.SubscriptionCreationTime, " ", "T", 1)

	return parsedTime, nil
}
