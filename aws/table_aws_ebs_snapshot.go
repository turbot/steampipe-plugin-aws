package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func tableAwsEBSSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_snapshot",
		Description: "AWS EBS Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("snapshot_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidSnapshot.NotFound", "InvalidSnapshotID.Malformed", "InvalidParameterValue"}),
			},
			Hydrate: getAwsEBSSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsEBSSnapshots,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "description",
					Require: plugin.Optional,
				},
				{
					Name:    "encrypted",
					Require: plugin.Optional,
				},
				{
					Name:    "owner_alias",
					Require: plugin.Optional,
				},
				{
					Name:    "owner_id",
					Require: plugin.Optional,
				},
				{
					Name:    "snapshot_id",
					Require: plugin.Optional,
				},
				{
					Name:    "state",
					Require: plugin.Optional,
				},
				{
					Name:    "progress",
					Require: plugin.Optional,
				},
				{
					Name:    "volume_id",
					Require: plugin.Optional,
				},
				{
					Name:    "volume_size",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "snapshot_id",
				Description: "The ID of the snapshot. Each snapshot receives a unique identifier when it is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the snapshot.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEBSSnapshotARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "state",
				Description: "The snapshot state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "volume_size",
				Description: "The size of the volume, in GiB.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "volume_id",
				Description: "The ID of the volume that was used to create the snapshot. Snapshots created by the CopySnapshot action have an arbitrary volume ID that should not be used for any purpose.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encrypted",
				Description: "Indicates whether the snapshot is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "start_time",
				Description: "The time stamp when the snapshot was initiated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description for the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The Amazon Resource Name (ARN) of the AWS Key Management Service (AWS KMS) customer master key (CMK) that was used to protect the volume encryption key for the parent volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_encryption_key_id",
				Description: "The data encryption key identifier for the snapshot. This value is a unique identifier that corresponds to the data encryption key that was used to encrypt the original volume or snapshot copy. Because data encryption keys are inherited by volumes created from snapshots, and vice versa, if snapshots share the same data encryption key identifier, then they belong to the same volume/snapshot lineage.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "progress",
				Description: "The progress of the snapshot, as a percentage.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_message",
				Description: "Encrypted Amazon EBS snapshots are copied asynchronously. If a snapshot copy operation fails this field displays error state details to help you diagnose why the error occurred.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_alias",
				Description: "The AWS owner alias, from an Amazon-maintained list (amazon). This is not the user-configured AWS account alias set using the IAM console.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The AWS account ID of the EBS snapshot owner.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_volume_permissions",
				Description: "The users and groups that have the permissions for creating volumes from the snapshot.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEBSSnapshotCreateVolumePermissions,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the snapshot.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(ec2SnapshotTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEBSSnapshotARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsEBSSnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listAwsEBSSnapshots", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeSnapshotsInput{
		MaxResults: aws.Int64(1000),
	}

	// Build filter for ebs snapshot
	filters := buildEbsSnapshotFilter(ctx, d, h, d.KeyColumnQuals)
	input.Filters = filters

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = types.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}
	// List call
	err = svc.DescribeSnapshotsPages(
		input,
		func(page *ec2.DescribeSnapshotsOutput, isLast bool) bool {
			for _, snapshot := range page.Snapshots {
				d.StreamListItem(ctx, snapshot)

				// Context can be cancelled due to manual cancellation or the limit has been hit
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

func getAwsEBSSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEBSSnapshot")

	region := d.KeyColumnQualString(matrixKeyRegion)
	snapshotID := d.KeyColumnQuals["snapshot_id"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeSnapshotsInput{
		SnapshotIds: []*string{aws.String(snapshotID)},
	}

	// Get call
	data, err := svc.DescribeSnapshots(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsEBSSnapshot__", "ERROR", err)
		return nil, err
	}

	if data.Snapshots != nil {
		return data.Snapshots[0], nil
	}
	return nil, nil
}

// getAwsEBSSnapshotCreateVolumePermissions :: Describes the users and groups that have the permissions for creating volumes from the snapshot
func getAwsEBSSnapshotCreateVolumePermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEBSSnapshotCreateVolumePermissions")
	snapshotData := h.Item.(*ec2.Snapshot)
	region := d.KeyColumnQualString(matrixKeyRegion)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	if *snapshotData.OwnerId != commonColumnData.AccountId {
		return nil, nil
	}
	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &ec2.DescribeSnapshotAttributeInput{
		SnapshotId: snapshotData.SnapshotId,
		Attribute:  aws.String(ec2.SnapshotAttributeNameCreateVolumePermission),
	}

	// Describe create volume permission
	resp, err := svc.DescribeSnapshotAttribute(params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getEBSSnapshotARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEBSSnapshotARN")
	region := d.KeyColumnQualString(matrixKeyRegion)
	snapshotData := h.Item.(*ec2.Snapshot)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// Get the resource arn
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + *snapshotData.OwnerId + ":snapshot/" + *snapshotData.SnapshotId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func ec2SnapshotTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	snapshot := d.HydrateItem.(*ec2.Snapshot)
	return ec2TagsToMap(snapshot.Tags)
}

//// UTILITY FUNCTION

// build ebs snapshot list call input filter
func buildEbsSnapshotFilter(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, equalQuals plugin.KeyColumnEqualsQualMap) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)

	filterQuals := map[string]string{
		"description": "description",
		"encrypted":   "encrypted",
		"owner_alias": "owner-alias",
		"snapshot_id": "snapshot-id",
		"state":       "status",
		"progress":    "progress",
		"volume_id":   "volume-id",
		"volume_size": "volume-size",
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
	if equalQuals["owner_id"] != nil {
		ownerFilter.Name = types.String("owner-id")
		ownerFilter.Values = []*string{types.String(equalQuals["owner_id"].GetStringValue())}
	} else {
		// Use this section later and compare the results
		getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
		c, err := getCommonColumnsCached(ctx, d, h)
		if err != nil {
			return filters
		}
		commonColumnData := c.(*awsCommonColumnData)
		ownerFilter.Name = types.String("owner-id")
		ownerFilter.Values = []*string{aws.String(commonColumnData.AccountId)}
	}

	filters = append(filters, &ownerFilter)
	return filters
}
