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

type integrationInfo = struct {
	types.Integration
	ApiId string
}

//// TABLE DEFINITION

func tableAwsAPIGatewayV2Integration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gatewayv2_integration",
		Description: "AWS API Gateway Version 2 Integration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"integration_id", "api_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException", "TooManyRequestsException"}),
			},
			Hydrate: getAPIGatewayV2Integration,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAPIGatewayV2API,
			Hydrate:       listAPIGatewayV2Integrations,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "integration_id",
				Description: "Represents the identifier of an integration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_id",
				Description: "Represents the identifier of an API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the integration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAPIGatewayV2IntegrationARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "description",
				Description: "Represents the description of an integration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_method",
				Description: "Specifies the integration's HTTP method type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_type",
				Description: "Represents an API method integration type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_uri",
				Description: "A string representation of a URI with a length between [1-2048]. For a Lambda integration, specify the URI of a Lambda function. For an HTTP integration, specify a fully-qualified URL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_gateway_managed",
				Description: "Specifies whether an integration is managed by API Gateway. If you created an API using using quick create, the resulting integration is managed by API Gateway. You can update a managed integration, but you can't delete it.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "connection_id",
				Description: "The ID of the VPC link for a private integration. Supported only for HTTP APIs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_type",
				Description: "Represents a connection type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_handling_strategy",
				Description: "Specifies how to handle response payload content type conversions. Supported only for WebSocket APIs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "credentials_arn",
				Description: "Specifies the credentials required for the integration, if any. For AWS integrations, three options are available. To specify an IAM Role for API Gateway to assume, use the role's Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_response_selection_expression",
				Description: "An expression used to extract information at runtime. See Selection Expressions(https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-websocket-api-selection-expressions.html#apigateway-websocket-api-apikey-selection-expressions for more information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "integration_subtype",
				Description: "A string with a length between [1-128].",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "passthrough_behavior",
				Description: "Represents passthrough behavior for an integration response. Supported only for WebSocket APIs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "payload_format_version",
				Description: "Specifies the format of the payload sent to an integration. Required for HTTP APIs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "template_selection_expression",
				Description: "The template selection expression for the integration. Supported only for WebSocket APIs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timeout_in_millis",
				Description: "Indicates custom timeout between 50 and 29,000 milliseconds for WebSocket APIs and between 50 and 30,000 milliseconds for HTTP APIs. The default timeout is 29 seconds for WebSocket APIs and 30 seconds for HTTP APIs.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "request_parameters",
				Description: "For HTTP API itegrations, without a specified integrationSubtype request parameters are a key-value map specifying how to transform HTTP requests before sending them to backend integrations. The key should follow the pattern <action>:<header|querystring|path>.<location>. The action can be append, overwrite or remove. For values, you can provide static values, or map request data, stage variables, or context variables that are evaluated at runtime. To learn more, see Transforming API requests and responses (https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-parameter-mapping.html).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "request_templates",
				Description: "Represents a map of Velocity templates that are applied on the request payload based on the value of the Content-Type header sent by the client. The content type value is the key in this map, and the template (as a String) is the value. Supported only for WebSocket APIs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "response_parameters",
				Description: "API requests and responses (https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-parameter-mapping.html).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tls_config",
				Description: "The TLS configuration for a private integration. If you specify a TLS configuration, private integration traffic uses the HTTPS protocol. Supported only for HTTP APIs.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IntegrationId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAPIGatewayV2IntegrationARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAPIGatewayV2Integrations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get API details
	api := h.Item.(types.Api)

	// Create Session
	svc, err := APIGatewayV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_integration.listAPIGatewayV2Integrations", "service_client_error", err)
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
	params := &apigatewayv2.GetIntegrationsInput{
		ApiId:      api.ApiId,
		MaxResults: aws.String(fmt.Sprint(maxLimit)),
	}

	for pagesLeft {
		result, err := svc.GetIntegrations(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gatewayv2_integration.listAPIGatewayV2Integrations", "api_error", err)
			return nil, err
		}

		for _, integration := range result.Items {
			d.StreamLeafListItem(ctx, integrationInfo{integration, *api.ApiId})

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

func getAPIGatewayV2Integration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := APIGatewayV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_integration.getAPIGatewayV2Integration", "service_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	api := d.KeyColumnQuals["api_id"].GetStringValue()
	key := d.KeyColumnQuals["integration_id"].GetStringValue()
	params := &apigatewayv2.GetIntegrationInput{
		ApiId:         aws.String(api),
		IntegrationId: aws.String(key),
	}

	item, err := svc.GetIntegration(ctx, params)

	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_integration.getAPIGatewayV2Integration", "api_error", err)
		return nil, err
	}

	if item != nil {
		integration := &types.Integration{
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

func getAPIGatewayV2IntegrationARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(integrationInfo)
	region := d.KeyColumnQualString(matrixKeyRegion)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/apis/" + data.ApiId + "/integrations/" + *data.IntegrationId

	return arn, nil
}
