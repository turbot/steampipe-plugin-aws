package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
	"github.com/aws/aws-sdk-go-v2/service/quicksight/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsQuickSightDataset(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_data_set",
		Description: "AWS QuickSight Dataset",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "dataset_id", Require: plugin.Required},
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Hydrate: getAwsQuickSightDataset,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "quicksight", "action": "DescribeDataSet"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsQuickSightDatasets,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "quicksight", "action": "ListDataSets"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "A display name for the dataset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dataset_id",
				Description: "The ID of the dataset.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataSetId"),
			},
			// As we have already a column "account_id" as a common column for all the tables, we have renamed the column to "quicksight_account_id"
			{
				Name:        "quicksight_account_id",
				Description: "The Amazon Web Services account ID where the dataset is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("quicksight_account_id"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time that this dataset was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_time",
				Description: "The last time that this dataset was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "import_mode",
				Description: "A value that indicates whether you want to import the data into SPICE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "row_level_permission_data_set",
				Description: "The row-level security configuration for the dataset.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RowLevelPermissionDataSet"),
			},
			{
				Name:        "column_groups",
				Description: "Groupings of columns that work together in certain Amazon QuickSight features.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQuickSightDataset,
			},
			{
				Name:        "column_level_permission_rules",
				Description: "A set of one or more definitions of a ColumnLevelPermissionRule.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQuickSightDataset,
			},
			{
				Name:        "data_set_usage_configuration",
				Description: "The usage configuration to apply to child datasets that reference this dataset as a source.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQuickSightDataset,
			},
			{
				Name:        "logical_table_map",
				Description: "Configures the combination and transformation of the data from the physical tables.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQuickSightDataset,
			},
			{
				Name:        "output_columns",
				Description: "Output columns for the dataset.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQuickSightDataset,
			},
			{
				Name:        "physical_table_map",
				Description: "Declares the physical tables that are available in the underlying data sources.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQuickSightDataset,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsQuickSightDatasets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_data_set.listAwsQuickSightDatasets", "connection_error", err)
		return nil, err
	}

	accountId := d.EqualsQuals["quicksight_account_id"].GetStringValue()

	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	if accountId == "" {
		accountId = commonColumnData.AccountId
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &quicksight.ListDataSetsInput{
		AwsAccountId: aws.String(accountId),
		MaxResults:   aws.Int32(maxLimit),
	}

	paginator := quicksight.NewListDataSetsPaginator(svc, input, func(o *quicksight.ListDataSetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_data_set.listAwsQuickSightDatasets", "api_error", err)
			return nil, err
		}

		for _, item := range output.DataSetSummaries {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsQuickSightDataset(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_data_set.getAwsQuickSightDataset", "connection_error", err)
		return nil, err
	}

	var datasetID string
	if h.Item != nil {
		datasetID = *h.Item.(types.DataSetSummary).DataSetId
	} else {
		datasetID = d.EqualsQuals["dataset_id"].GetStringValue()
	}

	accountId := d.EqualsQuals["quicksight_account_id"].GetStringValue()
	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	if accountId == "" {
		accountId = commonColumnData.AccountId
	}

	// Build the params
	params := &quicksight.DescribeDataSetInput{
		AwsAccountId: aws.String(accountId),
		DataSetId:    aws.String(datasetID),
	}

	// Get call
	data, err := svc.DescribeDataSet(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_data_set.getAwsQuickSightDataset", "api_error", err)
		return nil, err
	}

	return *data.DataSet, nil
}
