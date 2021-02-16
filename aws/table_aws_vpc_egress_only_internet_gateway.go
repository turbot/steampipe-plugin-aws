package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcEgressOnlyIGW(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_egress_only_internet_gateway",
		Description: "AWS VPC Egress Only Internet Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidEgressOnlyInternetGatewayId.NotFound", "InvalidEgressOnlyInternetGatewayId.Malformed"}),
			ItemFromKey:       egressOnlyIGWFromKey,
			Hydrate:           getVpcEgressOnlyInternetGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEgressOnlyInternetGateways,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the egress-only internet gateway",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EgressOnlyInternetGatewayId"),
			},
			{
				Name:        "attachments",
				Description: "Information about the attachment of the egress-only internet gateway",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to egress only internet gateway",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(egressOnlyIGWApiDataToTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(egressOnlyIGWApiDataToTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcEgressOnlyInternetGatewayTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// ITEM FROM KEY

func egressOnlyIGWFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	gatewayID := quals["id"].GetStringValue()
	item := &ec2.EgressOnlyInternetGateway{
		EgressOnlyInternetGatewayId: &gatewayID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcEgressOnlyInternetGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listVpcEgressOnlyInternetGateways", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeEgressOnlyInternetGatewaysPages(
		&ec2.DescribeEgressOnlyInternetGatewaysInput{},
		func(page *ec2.DescribeEgressOnlyInternetGatewaysOutput, isLast bool) bool {
			for _, egressOnlyInternetGateway := range page.EgressOnlyInternetGateways {
				d.StreamListItem(ctx, egressOnlyInternetGateway)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcEgressOnlyInternetGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcEgressOnlyInternetGateway")
	subnet := h.Item.(*ec2.EgressOnlyInternetGateway)
	defaultRegion := GetDefaultRegion()

	// get service
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeEgressOnlyInternetGatewaysInput{
		EgressOnlyInternetGatewayIds: []*string{subnet.EgressOnlyInternetGatewayId},
	}

	// Get call
	op, err := svc.DescribeEgressOnlyInternetGateways(params)
	if err != nil {
		logger.Debug("getVpcEgressOnlyInternetGateway__", "ERROR", err)
		return nil, err
	}

	if op.EgressOnlyInternetGateways != nil && len(op.EgressOnlyInternetGateways) > 0 {
		h.Item = op.EgressOnlyInternetGateways[0]
		return op.EgressOnlyInternetGateways[0], nil
	}
	return nil, nil
}

func getVpcEgressOnlyInternetGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcEgressOnlyInternetGatewayTurbotAkas")
	egw := h.Item.(*ec2.EgressOnlyInternetGateway)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get resource aka
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":egress-only-internet-gateway/" + *egw.EgressOnlyInternetGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func egressOnlyIGWApiDataToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	egw := d.HydrateItem.(*ec2.EgressOnlyInternetGateway)
	param := d.Param.(string)

	// Get resource title
	title := egw.EgressOnlyInternetGatewayId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if egw.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range egw.Tags {
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
