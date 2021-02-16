package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcDhcpOptions(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "aws_vpc_dhcp_options",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("dhcp_options_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidDhcpOptionID.NotFound"}),
			ItemFromKey:       dhcpOptionFromKey,
			Hydrate:           getVpcDhcpOption,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcDhcpOptions,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "dhcp_options_id",
				Description: "The ID of the set of DHCP options",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the DHCP options set",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_name",
				Description: "The domain name for instances. This value is used to complete unqualified DNS hostnames",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "domain-name"),
			},
			{
				Name:        "domain_name_servers",
				Description: "The IP addresses of up to four domain name servers, or AmazonProvidedDNS",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "domain-name-servers"),
			},
			{
				Name:        "netbios_name_servers",
				Description: "The IP addresses of up to four NetBIOS name servers",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "netbios-name-servers"),
			},
			{
				Name:        "netbios_node_type",
				Description: "The NetBIOS node type (1, 2, 4, or 8)",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "netbios-node-type"),
			},
			{
				Name:        "ntp_servers",
				Description: "The IP addresses of up to four Network Time Protocol (NTP) servers",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DhcpConfigurations").TransformP(dhcpConfigurationToStringSlice, "ntp-servers"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to vpc dhcp options",
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

//// ITEM FROM KEY

func dhcpOptionFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	dhcpOptionsID := quals["dhcp_options_id"].GetStringValue()
	item := &ec2.DhcpOptions{
		DhcpOptionsId: &dhcpOptionsID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcDhcpOptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listVpcDhcpOptions", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.DescribeDhcpOptionsPages(
		&ec2.DescribeDhcpOptionsInput{},
		func(page *ec2.DescribeDhcpOptionsOutput, lastPage bool) bool {
			for _, item := range page.DhcpOptions {
				plugin.Logger(ctx).Trace("listVpcDhcpOptions", "Data", item)
				d.StreamListItem(ctx, item)
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcDhcpOption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcDhcpOption")
	dhcpOptions := h.Item.(*ec2.DhcpOptions).DhcpOptionsId
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []*string{dhcpOptions},
	}

	// get call
	items, err := svc.DescribeDhcpOptions(params)
	if err != nil {
		logger.Debug("getVpcDhcpOption__", "ERROR", err)
		return nil, err
	}

	for _, item := range items.DhcpOptions {
		if *item.DhcpOptionsId == *dhcpOptions {
			return item, nil
		}
	}
	return nil, nil
}

func getVpcDhcpOptionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcDhcpOptionAkas")
	dhcpOption := h.Item.(*ec2.DhcpOptions)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":dhcp-options/" + *dhcpOption.DhcpOptionsId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func dhcpConfigurationToStringSlice(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	dhcpConfigurations := d.Value.([]*ec2.DhcpConfiguration)

	var values []*string
	for _, configuration := range dhcpConfigurations {
		if *configuration.Key == d.Param.(string) {
			values = mapString(configuration.Values)
		}
	}
	return values, nil
}

func vpcDhcpOptionsAPIDataToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpcDhcpOptions := d.HydrateItem.(*ec2.DhcpOptions)
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

func mapString(l []*ec2.AttributeValue) []*string {
	var values []*string
	for _, v := range l {
		values = append(values, v.Value)
	}
	return values
}
