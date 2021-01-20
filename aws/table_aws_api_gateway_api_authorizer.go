package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayAuthorizer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_authorizer",
		Description: "AWS API Gateway Authorizer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"rest_api_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       apiAuthorizerFromKey,
			Hydrate:           getRestAPIAuthorizer,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listRestAPI,
			Hydrate:       listRestAPIAuthorizers,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The identifier for the authorizer resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.Id"),
			},
			{
				Name:        "name",
				Description: "The name of the authorizer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.Name"),
			},
			{
				Name:        "rest_api_id",
				Description: "The id of the rest api",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RestAPIId"),
			},
			{
				Name:        "auth_type",
				Description: "Optional customer-defined field, used in OpenAPI imports and exports without functional impact",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.AuthType"),
			},
			{
				Name:        "authorizer_credentials",
				Description: "Specifies the required credentials as an IAM role for API Gateway to invoke the authorizer",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.AuthorizerCredentials"),
			},
			{
				Name:        "authorizer_uri",
				Description: "Specifies the authorizer's Uniform Resource Identifier (URI). For TOKEN or REQUEST authorizers, this must be a well-formed Lambda function URI",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.AuthorizerUri"),
			},
			{
				Name:        "identity_validation_expression",
				Description: "A validation expression for the incoming identity token. For TOKEN authorizers, this value is a regular expression",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.IdentityValidationExpression"),
			},
			{
				Name:        "identity_source",
				Description: "The identity source for which authorization is requested",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.IdentitySource"),
			},
			{
				Name:        "provider_arns",
				Description: "A list of the Amazon Cognito user pool ARNs for the COGNITO_USER_POOLS authorizer",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Authorizer.ProviderARNs"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Authorizer.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAPIGatewayAuthorizerAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type authorizerRowData = struct {
	Authorizer *apigateway.Authorizer
	RestAPIId  *string
}

//// BUILD HYDRATE INPUT

func apiAuthorizerFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	authorizerID := quals["id"].GetStringValue()
	RestAPIID := quals["rest_api_id"].GetStringValue()
	item := &authorizerRowData{
		RestAPIId: &RestAPIID,
		Authorizer: &apigateway.Authorizer{
			Id: &authorizerID,
		},
	}

	return item, nil
}

//// LIST FUNCTION

func listRestAPIAuthorizers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listRestAPIAuthorizers", "AWS_REGION", defaultRegion)

	restAPI := h.Item.(*apigateway.RestApi)

	// Create Session
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetAuthorizersInput{
		RestApiId: restAPI.Id,
	}

	op, err := svc.GetAuthorizers(params)
	if err != nil {
		return nil, err
	}

	for _, authorizer := range op.Items {
		d.StreamLeafListItem(ctx, &authorizerRowData{authorizer, restAPI.Id})
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRestAPIAuthorizer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getRestAPIAuthorizer")

	apiAuthorizer := h.Item.(*authorizerRowData)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetAuthorizerInput{
		AuthorizerId: apiAuthorizer.Authorizer.Id,
		RestApiId:    apiAuthorizer.RestAPIId,
	}

	authorizerData, err := svc.GetAuthorizer(params)
	if err != nil {
		logger.Debug("getRestAPIAuthorizer__", "ERROR", err)
		return nil, err
	}

	return &authorizerRowData{authorizerData, apiAuthorizer.RestAPIId}, nil
}

func getAPIGatewayAuthorizerAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAPIGatewayAuthorizerAkas")
	apiAuthorizer := h.Item.(*authorizerRowData)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + ":" + commonColumnData.AccountId + "::/restapis/" + *apiAuthorizer.RestAPIId + "/authorizer/" + *apiAuthorizer.Authorizer.Id}
	return akas, nil
}
