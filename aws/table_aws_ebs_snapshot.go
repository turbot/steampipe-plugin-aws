package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsEBSSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ebs_snapshot",
		Description: "AWS EBS Snapshot",
		List: &plugin.ListConfig{
			Hydrate: listAwsEBSSnapshots,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSnapshots"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidSnapshot.NotFound", "InvalidSnapshotID.Malformed", "InvalidUserID.Malformed"}),
			},
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
					Name:       "owner_alias",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "owner_id",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "snapshot_id",
					Require:    plugin.Optional,
					CacheMatch: "exact",
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
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsEBSSnapshotCreateVolumePermissions,
				Tags: map[string]string{"service": "ec2", "action": "DescribeSnapshotAttribute"},
			},
		},
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
				Name:        "restore_expiry_time",
				Description: "Only for archived snapshots that are temporarily restored. Indicates the date and time when a temporarily restored snapshot will be automatically re-archived.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description for the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "outpost_arn",
				Description: "The ARN of the Outpost on which the snapshot is stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_tier",
				Description: "The storage tier in which the snapshot is stored.",
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
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_snapshot.listAwsEBSSnapshots", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeSnapshotsInput{}

	// Limiting the results
	// You cannot specify limit and the snapshot ID in the same request.
	// It would be more efficient / faster for us to not paginate requests to the DescribeSnapshots API (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSnapshots.html#API_DescribeSnapshots_RequestParameters).
	// If the user specifies a limit parameter with a value ranging from 5 to 1000, it will be set accordingly. Otherwise, the DescribeSnapshots API will default to retrieving all available snapshots, which is faster than paginating with a value of 1000 per page.
	if d.EqualsQualString("snapshot_id") == "" {
		if d.QueryContext.Limit != nil {
			limit := int32(*d.QueryContext.Limit)
			if limit < 5 {
				input.MaxResults = aws.Int32(5)
			} else if limit > 1000 {
				input.MaxResults = aws.Int32(1000)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	// Build filter for ebs snapshot
	filters := buildEbsSnapshotFilter(ctx, d, h, d.EqualsQuals, input)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeSnapshotsPaginator(svc, input, func(o *ec2.DescribeSnapshotsPaginatorOptions) {
		if input.MaxResults != nil {
			o.Limit = *input.MaxResults
		}
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ebs_snapshot.listAwsEBSSnapshots", "api_error", err)
			return nil, err
		}

		for _, items := range output.Snapshots {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

// getAwsEBSSnapshotCreateVolumePermissions :: Describes the users and groups that have the permissions for creating volumes from the snapshot
func getAwsEBSSnapshotCreateVolumePermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	snapshotData := h.Item.(types.Snapshot)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	if *snapshotData.OwnerId != commonColumnData.AccountId {
		return nil, nil
	}
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_snapshot.getAwsEBSSnapshotCreateVolumePermissions", "connection_error", err)
		return nil, err
	}

	// Build params
	params := &ec2.DescribeSnapshotAttributeInput{
		SnapshotId: snapshotData.SnapshotId,
		Attribute:  types.SnapshotAttributeNameCreateVolumePermission,
	}

	// Describe create volume permission
	resp, err := svc.DescribeSnapshotAttribute(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ebs_snapshot.getAwsEBSSnapshotCreateVolumePermissions", "api_error", err)
		return nil, err
	}
	return resp, nil
}

func getEBSSnapshotARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	snapshotData := h.Item.(types.Snapshot)

	c, err := getCommonColumns(ctx, d, h)
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
	snapshot := d.HydrateItem.(types.Snapshot)
	var turbotTagsMap map[string]string
	if snapshot.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range snapshot.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}
	return turbotTagsMap, nil
}

//// UTILITY FUNCTION

// build ebs snapshot list call input filter
func buildEbsSnapshotFilter(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, equalQuals plugin.KeyColumnEqualsQualMap, input *ec2.DescribeSnapshotsInput) []types.Filter {
	filters := make([]types.Filter, 0)

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
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := equalQuals[columnName]
			if value.GetStringValue() != "" {
				filter.Values = []string{equalQuals[columnName].GetStringValue()}
			}
			filters = append(filters, filter)
		}
	}

	if equalQuals["owner_id"] != nil {
		input.OwnerIds = append(input.OwnerIds, equalQuals["owner_id"].GetStringValue())
	}

	/*
	 * By default, DescribeSnapshots API returns all the snapshots, including public & shared ones, which will be too many and can cause performance issues for the table.
	 * If the user did not provide any of the below filters in the where clause, then by default owner ID will be set to the caller account ID, and API will return snapshots from the same account. This will help in table performance.
	 */
	if equalQuals["owner_alias"] == nil && equalQuals["owner_id"] == nil && equalQuals["snapshot_id"] == nil {
		c, err := getCommonColumns(ctx, d, h)
		if err != nil {
			return filters
		}
		commonColumnData := c.(*awsCommonColumnData)
		input.OwnerIds = append(input.OwnerIds, commonColumnData.AccountId)
	}

	return filters
}
