package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsCodeDeployDeploymentGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codedeploy_deployment_group",
		Description: "AWS Code Deploy Deployment Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"application_name", "deployment_group_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidApplicationNameException", "ApplicationDoesNotExistException", "InvalidDeploymentGroupNameException", "DeploymentGroupDoesNotExistException"}),
			},
			Hydrate: getCodeDeployDeploymentGroup,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCodeDeployApplications,
			Hydrate:       listCodeDeployDeploymentGroups,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "deployment_group_name",
				Description: "The deployment group name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_group_id",
				Description: "The deployment group ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "application_name",
				Description: "The application name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentGroupArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "compute_platform",
				Description: "The destination platform type for the deployment (Lambda, Server, or ECS).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "deployment_config_name",
				Description: "The deployment configuration name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "outdated_instances_strategy",
				Description: "Indicates what happens when new EC2 instances are launched mid-deployment and do not receive the deployed application revision.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "service_role_arn",
				Description: "A service role Amazon Resource Name (ARN) that grants CodeDeploy permission to make calls to AWS services on your behalf.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentGroup,
			},

			// JSON columns
			{
				Name:        "alarm_configuration",
				Description: "A list of alarms associated with the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "auto_rollback_configuration",
				Description: "Information about the automatic rollback configuration associated with the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "auto_scaling_groups",
				Description: "A list of associated Auto Scaling groups.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "blue_green_deployment_configuration",
				Description: "Information about blue/green deployment options for a deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "deployment_style",
				Description: "Information about the type of deployment, either in-place or blue/green, you want to run and whether to route deployment traffic behind a load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "ec2_tag_filters",
				Description: "The Amazon EC2 tags on which to filter.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "ec2_tag_set",
				Description: "Information about groups of tags applied to an EC2 instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "ecs_services",
				Description: "The target Amazon ECS services in the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "last_attempted_deployment",
				Description: "Information about the most recent attempted deployment to the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "last_successful_deployment",
				Description: "Information about the most recent successful deployment to the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "load_balancer_info",
				Description: "Information about the load balancer to use in a deployment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "on_premises_instance_tag_filters",
				Description: "The on-premises instance tags on which to filter.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "on_premises_instance_tag_set",
				Description: "Information about groups of tags applied to an on-premises instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "target_revision",
				Description: "Information about the deployment group's target revision, including type and location.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},
			{
				Name:        "trigger_configurations",
				Description: "Information about triggers associated with the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroup,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroupArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodeDeployDeploymentGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	applicationName := h.Item.(*codedeploy.ApplicationInfo).ApplicationName

	// Create session
	svc, err := CodeDeployService(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.listCodeDeployDeploymentGroups", "service_creation_error", err)
		return nil, err
	}

	input := &codedeploy.ListDeploymentGroupsInput{
		ApplicationName: applicationName,
	}

	// List call
	err = svc.ListDeploymentGroupsPages(
		input,
		func(page *codedeploy.ListDeploymentGroupsOutput, isLast bool) bool {
			for _, group := range page.DeploymentGroups {
				item := &codedeploy.DeploymentGroupInfo{
					DeploymentGroupName: group,
					ApplicationName:     page.ApplicationName,
				}
				d.StreamListItem(ctx, item)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
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

func getCodeDeployDeploymentGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var applicationName, deploymentGroupName string
	if h.Item != nil {
		applicationName = *h.Item.(*codedeploy.DeploymentGroupInfo).ApplicationName
		deploymentGroupName = *h.Item.(*codedeploy.DeploymentGroupInfo).DeploymentGroupName
	} else {
		deploymentGroupName = d.KeyColumnQuals["deployment_group_name"].GetStringValue()
		applicationName = d.KeyColumnQuals["application_name"].GetStringValue()
	}

	if deploymentGroupName == "" || applicationName == "" {
		return nil, nil
	}

	// Build the params
	params := &codedeploy.GetDeploymentGroupInput{
		DeploymentGroupName: &deploymentGroupName,
		ApplicationName:     &applicationName,
	}

	// Create session
	svc, err := CodeDeployService(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.getCodeDeployDeploymentGroup", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.GetDeploymentGroup(params)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.getCodeDeployDeploymentGroup", "api_error", err)
		return nil, err
	}
	return data.DeploymentGroupInfo, nil
}

func getCodeDeployDeploymentGroupArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	applicationName := *h.Item.(*codedeploy.DeploymentGroupInfo).ApplicationName
	deploymentGroupName := *h.Item.(*codedeploy.DeploymentGroupInfo).DeploymentGroupName

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.getCodeDeployDeploymentGroupArn", "caching_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	//arn:aws:codedeploy:region:account-id:deploymentgroup:application-name/deployment-group-name
	tableArn := "arn:" + commonColumnData.Partition + ":codedeploy:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":deploymentgroup/" + applicationName + "/" + deploymentGroupName
	return tableArn, nil
}
