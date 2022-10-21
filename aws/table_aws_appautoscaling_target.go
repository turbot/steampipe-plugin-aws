package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
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
			Hydrate: listAwsApplicationAutoScalingTargets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "service_namespace",
					Require: plugin.Required,
				},
				{
					Name:    "resource_id",
					Require: plugin.Optional,
				},
				{
					Name:    "scalable_dimension",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
		plugin.Logger(ctx).Error("listAwsApplicationAutoScalingTargets", "connection_error", err)
		return nil, err
	}

	// The maximum number for MaxResults parameter is not defined by the API
	// We have set the MaxResults to 1000 based on our test
	input := &applicationautoscaling.DescribeScalableTargetsInput{
		MaxResults: aws.Int64(1000),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	input.ServiceNamespace = types.String(name)

	if equalQuals["resource_id"] != nil {
		input.ResourceIds = []*string{types.String(equalQuals["resource_id"].GetStringValue())}
	}
	if equalQuals["scalable_dimension"] != nil {
		input.ScalableDimension = types.String(equalQuals["scalable_dimension"].GetStringValue())
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = types.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeScalableTargetsPages(
		input,
		func(page *applicationautoscaling.DescribeScalableTargetsOutput, isLast bool) bool {
			for _, scalableTarget := range page.ScalableTargets {
				d.StreamListItem(ctx, scalableTarget)

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
