package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEfsMountTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_efs_mount_target",
		Description: "AWS EFS Mount Target",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("mount_target_id"),
			ShouldIgnoreError: isNotFoundError([]string{"MountTargetNotFound", "InvalidParameter"}),
			Hydrate:           getAwsEfsMountTarget,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listElasticFileSystem,
			Hydrate:       listAwsEfsMountTargets,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "mount_target_id",
				Description: "The ID of the mount target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "file_system_id",
				Description: "The ID of the file system for which the mount target is intended.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "life_cycle_state",
				Description: "Lifecycle state of the mount target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone_id",
				Description: "The unique and consistent identifier of the Availability Zone that the mount target resides in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone_name",
				Description: "The name of the Availability Zone in which the mount target is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "Address at which the file system can be mounted by using the mount target.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "network_interface_id",
				Description: "The ID of the network interface that Amazon EFS created when it created the mount target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "AWS account ID that owns the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_id",
				Description: "The ID of the mount target's subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The virtual private cloud (VPC) ID that the mount target is configured in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_groups",
				Description: "Specifies the security groups currently in effect for a mount target.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEfsMountTargetSecurityGroup,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MountTargetId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEfsMountTargetAkas,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsEfsMountTargets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create session
	svc, err := EfsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	data := h.Item.(*efs.FileSystemDescription)
	params := &efs.DescribeMountTargetsInput{
		FileSystemId: data.FileSystemId,
	}

	op, err := svc.DescribeMountTargets(params)
	if err != nil {
		return nil, err
	}

	for _, mounttarget := range op.MountTargets {
		d.StreamLeafListItem(ctx, mounttarget)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsEfsMountTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEfsMountTarget")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// create service
	svc, err := EfsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	mountTargetID := d.KeyColumnQuals["mount_target_id"].GetStringValue()

	params := &efs.DescribeMountTargetsInput{
		MountTargetId: aws.String(mountTargetID),
	}

	op, err := svc.DescribeMountTargets(params)
	if err != nil {
		return nil, err
	}

	if op.MountTargets != nil && len(op.MountTargets) > 0 {
		return op.MountTargets[0], nil
	}

	return nil, nil
}

func getAwsEfsMountTargetSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEfsMountTargetSecurityGroup")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// create service
	svc, err := EfsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	data := h.Item.(*efs.MountTargetDescription)
	params := &efs.DescribeMountTargetSecurityGroupsInput{
		MountTargetId: aws.String(*data.MountTargetId),
	}

	op, err := svc.DescribeMountTargetSecurityGroups(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func getAwsEfsMountTargetAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEfsMountTargetAkas")
	data := h.Item.(*efs.MountTargetDescription)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	aka := "arn:" + commonColumnData.Partition + ":elasticfilesystem:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":file-system/" + *data.FileSystemId + "/mount-target/" + *data.MountTargetId

	return aka, nil
}
