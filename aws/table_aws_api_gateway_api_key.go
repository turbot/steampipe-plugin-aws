package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayAPIKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_api_key",
		Description: "AWS API Gateway API Key",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       apiKeyFromKey,
			Hydrate:           getAPIKey,
		},
		List: &plugin.ListConfig{
			Hydrate: listAPIKeys,
		},
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

//// ITEM FROM KEY

func apiKeyFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	ID := quals["id"].GetStringValue()
	item := &apigateway.ApiKey{
		Id: &ID,
	}
	return item, nil
}

//// LIST FUNCTION

func listAPIKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listAPIKeys", "AWS_REGION", defaultRegion)

	// Create service
	svc, err := APIGatewayService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.GetApiKeysPages(
		&apigateway.GetApiKeysInput{},
		func(page *apigateway.GetApiKeysOutput, lastPage bool) bool {
			for _, items := range page.Items {
				d.StreamListItem(ctx, items)
			}
			return true
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAPIKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAPIKey")
	item := h.Item.(*apigateway.ApiKey)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := APIGatewayService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetApiKeyInput{
		ApiKey: item.Id,
	}

	detail, err := svc.GetApiKey(params)
	if err != nil {
		logger.Debug("getAPIKey__", "ERROR", err)
		return nil, err
	}
	return detail, nil
}

func getAwsAPIKeysAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsAPIKeysAkas")
	item := h.Item.(*apigateway.ApiKey)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + "::/apikeys/" + *item.Id}

	return akas, nil
}
