package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEmrInstanceGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_instance_group",
		Description: "AWS EMR Instance Group",
		List: &plugin.ListConfig{
			ParentHydrate: listEmrClusters,
			Hydrate:       listEmrInstanceGroups,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the instance group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The identifier of the instance group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the instance group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEmrInstanceGroupARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "cluster_id",
				Description: "The unique identifier for the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterID"),
			},
			{
				Name:        "instance_group_type",
				Description: "The type of the instance group. Valid values are MASTER, CORE or TASK.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "The EC2 instance type for all instances in the instance group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the instance group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.State"),
			},
			{
				Name:        "bid_price",
				Description: "The maximum price you are willing to pay for Spot Instances. If specified, indicates that the instance group uses Spot Instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "configurations_version",
				Description: "The version number of the requested configuration specification for this instance group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ebs_optimized",
				Description: "Indicates whether the instance group is EBS-optimized, or not.  An Amazon EBS-optimized instance uses an optimized configuration stack and provides additional, dedicated capacity for Amazon EBS I/O.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_successfully_applied_configurations_version",
				Description: "The version number of a configuration specification that was successfully applied for an instance group last time.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "market",
				Description: "The marketplace to provision instances for this group. Valid values are ON_DEMAND or SPOT.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "requested_instance_count",
				Description: "The target number of instances for the instance group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "running_instance_count",
				Description: "The number of instances currently running in this instance group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "autoscaling_policy",
				Description: "An automatic scaling policy for a core instance group or task instance group in an Amazon EMR cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AutoScalingPolicy"),
			},
			{
				Name:        "configurations",
				Description: "A list of configurations supplied for an EMR cluster instance group. Only availbale for Amazon EMR releases 4.x or later.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ebs_block_devices",
				Description: "The EBS block devices that are mapped to this instance group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_successfully_applied_configurations",
				Description: "A list of configurations that were successfully applied for an instance group last time.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "shrink_policy",
				Description: "Policy for customizing shrink operations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_change_reason",
				Description: "The status change reason details for the instance group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.StateChangeReason"),
			},
			{
				Name:        "status_timeline",
				Description: "The timeline of the instance group status over time.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.Timeline"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(EmrInstanceGroupTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEmrInstanceGroupARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type instanceGroupDetails = struct {
	types.InstanceGroup
	ClusterID string
}

//// LIST FUNCTION

func listEmrInstanceGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EmrClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_instance_group.listEmrInstanceGroups", "connection_error", err)
		return nil, err
	}

	// Get cluster details
	clusterID := h.Item.(types.ClusterSummary).Id

	input := &emr.ListInstanceGroupsInput{
		ClusterId: clusterID,
	}

	paginator := emr.NewListInstanceGroupsPaginator(svc, input, func(o *emr.ListInstanceGroupsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_emr_instance_group.listEmrInstanceGroups", err)
			return nil, err
		}

		for _, items := range output.InstanceGroups {
			d.StreamListItem(ctx, instanceGroupDetails{items, *clusterID})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getEmrInstanceGroupARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(instanceGroupDetails)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":emr:" + region + ":" + commonColumnData.AccountId + ":instance-group/" + *data.Id

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func EmrInstanceGroupTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(instanceGroupDetails)

	if *data.Name == "" {
		return data.Id, nil
	}
	return data.Name, nil
}
