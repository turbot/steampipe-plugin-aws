package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayV2Api(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gatewayv2_api",
		Description: "AWS API Gateway Version 2 API",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("api_id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			Hydrate:           getAPIGatewayV2API,
		},
		List: &plugin.ListConfig{
			Hydrate: listAPIGatewayV2API,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the API",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_id",
				Description: "The API ID",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_endpoint",
				Description: "The URI of the API, of the form {api-id}.execute-api.{region}.amazonaws.com",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protocol_type",
				Description: "The API protocol",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_key_selection_expression",
				Description: "An API key selection expression. Supported only for WebSocket APIs",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "route_selection_expression",
				Description: "The route selection expression for the API. For HTTP APIs, the routeSelectionExpression must be ${request.method} ${request.path}. If not provided, this will be the default for HTTP APIs",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_date",
				Description: "The timestamp when the API was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
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
				Hydrate:     getAPIGatewayV2APIAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAPIGatewayV2API(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAPIGatewayV2API", "AWS_REGION", region)

	// Create Session
	svc, err := APIGatewayV2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	pagesLeft := true
	params := &apigatewayv2.GetApisInput{}

	for pagesLeft {
		result, err := svc.GetApis(params)
		if err != nil {
			return nil, err
		}

		for _, apiGatewayV2Api := range result.Items {
			d.StreamListItem(ctx, apiGatewayV2Api)
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAPIGatewayV2API(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAPIGatewayV2API")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := APIGatewayV2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	id := d.KeyColumnQuals["api_id"].GetStringValue()
	params := &apigatewayv2.GetApiInput{
		ApiId: aws.String(id),
	}

	apiData, err := svc.GetApi(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAPIGatewayV2API__", "ERROR", err)
		return nil, err
	}

	if apiData != nil {
		api := &apigatewayv2.Api{
			Name:                      apiData.Name,
			ApiId:                     apiData.ApiId,
			ApiEndpoint:               apiData.ApiEndpoint,
			ProtocolType:              apiData.ProtocolType,
			ApiKeySelectionExpression: apiData.ApiKeySelectionExpression,
			RouteSelectionExpression:  apiData.RouteSelectionExpression,
			CreatedDate:               apiData.CreatedDate,
			Tags:                      apiData.Tags,
		}
		return api, nil
	}

	return nil, nil
}

func getAPIGatewayV2APIAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	apigatewayV2Api := h.Item.(*apigatewayv2.Api)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + "::/apis/" + *apigatewayV2Api.ApiId}

	return akas, nil
}
