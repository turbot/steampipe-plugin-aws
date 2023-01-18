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

func tableAwsVpcDhcpOptions(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_dhcp_options",
		Description: "AWS VPC DHCP Options",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("dhcp_options_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidDhcpOptionID.NotFound"}),
			},
			Hydrate: getVpcDhcpOption,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcDhcpOptions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "dhcp_options_id",
				Description: "The ID of the set of DHCP options.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the DHCP options set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_name",
				Description: "The domain name for instances. This value is used to complete unqualified DNS hostnames.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "domain-name"),
			},
			{
				Name:        "domain_name_servers",
				Description: "The IP addresses of up to four domain name servers, or AmazonProvidedDNS.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "domain-name-servers"),
			},
			{
				Name:        "netbios_name_servers",
				Description: "The IP addresses of up to four NetBIOS name servers.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "netbios-name-servers"),
			},
			{
				Name:        "netbios_node_type",
				Description: "The NetBIOS node type (1, 2, 4, or 8).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "netbios-node-type"),
			},
			{
				Name:        "ntp_servers",
				Description: "The IP addresses of up to four Network Time Protocol (NTP) servers.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "ntp-servers"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to vpc dhcp options.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(vpcDhcpOptionsAPIDataToTurbotData, "Title"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(vpcDhcpOptionsAPIDataToTurbotData, "Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcDhcpOptionAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcDhcpOptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_dhcp_options.listVpcDhcpOptions", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
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

	input := &ec2.DescribeDhcpOptionsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "owner_id", FilterName: "owner-id", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeDhcpOptionsPaginator(svc, input, func(o *ec2.DescribeDhcpOptionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_dhcp_options.listVpcDhcpOptions", "api_error", err)
			return nil, err
		}

		for _, items := range output.DhcpOptions {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVpcDhcpOption(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	dhcpOptionsID := d.KeyColumnQuals["dhcp_options_id"].GetStringValue()

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_dhcp_options.getVpcDhcpOption", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []string{dhcpOptionsID},
	}

	// get call
	items, err := svc.DescribeDhcpOptions(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_dhcp_options.getVpcDhcpOption", "api_error", err)
		return nil, err
	}

	for _, item := range items.DhcpOptions {
		if *item.DhcpOptionsId == dhcpOptionsID {
			return item, nil
		}
	}
	return nil, nil
}

func getVpcDhcpOptionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	dhcpOption := h.Item.(types.DhcpOptions)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_dhcp_options.getVpcDhcpOptionAkas", "common_data_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":dhcp-options/" + *dhcpOption.DhcpOptionsId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func dhcpConfigurationToStringSlice(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	dhcpConfigurations := d.Value.([]types.DhcpConfiguration)

	var values []*string
	for _, configuration := range dhcpConfigurations {
		if *configuration.Key == d.Param.(string) {
			values = mapString(configuration.Values)
		}
	}
	return values, nil
}

func vpcDhcpOptionsAPIDataToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpcDhcpOptions := d.HydrateItem.(types.DhcpOptions)
	param := d.Param.(string)

	// Get resource title
	title := *vpcDhcpOptions.DhcpOptionsId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if vpcDhcpOptions.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range vpcDhcpOptions.Tags {
			turbotTagsMap[*i.Key] = *i.Value
			if *i.Key == "Name" {
				title = *i.Value
			}
		}
	}

	if param == "Tags" {
		return turbotTagsMap, nil
	}

	return title, nil
}

func mapString(l []types.AttributeValue) []*string {
	var values []*string
	for _, v := range l {
		values = append(values, v.Value)
	}
	return values
}
