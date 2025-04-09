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

func tableAwsQuickSightGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_group",
		Description: "AWS QuickSight Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"group_name", "namespace"}),
			Hydrate:    getAwsQuickSightGroup,
			Tags:       map[string]string{"service": "quicksight", "action": "DescribeGroup"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsQuickSightNamespaces,
			Hydrate:       listAwsQuickSightGroups,
			Tags:          map[string]string{"service": "quicksight", "action": "ListGroups"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "namespace", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "group_name",
				Description: "The name of the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn"),
			},
			{
				Name:        "description",
				Description: "The group description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "principal_id",
				Description: "The principal ID of the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrincipalId"),
			},
			{
				Name:        "namespace",
				Description: "The namespace. Currently, you should set this to default.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("namespace"),
				Default:     "default",
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsQuickSightGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_group.listAwsQuickSightGroups", "connection_error", err)
		return nil, err
	}

	// Get namespace from parent or quals
	namespaceInfo := h.Item.(types.NamespaceInfoV2)
	if d.EqualsQuals["namespace"] != nil && d.EqualsQuals["namespace"].GetStringValue() != *namespaceInfo.Name {
		return nil, nil
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

	input := &quicksight.ListGroupsInput{
		AwsAccountId: aws.String(commonColumnData.AccountId),
		Namespace:    namespaceInfo.Name,
		MaxResults:   aws.Int32(maxLimit),
	}

	paginator := quicksight.NewListGroupsPaginator(svc, input, func(o *quicksight.ListGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_group.listAwsQuickSightGroups", "api_error", err)
			return nil, err
		}

		for _, item := range output.GroupList {
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

func getAwsQuickSightGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_group.getAwsQuickSightGroup", "connection_error", err)
		return nil, err
	}

	groupName := d.EqualsQuals["group_name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// Default namespace is default
	if namespace == "" {
		namespace = "default"
	}

	// Get AWS Account ID
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	params := &quicksight.DescribeGroupInput{
		AwsAccountId: aws.String(commonColumnData.AccountId),
		Namespace:    aws.String(namespace),
		GroupName:    aws.String(groupName),
	}

	// Get call
	data, err := svc.DescribeGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_group.getAwsQuickSightGroup", "api_error", err)
		return nil, err
	}

	return *data.Group, nil
}
