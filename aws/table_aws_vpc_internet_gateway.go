package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableAwsVpcInternetGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_internet_gateway",
		Description: "AWS VPC Internet Gateway",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("internet_gateway_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidInternetGatewayID.NotFound", "InvalidInternetGatewayID.Malformed"}),
			Hydrate:           getVpcInternetGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcInternetGateways,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcInternetGateways", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeInternetGatewaysInput{
		MaxResults: aws.Int64(1000),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = aws.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeInternetGatewaysPages(
		input,
		func(page *ec2.DescribeInternetGatewaysOutput, isLast bool) bool {
			for _, internetGateway := range page.InternetGateways {
				d.StreamListItem(ctx, internetGateway)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcInternetGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcInternetGateway")

	region := d.KeyColumnQualString(matrixKeyRegion)
	internetGatewayID := d.KeyColumnQuals["internet_gateway_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(internetGatewayID)},
	}

	// Get call
	op, err := svc.DescribeInternetGateways(params)
	if err != nil {
		plugin.Logger(ctx).Debug("[getVpcInternetGateway__", "ERROR", err)
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	internetGateway := h.Item.(*ec2.InternetGateway)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":internet-gateway/" + *internetGateway.InternetGatewayId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcInternetGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
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
