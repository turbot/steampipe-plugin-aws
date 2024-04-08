package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpcEgressOnlyIGW(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_egress_only_internet_gateway",
		Description: "AWS VPC Egress Only Internet Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidEgressOnlyInternetGatewayId.NotFound", "InvalidEgressOnlyInternetGatewayId.Malformed"}),
			},
			Hydrate: getVpcEgressOnlyInternetGateway,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeEgressOnlyInternetGateways"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEgressOnlyInternetGateways,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeEgressOnlyInternetGateways"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the egress-only internet gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EgressOnlyInternetGatewayId"),
			},
			{
				Name:        "attachments",
				Description: "Information about the attachment of the egress-only internet gateway.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to egress only internet gateway.",
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

//// LIST FUNCTION

func listVpcEgressOnlyInternetGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_egress_only_internet_gateway.listVpcEgressOnlyInternetGateways", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(255)

	// Reduce the basic request limit down if the user has only requested a small number of rows
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
	input := &ec2.DescribeEgressOnlyInternetGatewaysInput{
		MaxResults: &maxLimit,
	}

	paginator := ec2.NewDescribeEgressOnlyInternetGatewaysPaginator(svc, input, func(o *ec2.DescribeEgressOnlyInternetGatewaysPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_egress_only_internet_gateway.listVpcEgressOnlyInternetGateways", "api_error", err)
			return nil, err
		}

		for _, egressOnlyInternetGateway := range output.EgressOnlyInternetGateways {
			d.StreamListItem(ctx, egressOnlyInternetGateway)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVpcEgressOnlyInternetGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	gatewayID := d.EqualsQuals["id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_egress_only_internet_gateway.getVpcEgressOnlyInternetGateway", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeEgressOnlyInternetGatewaysInput{
		EgressOnlyInternetGatewayIds: []string{gatewayID},
	}

	// Get call
	op, err := svc.DescribeEgressOnlyInternetGateways(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_egress_only_internet_gateway.getVpcEgressOnlyInternetGateway", "api_error", err)
		return nil, err
	}

	if op.EgressOnlyInternetGateways != nil && len(op.EgressOnlyInternetGateways) > 0 {
		h.Item = op.EgressOnlyInternetGateways[0]
		return op.EgressOnlyInternetGateways[0], nil
	}
	return nil, nil
}

func getVpcEgressOnlyInternetGatewayTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region := d.EqualsQualString(matrixKeyRegion)
	egw := h.Item.(types.EgressOnlyInternetGateway)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_egress_only_internet_gateway.getVpcEgressOnlyInternetGatewayTurbotAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:egress-only-internet-gateway/%s", commonColumnData.Partition, region, commonColumnData.AccountId, *egw.EgressOnlyInternetGatewayId)

	return []string{arn}, nil
}

//// TRANSFORM FUNCTIONS

func egressOnlyIGWApiDataToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	egw := d.HydrateItem.(types.EgressOnlyInternetGateway)
	param := d.Param.(string)

	// Get resource title
	title := egw.EgressOnlyInternetGatewayId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if len(egw.Tags) > 0 {
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
