package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayAPIKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_api_key",
		Description: "AWS API Gateway API Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NotFoundException"}),
			},
			Hydrate: getAPIKey,
		},
		List: &plugin.ListConfig{
			Hydrate: listAPIKeys,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "customer_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	logger := plugin.Logger(ctx)

	// Create service
	svc, err := APIGatewayService(ctx, d)
	if err != nil {
		logger.Trace("listAPIKeys", "connection error", err)
		return nil, err
	}

	input := &apigateway.GetApiKeysInput{
		Limit: aws.Int64(500),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["customer_id"] != nil {
		input.CustomerId = types.String(equalQuals["customer_id"].GetStringValue())
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = types.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	// List call
	err = svc.GetApiKeysPages(
		input,
		func(page *apigateway.GetApiKeysOutput, lastPage bool) bool {
			for _, items := range page.Items {
				d.StreamListItem(ctx, items)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAPIKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAPIKey")

	// Create session
	svc, err := APIGatewayService(ctx, d)
	if err != nil {
		return nil, err
	}

	id := d.KeyColumnQuals["id"].GetStringValue()
	params := &apigateway.GetApiKeyInput{
		ApiKey: aws.String(id),
	}

	detail, err := svc.GetApiKey(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAPIKey__", "ERROR", err)
		return nil, err
	}
	return detail, nil
}

func getAwsAPIKeysAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsAPIKeysAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	item := h.Item.(*apigateway.ApiKey)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/apikeys/" + *item.Id}

	return akas, nil
}
