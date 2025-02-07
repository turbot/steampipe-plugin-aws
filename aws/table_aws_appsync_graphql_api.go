package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"

	appsyncEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsAppsyncGraphQLApi(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appsync_graphql_api",
		Description: "AWS AppSync GraphQL API",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("api_id"),
			Hydrate:    getAppsyncGraphqlApi,
			Tags:       map[string]string{"service": "appsync", "action": "GetGraphqlApi"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAppsyncGraphqlApis,
			Tags:    map[string]string{"service": "appsync", "action": "ListGraphqlApis"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "api_type",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(appsyncEndpoint.AWS_APPSYNC_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{

			{
				Name:        "name",
				Description: "The API name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_id",
				Description: "The API ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of AppSync GraphQL API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_type",
				Description: "The value that indicates whether the GraphQL API is a standard API ( GRAPHQL ) or merged API ( MERGED ).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "authentication_type",
				Description: "The authentication type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "merged_api_execution_role_arn",
				Description: "The Identity and Access Management service role ARN for a merged API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner",
				Description: "The account owner of the GraphQL API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_contact",
				Description: "The owner contact information for an API resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "visibility",
				Description: "Sets the value of the GraphQL API to public ( GLOBAL ) or private ( PRIVATE ).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "waf_web_acl_arn",
				Description: "The ARN of the WAF access control list (ACL) associated with this GraphqlApi, if one exists.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "query_depth_limit",
				Description: "The maximum depth a query can have in a single request.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "resolver_count_limit",
				Description: "The maximum number of resolvers that can be invoked in a single request.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "xray_enabled",
				Description: "A flag indicating whether to use X-Ray tracing for this GraphqlApi.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "log_config",
				Description: "The Amazon CloudWatch Logs configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "open_id_connect_config",
				Description: "The OpenID Connect configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OpenIDConnectConfig"),
			},
			{
				Name:        "additional_authentication_providers",
				Description: "A list of additional authentication providers for the GraphqlApi API.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dns",
				Description: "The DNS records for the API.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lambda_authorizer_config",
				Description: "Configuration for Lambda function authorization.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "uris",
				Description: "The URIs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_pool_config",
				Description: "The Amazon Cognito user pool configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "enhanced_metrics_config",
				Description: "The enhancedMetricsConfig object.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "introspection_config",
				Description: "Sets the value of the GraphQL API to enable ( ENABLED ) or disable ( DISABLED ) introspection. If no value is provided, the introspection configuration will be set to ENABLED by default.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAppsyncGraphqlApis(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := AppSyncClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appsync_graphql_api.listAppsyncGraphqlApis", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxResults := int32(25)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxResults {
			maxResults = int32(limit)
		}
	}

	input := appsync.ListGraphqlApisInput{
		MaxResults: maxResults,
	}

	if d.EqualsQualString("api_type") != "" {
		input.ApiType = types.GraphQLApiType(d.EqualsQualString("api_type"))
	}

	pageLeft := true

	for pageLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		res, err := svc.ListGraphqlApis(ctx, &input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_appsync_graphql_api.listAppsyncGraphqlApis", "api_error", err)
			return nil, err
		}

		for _, api := range res.GraphqlApis {
			d.StreamListItem(ctx, api)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if res.NextToken != nil {
			input.NextToken = res.NextToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAppsyncGraphqlApi(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	id := d.EqualsQualString("api_id")

	// Empty input check
	if id == "" {
		return nil, nil
	}

	// Create Session
	svc, err := AppSyncClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appsync_graphql_api.getAppsyncGraphqlApi", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &appsync.GetGraphqlApiInput{
		ApiId: aws.String(id),
	}

	rowData, err := svc.GetGraphqlApi(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appsync_graphql_api.getAppsyncGraphqlApi", "api_error", err)
		return nil, err
	}

	if rowData != nil {
		return rowData.GraphqlApi, nil
	}

	return nil, nil
}
