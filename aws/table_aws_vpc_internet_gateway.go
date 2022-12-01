package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsVpcInternetGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_internet_gateway",
		Description: "AWS VPC Internet Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("internet_gateway_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInternetGatewayID.NotFound", "InvalidInternetGatewayID.Malformed"}),
			},
			Hydrate: getVpcInternetGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcInternetGateways,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "internet_gateway_id",
				Description: "The ID of the internet gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the internet gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attachments",
				Description: "Any VPCs attached to the internet gateway.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "tags assigned to the internet gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getVpcInternetGatewayTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getVpcInternetGatewayTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcInternetGatewayTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcInternetGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_internet_gateway.listVpcInternetGateways", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = int32(5)
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeInternetGatewaysInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeInternetGatewaysPaginator(svc, input, func(o *ec2.DescribeInternetGatewaysPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_internet_gateway.listVpcInternetGateways", "api_error", err)
			return nil, err
		}

		for _, items := range output.InternetGateways {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcInternetGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	internetGatewayID := d.KeyColumnQuals["internet_gateway_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_internet_gateway.getVpcInternetGateway", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []string{internetGatewayID},
	}

	// Get call
	op, err := svc.DescribeInternetGateways(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_internet_gateway.getVpcInternetGateway", "api_error", err)
		return nil, err
	}

	if op.InternetGateways != nil && len(op.InternetGateways) > 0 {
		h.Item = op.InternetGateways[0]
		return op.InternetGateways[0], nil
	}
	return nil, nil
}

func getVpcInternetGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	internetGateway := h.Item.(types.InternetGateway)

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_internet_gateway.getVpcInternetGatewayTurbotAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":internet-gateway/" + *internetGateway.InternetGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcInternetGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	internetGateway := d.HydrateItem.(types.InternetGateway)
	param := d.Param.(string)

	// Get resource title
	title := internetGateway.InternetGatewayId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if internetGateway.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range internetGateway.Tags {
			turbotTagsMap[*i.Key] = *i.Value
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	if param == "Tags" {
		return turbotTagsMap, nil
	}

	return title, nil
}
