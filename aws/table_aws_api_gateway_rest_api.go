package aws

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayRestAPI(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_rest_api",
		Description: "AWS API Gateway Rest API ",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("api_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NotFoundException"}),
			},
			Hydrate: getRestAPI,
		},
		List: &plugin.ListConfig{
			Hydrate: listRestAPI,
		},
		GetMatrixItem: BuildRegionList,
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
	svc, err := APIGatewayService(ctx, d)
	if err != nil {
		logger.Trace("listRestAPI", "connection error", err)
		return nil, err
	}

	input := &apigateway.GetRestApisInput{
		Limit: aws.Int64(500),
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = types.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	// List call
	err = svc.GetRestApisPages(
		input,
		func(page *apigateway.GetRestApisOutput, lastPage bool) bool {
			for _, items := range page.Items {
				d.StreamListItem(ctx, items)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRestAPI(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRestAPI")

	// Create session
	svc, err := APIGatewayService(ctx, d)
	if err != nil {
		return nil, err
	}

	id := d.KeyColumnQuals["api_id"].GetStringValue()
	params := &apigateway.GetRestApiInput{
		RestApiId: aws.String(id),
	}

	detail, err := svc.GetRestApi(params)
	if err != nil {
		plugin.Logger(ctx).Debug("GetRestApi__", "ERROR", err)
		return nil, err
	}
	return detail, nil
}

func getAwsRestAPITurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRestAPITurbotData")
	region := d.KeyColumnQualString(matrixKeyRegion)
	item := h.Item.(*apigateway.RestApi)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/restapis/" + *item.Id}

	return akas, nil
}

//// TRANSFORM FUNCTION

// unmarshalJSON :: parse the yaml-encoded data and return the result
func unmarshalJSON(_ context.Context, d *transform.TransformData) (interface{}, error) {
	inputStr := types.SafeString(d.Value)
	var result interface{}
	if inputStr != "" {
		// Resource IAM policy for aws_api_gateway_rest_api is stored as stringified json object after removing the double quotes from end
		decoded, err := url.QueryUnescape("\"" + inputStr + "\"")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(decoded), &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
