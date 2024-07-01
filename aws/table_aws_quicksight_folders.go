package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"

	quicksightv1 "github.com/aws/aws-sdk-go/service/quicksight"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsQuicksightFolders(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_folders",
		Description: "AWS Quicksight Folders",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("folder_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException"}),
			},
			Hydrate: getQuicksightFolder,
			Tags:    map[string]string{"service": "quicksight", "action": "ListFolders"},
		},
		List: &plugin.ListConfig{
			Hydrate:    listQuicksightFolders,
			Tags:       map[string]string{"service": "quicksight", "action": "DescribeFolders"},
			KeyColumns: []*plugin.KeyColumn{
				//{Name: "account_id", Require: plugin.Optional},
			},
		},
		/*
			HydrateConfig: []plugin.HydrateConfig{
				{
					Func: listWorkspacesTags,
					Tags: map[string]string{"service": "workspaces", "action": "DescribeTags"},
				},
			},*/
		GetMatrixItemFunc: SupportedRegionMatrix(quicksightv1.EndpointsID),
		Columns: []*plugin.Column{
			{
				Name:        "folder_id",
				Description: "The id of the Folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the Folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The arn of the Folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "folder_type",
				Description: "The type of the Folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The last update time of the Folder.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_time",
				Description: "The creation time of the Folder.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "folder_type",
				Description: "The type of the Folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "folder_path",
				Description: "An array of ancestor ARN strings for the folder",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sharing_model",
				Description: "The sharing scope of the folder.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
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
		},
	}
}

//// LIST FUNCTION

func listQuicksightFolders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {


	// Create Session
	svc, err := QuicksightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_folders.listQuicksightFolders", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(25)
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

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	input := &quicksight.ListFoldersInput{
		MaxResults: aws.Int32(maxLimit),
		AwsAccountId: &commonColumnData.AccountId,
	}

	paginator := quicksight.NewListFoldersPaginator(svc, input, func(o *quicksight.ListFoldersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_folders.listQuicksightFolders", "api_error", err)
			return nil, err
		}

		for _, items := range output.FolderSummaryList {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getQuicksightFolder(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	FoldersId := d.EqualsQuals["folders_id"].GetStringValue()

	// check if folder id is empty
	if FoldersId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := QuicksightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_folders.getQuicksightFolders", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build the params
	params := &quicksight.DescribeFolderInput{
		FolderId: &FoldersId,
		AwsAccountId: &commonColumnData.AccountId,
	}

	// Get call
	data, err := svc.DescribeFolder(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_folders.getQuicksightFolder", "api_error", err)
		return nil, err
	}

	if data.Folder != nil {
		return data.Folder, nil
	}

	return nil, nil
}
