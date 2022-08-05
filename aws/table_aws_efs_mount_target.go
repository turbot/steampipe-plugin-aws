package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEfsMountTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_efs_mount_target",
		Description: "AWS EFS Mount Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("mount_target_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"MountTargetNotFound", "InvalidParameter"}),
			},
			Hydrate: getAwsEfsMountTarget,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listElasticFileSystem,
			Hydrate:       listAwsEfsMountTargets,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	// Create session
	svc, err := EfsService(ctx, d)
	if err != nil {
		return nil, err
	}

	data := h.Item.(*efs.FileSystemDescription)
	params := &efs.DescribeMountTargetsInput{
		FileSystemId: data.FileSystemId,
		MaxItems:     aws.Int64(100),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxItems {
			if *limit < 1 {
				params.MaxItems = aws.Int64(1)
			} else {
				params.MaxItems = limit
			}
		}
	}

	// List call
	pagesLeft := true
	for pagesLeft {
		result, err := svc.DescribeMountTargets(params)
		if err != nil {
			plugin.Logger(ctx).Error("listAwsEfsMountTargets", "DescribeMountTargets_error", err)
			return nil, err
		}
		for _, mountTarget := range result.MountTargets {
			d.StreamListItem(ctx, mountTarget)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				pagesLeft = false
			}
		}
		if result.NextMarker != nil {
			params.Marker = result.NextMarker
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsEfsMountTarget(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEfsMountTarget")

	// Create service
	svc, err := EfsService(ctx, d)
	if err != nil {
		return nil, err
	}

	mountTargetID := d.KeyColumnQuals["mount_target_id"].GetStringValue()

	params := &efs.DescribeMountTargetsInput{
		MountTargetId: aws.String(mountTargetID),
	}

	op, err := svc.DescribeMountTargets(params)
	if err != nil {
		plugin.Logger(ctx).Error("getAwsEfsMountTarget", "DescribeMountTargets_error", err)
		return nil, err
	}

	if op.MountTargets != nil && len(op.MountTargets) > 0 {
		return op.MountTargets[0], nil
	}

	return nil, nil
}

func getAwsEfsMountTargetSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEfsMountTargetSecurityGroup")

	// Create service
	svc, err := EfsService(ctx, d)
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*efs.MountTargetDescription)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	aka := "arn:" + commonColumnData.Partition + ":elasticfilesystem:" + region + ":" + commonColumnData.AccountId + ":file-system/" + *data.FileSystemId + "/mount-target/" + *data.MountTargetId

	return aka, nil
}
