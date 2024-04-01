package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/greengrassv2"
	"github.com/aws/aws-sdk-go-v2/service/greengrassv2/types"

	greengrassv1 "github.com/aws/aws-sdk-go/service/greengrassv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIotGreengrassDeployment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iot_greengrass_deployment",
		Description: "AWS IoT Greengrass Deployment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("deployment_id"),
			Hydrate:    getIoTGreengrassDeployment,
			Tags:       map[string]string{"service": "greengrassv2", "action": "GetDeployment"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIoTGreengrassDeployments,
			Tags:    map[string]string{"service": "greengrassv2", "action": "ListDeployments"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "parent_target_arn", Require: plugin.Optional},
				{Name: "target_arn", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getIoTGreengrassDeployment,
				Tags: map[string]string{"service": "greengrassv2", "action": "GetDeployment"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(greengrassv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "deployment_name",
				Description: "The name of the deployment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_id",
				Description: "The ID of the deployment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_status",
				Description: "The status of the deployment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_latest_for_target",
				Description: "Whether or not the deployment is the latest revision for its target.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "creation_timestamp",
				Description: "The time at which the deployment was created, expressed in ISO 8601 format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "parent_target_arn",
				Description: "The parent deployment's target ARN within a subdeployment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "revision_id",
				Description: "The revision number of the deployment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_arn",
				Description: "The ARN of the target IoT thing or thing group. When creating a subdeployment, the targetARN can only be a thing group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iot_job_arn",
				Description: "The ARN of the IoT job that applies the deployment to target devices.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTGreengrassDeployment,
			},
			{
				Name:        "iot_job_id",
				Description: "The ID of the IoT job that applies the deployment to target devices.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIoTGreengrassDeployment,
			},

			// JSON fields
			{
				Name:        "components",
				Description: "The components to deploy. This is a dictionary, where each key is the name of a component, and each key's value is the version and configuration to deploy for that component.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTGreengrassDeployment,
			},
			{
				Name:        "deployment_policies",
				Description: "The deployment policies for the deployment. These policies define how the deployment updates components and handles failure.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTGreengrassDeployment,
			},
			{
				Name:        "iot_job_configuration",
				Description: " The job configuration for the deployment configuration. The job configuration specifies the rollout, timeout, and stop configurations for the deployment configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTGreengrassDeployment,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeploymentName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIoTGreengrassDeployment,
			},
		}),
	}
}

//// LIST FUNCTION

func listIoTGreengrassDeployments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IoTGreengrassClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_greengrass_deployment.listIoTGreengrassDeployments", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &greengrassv2.ListDeploymentsInput{
		MaxResults: aws.Int32(maxLimit),
	}
	if d.EqualsQualString("target_arn") != "" {
		input.TargetArn = aws.String(d.EqualsQualString("target_arn"))
	}
	if d.EqualsQualString("parent_target_arn") != "" {
		input.ParentTargetArn = aws.String(d.EqualsQualString("parent_target_arn"))
	}

	paginator := greengrassv2.NewListDeploymentsPaginator(svc, input, func(o *greengrassv2.ListDeploymentsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iot_greengrass_deployment.listIoTGreengrassDeployments", "api_error", err)
			return nil, err
		}

		for _, item := range output.Deployments {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIoTGreengrassDeployment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	deploymentId := ""
	if h.Item != nil {
		t := h.Item.(types.Deployment)
		deploymentId = *t.DeploymentId
	} else {
		deploymentId = d.EqualsQualString("deployment_id")
	}

	if deploymentId == "" {
		return nil, nil
	}

	// Create service
	svc, err := IoTGreengrassClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_greengrass_deployment.getIoTGreengrassDeployment", "connection_error", err)
		return nil, err
	}

	params := &greengrassv2.GetDeploymentInput{
		DeploymentId: aws.String(deploymentId),
	}

	resp, err := svc.GetDeployment(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iot_greengrass_deployment.getIoTGreengrassDeployment", "api_error", err)
		return nil, err
	}

	return resp, nil
}
