package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"

	eventbridgeEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsEventBridgeBus(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eventbridge_bus",
		Description: "AWS EventBridge Bus",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameter", "ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getAwsEventBridgeBus,
			Tags:    map[string]string{"service": "events", "action": "DescribeEventBus"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEventBridgeBuses,
			Tags:    map[string]string{"service": "events", "action": "ListEventBuses"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEventBridgeBusTags,
				Tags: map[string]string{"service": "events", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(eventbridgeEndpoint.EVENTSServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the event bus.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the account permitted to write events to the current account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy",
				Description: "The policy that enables the external account to send events to your account.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy that enables the external account to send events to your account in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the bus.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEventBridgeBusTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEventBridgeBusTags,
				Transform:   transform.FromField("Tags").Transform(eventBridgeBusTagListToTurbotTags),
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

//// LIST FUNCTION

func listAwsEventBridgeBuses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.listAwsEventBridgeBuses", "get_client_error", err)
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

	pagesLeft := true
	params := &eventbridge.ListEventBusesInput{
		// Default to the maximum allowed
		Limit: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		params.NamePrefix = aws.String(equalQuals["name"].GetStringValue())
	}

	// For case when listAwsEventBridgeBuses is used as parent hydrate in aws_eventbridge_rule table
	if equalQuals["name"] == nil && equalQuals["event_bus_name"] != nil {
		params.NamePrefix = aws.String(equalQuals["event_bus_name"].GetStringValue())
	}

	// API doesn't support aws-go-sdk-v2 paginator as of date
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.ListEventBuses(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eventbridge_bus.listAwsEventBridgeBuses", "api_error", err)
			return nil, err
		}

		for _, bus := range output.EventBuses {
			d.StreamListItem(ctx, &eventbridge.DescribeEventBusOutput{
				Name:   bus.Name,
				Arn:    bus.Arn,
				Policy: bus.Policy,
			})
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if output.NextToken != nil {
			pagesLeft = true
			params.NextToken = output.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsEventBridgeBus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.getAwsEventBridgeBus", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	arn := d.EqualsQuals["arn"].GetStringValue()

	// Build the params
	params := &eventbridge.DescribeEventBusInput{
		Name: &arn,
	}

	// Get call
	data, err := svc.DescribeEventBus(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.getAwsEventBridgeBus", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getAwsEventBridgeBusTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	arn := h.Item.(*eventbridge.DescribeEventBusOutput).Arn

	// Create Session
	svc, err := EventBridgeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.getAwsEventBridgeBusTags", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &eventbridge.ListTagsForResourceInput{
		ResourceARN: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eventbridge_bus.getAwsEventBridgeBusTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func eventBridgeBusTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(*eventbridge.ListTagsForResourceOutput)

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
