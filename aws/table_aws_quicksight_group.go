package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
	"github.com/aws/aws-sdk-go-v2/service/quicksight/types"
	"github.com/aws/smithy-go"

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
			KeyColumns: []*plugin.KeyColumn{
				{Name:    "group_name", Require: plugin.Required},
				{Name:    "region", Require: plugin.Required},
				{Name:    "namespace", Require: plugin.Required},
				{Name:    "quicksight_account_id", Require: plugin.Optional},
			},
			Hydrate: getAwsQuickSightGroup,
			Tags:    map[string]string{"service": "quicksight", "action": "DescribeGroup"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsQuickSightNamespaces,
			Hydrate:       listAwsQuickSightGroups,
			Tags:          map[string]string{"service": "quicksight", "action": "ListGroups"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "quicksight_account_id", Require: plugin.Optional},
				{Name: "namespace", Require: plugin.Optional},
				{Name: "region", Require: plugin.Required},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "group_name",
				Description: "The name of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the group.",
				Type:        proto.ColumnType_STRING,
			},
			// As we have already a column "account_id" as a common column for all the tables, we have renamed the column to "quicksight_account_id"
			{
				Name:        "quicksight_account_id",
				Description: "The ID for the Amazon Web Services account that the group is in. Currently, you use the ID for the Amazon Web Services account that contains your Amazon QuickSight account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("quicksight_account_id"),
			},
			{
				Name:        "description",
				Description: "The group description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "principal_id",
				Description: "The principal ID of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The namespace. Currently, you should set this to default.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_members",
				Description: "The members of the group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQuickSightGroupMembers,
				Transform:   transform.FromValue(),
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupName"),
			},
		}),
	}
}

type QuickSightGroup struct {
	types.Group
	Namespace string
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

	input := &quicksight.ListGroupsInput{
		AwsAccountId: aws.String(accountId),
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
			//In the case of parent hydrate use the ignore config is not working as expected in the list config.
			//So we are using this workaround to ignore the ResourceNotFoundException
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "ResourceNotFoundException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_quicksight_group.listAwsQuickSightGroups", "api_error", err)
			return nil, err
		}

		for _, item := range output.GroupList {
			d.StreamListItem(ctx, QuickSightGroup{Group: item, Namespace: *namespaceInfo.Name})

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

	// Empty check for required parameters
	if groupName == "" || namespace == "" {
		return nil, nil
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

	params := &quicksight.DescribeGroupInput{
		AwsAccountId: aws.String(accountId),
		Namespace:    aws.String(namespace),
		GroupName:    aws.String(groupName),
	}

	// Get call
	data, err := svc.DescribeGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_group.getAwsQuickSightGroup", "api_error", err)
		return nil, err
	}

	return QuickSightGroup{Group: *data.Group, Namespace: namespace}, nil
}

func getAwsQuickSightGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(QuickSightGroup)

	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_group.getAwsQuickSightGroupMembers", "connection_error", err)
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

	params := &quicksight.ListGroupMembershipsInput{
		AwsAccountId: aws.String(accountId),
		Namespace:    aws.String(group.Namespace),
		GroupName:    group.GroupName,
	}

	// Get call
	paginator := quicksight.NewListGroupMembershipsPaginator(svc, params, func(o *quicksight.ListGroupMembershipsPaginatorOptions) {
		o.Limit = int32(100)
		o.StopOnDuplicateToken = true
	})

	var data []types.GroupMember
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_group.getAwsQuickSightGroupMembers", "api_error", err)
			return nil, err
		}
		data = append(data, output.GroupMemberList...)
	}

	return data, nil
}
