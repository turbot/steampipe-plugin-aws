package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling/types"

	applicationautoscalingv1 "github.com/aws/aws-sdk-go/service/applicationautoscaling"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		GetMatrixItemFunc: SupportedRegionMatrix(applicationautoscalingv1.EndpointsID),
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
	name := d.EqualsQuals["service_namespace"].GetStringValue()

	// Create Session
	svc, err := ApplicationAutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appautoscaling_target.listAwsApplicationAutoScalingTargets", "get_client_error", err)
		return nil, err
	}

	// The maximum number for MaxResults parameter is not defined by the API
	// We have set the MaxResults to 1000 based on our test
	input := &applicationautoscaling.DescribeScalableTargetsInput{
		MaxResults: aws.Int32(1000),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	input.ServiceNamespace = types.ServiceNamespace(name)

	if equalQuals["resource_id"] != nil {
		input.ResourceIds = []string{equalQuals["resource_id"].GetStringValue()}
	}
	if equalQuals["scalable_dimension"] != nil {
		input.ScalableDimension = types.ScalableDimension(equalQuals["scalable_dimension"].GetStringValue())
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxResults {
			if limit < 1 {
				input.MaxResults = aws.Int32(1)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	paginator := applicationautoscaling.NewDescribeScalableTargetsPaginator(svc, input, func(o *applicationautoscaling.DescribeScalableTargetsPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_neptune_db_cluster.listNeptuneDBClusters", "api_error", err)
			return nil, err
		}

		for _, scalableTarget := range output.ScalableTargets {
			d.StreamListItem(ctx, scalableTarget)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsApplicationAutoScalingTarget(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["service_namespace"].GetStringValue()
	id := d.EqualsQuals["resource_id"].GetStringValue()

	// create service
	svc, err := ApplicationAutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appautoscaling_target.getAwsApplicationAutoScalingTarget", "get_client_error", err)
		return nil, err
	}

	// Build the params
	params := &applicationautoscaling.DescribeScalableTargetsInput{
		ServiceNamespace: types.ServiceNamespace(name),
		ResourceIds:      []string{(id)},
	}

	// Get call
	op, err := svc.DescribeScalableTargets(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appautoscaling_target.getAwsApplicationAutoScalingTarget", "api_error", err)
		return nil, err
	}
	if op.ScalableTargets != nil && len(op.ScalableTargets) > 0 {
		return op.ScalableTargets[0], nil
	}

	return nil, nil
}
