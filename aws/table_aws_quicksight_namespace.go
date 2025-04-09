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
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getAwsQuickSightNamespace,
			Tags:       map[string]string{"service": "quicksight", "action": "DescribeNamespace"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsQuickSightNamespaces,
			Tags:    map[string]string{"service": "quicksight", "action": "ListNamespaces"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn"),
			},
			{
				Name:        "capacity_region",
				Description: "The region that hosts the namespace's capacity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CapacityRegion"),
			},
			{
				Name:        "creation_status",
				Description: "The creation status of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CreationStatus"),
			},
			{
				Name:        "identity_store",
				Description: "The identity store used for the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IdentityStore"),
			},
			{
				Name:        "iam_identity_center_application_arn",
				Description: "The ARN of the IAM Identity Center application that is integrated with Amazon QuickSight.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamIdentityCenterApplicationArn"),
			},
			{
				Name:        "iam_identity_center_instance_arn",
				Description: "The ARN of the IAM Identity Center instance that is integrated with Amazon QuickSight.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamIdentityCenterInstanceArn"),
			},
			{
				Name:        "namespace_error",
				Description: "The error type and message for a namespace error.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NamespaceError"),
			},
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

	input := &quicksight.ListNamespacesInput{
		AwsAccountId: aws.String(commonColumnData.AccountId),
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

func getAwsQuickSightNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_namespace.getAwsQuickSightNamespace", "connection_error", err)
		return nil, err
	}

	namespaceName := d.EqualsQuals["name"].GetStringValue()

	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	params := &quicksight.DescribeNamespaceInput{
		AwsAccountId: aws.String(commonColumnData.AccountId),
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
