package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	lambdav1 "github.com/aws/aws-sdk-go/service/lambda"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsLambdaLayerVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_layer_version",
		Description: "AWS Lambda Layer Version",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"layer_name", "version"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameter", "InvalidParameterValueException"}),
			},
			Hydrate: getLambdaLayerVersion,
		},
		List: &plugin.ListConfig{
			Hydrate:       listLambdaLayerVersions,
			ParentHydrate: listLambdaLayers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "layer_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lambdav1.EndpointsID),
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
				Type:        proto.ColumnType_TIMESTAMP,
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
	layerName := h.Item.(types.LayersListItem).LayerName

	// Create service
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_layer_version.listLambdaLayerVersions", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	equalQuals := d.EqualsQuals
	// Minimize the API call with the given layer name
	if equalQuals["layer_name"] != nil {
		if equalQuals["layer_name"].GetStringValue() != *layerName {
			return nil, nil
		}
	}

	maxItems := int32(50)
	input := lambda.ListLayerVersionsInput{
		LayerName: layerName,
	}

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

	input.MaxItems = aws.Int32(maxItems)
	paginator := lambda.NewListLayerVersionsPaginator(svc, &input, func(o *lambda.ListLayerVersionsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lambda_function.listAwsLambdaFunctions", "api_error", err)
			return nil, err
		}

		for _, version := range output.LayerVersions {
			d.StreamListItem(ctx, LayerVersionInfo{*layerName, lambda.GetLayerVersionOutput{
				CompatibleArchitectures: version.CompatibleArchitectures,
				CompatibleRuntimes:      version.CompatibleRuntimes,
				CreatedDate:             version.CreatedDate,
				Description:             version.Description,
				LayerVersionArn:         version.LayerVersionArn,
				LicenseInfo:             version.LicenseInfo,
				Version:                 version.Version,
			}})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLambdaLayerVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var layerName string
	var version int64
	if h.Item != nil {
		layerName = h.Item.(LayerVersionInfo).LayerName
		version = h.Item.(LayerVersionInfo).Version
	} else {
		layerName = d.EqualsQuals["layer_name"].GetStringValue()
		version = d.EqualsQuals["version"].GetInt64Value()
	}

	if strings.TrimSpace(layerName) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_layer_version.getLambdaLayerVersion", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &lambda.GetLayerVersionInput{
		LayerName:     aws.String(layerName),
		VersionNumber: version,
	}

	// Get call
	data, err := svc.GetLayerVersion(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_layer_version.getLambdaLayerVersion", "api_error", err)
		return nil, err
	}

	return LayerVersionInfo{layerName, *data}, nil
}

func getLambdaLayerVersionPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var layerName string
	var version int64
	if h.Item != nil {
		layerName = h.Item.(LayerVersionInfo).LayerName
		version = h.Item.(LayerVersionInfo).Version
	} else {
		layerName = d.EqualsQuals["layer_name"].GetStringValue()
		version = d.EqualsQuals["version"].GetInt64Value()
	}

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_layer_version.getLambdaLayerVersionPolicy", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &lambda.GetLayerVersionPolicyInput{
		LayerName:     aws.String(layerName),
		VersionNumber: version,
	}

	// Get call
	data, err := svc.GetLayerVersionPolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// If the function does not exist or does not have url config, the operation returns a 404 (ResourceNotFoundException) error.
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_lambda_layer_version.getLambdaLayerVersionPolicy", "api_error", err)
		return nil, err
	}

	return data, nil
}
