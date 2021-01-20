package aws

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayRestAPI(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_rest_api",
		Description: "AWS API Gateway Rest API ",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("api_id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       restAPIFromKey,
			Hydrate:           getRestAPI,
		},
		List: &plugin.ListConfig{
			Hydrate: listRestAPI,
		},
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
				Type:        proto.ColumnType_DATETIME,
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

//// ITEM FROM KEY

func restAPIFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	ID := quals["api_id"].GetStringValue()
	item := &apigateway.RestApi{
		Id: &ID,
	}
	return item, nil
}

//// LIST FUNCTION

func listRestAPI(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listRestAPI", "AWS_REGION", defaultRegion)

	// Create service
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.GetRestApisPages(
		&apigateway.GetRestApisInput{},
		func(page *apigateway.GetRestApisOutput, lastPage bool) bool {
			for _, items := range page.Items {
				d.StreamListItem(ctx, items)
			}
			return true
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRestAPI(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getRestAPI")
	item := h.Item.(*apigateway.RestApi)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetRestApiInput{
		RestApiId: item.Id,
	}

	detail, err := svc.GetRestApi(params)
	if err != nil {
		logger.Debug("GetRestApi__", "ERROR", err)
		return nil, err
	}
	return detail, nil
}

func getAwsRestAPITurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRestAPITurbotData")
	item := h.Item.(*apigateway.RestApi)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + "::/restapis/" + *item.Id}

	return akas, nil
}

//// TRANSFORM FUNCTION

// unmarshalJSON :: parse the yaml-encoded data and return the result
func unmarshalJSON(_ context.Context, d *transform.TransformData) (interface{}, error) {
	inputStr := types.SafeString(d.Value)
	var result interface{}
	if inputStr != "" {
		// Resource IAM policy for aws_api_gateway_rest_api is stored as stringfied json object after removing the double quotes from end
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
