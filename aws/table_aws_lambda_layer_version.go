package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsLambdaLayerVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_layer_version",
		Description: "AWS Lambda Layer Version",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"layer_name", "version"}),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "InvalidParameter", "InvalidParameterValueException"}),
			Hydrate:           getLambdaLayerVersion,
		},
		List: &plugin.ListConfig{
			Hydrate:           listLambdaLayerVersions,
			ParentHydrate:     listLambdaLayers,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "layer_name",
				Description: "The name of the layer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "layer_arn",
				Description: "The ARN of the layer.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getLambdaLayerVersion,
			},
			{
				Name:        "layer_version_arn",
				Description: "The ARN of the layer version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_date",
				Description: "The date that the version was created, in ISO 8601 format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_info",
				Description: "The layer's open-source license.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "revision_id",
				Description: "A unique identifier for the current revision of the policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getLambdaLayerVersionPolicy,
			},
			{
				Name:        "version",
				Description: "The version number.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "compatible_architectures",
				Description: "A list of compatible instruction set architectures.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "compatible_runtimes",
				Description: "The layer's compatible runtimes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "content",
				Description: "Details about the layer version.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLambdaLayerVersion,
			},
			{
				Name:        "policy",
				Description: "The policy document.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLambdaLayerVersionPolicy,
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy document in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLambdaLayerVersionPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
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
				Transform:   transform.FromField("LayerVersionArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type LayerVersionInfo struct {
	LayerName string
	lambda.GetLayerVersionOutput
}

//// LIST FUNCTION

func listLambdaLayerVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listLambdaLayerVersions")

	// Create service
	svc, err := LambdaService(ctx, d)
	if err != nil {
		logger.Error("listLambdaLayerVersions", "error_LambdaService", err)
		return nil, err
	}

	layerName := h.Item.(*lambda.LayersListItem).LayerName

	// Set MaxItems to the maximum number allowed
	input := lambda.ListLayerVersionsInput{
		LayerName: layerName,
		MaxItems:  types.Int64(50),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			input.MaxItems = limit
		}
	}

	err = svc.ListLayerVersionsPages(
		&input,
		func(page *lambda.ListLayerVersionsOutput, lastPage bool) bool {
			for _, version := range page.LayerVersions {
				d.StreamListItem(ctx, LayerVersionInfo{*layerName, lambda.GetLayerVersionOutput{
					CompatibleArchitectures: version.CompatibleArchitectures,
					CompatibleRuntimes:      version.CompatibleRuntimes,
					CreatedDate:             version.CreatedDate,
					Description:             version.Description,
					LayerVersionArn:         version.LayerVersionArn,
					LicenseInfo:             version.LicenseInfo,
					Version:                 version.Version,
				}})
			}
			return !lastPage
		},
	)

	if err != nil {
		logger.Error("listLambdaLayerVersions", "error_ListLayerVersionsPages", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLambdaLayerVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getLambdaLayerVersion")

	var layerName string
	var version int64
	if h.Item != nil {
		layerName = h.Item.(LayerVersionInfo).LayerName
		version = *h.Item.(LayerVersionInfo).Version
	} else {
		layerName = d.KeyColumnQuals["layer_name"].GetStringValue()
		version = d.KeyColumnQuals["version"].GetInt64Value()
	}

	// Create Session
	svc, err := LambdaService(ctx, d)
	if err != nil {
		logger.Error("getLambdaLayerVersion", "error_LambdaService", err)
		return nil, err
	}

	// Build the params
	params := &lambda.GetLayerVersionInput{
		LayerName:     &layerName,
		VersionNumber: &version,
	}

	// Get call
	data, err := svc.GetLayerVersion(params)
	if err != nil {
		logger.Error("getLambdaLayerVersion", "error_GetLayerVersion", err)
		return nil, err
	}

	return LayerVersionInfo{layerName, *data}, nil
}

func getLambdaLayerVersionPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getLambdaLayerVersionPolicy")

	var layerName string
	var version int64
	if h.Item != nil {
		layerName = h.Item.(LayerVersionInfo).LayerName
		version = *h.Item.(LayerVersionInfo).Version
	} else {
		layerName = d.KeyColumnQuals["layer_name"].GetStringValue()
		version = d.KeyColumnQuals["version"].GetInt64Value()
	}

	// Create Session
	svc, err := LambdaService(ctx, d)
	if err != nil {
		logger.Error("getLambdaLayerVersionPolicy", "error_LambdaService", err)
		return nil, err
	}

	// Build the params
	params := &lambda.GetLayerVersionPolicyInput{
		LayerName:     &layerName,
		VersionNumber: &version,
	}

	// Get call
	data, err := svc.GetLayerVersionPolicy(params)
	if err != nil {
		logger.Error("getLambdaLayerVersionPolicy", "error_GetLayerVersionPolicy", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return data, nil
}
