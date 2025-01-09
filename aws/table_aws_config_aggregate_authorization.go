package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"

	cognitoidentityEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsConfigAggregateAuthorization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_aggregate_authorization",
		Description: "AWS Config Aggregate Authorization",
		List: &plugin.ListConfig{
			Hydrate: listConfigAggregateAuthorizations,
			Tags:    map[string]string{"service": "config", "action": "DescribeAggregationAuthorizations"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getConfigAggregateAuthorizationsTags,
				Tags: map[string]string{"service": "config", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cognitoidentityEndpoint.COGNITO_IDENTITYServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the aggregation object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AggregationAuthorizationArn"),
			},
			{
				Name:        "authorized_account_id",
				Description: "The 12-digit account ID of the account authorized to aggregate data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "authorized_aws_region",
				Description: "The region authorized to collect aggregated data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time stamp when the aggregation authorization was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			{
				Name:        "tags_src",
				Description: "A list of tags attached to the Cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConfigAggregateAuthorizationsTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConfigAggregateAuthorizationsTags,
				Transform:   transform.FromField("Tags").Transform(configAggregateAuthorizationsTagListToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AggregationAuthorizationArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AggregationAuthorizationArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listConfigAggregateAuthorizations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_aggregate_authorization.listConfigAggregateAuthorizations", "get_client_error", err)
		return nil, err
	}

	input := &configservice.DescribeAggregationAuthorizationsInput{
		Limit: int32(100),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(input.Limit) {
			input.Limit = int32(*limit)
		}
	}

	paginator := configservice.NewDescribeAggregationAuthorizationsPaginator(svc, input, func(o *configservice.DescribeAggregationAuthorizationsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_config_aggregate_authorization.listConfigAggregateAuthorizations", "api_error", err)
			return nil, err
		}
		for _, aurhorization := range output.AggregationAuthorizations {
			d.StreamListItem(ctx, aurhorization)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

func getConfigAggregateAuthorizationsTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	auth := h.Item.(types.AggregationAuthorization)

	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_aggregate_authorization.getConfigAggregateAuthorizationsTags", "api_error", err)
		return nil, err
	}

	// Build the params
	params := &configservice.ListTagsForResourceInput{
		ResourceArn: auth.AggregationAuthorizationArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_aggregate_authorization.getConfigAggregateAuthorizationsTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func configAggregateAuthorizationsTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(tagList) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
