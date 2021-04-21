package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/turbot/go-kit/types"
)

//// TABLE DEFINITION

func tableAwsAppAutoScalingTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appautoscaling_target",
		Description: "AWS ApplicationAutoScaling Target",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"service_namespace", "resource_id"}),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException"}),
			Hydrate:           getAwsApplicationAutoScalingTarget,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_namespace",
				Description: "The namespace of the AWS service that provides the resource, or a custom-resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},
			{
				Name:        "resource_id",
				Description: "The identifier of the resource associated with the scalable target.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},
			{
				Name:        "scalable_dimension",
				Description: "The scalable dimension associated with the scalable target. This string consists of the service namespace, resource type, and scaling property.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},
			{
				Name:        "creation_time",
				Description: " The Unix timestamp for when the scalable target was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},
			{
				Name:        "min_capacity",
				Description: "The minimum value to scale to in response to a scale-in activity.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},
			{
				Name:        "max_capacity",
				Description: "The maximum value to scale to in response to a scale-out activity.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},
			{
				Name:        "role_arn",
				Description: "The ARN of an IAM role that allows Application Auto Scaling to modify the scalable target on your behalf.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},
			{
				Name:        "suspended_state",
				Description: "Specifies whether the scaling activities for a scalable target are in a suspended state.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsApplicationAutoScalingTarget,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceNamespace"),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func getAwsApplicationAutoScalingTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsApplicationAutoScalingTarget")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	name := d.KeyColumnQuals["service_namespace"].GetStringValue()
	ID := d.KeyColumnQuals["resource_id"].GetStringValue()

	// create service
	svc, err := ApplicationAutoScalingService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &applicationautoscaling.DescribeScalableTargetsInput{
		ServiceNamespace: &name,
		ResourceIds:      []*string{types.String("table/" + ID)},
	}

	// Get call
	op, err := svc.DescribeScalableTargets(params)
	if err != nil {
		return nil, err
	}
	if len(op.ScalableTargets) > 0 {
		return op.ScalableTargets[0], nil
	}

	return nil, nil
}
