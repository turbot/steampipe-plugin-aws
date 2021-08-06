package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
)

//// TABLE DEFINITION

func tableAwsAppAutoScalingTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appautoscaling_target",
		Description: "AWS Application Auto Scaling Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"service_namespace", "resource_id"}),
			Hydrate:    getAwsApplicationAutoScalingTarget,
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("service_namespace"),
			Hydrate:    listAwsApplicationAutoScalingTargets,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_namespace",
				Description: "The namespace of the AWS service that provides the resource, or a custom-resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_id",
				Description: "The identifier of the resource associated with the scalable target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scalable_dimension",
				Description: "The scalable dimension associated with the scalable target. This string consists of the service namespace, resource type, and scaling property.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The Unix timestamp for when the scalable target was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "min_capacity",
				Description: "The minimum value to scale to in response to a scale-in activity.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_capacity",
				Description: "The maximum value to scale to in response to a scale-out activity.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "role_arn",
				Description: "The ARN of an IAM role that allows Application Auto Scaling to modify the scalable target on your behalf.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "suspended_state",
				Description: "Specifies whether the scaling activities for a scalable target are in a suspended state.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsApplicationAutoScalingTargets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["service_namespace"].GetStringValue()

	// Create Session
	svc, err := ApplicationAutoScalingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeScalableTargetsPages(
		&applicationautoscaling.DescribeScalableTargetsInput{
			ServiceNamespace: &name,
		},
		func(page *applicationautoscaling.DescribeScalableTargetsOutput, isLast bool) bool {
			for _, scalableTarget := range page.ScalableTargets {
				d.StreamListItem(ctx, scalableTarget)

			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsApplicationAutoScalingTarget(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsApplicationAutoScalingTarget")
	
	name := d.KeyColumnQuals["service_namespace"].GetStringValue()
	id := d.KeyColumnQuals["resource_id"].GetStringValue()

	// create service
	svc, err := ApplicationAutoScalingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &applicationautoscaling.DescribeScalableTargetsInput{
		ServiceNamespace: &name,
		ResourceIds:      []*string{types.String(id)},
	}

	// Get call
	op, err := svc.DescribeScalableTargets(params)
	if err != nil {
		return nil, err
	}
	if op.ScalableTargets != nil && len(op.ScalableTargets) > 0 {
		return op.ScalableTargets[0], nil
	}

	return nil, nil
}
