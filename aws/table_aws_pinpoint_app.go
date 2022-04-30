package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsPinpointApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_pinpoint_app",
		Description: "AWS Pinpoint App",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			Hydrate:           getPinpointApp,
		},
		List: &plugin.ListConfig{
			Hydrate: listPinpointApps,
		},
		GetMatrixItem: BuildPinpointRegionList,
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
	plugin.Logger(ctx).Trace("listPinpointApps")

	// Create Session
	svc, err := PinpointService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &pinpoint.GetAppsInput{
		PageSize: aws.String("100"),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < 100 {
			if *limit < 1 {
				input.PageSize = aws.String("1")
			} else {
				input.PageSize = aws.String(string(*limit))
			}
		}
	}

	for {
		apps, err := svc.GetApps(input)
		if err != nil {
			plugin.Logger(ctx).Error("listPinpointApps", "list", err)
			return nil, err
		}

		if apps.ApplicationsResponse == nil {
			break
		}

		for _, app := range apps.ApplicationsResponse.Item {
			d.StreamListItem(ctx, app)

			// Check if context has been cancelled or if the limit has been reached (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if apps.ApplicationsResponse.NextToken != nil {
			input.Token = apps.ApplicationsResponse.NextToken
		} else {
			break
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getPinpointApp(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getPinpointApp")
	appId := d.KeyColumnQuals["id"].GetStringValue()
	plugin.Logger(ctx).Error("getPinpointApp", "ApplicationId", appId)

	if appId == "" {
		return nil, nil
	}

	// Create service
	svc, err := PinpointService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getPinpointApp", "connection", err)
		return nil, err
	}

	params := &pinpoint.GetAppInput{}
	params.SetApplicationId(appId)

	op, err := svc.GetApp(params)
	plugin.Logger(ctx).Error("getPinpointApp", "Ok", 1111111)
	if err != nil {
		plugin.Logger(ctx).Error("getPinpointApp", "get", err)
		return nil, err
	}
	plugin.Logger(ctx).Error("getPinpointApp", "Ok", 222222)
	plugin.Logger(ctx).Error("getPinpointApp", "Ok", *op.ApplicationResponse.Arn)
	if op == nil {
		return nil, nil
	}

	return op.ApplicationResponse, nil
}

//// TRANSFORM FUNCTIONS

// func getPinpointAppTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	app := d.HydrateItem.(*pinpoint.ApplicationResponse)

// 	if app.Tags != nil {
// 		turbotTagsMap := map[string]string{}
// 		for _, i := range app.Tags {
// 			turbotTagsMap[*i.Key] = *i.Value
// 		}
// 		return turbotTagsMap, nil
// 	}
// 	return nil, nil
// }
