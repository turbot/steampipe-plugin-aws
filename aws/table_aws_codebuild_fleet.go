package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/aws/aws-sdk-go-v2/service/codebuild/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeBuildFleet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codebuild_fleet",
		Description: "AWS CodeBuild Fleet",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException", "ResourceNotFoundException"}),
			},
			Hydrate: getCodeBuildFleet,
			Tags:    map[string]string{"service": "codebuild", "action": "BatchGetFleets"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeBuildFleets,
			Tags:    map[string]string{"service": "codebuild", "action": "ListFleets"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CODEBUILD_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "base_capacity",
				Description: "The base capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "created",
				Description: "When the compute fleet was created, expressed in Unix time format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_modified",
				Description: "When the compute fleet's settings were last modified, expressed in Unix time format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The current status of the compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_reason",
				Description: "The reason for the current status of the compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compute_type",
				Description: "The compute type of the compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compute_fleet_type",
				Description: "The type of compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_config",
				Description: "Information about the VPC configuration that AWS CodeBuild uses to access resources in a VPC.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "environment_type",
				Description: "The environment type of the compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "current_capacity",
				Description: "The current capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "desired_capacity",
				Description: "The desired capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_capacity",
				Description: "The maximum capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_capacity",
				Description: "The minimum capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "fleet_service_role",
				Description: "The service role ARN for the compute fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tag key and value pairs associated with this compute fleet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: "A map of tags key and value pairs associated with this compute fleet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(codeBuildFleetTurbotTags),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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

func listCodeBuildFleets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service client
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_fleet.listCodeBuildFleets", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Build the params
	params := &codebuild.ListFleetsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := codebuild.NewListFleetsPaginator(svc, params, func(o *codebuild.ListFleetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codebuild_fleet.listCodeBuildFleets", "api_error", err)
			return nil, err
		}

		for _, fleet := range output.Fleets {
			d.StreamListItem(ctx, fleet)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeBuildFleet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		data := h.Item.(types.Fleet)
		arn = *data.Arn
	} else {
		arn = d.EqualsQualString("arn")
	}

	// Empty check
	if arn == "" {
		return nil, nil
	}

	// Create service client
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_fleet.getCodeBuildFleet", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &codebuild.BatchGetFleetsInput{
		Names: []string{arn},
	}

	// Get call
	data, err := svc.BatchGetFleets(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_fleet.getCodeBuildFleet", "api_error", err)
		return nil, err
	}

	if len(data.Fleets) > 0 {
		return data.Fleets[0], nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func codeBuildFleetTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.Fleet)

	if data.Tags == nil {
		return nil, nil
	}

	// Turn the tags into a map
	tags := make(map[string]string)
	for _, tag := range data.Tags {
		tags[*tag.Key] = *tag.Value
	}

	return tags, nil
}
