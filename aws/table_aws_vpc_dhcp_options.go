package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcDhcpOptions(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_dhcp_options",
		Description: "AWS VPC DHCP Options",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("dhcp_options_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidDhcpOptionID.NotFound"}),
			Hydrate:           getVpcDhcpOption,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcDhcpOptions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcDhcpOptions", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeDhcpOptionsInput{
		MaxResults: aws.Int64(100),
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

	err = svc.DescribeDhcpOptionsPages(
		input,
		func(page *ec2.DescribeDhcpOptionsOutput, lastPage bool) bool {
			for _, item := range page.DhcpOptions {
				plugin.Logger(ctx).Trace("listVpcDhcpOptions", "Data", item)
				d.StreamListItem(ctx, item)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcDhcpOption(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcDhcpOption")

	region := d.KeyColumnQualString(matrixKeyRegion)
	dhcpOptionsID := d.KeyColumnQuals["dhcp_options_id"].GetStringValue()

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []*string{aws.String(dhcpOptionsID)},
	}

	// get call
	items, err := svc.DescribeDhcpOptions(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcDhcpOption__", "ERROR", err)
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
	plugin.Logger(ctx).Trace("getVpcDhcpOptionAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	dhcpOption := h.Item.(*ec2.DhcpOptions)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
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
