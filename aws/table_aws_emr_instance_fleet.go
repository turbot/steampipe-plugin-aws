package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEmrInstanceFleet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_instance_fleet",
		Description: "AWS EMR Instance Fleet",
		List: &plugin.ListConfig{
			ParentHydrate: listEmrClusters,
			Hydrate:       listEmrInstanceFleets,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the instance fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The identifier of the instance fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the instance fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEmrInstanceFleetARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "cluster_id",
				Description: "The unique identifier for the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterID"),
			},
			{
				Name:        "instance_fleet_type",
				Description: "The type of the instance fleet. Valid values are MASTER, CORE or TASK.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the instance fleet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.State"),
			},
			{
				Name:        "provisioned_on_demand_capacity",
				Description: "The number of On-Demand units that have been provisioned for the instance fleet to fulfill TargetOnDemandCapacity.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "provisioned_spot_capacity",
				Description: "The number of Spot units that have been provisioned for this instance fleet to fulfill TargetSpotCapacity.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "target_on_demand_capacity",
				Description: "The target capacity of On-Demand units for the instance fleet, which determines how many On-Demand Instances to provision.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "target_spot_capacity",
				Description: "The target capacity of Spot units for the instance fleet, which determines how many Spot Instances to provision.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "instance_type_specifications",
				Description: "An array of specifications for the instance types that comprise an instance fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "launch_specifications",
				Description: "Describes the launch specification for an instance fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_change_reason",
				Description: "Provides status change reason details for the instance fleet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.StateChangeReason"),
			},
			{
				Name:        "status_timeline",
				Description: "Provides historical timestamps for the instance fleet, including the time of creation, the time it became ready to run jobs, and the time of termination.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.Timeline"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(emrInstanceFleetTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEmrInstanceFleetARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type instanceFleetDetails = struct {
	emr.InstanceFleet
	ClusterID string
}

//// LIST FUNCTION

func listEmrInstanceFleets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EmrService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get cluster details
	clusterID := h.Item.(*emr.ClusterSummary).Id

	// List call
	err = svc.ListInstanceFleetsPages(
		&emr.ListInstanceFleetsInput{
			ClusterId: clusterID,
		},
		func(page *emr.ListInstanceFleetsOutput, isLast bool) bool {
			for _, instanceFleet := range page.InstanceFleets {
				d.StreamListItem(ctx, instanceFleetDetails{*instanceFleet, *clusterID})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "InvalidRequestException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listEmrInstanceFleets", "ListInstanceFleetsPages-err", err)
		return nil, err
	}
	return nil, nil
}

func getEmrInstanceFleetARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEmrInstanceFleetARN")
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(instanceFleetDetails)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":emr:" + region + ":" + commonColumnData.AccountId + ":instance-fleet/" + *data.Id

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func emrInstanceFleetTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(instanceFleetDetails)

	if *data.Name == "" {
		return data.Id, nil
	}
	return data.Name, nil
}
