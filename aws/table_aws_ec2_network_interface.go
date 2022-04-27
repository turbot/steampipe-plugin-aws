package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2NetworkInterface(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_network_interface",
		Description: "AWS EC2 Network Interface",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("network_interface_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorWithContext([]string{"InvalidNetworkInterfaceID.NotFound", "InvalidNetworkInterfaceID.Unavailable", "InvalidNetworkInterfaceID.Malformed"}),
			},
			Hydrate: getEc2NetworkInterface,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2NetworkInterfaces,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "association_id", Require: plugin.Optional},
				{Name: "association_allocation_id", Require: plugin.Optional},
				{Name: "association_ip_owner_id", Require: plugin.Optional},
				{Name: "association_public_ip", Require: plugin.Optional},
				{Name: "association_public_dns_name", Require: plugin.Optional},
				{Name: "attachment_id", Require: plugin.Optional},
				{Name: "attachment_time", Require: plugin.Optional},
				{Name: "delete_on_instance_termination", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "attached_instance_id", Require: plugin.Optional},
				{Name: "attached_instance_owner_id", Require: plugin.Optional},
				{Name: "attachment_status", Require: plugin.Optional},
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "mac_address", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "private_ip_address", Require: plugin.Optional},
				{Name: "private_dns_name", Require: plugin.Optional},
				{Name: "requester_id", Require: plugin.Optional},
				{Name: "requester_managed", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "source_dest_check", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "status", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "network_interface_id",
				Description: "The ID of the network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "interface_type",
				Description: "The type of network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The AWS account ID of the owner of the network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_allocation_id",
				Description: "Allocation id for the association. Association can be an Elastic IP address (IPv4 only), or a Carrier IP address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Association.AllocationId"),
			},
			{
				Name:        "association_carrier_ip",
				Description: "The carrier IP address associated with the network interface.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Association.CarrierIp"),
			},
			{
				Name:        "association_customer_owned_ip",
				Description: "The customer-owned IP address associated with the network interface.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Association.CustomerOwnedIp"),
			},
			{
				Name:        "association_id",
				Description: "The association ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Association.AssociationId"),
			},
			{
				Name:        "association_ip_owner_id",
				Description: "The ID of the Elastic IP address owner.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Association.IpOwnerId"),
			},
			{
				Name:        "association_public_dns_name",
				Description: "The public DNS name of the association.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Association.PublicDnsName"),
			},
			{
				Name:        "association_public_ip",
				Description: "The address of the Elastic IP address bound to the network interface.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Association.PublicIp"),
			},
			{
				Name:        "attached_instance_id",
				Description: "The ID of the attached instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attachment.InstanceId"),
			},
			{
				Name:        "attached_instance_owner_id",
				Description: "The AWS account ID of the owner of the attached instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attachment.InstanceOwnerId"),
			},
			{
				Name:        "attachment_id",
				Description: "The ID of the network interface attachment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attachment.AttachmentId"),
			},
			{
				Name:        "attachment_status",
				Description: "The attachment state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Attachment.Status"),
			},
			{
				Name:        "attachment_time",
				Description: "The timestamp indicating when the attachment initiated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Attachment.AttachTime"),
			},
			{
				Name:        "delete_on_instance_termination",
				Description: "Indicates whether the network interface is deleted when the instance is terminated.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Attachment.DeleteOnTermination"),
			},
			{
				Name:        "device_index",
				Description: "The device index of the network interface attachment on the instance.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Attachment.DeviceIndex"),
			},

			{
				Name:        "mac_address",
				Description: "The MAC address of the interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_dns_name",
				Description: "The private DNS name",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_ip_address",
				Description: "The IPv4 address of the network interface within the subnet.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "requester_id",
				Description: "The ID of the entity that launched the instance on your behalf (for example, AWS Management Console or Auto Scaling).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "requester_managed",
				Description: "Indicates whether the network interface is being managed by AWS.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "source_dest_check",
				Description: "Indicates whether traffic to or from the instance is validated.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "groups",
				Description: "Any security groups for the network interface.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ipv6_addresses",
				Description: "The IPv6 addresses associated with the network interface.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_ip_addresses",
				Description: "The IPv4 address of the network interface within the subnet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the network interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagSet"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NetworkInterfaceId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2NetworkInterfaceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2NetworkInterfaceAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2NetworkInterfaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2NetworkInterfaces", "AWS_REGION", region)

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeNetworkInterfacesInput{
		MaxResults: aws.Int64(1000),
	}

	filters := buildec2NetworkInterfaceFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	// Limiting the results
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
	err = svc.DescribeNetworkInterfacesPages(
		input,
		func(page *ec2.DescribeNetworkInterfacesOutput, isLast bool) bool {
			for _, networkInterface := range page.NetworkInterfaces {
				d.StreamListItem(ctx, networkInterface)

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

func getEc2NetworkInterface(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2NetworkInterface")

	region := d.KeyColumnQualString(matrixKeyRegion)
	networkInterfaceID := d.KeyColumnQuals["network_interface_id"].GetStringValue()

	// create service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{aws.String(networkInterfaceID)},
	}

	op, err := svc.DescribeNetworkInterfaces(params)
	if err != nil {
		return nil, err
	}

	if op.NetworkInterfaces != nil && len(op.NetworkInterfaces) > 0 {
		return op.NetworkInterfaces[0], nil
	}
	return nil, nil
}

func getAwsEc2NetworkInterfaceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2NetworkInterfaceTurbotData")
	region := d.KeyColumnQualString(matrixKeyRegion)
	networkInterface := h.Item.(*ec2.NetworkInterface)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":network-interface/" + *networkInterface.NetworkInterfaceId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getEc2NetworkInterfaceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.NetworkInterface)

	// Get resource tags
	var turbotTags map[string]string
	if data.TagSet != nil {
		turbotTags = map[string]string{}
		for _, i := range data.TagSet {
			turbotTags[*i.Key] = *i.Value
		}
	}
	return turbotTags, nil
}

//// UTILITY FUNCTION
// Build ec2 network interface list call input filter
func buildec2NetworkInterfaceFilter(quals plugin.KeyColumnQualMap) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)

	filterQuals := map[string]string{
		"association_id":                 "association.association-id",
		"association_allocation_id":      "association.allocation-id",
		"association_ip_owner_id":        "association.ip-owner-id",
		"association_public_ip":          "association.public-ip",
		"association_public_dns_name":    "association.public-dns-name",
		"attachment_id":                  "attachment.attachment-id",
		"attachment_time":                "attachment.attach-time",
		"attached_instance_id":           "attachment.instance-id",
		"attached_instance_owner_id":     "attachment.instance-owner-id",
		"attachment_status":              "attachment.status",
		"availability_zone":              "availability-zone",
		"delete_on_instance_termination": "attachment.delete-on-termination",
		"description":                    "description",
		"mac_address":                    "mac-address",
		"owner_id":                       "owner-id",
		"private_ip_address":             "private-ip-address",
		"private_dns_name":               "private-dns-name",
		"source_dest_check":              "source-dest-check",
		"requester_id":                   "requester-id",
		"requester_managed":              "requester-managed",
		"status":                         "status",
	}

	columnsBool := []string{"delete_on_instance_termination", "source_dest_check", "requester_managed"}
	columnIpAddr := []string{"association_ip_owner_id", "association_public_ip", "private_ip_address"}
	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := ec2.Filter{
				Name: aws.String(filterName),
			}
			if strings.Contains(fmt.Sprint(columnsBool), columnName) { //check Bool columns
				value := getQualsValueByColumn(quals, columnName, "boolean")
				filter.Values = []*string{aws.String(fmt.Sprint(value))}
			} else if strings.Contains(fmt.Sprint(columnIpAddr), columnName) {
				value := getQualsValueByColumn(quals, columnName, "ipaddr")
				filter.Values = []*string{aws.String(fmt.Sprint(value))}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				val, ok := value.(string)
				if ok {
					filter.Values = []*string{aws.String(val)}
				} else {
					v := value.([]*string)
					filter.Values = v
				}
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
