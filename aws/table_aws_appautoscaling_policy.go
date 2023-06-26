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

func tableAwsAppAutoScalingPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appautoscaling_policy",
		Description: "AWS Application Auto Scaling Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"service_namespace", "policy_name"}),
			Hydrate:    getAwsApplicationAutoScalingPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsApplicationAutoScalingPolicies,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "service_namespace",
					Require: plugin.Required,
				},
				{
					Name:    "policy_ARN",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(applicationautoscalingv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "policy_ARN",
				Description: "The Amazon Resource Name (ARN) of the appautoscaling policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_name",
				Description: "The name of the scaling policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_namespace",
				Description: "The namespace of the AWS service that provides the resource, or a custom-resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_id",
				Description: "The identifier of the resource associated with the scaling policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scalable_dimension",
				Description: "The scalable dimension associated with the scaling policy. This string consists of the service namespace, resource type, and scaling property.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_type",
				Description: "The policy type. Currently supported values are TargetTrackingScaling and StepScaling",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_tracking_scaling_policy_configuration",
				Description: "The target tracking scaling policy configuration (if policy type is TargetTrackingScaling).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "step_scaling_policy_configuration",
				Description: "The step tracking scaling policy configuration (if policy type is StepScaling).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "alarms",
				Description: "The CloudWatch alarms associated with the scaling policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "creation_time",
				Description: "The Unix timestamp for when the scaling policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
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

func listAwsApplicationAutoScalingPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["service_namespace"].GetStringValue()

	// Create Session
	svc, err := ApplicationAutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appautoscaling_policy.listAwsApplicationAutoScalingPolicies", "get_client_error", err)
		return nil, err
	}

	// The maximum number for MaxResults parameter is not defined by the API
	// We have set the MaxResults to 1000 based on our test
	input := &applicationautoscaling.DescribeScalingPoliciesInput{
		MaxResults: aws.Int32(1000),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	input.ServiceNamespace = types.ServiceNamespace(name)

	if equalQuals["policy_name"] != nil {
		input.PolicyNames = []string{equalQuals["policy_name"].GetStringValue()}
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

	paginator := applicationautoscaling.NewDescribeScalingPoliciesPaginator(svc, input, func(o *applicationautoscaling.DescribeScalingPoliciesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_appautoscaling_policy.listAwsApplicationAutoScalingPolicies", "api_error", err)
			return nil, err
		}

		for _, scalingPolicy := range output.ScalingPolicies {
			d.StreamListItem(ctx, scalingPolicy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsApplicationAutoScalingPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["service_namespace"].GetStringValue()
	id := d.EqualsQuals["resource_id"].GetStringValue()

	// create service
	svc, err := ApplicationAutoScalingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appautoscaling_policy.getAwsApplicationAutoScalingPolicy", "get_client_error", err)
		return nil, err
	}

	// Build the params
	params := &applicationautoscaling.DescribeScalingPoliciesInput{
		ServiceNamespace: types.ServiceNamespace(name),
		ResourceId:       &id,
	}

	// Get call
	op, err := svc.DescribeScalingPolicies(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_appautoscaling_policy.getAwsApplicationAutoScalingPolicy", "api_error", err)
		return nil, err
	}
	if op.ScalingPolicies != nil && len(op.ScalingPolicies) > 0 {
		return op.ScalingPolicies[0], nil
	}

	return nil, nil
}
