package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayV2Api(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gatewayv2_api",
		Description: "AWS API Gateway Version 2 API",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("api_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getAPIGatewayV2API,
		},
		List: &plugin.ListConfig{
			Hydrate: listAPIGatewayV2API,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Name:        "disable_execute_api_endpoint",
				Description: "Specifies whether clients can invoke your API by using the default execute-api endpoint.",
				Type:        proto.ColumnType_BOOL,
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
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := APIGatewayV2Client(ctx, d)
	if err != nil {
		logger.Error("aws_api_gatewayv2_api.listAPIGatewayV2API", "connection error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
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

	pagesLeft := true
	params := &apigatewayv2.GetApisInput{
		MaxResults: aws.String(fmt.Sprint(maxLimit)),
	}

	for pagesLeft {
		result, err := svc.GetApis(ctx, params)
		if err != nil {
			logger.Error("aws_api_gatewayv2_api.listAPIGatewayV2API", "api_error", err)
			return nil, err
		}

		for _, apiGatewayV2Api := range result.Items {
			d.StreamListItem(ctx, apiGatewayV2Api)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
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

	// Create Session
	svc, err := APIGatewayV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_api.getAPIGatewayV2API", "service_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	id := d.KeyColumnQuals["api_id"].GetStringValue()
	params := &apigatewayv2.GetApiInput{
		ApiId: aws.String(id),
	}

	apiData, err := svc.GetApi(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_api.getAPIGatewayV2API", "api_error", err)
		return nil, err
	}

	if apiData != nil {
		api := &types.Api{
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	id := ""

	switch h.Item.(type) {
	case *types.Api:
		id = *h.Item.(*types.Api).ApiId
	case types.Api:
		id = *h.Item.(types.Api).ApiId
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/apis/" + id}

	return akas, nil
}
