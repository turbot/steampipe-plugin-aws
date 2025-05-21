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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getCodeBuildFleet,
			Tags:    map[string]string{"service": "codebuild", "action": "BatchGetFleets"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeBuildFleets,
			Tags:    map[string]string{"service": "codebuild", "action": "ListFleets"},
			// According to the supported endpoints listed here: https://docs.aws.amazon.com/general/latest/gr/codebuild.html,
			// AWS CodeBuild is available in all regions as mentioned in the above doc. However, the specific resource type "fleet" is not supported in every region.
			// If you attempt to perform the ListFleets operation in a region where fleet support is unavailable,
			// the API returns an InvalidInputException with the message:
			// "An error occurred (InvalidInputException) when calling the ListFleets operation: Unknown operation ListFleets"
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CODEBUILD_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the compute fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildFleet,
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
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "created",
				Description: "When the compute fleet was created, expressed in Unix time format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "last_modified",
				Description: "When the compute fleet's settings were last modified, expressed in Unix time format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "status",
				Description: "The current status of the compute fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "status_reason",
				Description: "The reason for the current status of the compute fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "compute_type",
				Description: "The compute type of the compute fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "compute_fleet_type",
				Description: "The type of compute fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "vpc_config",
				Description: "Information about the VPC configuration that AWS CodeBuild uses to access resources in a VPC.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "environment_type",
				Description: "The environment type of the compute fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "current_capacity",
				Description: "The current capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "desired_capacity",
				Description: "The desired capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "max_capacity",
				Description: "The maximum capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "min_capacity",
				Description: "The minimum capacity of the compute fleet.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "fleet_service_role",
				Description: "The service role ARN for the compute fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildFleet,
			},
			{
				Name:        "tags_src",
				Description: "A list of tag key and value pairs associated with this compute fleet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
				Hydrate:     getCodeBuildFleet,
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
			{
				Name:        "tags",
				Description: "A map of tags key and value pairs associated with this compute fleet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(codeBuildFleetTurbotTags),
				Hydrate:     getCodeBuildFleet,
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

	// Unsupported region check
	if svc == nil {
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
			d.StreamListItem(ctx, types.Fleet{Arn: &fleet})

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

	if svc == nil {
		return nil, nil
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
