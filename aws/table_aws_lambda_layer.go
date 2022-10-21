package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsLambdaLayer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_layer",
		Description: "AWS Lambda Layer",
		List: &plugin.ListConfig{
			Hydrate: listLambdaLayers,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "layer_name",
				Description: "The name of the layer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "layer_arn",
				Description: "The Amazon Resource Name (ARN) of the function layer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_date",
				Description: "The date that the version was created, in ISO 8601 format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LatestMatchingVersion.CreatedDate"),
			},
			{
				Name:        "description",
				Description: "The description of the version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LatestMatchingVersion.Description"),
			},
			{
				Name:        "layer_version_arn",
				Description: "The ARN of the layer version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LatestMatchingVersion.LayerVersionArn"),
			},
			{
				Name:        "license_info",
				Description: "The layer's open-source license.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LatestMatchingVersion.LicenseInfo"),
			},
			{
				Name:        "version",
				Description: "The version number.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("LatestMatchingVersion.Version"),
			},
			{
				Name:        "compatible_architectures",
				Description: "A list of compatible instruction set architectures.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LatestMatchingVersion.CompatibleArchitectures"),
			},
			{
				Name:        "compatible_runtimes",
				Description: "The layer's compatible runtimes.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LatestMatchingVersion.CompatibleRuntimes"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LayerName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LayerArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listLambdaLayers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_layer.listLambdaLayers", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// Set MaxItems to the maximum number allowed
	maxItems := int32(50)
	input := lambda.ListLayersInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	paginator := lambda.NewListLayersPaginator(svc, &input, func(o *lambda.ListLayersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lambda_function.listAwsLambdaFunctions", "api_error", err)
			return nil, err
		}

		for _, layer := range output.Layers {
			d.StreamListItem(ctx, layer)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
