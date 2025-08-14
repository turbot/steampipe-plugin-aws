package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsAppsyncApi(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appsync_api",
		Description: "AWS AppSync API",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("api_id"),
			Hydrate:    getAppsyncApi,
			Tags:       map[string]string{"service": "appsync", "action": "GetApi"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAppsyncApis,
			Tags:    map[string]string{"service": "appsync", "action": "ListApis"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_APPSYNC_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "api_id",
				Description: "The API ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the API.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApiArn"),
			},
			{
				Name:        "created",
				Description: "The date and time that the API was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "dns",
				Description: "The DNS records for the API. This will include an HTTP and a real-time endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "event_config",
				Description: "The Event API configuration. This includes the default authorization configuration for connecting, publishing, and subscribing to an Event API.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "owner_contact",
				Description: "The owner contact information for the API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "waf_web_acl_arn",
				Description: "The Amazon Resource Name (ARN) of the WAF web access control list (web ACL) associated with this API, if one exists.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "xray_enabled",
				Description: "A flag indicating whether to use X-Ray tracing for this API.",
				Type:        proto.ColumnType_BOOL,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAppsyncApi,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: "A map with keys of TagKey objects and values of TagValue objects.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppsyncApi,
				Transform:   transform.FromField("ApiArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAppsyncApis(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := AppSyncClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appsync_api.listAppsyncApis", "connection_error", err)
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

	// Using the actual ListApis API
	input := appsync.ListApisInput{
		MaxResults: maxResults,
	}

	paginator := appsync.NewListApisPaginator(svc, &input, func(o *appsync.ListApisPaginatorOptions) {
		o.Limit = maxResults
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		res, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_appsync_api.listAppsyncApis", "api_error", err)
			return nil, err
		}

		for _, api := range res.Apis {
			d.StreamListItem(ctx, api)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAppsyncApi(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var apiId string
	if h.Item != nil {
		// If we have an item from the list, extract the API ID
		if api, ok := h.Item.(types.Api); ok {
			apiId = *api.ApiId
		}
	} else {
		// If this is a get call, use the key column
		apiId = d.EqualsQualString("api_id")
	}

	if apiId == "" {
		return nil, nil
	}

	// Create service
	svc, err := AppSyncClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appsync_api.getAppsyncApi", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Using the actual GetApi API
	params := &appsync.GetApiInput{
		ApiId: aws.String(apiId),
	}

	rowData, err := svc.GetApi(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appsync_api.getAppsyncApi", "api_error", err)
		return nil, err
	}

	if rowData != nil && rowData.Api != nil {
		return rowData.Api, nil
	}

	return nil, nil
}
