package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/aws/aws-sdk-go-v2/service/opensearch/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsOpenSearchReservedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_opensearch_reserved_instance",
		Description: "AWS OpenSearch Reserved Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("reserved_instance_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getOpenSearchReservedInstance,
			Tags:    map[string]string{"service": "es", "action": "DescribeReservedInstances"},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenSearchReservedInstances,
			Tags:    map[string]string{"service": "es", "action": "DescribeReservedInstances"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ES_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "reserved_instance_id",
				Description: "The unique identifier for the reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reservation_name",
				Description: "The customer-specified identifier to track this reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "The OpenSearch instance type offered by the Reserved Instance offering.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the Reserved Instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The date and time when the reservation was purchased.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "duration",
				Description: "The duration, in seconds, for which the OpenSearch instance is reserved.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "fixed_price",
				Description: "The upfront fixed charge you will pay to purchase the specific Reserved Instance offering.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "usage_price",
				Description: "The hourly rate at which you're charged for the domain using this Reserved Instance.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "currency_code",
				Description: "The currency code for the offering.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_count",
				Description: "The number of OpenSearch instances that have been reserved.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "payment_option",
				Description: "The payment option as defined in the Reserved Instance offering.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reserved_instance_offering_id",
				Description: "The unique identifier of the Reserved Instance offering.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "billing_subscription_id",
				Description: "The unique identifier of the billing subscription.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "recurring_charges",
				Description: "The recurring charge to your account, regardless of whether you create any domains using the Reserved Instance offering.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservationName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOpenSearchReservedInstanceARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listOpenSearchReservedInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := OpenSearchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_opensearch_reserved_instance.listOpenSearchReservedInstances", "client_error", err)
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

	input := &opensearch.DescribeReservedInstancesInput{
		MaxResults: maxLimit,
	}

	// API doesn't support aws-go-sdk-v2 paginator as of date
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.DescribeReservedInstances(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_opensearch_reserved_instance.listOpenSearchReservedInstances", "api_error", err)
			return nil, err
		}

		for _, reservedInstance := range output.ReservedInstances {
			d.StreamListItem(ctx, reservedInstance)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Check if there are more results
		if output.NextToken == nil {
			break
		}

		// Set the next token for the next iteration
		input.NextToken = output.NextToken
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenSearchReservedInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	reservedInstanceID := d.EqualsQuals["reserved_instance_id"].GetStringValue()

	// Empty check
	if reservedInstanceID == "" {
		return nil, nil
	}

	// Create service
	svc, err := OpenSearchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_opensearch_reserved_instance.getOpenSearchReservedInstance", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &opensearch.DescribeReservedInstancesInput{
		ReservedInstanceId: aws.String(reservedInstanceID),
	}

	result, err := svc.DescribeReservedInstances(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_opensearch_reserved_instance.getOpenSearchReservedInstance", "api_error", err)
		return nil, err
	}

	if len(result.ReservedInstances) > 0 {
		return result.ReservedInstances[0], nil
	}

	return nil, nil
}

func getOpenSearchReservedInstanceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	reservedInstance := h.Item.(types.ReservedInstance)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	arn := "arn:" + commonColumnData.Partition + ":es:" + region + ":" + commonColumnData.AccountId + ":reserved-instance/" + *reservedInstance.ReservedInstanceId

	return arn, nil
}
