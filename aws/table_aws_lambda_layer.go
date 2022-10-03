package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turbot/go-kit/types"
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
	logger := plugin.Logger(ctx)
	logger.Trace("listLambdaLayers")

	// Create service
	svc, err := LambdaService(ctx, d)
	if err != nil {
		logger.Error("listLambdaLayers", "error_LambdaService", err)
		return nil, err
	}

	// Set MaxItems to the maximum number allowed
	input := lambda.ListLayersInput{
		MaxItems: types.Int64(50),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = aws.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	err = svc.ListLayersPages(
		&input,
		func(page *lambda.ListLayersOutput, lastPage bool) bool {
			for _, layer := range page.Layers {
				d.StreamListItem(ctx, layer)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	if err != nil {
		logger.Error("listLambdaLayers", "error_ListLayersPages", err)
		return nil, err
	}

	return nil, nil
}
