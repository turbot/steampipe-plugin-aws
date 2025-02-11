package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayAPIKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_api_key",
		Description: "AWS API Gateway API Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getAPIKey,
			Tags:    map[string]string{"service": "apigateway", "action": "GetApiKey"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAPIKeys,
			Tags:    map[string]string{"service": "apigateway", "action": "GetApiKeys"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "customer_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_APIGATEWAY_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the API Key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The identifier of the API Key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "Specifies whether the API Key can be used by callers",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "created_date",
				Description: "The timestamp when the API Key was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_date",
				Description: "The timestamp when the API Key was last updated",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "customer_id",
				Description: "An AWS Marketplace customer identifier , when integrating with the AWS SaaS Marketplace",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the API Key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The value of the API Key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stage_keys",
				Description: "A list of Stage resources that are associated with the ApiKey resource",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to API key",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
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
				Hydrate:     getAwsAPIKeysAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAPIKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_api_key.listAPIKeys", "service_client_error", err)
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

	input := &apigateway.GetApiKeysInput{
		Limit: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["customer_id"] != nil {
		input.CustomerId = aws.String(equalQuals["customer_id"].GetStringValue())
	}

	paginator := apigateway.NewGetApiKeysPaginator(svc, input, func(o *apigateway.GetApiKeysPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_rest_api.listAPIKeys", "api_error", err)
			return nil, err
		}

		for _, items := range output.Items {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAPIKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getAPIKey", "service_client_error", err)
		return nil, err
	}

	id := d.EqualsQuals["id"].GetStringValue()
	params := &apigateway.GetApiKeyInput{
		ApiKey: aws.String(id),
	}

	detail, err := svc.GetApiKey(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getAPIKey", "api_error", err)
		return nil, err
	}
	return detail, nil
}

func getAwsAPIKeysAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	id := ""

	switch h.Item.(type) {
	case *apigateway.GetApiKeyOutput:
		id = *h.Item.(*apigateway.GetApiKeyOutput).Id
	case types.ApiKey:
		id = *h.Item.(types.ApiKey).Id
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/apikeys/" + id}

	return akas, nil
}
