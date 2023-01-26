package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pipes"
	"github.com/aws/aws-sdk-go-v2/service/pipes/types"

	pipesv1 "github.com/aws/aws-sdk-go/service/pipes"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsPipes(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pipes_pipe",
		Description: "AWS Pipes Pipe",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getAwsPipe,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsPipes,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "current_state", Require: plugin.Optional},
				{Name: "desired_state", Require: plugin.Optional},
				{Name: "source_prefix", Require: plugin.Optional},
				{Name: "target_prefix", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(pipesv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the pipe.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the pipe.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time the pipe was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "current_state",
				Description: "The state the pipe is in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the pipe.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsPipe,
			},
			{
				Name:        "desired_state",
				Description: "The state the pipe should be in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enrichment",
				Description: "The ARN of the enrichment resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_time",
				Description: "When the pipe was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "role_arn",
				Description: "The ARN of the role that allows the pipe to send data to the target.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsPipe,
			},
			{
				Name:        "source",
				Description: "The ARN of the source resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_prefix",
				Description: "The prefix matching the pipe source.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("source_prefix"),
			},
			{
				Name:        "state_reason",
				Description: "The reason the pipe is in its current state.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsPipe,
			},
			{
				Name:        "target",
				Description: "The ARN of the target resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_prefix",
				Description: "The prefix matching the pipe target.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("target_prefix"),
			},
			{
				Name:        "enrichment_parameters",
				Description: "The parameters required to set up enrichment on your pipe.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsPipe,
			},
			{
				Name:        "target_parameters",
				Description: "The parameters required to set up a target for your pipe.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsPipe,
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
				Hydrate:     getAwsPipe,
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

func listAwsPipes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := PipesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pipes_pipe.listAwsPipes", "client_error", err)
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

	pagesLeft := true
	params := &pipes.ListPipesInput{
		// Default to the maximum allowed
		Limit: aws.Int32(maxLimit),
	}

	if d.EqualsQualString("current_state") != "" {
		params.CurrentState = types.PipeState(d.EqualsQualString("current_state"))
	}
	if d.EqualsQualString("desired_state") != "" {
		params.DesiredState = types.RequestedPipeState(d.EqualsQualString("desired_state"))
	}
	if d.EqualsQualString("source_prefix") != "" {
		params.SourcePrefix = aws.String(d.EqualsQualString("source_prefix"))
	}
	if d.EqualsQualString("target_prefix") != "" {
		params.TargetPrefix = aws.String(d.EqualsQualString("target_prefix"))
	}

	// API doesn't support aws-go-sdk-v2 paginator as of date
	for pagesLeft {
		output, err := svc.ListPipes(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_pipes_pipe.listAwsPipes", "api_error", err)
			return nil, err
		}

		for _, item := range output.Pipes {
			d.StreamListItem(ctx, item)
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

func getAwsPipe(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := PipesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pipes_pipe.getAwsPipe", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if d.EqualsQualString("name") != "" {
		name = d.EqualsQualString("name")
	} else {
		pipe := h.Item.(types.Pipe)
		name = *pipe.Name
	}

	// Build the params
	params := &pipes.DescribePipeInput{
		Name: &name,
	}

	// Get call
	data, err := svc.DescribePipe(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pipes_pipe.getAwsPipe", "api_error", err)
		return nil, err
	}

	return data, nil
}
