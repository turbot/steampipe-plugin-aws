package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsCodeDeployDeploymentConfig(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codedeploy_deployment_config",
		Description: "AWS Code Deploy Deployment Config",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"deployment_config_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidDeploymentConfigNameException", "DeploymentConfigDoesNotExistException"}),
			},
			Hydrate: getCodeDeployDeploymentConfig,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeDeployDeploymentConfigs,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "deployment_config_name",
				Description: "The deployment configuration name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_config_id",
				Description: "The deployment configuration ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentConfig,
			},
			{
				Name:        "compute_platform",
				Description: "The destination platform type for the deployment (Lambda, Server, or ECS).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentConfig,
			},
			{
				Name:        "create_time",
				Description: "The time at which the deployment configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeDeployDeploymentConfig,
			},

			// JSON columns
			{
				Name:        "minimum_healthy_hosts",
				Description: "Information about the number or percentage of minimum healthy instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentConfig,
			},
			{
				Name:        "traffic_routing_config",
				Description: "The configuration that specifies how the deployment traffic is routed.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentConfig,
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
				Hydrate:     getCodeDeployDeploymentConfigArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodeDeployDeploymentConfigs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create session
	svc, err := CodeDeployService(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_config.listCodeDeployDeploymentConfigs", "service_creation_error", err)
		return nil, err
	}

	input := &codedeploy.ListDeploymentConfigsInput{}

	// List call
	err = svc.ListDeploymentConfigsPages(
		input,
		func(page *codedeploy.ListDeploymentConfigsOutput, isLast bool) bool {
			for _, config := range page.DeploymentConfigsList {
				item := &codedeploy.DeploymentConfigInfo{
					DeploymentConfigName: config,
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

func getCodeDeployDeploymentConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var deploymentConfigName string
	if h.Item != nil {
		deploymentConfigName = *h.Item.(*codedeploy.DeploymentConfigInfo).DeploymentConfigName
	} else {
		deploymentConfigName = d.KeyColumnQuals["deployment_config_name"].GetStringValue()
	}

	if deploymentConfigName == "" {
		return nil, nil
	}

	// Build the params
	params := &codedeploy.GetDeploymentConfigInput{
		DeploymentConfigName: &deploymentConfigName,
	}

	// Create session
	svc, err := CodeDeployService(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_config.getCodeDeployDeploymentConfig", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.GetDeploymentConfig(params)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_config.getCodeDeployDeploymentConfig", "api_error", err)
		return nil, err
	}
	return data.DeploymentConfigInfo, nil
}

func getCodeDeployDeploymentConfigArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	logger := plugin.Logger(ctx)
	deploymentConfigName := *h.Item.(*codedeploy.DeploymentConfigInfo).DeploymentConfigName

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_config.getCodeDeployDeploymentConfigArn", "caching_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// arn:aws:codedeploy:region:account-id:deploymentconfig:deployment-configuration-name
	tableArn := "arn:" + commonColumnData.Partition + ":codedeploy:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":deploymentconfig/" + deploymentConfigName
	return tableArn, nil
}
