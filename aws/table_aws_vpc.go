package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableAwsVpc(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc",
		Description: "AWS VPC",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("vpc_id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException", "InvalidVpcID.NotFound"}),
			Hydrate:           getVpc,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cidr_block", Require: plugin.Optional},
				{Name: "dhcp_options_id", Require: plugin.Optional},
				{Name: "is_default", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the vpc.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "cidr_block",
				Description: "The primary IPv4 CIDR block for the VPC.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "state",
				Description: "Contains the current state of the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_default",
				Description: "Indicates whether the VPC is the default VPC.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "dhcp_options_id",
				Description: "Contains the ID of the set of DHCP options, associated with the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_tenancy",
				Description: "The allowed tenancy of instances launched into the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "Contains ID of the AWS account that owns the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_block_association_set",
				Description: "Information about the IPv4 CIDR blocks associated with the VPC.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ipv6_cidr_block_association_set",
				Description: "Information about the IPv6 CIDR blocks associated with the VPC.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached with the VPC.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getVpcTurbotTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcs", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeVpcsInput{
		MaxResults: aws.Int64(1000),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "cidr_block", FilterName: "cidr", ColumnType: "cidr"},
		{ColumnName: "dhcp_options_id", FilterName: "dhcp-options-id", ColumnType: "string"},
		{ColumnName: "is_default", FilterName: "is-default", ColumnType: "boolean"},
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
		{ColumnName: "state", FilterName: "state", ColumnType: "string"},
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
	err = svc.DescribeVpcsPages(
		input,
		func(page *ec2.DescribeVpcsOutput, isLast bool) bool {
			for _, vpc := range page.Vpcs {
				d.StreamListItem(ctx, vpc)

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

func getVpc(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpc")

	vpcID := d.KeyColumnQuals["vpc_id"].GetStringValue()
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace(" getVpc", "AWS_REGION", region)

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcsInput{
		VpcIds: []*string{&vpcID},
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

func getVpcARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcARN")
	vpc := h.Item.(*ec2.Vpc)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":vpc/" + *vpc.VpcId

	return arn, nil
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

	if vpc.Tags != nil {
		for _, i := range vpc.Tags {
			if *i.Key == "Name" {
				return *i.Value, nil
			}
		}
	}

	return vpc.VpcId, nil
}

//// UTILITY FUNCTION

//Build vpc resources filter parameter

type VpcFilterKeyMap struct {
	ColumnName string
	FilterName string
	ColumnType string
}

func buildVpcResourcesFilterParameter(keyMap []VpcFilterKeyMap, quals plugin.KeyColumnQualMap) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)
	for _, keyDetail := range keyMap {
		if quals[keyDetail.ColumnName] != nil {
			filter := &ec2.Filter{
				Name: aws.String(keyDetail.FilterName),
			}

			value := getQualsValueByColumn(quals, keyDetail.ColumnName, keyDetail.ColumnType)
			switch keyDetail.ColumnType {
			case "string":
				val, ok := value.(string)
				if ok {
					filter.Values = []*string{&val}
				} else {
					filter.Values = value.([]*string)
				}
			case "int64":
				val, ok := value.(int64)
				if ok {
					filter.Values = []*string{aws.String(fmt.Sprint(val))}
				} else {
					filter.Values = value.([]*string)
				}
			case "double":
				val, ok := value.(float64)
				if ok {
					filter.Values = []*string{aws.String(fmt.Sprint(val))}
				} else {
					filter.Values = value.([]*string)
				}
			case "ipaddr":
				val, ok := value.(string)
				if ok {
					filter.Values = []*string{&val}
				} else {
					filter.Values = value.([]*string)
				}
			case "cidr": // Ip address with mask
				val, ok := value.(string)
				if ok {
					filter.Values = []*string{&val}
				} else {
					filter.Values = value.([]*string)
				}
			case "time":
				val, ok := value.(time.Time)
				if ok {
					v := val.Format(time.RFC3339) 
					filter.Values = []*string{&v}
				} else {
					filter.Values = value.([]*string)
				}
			case "boolean":
				filter.Values = []*string{aws.String(fmt.Sprint(value))}
			}

			filters = append(filters, filter)
		}
	}

	return filters
}
