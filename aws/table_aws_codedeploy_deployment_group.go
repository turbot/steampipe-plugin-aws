package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy/types"

	codedeployv1 "github.com/aws/aws-sdk-go/service/codedeploy"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeDeployDeploymentGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codedeploy_deployment_group",
		Description: "AWS CodeDeploy Deployment Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"deployment_group_name", "application_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ApplicationDoesNotExistException", "DeploymentGroupDoesNotExistException"}),
			},
			Hydrate: getDeploymentGroup,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCodeDeployApplications,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ApplicationDoesNotExistException"}),
			},
			Hydrate: listDeploymentGroups,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codedeployv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentGroupArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "alarm_configuration",
				Description: "A list of alarms associated with the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "application_name",
				Description: "The application name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_rollback_configuration",
				Description: "Information about the automatic rollback configuration associated with the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "auto_scaling_groups",
				Description: "A list of associated Auto Scaling groups.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "blue_green_deployment_configuration",
				Description: "Information about blue/green deployment options for a deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "compute_platform",
				Description: "The destination platform type for the deployment (Lambda, Server, or ECS).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "deployment_config_name",
				Description: "The deployment configuration name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "deployment_group_id",
				Description: "The deployment group ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "deployment_group_name",
				Description: "The deployment group name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "deployment_style",
				Description: "Information about the type of deployment, either in-place or blue/green, you want to run and whether to route deployment traffic behind a load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "ec2_tag_filters",
				Description: "The Amazon EC2 tags on which to filter. The deployment group includes EC2 instances with any of the specified tags.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "ec2_tag_set",
				Description: "Information about groups of tags applied to an Amazon EC2 instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "ecs_services",
				Description: "The target Amazon ECS services in the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "last_attempted_deployment",
				Description: "Information about the most recent attempted deployment to the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "last_successful_deployment",
				Description: "Information about the most recent successful deployment to the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "load_balancer_info",
				Description: "Information about the load balancer to use in a deployment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "on_premises_instance_tag_filters",
				Description: "The on-premises instance tags on which to filter.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "on_premises_tag_set",
				Description: "Information about groups of tags applied to an on-premises instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "outdated_instances_strategy",
				Description: "Indicates what happens when new Amazon EC2 instances are launched mid-deployment and do not receive the deployed application revision.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "service_role_arn",
				Description: "A service role Amazon Resource Name (ARN) that grants CodeDeploy permission to make calls to Amazon Web Services services on your behalf.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "target_revision",
				Description: "Information about the deployment group's target revision, including type and location.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},
			{
				Name:        "trigger_configurations",
				Description: "Information about triggers associated with the deployment group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDeploymentGroup,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeploymentGroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentGroupTags,
				Transform:   transform.From(codeDeployDeploymentGroupTurbotTags),
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

func listDeploymentGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	application := h.Item.(*types.ApplicationInfo)

	// Create session
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.listDeploymentGroups", "service_creation_error", err)
		return nil, err
	}
	applicationName := d.EqualsQuals["application_name"].GetStringValue()
	input := codedeploy.ListDeploymentGroupsInput{
		ApplicationName: aws.String(*application.ApplicationName),
	}

	if applicationName != "" && applicationName !=  *application.ApplicationName {
		return nil,nil
	}

	paginator := codedeploy.NewListDeploymentGroupsPaginator(svc, &input, func(o *codedeploy.ListDeploymentGroupsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codedeploy_deployment_group.listDeploymentGroups", "api_error", err)
			return nil, err
		}

		for _, deploymentgroup := range output.DeploymentGroups {
			item := &types.DeploymentGroupInfo{
				DeploymentGroupName: aws.String(deploymentgroup),
				ApplicationName: output.ApplicationName,
			}
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDeploymentGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name string
	var appname string
	if h.Item != nil {
		name = *h.Item.(*types.DeploymentGroupInfo).DeploymentGroupName
		appname = *h.Item.(*types.DeploymentGroupInfo).ApplicationName
	} else {
		name = d.EqualsQuals["deployment_group_name"].GetStringValue()
		appname = d.EqualsQuals["application_name"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	if appname == "" {
		return nil, nil
	}

	// Build the params
	params := &codedeploy.GetDeploymentGroupInput{
		DeploymentGroupName: aws.String(name),
		ApplicationName: aws.String(appname),
	}

	// Create session
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.getDeploymentGroup", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.GetDeploymentGroup(ctx, params)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.getDeploymentGroup", "api_error", err)
		return nil, err
	}
	return data.DeploymentGroupInfo, nil
}
func getCodeDeployDeploymentGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var name string
	if h.Item != nil {
		name = *h.Item.(*types.DeploymentGroupInfo).DeploymentGroupName
	} else {
		name = d.EqualsQuals["deployment_group_name"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &codedeploy.ListTagsForResourceInput{
		ResourceArn: aws.String(CodeDeployDeploymentGroupArn(ctx, d, h)),
	}

	// Create session
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.getCodeDeployDeploymentGroupTags", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.getCodeDeployDeploymentGroupTags", "api_error", err)
		return nil, err
	}
	return data, nil

}

func getCodeDeployDeploymentGroupArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return CodeDeployDeploymentGroupArn(ctx, d, h), nil
}

func CodeDeployDeploymentGroupArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) string {
	name := *h.Item.(*types.DeploymentGroupInfo).DeploymentGroupName
	appname := *h.Item.(*types.DeploymentGroupInfo).ApplicationName
	region := d.EqualsQualString(matrixKeyRegion)
	logger := plugin.Logger(ctx)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_group.CodeDeployDeploymentGroupArn", "caching_error", err)
		return ""
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	//arn:aws:codedeploy:region:account-id:deploymentgroup:application-name/deployment-group-name
	tableArn := "arn:" + commonColumnData.Partition + ":codedeploy:" + region + ":" + commonColumnData.AccountId + ":deploymentgroup:" + appname + "/" + name
	return tableArn
}

func codeDeployDeploymentGroupTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*codedeploy.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}