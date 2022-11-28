package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2Ami(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_ami",
		Description: "AWS EC2 AMI",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("image_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAMIID.NotFound", "InvalidAMIID.Unavailable", "InvalidAMIID.Malformed"}),
			},
			Hydrate: getEc2Ami,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2Amis,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "architecture", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "ena_support", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "hypervisor", Require: plugin.Optional},
				{Name: "image_type", Require: plugin.Optional},
				{Name: "public", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "kernel_id", Require: plugin.Optional},
				{Name: "platform", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "ramdisk_id", Require: plugin.Optional},
				{Name: "root_device_name", Require: plugin.Optional},
				{Name: "root_device_type", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "sriov_net_support", Require: plugin.Optional},
				{Name: "virtualization_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Default:     "self",
			},
			{
				Name:        "imds_support",
				Description: "Indicates that IMDSv2 is specified in the AMI",
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
				Transform:   transform.FromField("Platform").NullIfZero(),
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
				Name:        "launch_permissions",
				Description: "The users and groups that have the permissions for creating instances from the AMI.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2AmiLaunchPermissionData,
				Transform:   transform.FromField("LaunchPermissions"),
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

func listEc2Amis(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami.listEc2Amis", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeImagesInput{}

	filters := buildAmisWithOwnerFilter(d.Quals, "AMI", ctx, d, h)
	if len(filters) != 0 {
		input.Filters = filters
	}

	resp, err := svc.DescribeImages(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami.listEc2Amis", "api_error", err)
		return nil, err
	}
	for _, image := range resp.Images {
		d.StreamListItem(ctx, image)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2Ami(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	imageID := d.KeyColumnQuals["image_id"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami.getEc2Ami", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeImagesInput{
		ImageIds: []string{imageID},
	}

	op, err := svc.DescribeImages(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami.getEc2Ami", "api_error", err)
		return nil, err
	}

	if op.Images != nil && len(op.Images) > 0 {
		return op.Images[0], nil
	}
	return nil, nil
}

type LaunchPermissions struct {
	Group                 *string
	OrganizationArn       *string
	OrganizationalUnitArn *string
	UserId                *string
}

func getAwsEc2AmiLaunchPermissionData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	image := h.Item.(types.Image)

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami.getAwsEc2AmiLaunchPermissionData", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeImageAttributeInput{
		ImageId:   image.ImageId,
		Attribute: types.ImageAttributeNameLaunchPermission,
	}

	imageData, err := svc.DescribeImageAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami.getAwsEc2AmiLaunchPermissionData", "api_error", err)
		return nil, err
	}

	return imageData, nil
}

func getAwsEc2AmiAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	image := h.Item.(types.Image)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + *image.OwnerId + ":image/" + *image.ImageId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getEc2AmiTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	image := d.HydrateItem.(types.Image)
	var turbotTagsMap map[string]string
	if image.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range image.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func getEc2AmiTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.Image)

	title := data.ImageId
	if data.Name != nil {
		title = data.Name
	}

	return title, nil
}
