package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type integrationInfo = struct {
	apigatewayv2.Integration
	ApiId string
}

//// TABLE DEFINITION
func tableAwsAPIGatewayV2Integration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gatewayv2_integration",
		Description: "AWS API Gateway Version 2 Integration",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"integration_id", "api_id"}),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException", "TooManyRequestsException"}),
			Hydrate:           getAPIGatewayV2Integration,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAPIGatewayV2API,
			Hydrate:       listAPIGatewayV2Integrations,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "integration_id",
				Description: "The identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_id",
				Description: "TODO",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A string with a length between [0-1024].",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_method",
				Description: "A string with a length between [1-64].",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_type",
				Description: "Represents an API method integration type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_uri",
				Description: "A string representation of a URI with a length between [1-2048].",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_id",
				Description: "A string with a length between [1-1024].",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_type",
				Description: "Represents a connection type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timeout_in_millis",
				Description: "",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "integration_subtype",
				Description: "A string with a length between [1-128].",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_gateway_managed",
				Description: "Specifies whether an integration is managed by API Gateway. If you created an API using using quick create, the resulting integration is managed by API Gateway. You can update a managed integration, but you can't delete it.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "content_handling_strategy",
				Description: "Specifies how to handle response payload content type conversions. Supported only for WebSocket APIs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "credentials_arn",
				Description: "Represents an Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_response_selection_expression",
				Description: "An expression used to extract information at runtime. See Selection Expressions(https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-websocket-api-selection-expressions.html#apigateway-websocket-api-apikey-selection-expressions for more information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "passthrough_behavior",
				Description: "Represents passthrough behavior for an integration response. Supported only for WebSocket APIs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "payload_format_version",
				Description: "A string with a length between [1-64].",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "request_parameters",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "request_templates",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "response_parameters",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "template_selection_expression",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tls_config",
				Description: "The TLS configuration for a private integration. If you specify a TLS configuration, private integration traffic uses the HTTPS protocol. Supported only for HTTP APIs.",
				Type:        proto.ColumnType_JSON,
			},
			// Standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAPIGatewayV2IntegrationAkas,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IntegrationId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAPIGatewayV2Integrations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAPIGatewayV2Integrations", "AWS_REGION", region)
	api := h.Item.(*apigatewayv2.Api)

	// Create Session
	svc, err := APIGatewayV2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	pagesLeft := true
	params := &apigatewayv2.GetIntegrationsInput{
		ApiId: api.ApiId,
	}

	for pagesLeft {
		result, err := svc.GetIntegrations(params)
		if err != nil {
			return nil, err
		}

		for _, integration := range result.Items {
			d.StreamLeafListItem(ctx, integrationInfo{*integration, *api.ApiId})
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

func getAPIGatewayV2Integration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAPIGatewayV2Integration")
	api := d.KeyColumnQuals["api_id"].GetStringValue()
	key := d.KeyColumnQuals["integration_id"].GetStringValue()

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

	params := &apigatewayv2.GetIntegrationInput{
		ApiId:         aws.String(api),
		IntegrationId: aws.String(key),
	}

	item, err := svc.GetIntegration(params)

	if err != nil {
		plugin.Logger(ctx).Debug("getAPIGatewayV2API__", "ERROR", err)
		return nil, err
	}

	if item != nil {
		integration := &apigatewayv2.Integration{
			ApiGatewayManaged:                      item.ApiGatewayManaged,
			ConnectionId:                           item.ConnectionId,
			ConnectionType:                         item.ConnectionType,
			ContentHandlingStrategy:                item.ContentHandlingStrategy,
			CredentialsArn:                         item.CredentialsArn,
			Description:                            item.Description,
			IntegrationId:                          item.IntegrationId,
			IntegrationMethod:                      item.IntegrationMethod,
			IntegrationResponseSelectionExpression: item.IntegrationResponseSelectionExpression,
			IntegrationSubtype:                     item.IntegrationSubtype,
			IntegrationType:                        item.IntegrationType,
			IntegrationUri:                         item.IntegrationUri,
			PassthroughBehavior:                    item.PassthroughBehavior,
			PayloadFormatVersion:                   item.PayloadFormatVersion,
			RequestParameters:                      item.RequestParameters,
			RequestTemplates:                       item.RequestTemplates,
			ResponseParameters:                     item.ResponseParameters,
			TemplateSelectionExpression:            item.TemplateSelectionExpression,
			TimeoutInMillis:                        item.TimeoutInMillis,
			TlsConfig:                              item.TlsConfig,
		}
		return integrationInfo{*integration, api}, nil
	}

	return nil, nil
}

func getAPIGatewayV2IntegrationAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(integrationInfo)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + "::/apis/" + data.ApiId + "/integrations/" + *data.IntegrationId}

	return akas, nil
}
