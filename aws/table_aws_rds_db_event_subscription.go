package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
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
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"SubscriptionNotFound"}),
			},
			Hydrate: getRDSDBEventSubscription,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBEventSubscriptions,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listRDSDBEventSubscriptions")

	// Create Session
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &rds.DescribeEventSubscriptionsInput{
		MaxRecords: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxRecords {
			if *limit < 20 {
				input.MaxRecords = aws.Int64(20)
			} else {
				input.MaxRecords = limit
			}
		}
	}

	// List call
	err = svc.DescribeEventSubscriptionsPages(
		input,
		func(page *rds.DescribeEventSubscriptionsOutput, isLast bool) bool {
			for _, eventSubscription := range page.EventSubscriptionsList {
				d.StreamListItem(ctx, eventSubscription)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBEventSubscription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	subscriptionId := d.KeyColumnQuals["cust_subscription_id"].GetStringValue()

	// Create service
	svc, err := RDSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &rds.DescribeEventSubscriptionsInput{
		SubscriptionName: aws.String(subscriptionId),
	}

	op, err := svc.DescribeEventSubscriptions(params)
	if err != nil {
		return nil, err
	}

	if op.EventSubscriptionsList != nil && len(op.EventSubscriptionsList) > 0 {
		return op.EventSubscriptionsList[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTION

func convertStringToRFC3339Timestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*rds.EventSubscription)
	parsedTime := strings.Replace(*data.SubscriptionCreationTime, " ", "T", 1)

	return parsedTime, nil
}
