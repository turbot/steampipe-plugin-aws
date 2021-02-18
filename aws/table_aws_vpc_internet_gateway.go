package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcInternetGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_internet_gateway",
		Description: "AWS VPC Internet Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("internet_gateway_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInternetGatewayID.NotFound", "InvalidInternetGatewayID.Malformed"}),
			ItemFromKey:       internetGatewayFromKey,
			Hydrate:           getVpcInternetGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcInternetGateways,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "internet_gateway_id",
				Description: "The ID of the internet gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the internet gateway",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attachments",
				Description: "Any VPCs attached to the internet gateway",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "tags assigned to the internet gateway",
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

//// ITEM FROM KEY

func internetGatewayFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	internetGatewayID := quals["internet_gateway_id"].GetStringValue()
	item := &ec2.InternetGateway{
		InternetGatewayId: &internetGatewayID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcInternetGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listVpcInternetGateways", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeInternetGatewaysPages(
		&ec2.DescribeInternetGatewaysInput{},
		func(page *ec2.DescribeInternetGatewaysOutput, isLast bool) bool {
			for _, internetGateway := range page.InternetGateways {
				d.StreamListItem(ctx, internetGateway)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcInternetGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcInternetGateway")
	internetGateway := h.Item.(*ec2.InternetGateway)
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{internetGateway.InternetGatewayId},
	}

	// Get call
	op, err := svc.DescribeInternetGateways(params)
	if err != nil {
		logger.Debug("[getVpcInternetGateway__", "ERROR", err)
		return nil, err
	}

	if op.InternetGateways != nil && len(op.InternetGateways) > 0 {
		h.Item = op.InternetGateways[0]
		return op.InternetGateways[0], nil
	}
	return nil, nil
}

func getVpcInternetGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcInternetGatewayTurbotAkas")
	internetGateway := h.Item.(*ec2.InternetGateway)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":internet-gateway/" + *internetGateway.InternetGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcInternetGatewayTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	internetGateway := d.HydrateItem.(*ec2.InternetGateway)
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
