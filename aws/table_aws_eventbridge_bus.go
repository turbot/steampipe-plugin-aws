package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsEventBridgeBus(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eventbridge_bus",
		Description: "AWS EventBridge Bus",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameter", "ResourceNotFoundException", "ValidationException"}),
			Hydrate:           getAwsEventBridgeBus,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEventBridgeBuses,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name_prefix", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
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
				Name:        "name_prefix",
				Description: "Specifying this limits the results to only those event buses with names that start with the specified prefix.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("name_prefix"),
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
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsEventBridgeBuses")

	// Create session
	svc, err := EventBridgeService(ctx, d)
	if err != nil {
		logger.Error("listAwsEventBridgeBuses", "error_EventBridgeService", err)
		return nil, err
	}

	// List call
	input := eventbridge.ListEventBusesInput{
		// Default to the maximum allowed
		Limit: aws.Int64(100),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["name_prefix"] != nil {
		input.NamePrefix = aws.String(equalQuals["name_prefix"].GetStringValue())
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = aws.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	for {
		response, err := svc.ListEventBuses(&input)
		if err != nil {
			logger.Error("listAwsEventBridgeBuses", "error_ListEventBuses", err)
			return nil, err
		}

		for _, bus := range response.EventBuses {
			d.StreamListItem(ctx, &eventbridge.DescribeEventBusOutput{
				Name:   bus.Name,
				Arn:    bus.Arn,
				Policy: bus.Policy,
			})
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				break
			}
		}

		if response.NextToken == nil {
			break
		}
		input.NextToken = response.NextToken
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsEventBridgeBus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEventBridgeBus")

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// Create Session
	svc, err := EventBridgeService(ctx, d)
	if err != nil {
		logger.Error("getAwsEventBridgeBus", "error_EventBridgeService", err)
		return nil, err
	}
	// Build the params
	params := &eventbridge.DescribeEventBusInput{
		Name: &arn,
	}

	// Get call
	data, err := svc.DescribeEventBus(params)
	if err != nil {
		logger.Error("getAwsEventBridgeBus", "error_DescribeEventBus", err)
		return nil, err
	}

	return data, nil
}

func getAwsEventBridgeBusTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsEventBridgeBusTags")

	arn := h.Item.(*eventbridge.DescribeEventBusOutput).Arn

	// Create Session
	svc, err := EventBridgeService(ctx, d)
	if err != nil {
		logger.Error("getAwsEventBridgeBusTags", "error_EventBridgeService", err)
		return nil, err
	}

	// Build the params
	params := &eventbridge.ListTagsForResourceInput{
		ResourceARN: arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Error("getAwsEventBridgeBusTags", "error_ListTagsForResource", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func eventBridgeBusTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("eventBridgeBusTagListToTurbotTags")
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
