package aws

import (
	"context"
	"fmt"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2AmiShared(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_ami_shared",
		Description: "AWS EC2 AMI - All public, private, and shared AMIs",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("image_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidAMIID.NotFound", "InvalidAMIID.Unavailable", "InvalidAMIID.Malformed"}),
			Hydrate:           getEc2Ami,
		},
		List: &plugin.ListConfig{
			Hydrate: listAmisByOwner,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Required},
				{Name: "architecture", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "ena_support", Require: plugin.Optional},
				{Name: "hypervisor", Require: plugin.Optional},
				{Name: "image_type", Require: plugin.Optional},
				{Name: "public", Require: plugin.Optional},
				{Name: "kernel_id", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "platform", Require: plugin.Optional},
				{Name: "ramdisk_id", Require: plugin.Optional},
				{Name: "root_device_name", Require: plugin.Optional},
				{Name: "root_device_type", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "sriov_net_support", Require: plugin.Optional},
				{Name: "virtualization_type", Require: plugin.Optional},
			},
			ShouldIgnoreError: isNotFoundError([]string{"InvalidAMIID.NotFound", "InvalidAMIID.Unavailable", "InvalidAMIID.Malformed"}),
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the AMI that was provided during image creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_id",
				Description: "The ID of the AMI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the AMI. If the state is available, the image is successfully registered and can be used to launch an instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_type",
				Description: "The type of image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_location",
				Description: "The location of the AMI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date and time when the image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "architecture",
				Description: "The architecture of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the AMI that was provided during image creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ena_support",
				Description: "Specifies whether enhanced networking with ENA is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "hypervisor",
				Description: "The hypervisor type of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_owner_alias",
				Description: "The AWS account alias (for example, amazon, self) or the AWS account ID of the AMI owner.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kernel_id",
				Description: "The kernel associated with the image, if any. Only applicable for machine images.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The AWS account ID of the image owner.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform",
				Description: "This value is set to windows for Windows AMIs; otherwise, it is blank.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_details",
				Description: "The platform details associated with the billing code of the AMI. For more information, see Obtaining Billing Information (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-billing-info.html) in the Amazon Elastic Compute Cloud User Guide.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public",
				Description: "Indicates whether the image has public launch permissions. The value is true if this image has public launch permissions or false if it has only implicit and explicit launch permissions.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "ramdisk_id",
				Description: "The RAM disk associated with the image, if any. Only applicable for machine images.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "root_device_name",
				Description: "The device name of the root device volume (for example, /dev/sda1).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "root_device_type",
				Description: "The type of root device used by the AMI. The AMI can use an EBS volume or an instance store volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sriov_net_support",
				Description: "Specifies whether enhanced networking with the Intel 82599 Virtual Function interface is enabled.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "usage_operation",
				Description: "The operation of the Amazon EC2 instance and the billing code that is associated with the AMI. For the list of UsageOperation codes, see Platform Details and [Usage Operation Billing Codes](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-billing-info.html#billing-info) in the Amazon Elastic Compute Cloud User Guide.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtualization_type",
				Description: "The type of virtualization of the AMI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "block_device_mappings",
				Description: "Any block device mapping entries.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product_codes",
				Description: "Any product codes associated with the AMI.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the AMI.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2AmiTurbotTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2AmiTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2AmiAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAmisByOwner(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listAmisByOwner", "AWS_REGION", region)

	owner_id := d.KeyColumnQuals["owner_id"].GetStringValue()

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeImagesInput{
		Owners: []*string{aws.String(owner_id)},
	}

	filters := buildAmisByOwnerFilterFilter(d.KeyColumnQuals, "SHARED_AMI")
	equalQuals := d.KeyColumnQuals
	if equalQuals["ena_support"] != nil {
		filters = append(filters, &ec2.Filter{Name: aws.String("ena-support"), Values: []*string{aws.String(fmt.Sprint(equalQuals["ena_support"].GetBoolValue()))}})
	}
	if equalQuals["public"] != nil {
		filters = append(filters, &ec2.Filter{Name: aws.String("is-public"), Values: []*string{aws.String(fmt.Sprint(equalQuals["public"].GetBoolValue()))}})
	}

	if len(filters) != 0 {
		input.Filters = filters
	}

	// There is no MaxResult property in param, through which we can limit the number of results
	resp, err := svc.DescribeImages(input)
	for _, image := range resp.Images {
		d.StreamListItem(ctx, image)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// UTILITY FUNCTION
// build amis list call input filter
func buildAmisByOwnerFilterFilter(equalQuals plugin.KeyColumnEqualsQualMap, amiType string) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)

	filterQuals := map[string]string{
		"architecture":        "architecture",
		"description":         "description",
		"hypervisor":          "hypervisor",
		"image_id ":           "image-id ",
		"image_type":          "image-type",
		"kernel_id":           "kernel-id",
		"name":                "name",
		"platform":            "platform",
		"ramdisk_id":          "ramdisk-id",
		"root_device_name":    "root-device-name",
		"root_device_type":    "root-device-type",
		"state":               "state",
		"sriov_net_support":   "sriov-net-support",
		"virtualization_type": "virtualization-type",
	}

	for columnName, filterName := range filterQuals {
		if equalQuals[columnName] != nil {
			filter := ec2.Filter{
				Name: types.String(filterName),
			}
			value := equalQuals[columnName]
			if value.GetStringValue() != "" {
				filter.Values = []*string{types.String(equalQuals[columnName].GetStringValue())}
			} else if value.GetListValue() != nil {
				filter.Values = getListValues(value.GetListValue())
			}
			filters = append(filters, &filter)
		}
	}

	ownerFilter := ec2.Filter{}
	if equalQuals["owner_id"].GetStringValue() != "SHARED_AMI" {
		if equalQuals["owner_id"] != nil {
			ownerFilter.Name = types.String("owner-id")
			ownerFilter.Values = []*string{types.String(equalQuals["owner_id"].GetStringValue())}
		} else {
			ownerFilter.Name = types.String("owner-id")
			ownerFilter.Values = []*string{types.String("self")}
		}

		filters = append(filters, &ownerFilter)
	}
	return filters
}
