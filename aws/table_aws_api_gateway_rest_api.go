package aws

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	apigatewayv1 "github.com/aws/aws-sdk-go/service/apigateway"

	go_kit_packs "github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayRestAPI(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_rest_api",
		Description: "AWS API Gateway Rest API ",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("api_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getRestAPI,
		},
		List: &plugin.ListConfig{
			Hydrate: listRestAPI,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(apigatewayv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The API's name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "api_id",
				Description: "The API's identifier. This identifier is unique across all of APIs in API Gateway",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "version",
				Description: "A version identifier for the API",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Version"),
			},
			{
				Name:        "api_key_source",
				Description: "The source of the API key for metering requests according to a usage plan",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApiKeySource"),
			},
			{
				Name:        "created_date",
				Description: "The timestamp when the API was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreatedDate"),
			},
			{
				Name:        "description",
				Description: "The API's description",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "minimum_compression_size",
				Description: "A nullable integer that is used to enable compression (with non-negative between 0 and 10485760 (10M) bytes, inclusive) or disable compression (with a null value) on an API. When compression is enabled, compression or decompression is not applied on the payload if the payload size is smaller than this value",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MinimumCompressionSize"),
			},
			{
				Name:        "policy",
				Description: "A stringified JSON policy document that applies to this RestApi regardless of the caller and Method configuration",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policy").Transform(unmarshalJSON).Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policy").Transform(unmarshalJSON).Transform(policyToCanonical),
			},
			{
				Name:        "binary_media_types",
				Description: "The list of binary media types supported by the RestApi. By default, the RestApi supports only UTF-8-encoded text payloads",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BinaryMediaTypes"),
			},
			{
				Name:        "endpoint_configuration_types",
				Description: "The endpoint configuration of this RestApi showing the endpoint types of the API",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointConfiguration.Types"),
			},
			{
				Name:        "endpoint_configuration_vpc_endpoint_ids",
				Description: "The endpoint configuration of this RestApi showing the endpoint types of the API",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointConfiguration.VpcEndpointIds"),
			},
			{
				Name:        "warnings",
				Description: "The warning messages reported when failonwarnings is turned on during API import",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Warnings"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
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
				Hydrate:     getAwsRestAPITurbotData,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listRestAPI(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create service
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		logger.Error("aws_api_gateway_rest_api.listRestAPI", "connection error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &apigateway.GetRestApisInput{
		Limit: aws.Int32(maxLimit),
	}

	paginator := apigateway.NewGetRestApisPaginator(svc, input, func(o *apigateway.GetRestApisPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_rest_api.listRestAPI", "api_error", err)
			return nil, err
		}

		for _, items := range output.Items {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRestAPI(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getRestAPI", "connection_error", err)
		return nil, err
	}

	id := d.KeyColumnQuals["api_id"].GetStringValue()
	params := &apigateway.GetRestApiInput{
		RestApiId: aws.String(id),
	}

	detail, err := svc.GetRestApi(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "NotFoundException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getRestAPI", "api_error", err)
		return nil, err
	}
	return detail, nil
}

func getAwsRestAPITurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	id := ""

	switch h.Item.(type) {
	case *apigateway.GetRestApiOutput:
		id = *h.Item.(*apigateway.GetRestApiOutput).Id
	case types.RestApi:
		id = *h.Item.(types.RestApi).Id
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/restapis/" + id}

	return akas, nil
}

//// TRANSFORM FUNCTION

// unmarshalJSON :: parse the yaml-encoded data and return the result
func unmarshalJSON(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	inputStr := go_kit_packs.SafeString(d.Value)
	var result interface{}
	if inputStr != "" {
		// Resource IAM policy for aws_api_gateway_rest_api is stored as stringified json object after removing the double quotes from end
		decoded, err := url.QueryUnescape("\"" + inputStr + "\"")
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_rest_api.unmarshalJSON", err)
			return nil, err
		}

		err = json.Unmarshal([]byte(decoded), &result)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_rest_api.unmarshalJSON", err)
			return nil, err
		}
	}

	return result, nil
}
