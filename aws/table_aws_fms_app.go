package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fms"
	"github.com/aws/aws-sdk-go-v2/service/fms/types"

	fmsv1 "github.com/aws/aws-sdk-go/service/fms"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsFMSAppList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_fms_app_list",
		Description: "AWS FMS App List",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("list_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getFmsAppList,
		},
		List: &plugin.ListConfig{
			Hydrate: listFmsAppLists,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(fmsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "list_name",
				Description: "The name of the applications list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ListName", "AppsList.ListName"),
			},
			{
				Name:        "list_id",
				Description: "The ID of the applications list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ListId", "AppsList.ListId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the applications list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ListArn", "AppsListArn"),
			},
			{
				Name:        "create_time",
				Description: "The time that the Firewall Manager applications list was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getFmsAppList,
			},
			{
				Name:        "last_update_time",
				Description: "The time that the Firewall Manager applications list was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getFmsAppList,
			},
			{
				Name:        "list_update_token",
				Description: "A unique identifier for each update to the list. When you update the list, the update token must match the token of the current version of the application list.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFmsAppList,
			},
			{
				Name:        "previous_apps_list",
				Description: "A map of previous version numbers to their corresponding App object arrays.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFmsAppList,
			},
			{
				Name:        "apps_list",
				Description: "An array of applications in the Firewall Manager applications list.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFmsAppList,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ListName", "AppsList.ListName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ListArn", "AppsListArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listFmsAppLists(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := FMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_app_list.listFmsAppLists", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(100)
	input := fms.ListAppsListsInput{
		MaxResults: aws.Int32(maxItems),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	paginator := fms.NewListAppsListsPaginator(svc, &input, func(o *fms.ListAppsListsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_fms_app_list.listFmsAppLists", "api_error", err)
			return nil, err
		}

		for _, app := range output.AppsLists {
			d.StreamListItem(ctx, app)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getFmsAppList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	listId := ""

	if h.Item != nil {
		data := h.Item.(types.AppsListDataSummary)
		listId = *data.ListId
	} else {
		listId = d.EqualsQualString("list_id")
	}

	if listId == "" {
		return nil, nil
	}
	// Create service
	svc, err := FMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_app_list.getFmsAppList", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &fms.GetAppsListInput{
		ListId: &listId,
	}

	op, err := svc.GetAppsList(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_app_list.getFmsAppList", "api_error", err)
		return nil, err
	}

	if op != nil {
		return op, nil
	}

	return nil, nil
}
