package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"

	efsEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEfsMountTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_efs_mount_target",
		Description: "AWS EFS Mount Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("mount_target_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"MountTargetNotFound", "InvalidParameter"}),
			},
			Hydrate: getAwsEfsMountTarget,
			Tags:    map[string]string{"service": "elasticfilesystem", "action": "DescribeMountTargets"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listElasticFileSystem,
			Hydrate:       listAwsEfsMountTargets,
			Tags:          map[string]string{"service": "elasticfilesystem", "action": "DescribeMountTargets"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEfsMountTargetSecurityGroup,
				Tags: map[string]string{"service": "elasticfilesystem", "action": "DescribeMountTargetSecurityGroups"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(efsEndpoint.ELASTICFILESYSTEMServiceID),
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
	svc, err := EFSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_mount_target.listAwsEfsMountTargets", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxLimit := int32(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}

	data := h.Item.(types.FileSystemDescription)
	params := &efs.DescribeMountTargetsInput{
		FileSystemId: data.FileSystemId,
		MaxItems:     aws.Int32(maxLimit),
	}

	// List call
	pagesLeft := true
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		result, err := svc.DescribeMountTargets(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_efs_mount_target.listAwsEfsMountTargets", "api_error", err)
			return nil, err
		}
		for _, mountTarget := range result.MountTargets {
			d.StreamListItem(ctx, mountTarget)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
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
	mountTargetID := d.EqualsQuals["mount_target_id"].GetStringValue()

	if strings.TrimSpace(mountTargetID) == "" {
		return nil, nil
	}

	// Create service
	svc, err := EFSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_mount_target.getAwsEfsMountTarget", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &efs.DescribeMountTargetsInput{
		MountTargetId: aws.String(mountTargetID),
	}

	op, err := svc.DescribeMountTargets(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_mount_target.getAwsEfsMountTarget", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.MountTargets) > 0 {
		return op.MountTargets[0], nil
	}

	return nil, nil
}

func getAwsEfsMountTargetSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.MountTargetDescription)

	// Create service
	svc, err := EFSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_mount_target.getAwsEfsMountTargetSecurityGroup", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &efs.DescribeMountTargetSecurityGroupsInput{
		MountTargetId: aws.String(*data.MountTargetId),
	}

	op, err := svc.DescribeMountTargetSecurityGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_mount_target.getAwsEfsMountTargetSecurityGroup", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func getAwsEfsMountTargetAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.MountTargetDescription)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_mount_target.getAwsEfsMountTargetAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	aka := "arn:" + commonColumnData.Partition + ":elasticfilesystem:" + region + ":" + commonColumnData.AccountId + ":file-system/" + *data.FileSystemId + "/mount-target/" + *data.MountTargetId

	return aka, nil
}
