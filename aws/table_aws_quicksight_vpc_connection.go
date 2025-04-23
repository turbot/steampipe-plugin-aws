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

func tableAwsQuickSightVpcConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_quicksight_vpc_connection",
		Description: "AWS QuickSight VPC Connection",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vpc_connection_id", Require: plugin.Required},
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Hydrate: getAwsQuickSightVpcConnection,
			Tags:    map[string]string{"service": "quicksight", "action": "DescribeVPCConnection"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsQuickSightVpcConnections,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "quicksight_account_id", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "quicksight", "action": "ListVPCConnections"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_QUICKSIGHT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "A display name for the VPC connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_connection_id",
				Description: "The ID of the VPC connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VPCConnectionId"),
			},
			// As we have already a column "account_id" as a common column for all the tables, we have renamed the column to "quicksight_account_id"
			{
				Name:        "quicksight_account_id",
				Description: "The Amazon Web Services account ID where the VPC connection is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("quicksight_account_id"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the VPC connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time that this VPC connection was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_time",
				Description: "The last time that this VPC connection was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The status of the VPC connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_status",
				Description: "The availability status of the VPC connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The ARN of the IAM role used to create the VPC connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_group_ids",
				Description: "The security group IDs used for the VPC connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dns_resolvers",
				Description: "The DNS resolvers used for the VPC connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interfaces",
				Description: "Information about the network interfaces for the VPC connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC associated with the connection.",
				Type:        proto.ColumnType_STRING,
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

func listAwsQuickSightVpcConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_vpc_connection.listAwsQuickSightVpcConnections", "connection_error", err)
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

	input := &quicksight.ListVPCConnectionsInput{
		AwsAccountId: aws.String(accountId),
		MaxResults:   aws.Int32(maxLimit),
	}

	paginator := quicksight.NewListVPCConnectionsPaginator(svc, input, func(o *quicksight.ListVPCConnectionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_quicksight_vpc_connection.listAwsQuickSightVpcConnections", "api_error", err)
			return nil, err
		}

		for _, item := range output.VPCConnectionSummaries {
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

func getAwsQuickSightVpcConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := QuickSightClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_vpc_connection.getAwsQuickSightVpcConnection", "connection_error", err)
		return nil, err
	}

	var vpcConnectionId string
	if h.Item != nil {
		vpcConnectionId = *h.Item.(types.VPCConnectionSummary).VPCConnectionId
	} else {
		vpcConnectionId = d.EqualsQuals["vpc_connection_id"].GetStringValue()
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
	params := &quicksight.DescribeVPCConnectionInput{
		AwsAccountId:    aws.String(accountId),
		VPCConnectionId: aws.String(vpcConnectionId),
	}

	// Get call
	data, err := svc.DescribeVPCConnection(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_quicksight_vpc_connection.getAwsQuickSightVpcConnection", "api_error", err)
		return nil, err
	}

	return *data.VPCConnection, nil
}
