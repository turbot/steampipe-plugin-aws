package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsQuickSightNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_namespace",
		Description: "AWS QuickSight Namespace",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Required},
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Hydrate: getAwsQuickSightNamespace,
			Tags:    map[string]string{"service": "quicksight", "action": "DescribeNamespace"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsQuickSightNamespaces,
			KeyColumns: []*plugin.KeyColumn{
				// Namespaces can only be created in the QuickSight identity region.
				// When listing namespaces from other regions, the API returns the same results,
				// which can lead to duplicate entries.
				// Simply requiring the "region" as a qualifier still allows incorrect results to be returned.
				//
				// Example: For the query "SELECT * FROM aws_quicksight_namespace WHERE region = 'us-east-1';"
				// (where 'us-east-1' is a non-identity region), the API still returns namespaces that exist in the identity region, in my case identity region is "ap-south-1".
				//
				// Instead, we’ve added a check in the list function to return results only when the queried region
				// matches the QuickSight capacity (identity) region.
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "quicksight", "action": "ListNamespaces"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			// As we have already a column "account_id" as a common column for all the tables, we have renamed the column to "quicksight_account_id"
			{
				Name:        "quicksight_account_id",
				Description: "The ID for the Amazon Web Services account that contains the Amazon QuickSight namespace that you want to describe.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("quicksight_account_id"),
			},
			{
				Name:        "capacity_region",
				Description: "The region that hosts the namespace's capacity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_status",
				Description: "The creation status of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_store",
				Description: "The identity store used for the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace_error",
				Description: "The error type and message for a namespace error.",
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

func listAwsQuickSightNamespaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_namespace.listAwsQuickSightNamespaces", "connection_error", err)
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

	input := &quicksight.ListNamespacesInput{
		AwsAccountId: aws.String(accountId),
		MaxResults:   aws.Int32(maxLimit),
	}

	paginator := quicksight.NewListNamespacesPaginator(svc, input, func(o *quicksight.ListNamespacesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_namespace.listAwsQuickSightNamespaces", "api_error", err)
			return nil, err
		}

		for _, item := range output.Namespaces {
			// The API returns the same result regardless of the region.
			// Adding a check for the capacity region helps prevent duplicate entries.
			if d.EqualsQuals[matrixKeyRegion].GetStringValue() == *item.CapacityRegion {
				d.StreamListItem(ctx, item)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsQuickSightNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_namespace.getAwsQuickSightNamespace", "connection_error", err)
		return nil, err
	}

	namespaceName := d.EqualsQuals["name"].GetStringValue()

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

	params := &quicksight.DescribeNamespaceInput{
		AwsAccountId: aws.String(accountId),
		Namespace:    aws.String(namespaceName),
	}

	// Get call
	data, err := svc.DescribeNamespace(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_namespace.getAwsQuickSightNamespace", "api_error", err)
		return nil, err
	}

	return *data.Namespace, nil
}
