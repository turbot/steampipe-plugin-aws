package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"

	redshiftv1 "github.com/aws/aws-sdk-go/service/redshift"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRedshiftEventSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_event_subscription",
		Description: "AWS Redshift Event Subscription",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cust_subscription_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"SubscriptionNotFound"}),
			},
			Hydrate: getRedshiftEventSubscription,
			Tags:    map[string]string{"service": "redshift", "action": "DescribeEventSubscriptions"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftEventSubscriptions,
			Tags:    map[string]string{"service": "redshift", "action": "DescribeEventSubscriptions"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(redshiftv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cust_subscription_id",
				Description: "The name of the Amazon Redshift event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_aws_id",
				Description: "The AWS customer account associated with the Amazon Redshift event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "A boolean value indicating whether the subscription is enabled or disabled",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "severity",
				Description: "The event severity specified in the Amazon Redshift event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sns_topic_arn",
				Description: "The Amazon Resource Name (ARN) of the Amazon SNS topic used by the event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_type",
				Description: "The source type of the events returned by the Amazon Redshift event notification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the Amazon Redshift event notification subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_creation_time",
				Description: "The date and time the Amazon Redshift event notification subscription was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "event_categories_list",
				Description: "The list of Amazon Redshift event categories specified in the event notification subscription.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_ids_list",
				Description: "A list of the sources that publish events to the Amazon Redshift event notification subscription.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the event subscription.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CustSubscriptionId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(redshiftEventSubListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRedshiftEventSubscriptionAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listRedshiftEventSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_event_subscription.listRedshiftEventSubscriptions", "connection_error", err)
		return nil, err
	}

	input := &redshift.DescribeEventSubscriptionsInput{
		MaxRecords: aws.Int32(100),
	}
	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxRecords {
			if limit < 20 {
				input.MaxRecords = aws.Int32(20)
			} else {
				input.MaxRecords = aws.Int32(limit)
			}
		}
	}

	// List call
	paginator := redshift.NewDescribeEventSubscriptionsPaginator(svc, input, func(o *redshift.DescribeEventSubscriptionsPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_redshift_event_subscription.listRedshiftEventSubscriptions", "api_error", err)
			return nil, err
		}

		for _, items := range output.EventSubscriptionsList {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedshiftEventSubscription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_event_subscription.getRedshiftEventSubscription", "connection_error", err)
		return nil, err
	}

	name := d.EqualsQuals["cust_subscription_id"].GetStringValue()

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &redshift.DescribeEventSubscriptionsInput{
		SubscriptionName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeEventSubscriptions(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_event_subscription.getRedshiftEventSubscription", "api_error", err)
		return nil, err
	}

	if len(data.EventSubscriptionsList) > 0 {
		return data.EventSubscriptionsList[0], nil
	}
	return nil, nil
}

func getRedshiftEventSubscriptionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	parameterData := h.Item.(types.EventSubscription)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_event_subscription.getRedshiftEventSubscriptionAkas", "getCommonColumns_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":redshift:" + region + ":" + commonColumnData.AccountId + ":eventsubscription"

	if strings.HasPrefix(*parameterData.CustSubscriptionId, ":") {
		aka = aka + *parameterData.CustSubscriptionId
	} else {
		aka = aka + ":" + *parameterData.CustSubscriptionId
	}

	return []string{aka}, nil
}

//// TRANSFORM FUNCTION

func redshiftEventSubListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(types.EventSubscription)

	if len(tagList.Tags) > 0 {
		turbotTagsMap := map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
