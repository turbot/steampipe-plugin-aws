package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEBSVolume(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_volume",
		Description: "AWS EBS Volume",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("volume_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidVolume.NotFound", "InvalidParameterValue"}),
			},
			Hydrate: getEBSVolume,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVolumes"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEBSVolume,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVolumes"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "encrypted", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "fast_restored", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "multi_attach_enabled", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "size", Require: plugin.Optional},
				{Name: "snapshot_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "volume_id", Require: plugin.Optional},
				{Name: "volume_type", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getVolumeAutoEnableIOData,
				Tags: map[string]string{"service": "ec2", "action": "DescribeVolumeAttribute"},
			},
			{
				Func: getVolumeProductCodes,
				Tags: map[string]string{"service": "ec2", "action": "DescribeVolumeAttribute"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "volume_id",
				Description: "The ID of the volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the volume.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEBSVolumeARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "volume_type",
				Description: "The volume type. This can be gp2 for General Purpose SSD, io1 or io2 for Provisioned IOPS SSD, st1 for Throughput Optimized HDD, sc1 for Cold HDD, or standard for Magnetic volumes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The volume state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time stamp when volume creation was initiated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "auto_enable_io",
				Description: "The state of autoEnableIO attribute.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getVolumeAutoEnableIOData,
				Transform:   transform.FromField("AutoEnableIO.Value"),
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone for the volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encrypted",
				Description: "Indicates whether the volume is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "fast_restored",
				Description: "Indicates whether the volume was created using fast snapshot restore.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "iops",
				Description: "The number of I/O operations per second (IOPS) that the volume supports.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "throughput",
				Description: "The throughput that the volume supports, in MiB/s..",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kms_key_id",
				Description: "The Amazon Resource Name (ARN) of the AWS Key Management Service (AWS KMS) customer master key (CMK) that was used to protect the volume encryption key for the volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "multi_attach_enabled",
				Description: "Indicates whether Amazon EBS Multi-Attach is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "The size of the volume, in GiBs.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "snapshot_id",
				Description: "The snapshot from which the volume was created, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attachments",
				Description: "Information about the volume attachments.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product_codes",
				Description: "A list of product codes.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVolumeProductCodes,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the volume.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEBSVolumeTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEBSVolumeTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEBSVolumeARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEBSVolume(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_volume.listEBSVolume", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filters := buildEbsVolumeFilter(d.Quals)

	if len(filters) != 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeVolumesPaginator(svc, input, func(o *ec2.DescribeVolumesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ebs_volume.listEBSVolume", "api_error", err)
			return nil, err
		}

		for _, items := range output.Volumes {
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

func getEBSVolume(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	volumeID := d.EqualsQuals["volume_id"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_volume.getEBSVolume", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVolumesInput{
		VolumeIds: []string{volumeID},
	}

	// Get call
	op, err := svc.DescribeVolumes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_volume.getEBSVolume", "api_error", err)
		return nil, err
	}

	if len(op.Volumes) > 0 {
		h.Item = op.Volumes[0]
		return op.Volumes[0], nil
	}
	return nil, nil
}

func getVolumeAutoEnableIOData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	volume := h.Item.(types.Volume)

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_volume.getVolumeAutoEnableIOData", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVolumeAttributeInput{
		VolumeId:  volume.VolumeId,
		Attribute: types.VolumeAttributeNameAutoEnableIO,
	}

	volumeAttributes, err := svc.DescribeVolumeAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_volume.getVolumeAutoEnableIOData", "api_error", err)
		return nil, err
	}

	return volumeAttributes, nil
}

func getVolumeProductCodes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	volume := h.Item.(types.Volume)

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_volume.getVolumeProductCodes", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeVolumeAttributeInput{
		VolumeId:  volume.VolumeId,
		Attribute: types.VolumeAttributeNameProductCodes,
	}

	volumeAttributes, err := svc.DescribeVolumeAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_volume.getVolumeProductCodes", "api_error", err)
		return nil, err
	}

	return volumeAttributes, nil
}

func getEBSVolumeARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	volume := h.Item.(types.Volume)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":volume/" + *volume.VolumeId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getEBSVolumeTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volume := d.HydrateItem.(types.Volume)
	var turbotTagsMap map[string]string
	if volume.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range volume.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func getEBSVolumeTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volume := d.HydrateItem.(types.Volume)

	title := volume.VolumeId
	if volume.Tags != nil {
		for _, i := range volume.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	return title, nil
}

//// UTILITY FUNCTION

// Build ebs volume list call input filter
func buildEbsVolumeFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"availability_zone":    "availability-zone",
		"encrypted":            "encrypted",
		"fast_restored":        "fast-restored",
		"multi_attach_enabled": "multi-attach-enabled",
		"size":                 "size",
		"snapshot_id":          "snapshot-id",
		"state":                "status",
		"volume_id":            "volume-id",
		"volume_type":          "volume-type",
	}

	columnsBool := []string{"encrypted", "fast_restored", "multi_attach_enabled"}
	columnsInt := []string{"size"}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			if strings.Contains(fmt.Sprint(columnsBool), columnName) { //check Bool columns
				value := getQualsValueByColumn(quals, columnName, "boolean")
				filter.Values = []string{fmt.Sprint(value)}
			} else if strings.Contains(fmt.Sprint(columnsInt), columnName) { //check Int columns
				value := getQualsValueByColumn(quals, columnName, "int64")
				filter.Values = []string{fmt.Sprint(value)}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				if value != nil {
					val, ok := value.(string)
					if ok {
						filter.Values = []string{val}
					}
				}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
