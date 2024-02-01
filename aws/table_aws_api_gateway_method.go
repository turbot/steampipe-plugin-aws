package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/mitchellh/mapstructure"

	apigatewayv1 "github.com/aws/aws-sdk-go/service/apigateway"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayMethod(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_method",
		Description: "AWS API Gateway Method",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"resource_id", "rest_api_id", "http_method"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getApiGatewayMethod,
			Tags:    map[string]string{"service": "apigateway", "action": "GetMethod"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listRestAPI,
			Hydrate:       listApiGatewayMethods,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "rest_api_id", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "apigateway", "action": "GetResources"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(apigatewayv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "rest_api_id",
				Description: "The string identifier of the associated RestApi.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_id",
				Description: "The Resource identifier for the Method resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "http_method",
				Description: "The method's HTTP verb.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The full path for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path_part",
				Description: "The last path segment for this resource.",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "api_key_required",
				Description: "A boolean flag specifying whether a valid ApiKey is required to invoke this method.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "authorization_type",
				Description: "The method's authorization type. Valid values are NONE for open access, AWS_IAM for using AWS IAM permissions, CUSTOM for using a custom authorizer, or COGNITO_USER_POOLS for using a Cognito user pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "authorizer_id",
				Description: "The identifier of an Authorizer to use on this method. The authorizationType must be CUSTOM.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operation_name",
				Description: "A human-friendly operation identifier for the method. For example, you can assign the operationName of ListPets for the GET /pets method in the PetStore example.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "request_validator_id",
				Description: "The identifier of a RequestValidator for request validation.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON columns
			{
				Name:        "authorization_scopes",
				Description: "A list of authorization scopes configured on the method. The scopes are used with a COGNITO_USER_POOLS authorizer to authorize the method invocation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "method_integration",
				Description: "Gets the method's integration responsible for passing the client-submitted request to the back end and performing necessary transformations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "method_responses",
				Description: "Gets a method response associated with a given HTTP status code.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "request_models",
				Description: "A key-value map specifying data schemas, represented by Model resources, of the request payloads of given content types.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "request_parameters",
				Description: "A key-value map defining required or optional method request parameters that can be accepted by API Gateway.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("HttpMethod"),
			},
		}),
	}
}

type MethodInfo struct {
	ResourceId *string
	RestApiId  *string
	HttpMethod *string
	Path       *string
	PathPart   *string
	types.Method
}

//// LIST FUNCTION

func listApiGatewayMethods(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	restAPI := h.Item.(types.RestApi)
	restApiId := d.EqualsQualString("rest_api_id")
	if restApiId != "" {
		if restApiId != *restAPI.Id {
			return nil, nil
		}
	}

	// Create service
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		logger.Error("aws_api_gateway_method.listApiGatewayMethods", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 25 {
				maxLimit = 25
			} else {
				maxLimit = limit
			}
		}
	}

	input := &apigateway.GetResourcesInput{
		RestApiId: restAPI.Id,
		Embed:     []string{"methods"},
		Limit:     aws.Int32(maxLimit),
	}

	// Currently, there is no dedicated list API to retrieve all methods associated with a REST API and its resources directly.
	// The 'GetResources' API can be used to obtain all available methods for the resources by providing 'Embed' with the value '[]string{"methods"}' as an input parameter. The response format of this API mirrors that of the 'GetMethod' API.
	// To simplify the process and avoid requiring users to provide three parameters (RestAPI ID, Resource ID, and HTTP Method) for the 'GetMethod' API call, the 'GetResources' API is employed here for greater efficiency and ease of use.
	paginator := apigateway.NewGetResourcesPaginator(svc, input, func(o *apigateway.GetResourcesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_method.listApiGatewayMethods", "api_error", err)
			return nil, err
		}

		for _, resources := range output.Items {
			for httpMethod, item := range resources.ResourceMethods {
				d.StreamListItem(ctx, &MethodInfo{resources.Id, restAPI.Id, aws.String(httpMethod), resources.Path, resources.PathPart, item})

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getApiGatewayMethod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	restApiId := d.EqualsQualString("rest_api_id")
	resourceId := d.EqualsQualString("resource_id")
	httpMethod := d.EqualsQualString("http_method")

	if restApiId == "" || resourceId == "" || httpMethod == "" {
		return nil, nil
	}

	// Create session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_method.getApiGatewayMethod", "connection_error", err)
		return nil, err
	}

	params := &apigateway.GetMethodInput{
		RestApiId:  aws.String(restApiId),
		ResourceId: aws.String(resourceId),
		HttpMethod: aws.String(httpMethod),
	}

	detail, err := svc.GetMethod(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "NotFoundException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_api_gateway_method.getApiGatewayMethod", "api_error", err)
		return nil, err
	}

	var methodDetails types.Method
	err = mapstructure.Decode(detail, &methodDetails)
	if err != nil {
		return nil, err
	}

	return &MethodInfo{aws.String(resourceId), aws.String(restApiId), aws.String(httpMethod), nil, nil, methodDetails}, nil
}
