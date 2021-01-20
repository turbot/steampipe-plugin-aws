package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpc(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc",
		Description: "AWS VPC",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("vpc_id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       vpcFromKey,
			Hydrate:           getVpc,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcs,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_block",
				Description: "The primary IPv4 CIDR block for the VPC",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "state",
				Description: "Contains the current state of the VPC",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_default",
				Description: "Indicates whether the VPC is the default VPC",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "dhcp_options_id",
				Description: "Contains the ID of the set of DHCP options, associated with the VPC",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_tenancy",
				Description: "The allowed tenancy of instances launched into the VPC",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "Contains ID of the AWS account that owns the VPC",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_block_association_set",
				Description: "Information about the IPv4 CIDR blocks associated with the VPC",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ipv6_cidr_block_association_set",
				Description: "Information about the IPv6 CIDR blocks associated with the VPC",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_raw",
				Description: "A list of tags that are attached to the vpc",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Default:     transform.FromField("VpcId"),
				Transform:   transform.From(getVpcTurbotTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsVpcTurbotData,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// ITEM FROM KEY

func vpcFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	vpcID := quals["vpc_id"].GetStringValue()
	item := &ec2.Vpc{
		VpcId: &vpcID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("[TRACE] listVpcs", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeVpcsPages(
		&ec2.DescribeVpcsInput{},
		func(page *ec2.DescribeVpcsOutput, isLast bool) bool {
			for _, vpc := range page.Vpcs {
				d.StreamListItem(ctx, vpc)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpc(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpc")
	vpc := h.Item.(*ec2.Vpc)
	defaultRegion := GetDefaultRegion()

	// get service
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcsInput{
		VpcIds: []*string{vpc.VpcId},
	}

	// Get call
	op, err := svc.DescribeVpcs(params)
	if err != nil {
		logger.Debug("getVpc__", "ERROR", err)
		return nil, err
	}

	if op.Vpcs != nil && len(op.Vpcs) > 0 {
		return op.Vpcs[0], nil
	}
	return nil, nil
}

func getAwsVpcTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsVpcTurbotData")
	vpc := h.Item.(*ec2.Vpc)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":vpc/" + *vpc.VpcId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpc := d.HydrateItem.(*ec2.Vpc)
	var turbotTagsMap map[string]string
	if vpc.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range vpc.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

func getVpcTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpc := d.HydrateItem.(*ec2.Vpc)
	vpcData := d.HydrateResults
	var title string
	if vpc.Tags != nil {
		for _, i := range vpc.Tags {
			if *i.Key == "Name" {
				title = *i.Value
			}
		}
	}

	if title == "" {
		if vpcData["getVpc"] != nil {
			title = *vpcData["getVpc"].(*ec2.Vpc).VpcId
		} else {
			title = *vpcData["listVpcs"].(*ec2.Vpc).VpcId
		}
	}
	return title, nil
}
