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

func tableAwsQuickSightDatasource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_datasource",
		Description: "AWS QuickSight Data Source",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("data_source_id"),
			Hydrate:    getAwsQuickSightDatasource,
			Tags:       map[string]string{"service": "quicksight", "action": "DescribeDataSource"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsQuickSightDatasources,
			Tags:    map[string]string{"service": "quicksight", "action": "ListDataSources"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "A display name for the data source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_source_id",
				Description: "The ID of the data source.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataSourceId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the data source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time that this data source was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_time",
				Description: "The last time that this data source was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The status of the data source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the data source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_connection_properties",
				Description: "The VPC connection properties of the data source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ssl_properties",
				Description: "Secure Socket Layer (SSL) properties that apply when QuickSight connects to your underlying data source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "error_info",
				Description: "Error information from the last update or the creation of the data source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "data_source_parameters",
				Description: "The parameters that Amazon QuickSight uses to connect to your underlying data source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "alternate_data_source_parameters",
				Description: "A list of alternate data source parameters.",
				Type:        proto.ColumnType_JSON,
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

func listAwsQuickSightDatasources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_datasource.listAwsQuickSightDatasources", "connection_error", err)
		return nil, err
	}

	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &quicksight.ListDataSourcesInput{
		AwsAccountId: aws.String(commonColumnData.AccountId),
		MaxResults:   aws.Int32(maxLimit),
	}

	paginator := quicksight.NewListDataSourcesPaginator(svc, input, func(o *quicksight.ListDataSourcesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_datasource.listAwsQuickSightDatasources", "api_error", err)
			return nil, err
		}

		for _, item := range output.DataSources {
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

func getAwsQuickSightDatasource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_datasource.getAwsQuickSightDatasource", "connection_error", err)
		return nil, err
	}

	var dataSourceID string
	if h.Item != nil {
		dataSourceID = *h.Item.(types.DataSource).DataSourceId
	} else {
		dataSourceID = d.EqualsQuals["data_source_id"].GetStringValue()
	}

	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build the params
	params := &quicksight.DescribeDataSourceInput{
		AwsAccountId: aws.String(commonColumnData.AccountId),
		DataSourceId: aws.String(dataSourceID),
	}

	// Get call
	data, err := svc.DescribeDataSource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_datasource.getAwsQuickSightDatasource", "api_error", err)
		return nil, err
	}

	return *data.DataSource, nil
}
