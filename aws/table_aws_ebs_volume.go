package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEBSVolume(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_volume",
		Description: "AWS EBS Volume",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("volume_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVolume.NotFound", "InvalidParameterValue"}),
			Hydrate:           getEBSVolume,
		},
		List: &plugin.ListConfig{
			Hydrate: listEBSVolume,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "volume_id",
				Description: "The ID of the volume.",
				Type:        proto.ColumnType_STRING,
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

			/// Standard columns
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
				Hydrate:     getEBSVolumeAkas,
			},
		}),
	}
}

//// LIST FUNCTION

func listEBSVolume(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEBSVolume", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeVolumesPages(
		&ec2.DescribeVolumesInput{},
		func(page *ec2.DescribeVolumesOutput, isLast bool) bool {
			for _, volume := range page.Volumes {
				d.StreamListItem(ctx, volume)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEBSVolume(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEBSVolume")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	volumeID := d.KeyColumnQuals["volume_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(volumeID)},
	}

	// Get call
	op, err := svc.DescribeVolumes(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getEBSVolume__", "ERROR", err)
		return nil, err
	}

	if len(op.Volumes) > 0 {
		h.Item = op.Volumes[0]
		return op.Volumes[0], nil
	}
	return nil, nil
}

func getVolumeAutoEnableIOData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVolumeAutoEnableIOData")
	volume := h.Item.(*ec2.Volume)
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVolumeAttributeInput{
		VolumeId:  volume.VolumeId,
		Attribute: aws.String("autoEnableIO"),
	}

	volumeAttributes, err := svc.DescribeVolumeAttribute(params)
	if err != nil {
		return nil, err
	}

	return volumeAttributes, nil
}

func getVolumeProductCodes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVolumeProductCodes")
	volume := h.Item.(*ec2.Volume)
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeVolumeAttributeInput{
		VolumeId:  volume.VolumeId,
		Attribute: aws.String("productCodes"),
	}

	volumeAttributes, err := svc.DescribeVolumeAttribute(params)
	if err != nil {
		return nil, err
	}

	return volumeAttributes, nil
}

func getEBSVolumeAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEBSVolumeTurbotTags")
	volume := h.Item.(*ec2.Volume)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":volume/" + *volume.VolumeId}

	// Mapping all the turbot defined properties
	return map[string]interface{}{
		"Akas": akas,
	}, nil
}

//// TRANSFORM FUNCTIONS

func getEBSVolumeTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volume := d.HydrateItem.(*ec2.Volume)
	return ec2TagsToMap(volume.Tags)
}

func getEBSVolumeTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volume := d.HydrateItem.(*ec2.Volume)

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
