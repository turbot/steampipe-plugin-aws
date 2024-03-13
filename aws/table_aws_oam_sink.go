package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	oamv1 "github.com/aws/aws-sdk-go/service/oam"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/oam"
	"github.com/aws/aws-sdk-go-v2/service/oam/types"
)

func tableAwsOAMSink(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_oam_sink",
		Description: "AWS OAM Sink",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue", "ResourceNotFoundException"}),
			},
			Hydrate: getAwsOAMSink,
			Tags:    map[string]string{"service": "oam", "action": "GetSink"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsOAMSinks,
			Tags:    map[string]string{"service": "oam", "action": "ListSinks"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listAwsOAMSinkTags,
				Tags: map[string]string{"service": "oam", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(oamv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the sink.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The random ID string that Amazon Web Service generates as part of the sink ARN.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the sink.",
				Type:        proto.ColumnType_STRING,
			},

			/// Standard columns for all tables
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
				Hydrate:     listAwsOAMSinkTags,
				Transform:   transform.FromValue(),
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

func listAwsOAMSinks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := OAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_sink.listAwsOAMSinks", "connection_error", err)
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
			maxLimit = limit
		}
	}

	input := &oam.ListSinksInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := oam.NewListSinksPaginator(svc, input, func(o *oam.ListSinksPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_oam_sink.listAwsOAMSinks", "api_error", err)
			return nil, err
		}

		for _, sink := range output.Items {
			d.StreamListItem(ctx, sink)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsOAMSink(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQualString("arn")
	if arn == "" {
		return nil, nil
	}

	// Create Client
	svc, err := OAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_sink.getAwsOAMSink", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &oam.GetSinkInput{
		Identifier: aws.String(arn),
	}

	// Get call
	resp, err := svc.GetSink(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_sink.getAwsOAMSink", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func listAwsOAMSinkTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := getOAMSinkArn(h.Item)

	// Create Client
	svc, err := OAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_sink.listAwsOAMSinkTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &oam.ListTagsForResourceInput{
		ResourceArn: aws.String(arn),
	}

	// Get call
	resp, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_sink.listAwsOAMSinkTags", "api_error", err)
		return nil, err
	}

	if resp != nil {
		return resp.Tags, nil
	}

	return nil, nil
}

//// UTILITY FUNCTION

func getOAMSinkArn(item interface{}) string {
	switch item := item.(type) {
	case types.ListSinksItem:
		return *item.Arn
	case *oam.GetSinkOutput:
		return *item.Arn
	}
	return ""
}
