package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"

	emrv1 "github.com/aws/aws-sdk-go/service/emr"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEmrInstanceFleet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_instance_fleet",
		Description: "AWS EMR Instance Fleet",
		List: &plugin.ListConfig{
			ParentHydrate: listEmrClusters,
			Hydrate:       listEmrInstanceFleets,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidRequestException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(emrv1.EndpointsID),
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
	types.InstanceFleet
	ClusterID string
}

//// LIST FUNCTION

func listEmrInstanceFleets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EMRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_instance_fleet.listEmrInstanceFleets", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Get cluster details
	clusterID := h.Item.(types.ClusterSummary).Id
	input := &emr.ListInstanceFleetsInput{
		ClusterId: clusterID,
	}

	paginator := emr.NewListInstanceFleetsPaginator(svc, input, func(o *emr.ListInstanceFleetsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				// Error: operation error EMR: ListInstanceFleets, https response error StatusCode: 400, RequestID: 560c660e-9fd8-4457-9cfd-fa79427912d4, InvalidRequestException: Instance fleets and instance groups are mutually exclusive. The EMR cluster specified in the request uses instance groups. The ListInstanceFleets operation does not support clusters that use instance groups. Use the ListInstanceGroups operation instead. (SQLSTATE HV000)
				if ae.ErrorCode() == "InvalidRequestException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_emr_instance_fleet.listEmrInstanceFleets", "api_error", err)
			return nil, err
		}

		for _, instanceFleet := range output.InstanceFleets {
			d.StreamListItem(ctx, instanceFleetDetails{instanceFleet, *clusterID})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getEmrInstanceFleetARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(instanceFleetDetails)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_instance_fleet.getEmrInstanceFleetARN", "common_data_error", err)
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
