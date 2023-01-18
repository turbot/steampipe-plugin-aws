package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayAuthorizer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_authorizer",
		Description: "AWS API Gateway Authorizer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"rest_api_id", "id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getRestAPIAuthorizer,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listRestAPI,
			Hydrate:       listRestAPIAuthorizers,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	Authorizer any
	RestAPIId  *string
}

//// LIST FUNCTION

func listRestAPIAuthorizers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get Rest API details
	restAPI := h.Item.(types.RestApi)

	// Create Session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_authorizer.listRestAPIAuthorizers", "service_client_error", err)
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

	params := &apigateway.GetAuthorizersInput{
		Limit:     aws.Int32(maxLimit),
		RestApiId: restAPI.Id,
	}

	op, err := svc.GetAuthorizers(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_authorizer.listRestAPIAuthorizers", "api_error", err)
		return nil, err
	}

	for _, authorizer := range op.Items {
		d.StreamLeafListItem(ctx, &authorizerRowData{authorizer, restAPI.Id})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRestAPIAuthorizer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_authorizer.getRestAPIAuthorizer", "service_client_error", err)
		return nil, err
	}

	authorizerID := d.KeyColumnQuals["id"].GetStringValue()
	RestAPIID := d.KeyColumnQuals["rest_api_id"].GetStringValue()

	params := &apigateway.GetAuthorizerInput{
		AuthorizerId: aws.String(authorizerID),
		RestApiId:    aws.String(RestAPIID),
	}

	authorizerData, err := svc.GetAuthorizer(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "NotFoundException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_api_gateway_authorizer.getRestAPIAuthorizer", "api_error", err)
		return nil, err
	}

	return &authorizerRowData{authorizerData, aws.String(RestAPIID)}, nil
}

func getAPIGatewayAuthorizerAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	authorizer := h.Item.(*authorizerRowData).Authorizer

	id := ""
	restApiId := h.Item.(*authorizerRowData).RestAPIId
	switch item := authorizer.(type) {
	case *apigateway.GetAuthorizerOutput:
		id = *item.Id
	case types.Authorizer:
		id = *item.Id
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + ":" + commonColumnData.AccountId + "::/restapis/" + *restApiId + "/authorizer/" + id}
	return akas, nil
}
