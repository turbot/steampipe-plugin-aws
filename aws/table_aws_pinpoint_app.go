package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"

	pinpointEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsPinpointApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pinpoint_app",
		Description: "AWS Pinpoint App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getPinpointApp,
			Tags:    map[string]string{"service": "mobiletargeting", "action": "GetApp"},
		},
		List: &plugin.ListConfig{
			Hydrate: listPinpointApps,
			Tags:    map[string]string{"service": "mobiletargeting", "action": "GetApps"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getPinpointApplicationSettings,
				Tags: map[string]string{"service": "mobiletargeting", "action": "GetApplicationSettings"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(pinpointEndpoint.AWS_PINPOINT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The display name of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_date",
				Description: "The date and time, in ISO 8601 format, when the application's settings were last modified.",
				Hydrate:     getPinpointApplicationSettings,
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "campaign_hook",
				Description: "The settings for the AWS Lambda function to invoke by default as a code hook for campaigns in the application.",
				Hydrate:     getPinpointApplicationSettings,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "limits",
				Description: "The default sending limits for campaigns in the application.",
				Hydrate:     getPinpointApplicationSettings,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "quiet_time",
				Description: "The default quiet time for campaigns in the application.",
				Hydrate:     getPinpointApplicationSettings,
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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

func listPinpointApps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := PinpointClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pinpoint_app.listPinpointApps", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Page size must be greater than 0 and less than or equal to 1000
	input := &pinpoint.GetAppsInput{
		PageSize: aws.String("1000"),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit <= 1000 {
			if *limit < 1 {
				input.PageSize = aws.String("1")
			} else {
				input.PageSize = aws.String(fmt.Sprint(*limit))
			}
		}
	}

	// API doesn't support aws-g0-sdk-v2 paginator as of date
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		apps, err := svc.GetApps(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_pinpoint_app.listPinpointApps", "api_error", err)
			return nil, err
		}

		if apps.ApplicationsResponse == nil {
			break
		}

		for _, app := range apps.ApplicationsResponse.Item {
			d.StreamListItem(ctx, app)

			// Check if context has been cancelled or if the limit has been reached (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if apps.ApplicationsResponse.NextToken != nil {
			input.Token = apps.ApplicationsResponse.NextToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getPinpointApp(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	appId := d.EqualsQuals["id"].GetStringValue()
	// Empty check
	if appId == "" {
		return nil, nil
	}

	// Create service
	svc, err := PinpointClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pinpoint_app.getPinpointApp", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &pinpoint.GetAppInput{
		ApplicationId: &appId,
	}

	op, err := svc.GetApp(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pinpoint_app.getPinpointApp", "api_error", err)
		return nil, err
	}

	if op == nil {
		return nil, nil
	}

	return *op.ApplicationResponse, nil
}

func getPinpointApplicationSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	application := h.Item.(types.ApplicationResponse)

	// Create service
	svc, err := PinpointClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pinpoint_app.getPinpointApplicationSettings", "connection_error", err)
		return nil, err
	}
	params := &pinpoint.GetApplicationSettingsInput{
		ApplicationId: application.Id,
	}

	op, err := svc.GetApplicationSettings(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_pinpoint_app.getPinpointApplicationSettings", "api_error", err)
		return nil, err
	}
	if op == nil {
		return nil, nil
	}

	return op.ApplicationSettingsResource, nil
}
