package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpc(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc",
		Description: "AWS VPC",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vpc_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException", "InvalidVpcID.NotFound"}),
			},
			Hydrate: getVpc,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVpcs"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcs,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVpcs"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cidr_block", Require: plugin.Optional},
				{Name: "dhcp_options_id", Require: plugin.Optional},
				{Name: "is_default", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
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

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc.listVpcs", "connection error", err)
		return nil, err
	}

	input := &ec2.DescribeVpcsInput{
		MaxResults: aws.Int32(1000),
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

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				input.MaxResults = aws.Int32(5)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}
	paginator := ec2.NewDescribeVpcsPaginator(svc, input, func(o *ec2.DescribeVpcsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc.listVpcs", "api_error", err)
			return nil, err
		}

		for _, items := range output.Vpcs {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpc(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	vpcID := d.EqualsQuals["vpc_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc.getVpc", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	}

	// Get call
	op, err := svc.DescribeVpcs(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc.getVpc", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.Vpcs) > 0 {
		return op.Vpcs[0], nil
	}
	return nil, nil
}

func getVpcARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpc := h.Item.(types.Vpc)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc.getVpcARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":vpc/" + *vpc.VpcId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getVpcTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpc := d.HydrateItem.(types.Vpc)

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
	vpc := d.HydrateItem.(types.Vpc)

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

// Build vpc resources filter parameter

type VpcFilterKeyMap struct {
	ColumnName string
	FilterName string
	ColumnType string
}

func buildVpcResourcesFilterParameter(keyMap []VpcFilterKeyMap, quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)
	for _, keyDetail := range keyMap {
		if quals[keyDetail.ColumnName] != nil {
			filter := &types.Filter{
				Name: aws.String(keyDetail.FilterName),
			}

			value := getQualsValueByColumn(quals, keyDetail.ColumnName, keyDetail.ColumnType)
			switch keyDetail.ColumnType {
			case "string":
				val, ok := value.(string)
				if ok {
					filter.Values = []string{val}
				}
			case "int64":
				val, ok := value.(int64)
				if ok {
					filter.Values = []string{fmt.Sprint(val)}
				}
			case "double":
				val, ok := value.(float64)
				if ok {
					filter.Values = []string{fmt.Sprint(val)}
				}
			case "ipaddr":
				val, ok := value.(string)
				if ok {
					filter.Values = []string{fmt.Sprint(val)}
				}
			case "cidr": // Ip address with mask
				val, ok := value.(string)
				if ok {
					filter.Values = []string{fmt.Sprint(val)}
				}
			case "time":
				val, ok := value.(time.Time)
				if ok {
					v := val.Format(time.RFC3339)
					filter.Values = []string{fmt.Sprint(v)}
				}
			case "boolean":
				filter.Values = []string{fmt.Sprint(value)}
			}

			filters = append(filters, *filter)
		}
	}

	return filters
}
