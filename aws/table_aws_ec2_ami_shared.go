package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	go_kit_pack "github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2AmiShared(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_ami_shared",
		Description: "AWS EC2 AMI - All public, private, and shared AMIs",
		List: &plugin.ListConfig{
			Hydrate: listAmisByOwner,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeImages"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "owner_id", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "architecture", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "ena_support", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "hypervisor", Require: plugin.Optional},
				{Name: "image_id", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "image_type", Require: plugin.Optional},
				{Name: "public", Require: plugin.Optional, Operators: []string{"=", "<>"}},
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
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAMIID.NotFound", "InvalidAMIID.Unavailable", "InvalidAMIID.Malformed"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
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
				Hydrate:     getImageOwnerAlias,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "imds_support",
				Description: "If v2.0, it indicates that IMDSv2 is specified in the AMI.",
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

func listAmisByOwner(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	owner_id := d.EqualsQuals["owner_id"].GetStringValue()
	image_id := d.EqualsQuals["image_id"].GetStringValue()

	// check if owner_id and image_id is empty
	if owner_id == "" && image_id == "" {
		return nil, errors.New("please provide either owner_id or image_id")
	}

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami_shared.listAmisByOwner", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeImagesInput{}

	if owner_id != "" {
		input.Owners = []string{owner_id}
	}
	if image_id != "" {
		input.ImageIds = []string{image_id}
	}

	filters := buildSharedAmisWithOwnerFilter(d.Quals, ctx, d, h)

	if len(filters) != 0 {
		input.Filters = filters
	}

	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	// There is no MaxResult property in param, through which we can limit the number of results
	resp, err := svc.DescribeImages(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ami_shared.listAmisByOwner", "api_error", err)
		return nil, err
	}
	for _, image := range resp.Images {
		d.StreamListItem(ctx, image)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

func getImageOwnerAlias(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	image := h.Item.(types.Image)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	if image.ImageOwnerAlias == nil && commonColumnData.AccountId != *image.OwnerId {
		return *image.OwnerId, nil
	} else if image.ImageOwnerAlias == nil {
		return "self", nil
	} else {
		return *image.ImageOwnerAlias, nil
	}
}

// // UTILITY FUNCTION
// Build AMI's list call input filter
func buildSharedAmisWithOwnerFilter(quals plugin.KeyColumnQualMap, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"architecture":        "architecture",
		"description":         "description",
		"ena_support":         "ena-support",
		"hypervisor":          "hypervisor",
		"image_type":          "image-type",
		"kernel_id":           "kernel-id",
		"name":                "name",
		"platform":            "platform",
		"public":              "is-public",
		"ramdisk_id":          "ramdisk-id",
		"root_device_name":    "root-device-name",
		"root_device_type":    "root-device-type",
		"state":               "state",
		"sriov_net_support":   "sriov-net-support",
		"virtualization_type": "virtualization-type",
	}

	columnsBool := []string{"ena_support", "public"}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: go_kit_pack.String(filterName),
			}

			//check Bool columns
			if strings.Contains(fmt.Sprint(columnsBool), columnName) {
				value := getQualsValueByColumn(quals, columnName, "boolean")
				filter.Values = []string{fmt.Sprint(value)}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				val, ok := value.(string)
				if ok {
					filter.Values = []string{val}
				}
			}
			filters = append(filters, filter)
		}
	}

	return filters
}
