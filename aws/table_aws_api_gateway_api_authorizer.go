package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayAuthorizer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_authorizer",
		Description: "AWS API Gateway Authorizer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"rest_api_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			Hydrate:           getRestAPIAuthorizer,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listRestAPI,
			Hydrate:       listRestAPIAuthorizers,
		},
		GetMatrixItem: BuildRegionList,
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

//// LIST FUNCTION

func listRestAPIAuthorizers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get Rest API details
	restAPI := h.Item.(*apigateway.RestApi)

	// Create Session
	svc, err := APIGatewayService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetAuthorizersInput{
		Limit:     aws.Int64(500),
		RestApiId: restAPI.Id,
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.Limit {
			if *limit < 1 {
				params.Limit = types.Int64(1)
			} else {
				params.Limit = limit
			}
		}
	}

	op, err := svc.GetAuthorizers(params)
	if err != nil {
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
	plugin.Logger(ctx).Trace("getRestAPIAuthorizer")

	// Create Session
	svc, err := APIGatewayService(ctx, d)
	if err != nil {
		return nil, err
	}

	authorizerID := d.KeyColumnQuals["id"].GetStringValue()
	RestAPIID := d.KeyColumnQuals["rest_api_id"].GetStringValue()

	params := &apigateway.GetAuthorizerInput{
		AuthorizerId: aws.String(authorizerID),
		RestApiId:    aws.String(RestAPIID),
	}

	authorizerData, err := svc.GetAuthorizer(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getRestAPIAuthorizer__", "ERROR", err)
		return nil, err
	}

	return &authorizerRowData{authorizerData, aws.String(RestAPIID)}, nil
}

func getAPIGatewayAuthorizerAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAPIGatewayAuthorizerAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	apiAuthorizer := h.Item.(*authorizerRowData)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + ":" + commonColumnData.AccountId + "::/restapis/" + *apiAuthorizer.RestAPIId + "/authorizer/" + *apiAuthorizer.Authorizer.Id}
	return akas, nil
}
