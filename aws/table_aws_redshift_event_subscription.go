package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsRedshiftEventSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_event_subscription",
		Description: "AWS Redshift Event Subscription",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cust_subscription_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"SubscriptionNotFound"}),
			},
			Hydrate: getAwsRedshiftEventSubscription,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsRedshiftEventSubscriptions,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Hydrate:     getAwsRedshiftEventSubscriptionAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsRedshiftEventSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsRedshiftEventSubscriptions", "AWS_REGION")

	// Create session
	svc, err := RedshiftService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &redshift.DescribeEventSubscriptionsInput{
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
		func(page *redshift.DescribeEventSubscriptionsOutput, isLast bool) bool {
			for _, parameter := range page.EventSubscriptionsList {
				d.StreamListItem(ctx, parameter)

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

func getAwsRedshiftEventSubscription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsRedshiftEventSubscription")

	// Create Session
	svc, err := RedshiftService(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["cust_subscription_id"].GetStringValue()

	// Build the params
	params := &redshift.DescribeEventSubscriptionsInput{
		SubscriptionName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeEventSubscriptions(params)
	if err != nil {
		return nil, err
	}

	if data.EventSubscriptionsList != nil && len(data.EventSubscriptionsList) > 0 {
		return data.EventSubscriptionsList[0], nil
	}
	return nil, nil
}

func getAwsRedshiftEventSubscriptionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRedshiftEventSubscriptionAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	parameterData := h.Item.(*redshift.EventSubscription)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
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
	plugin.Logger(ctx).Trace("redshiftEventSubListToTurbotTags")

	tagList := d.HydrateItem.(*redshift.EventSubscription)

	if tagList.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
