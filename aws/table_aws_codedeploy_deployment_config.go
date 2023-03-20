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

func tableAwsCodeDeployDeploymentConfig(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codedeploy_deployment_config",
		Description: "AWS CodeDeploy Deployment Config",
		Get: &plugin.GetConfig{
			KeyColumns:   plugin.SingleColumn("deployment_config_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DeploymentConfigDoesNotExistException"}),
			},
			Hydrate: getCodeDeployDeploymentConfig,
		},
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DeploymentConfigDoesNotExistException"}),
			},
			Hydrate: listCodeDeployDeploymentConfigs,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codedeployv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentConfigArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "compute_platform",
				Description: "The destination platform type for the deployment (Lambda, Server, or ECS).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentConfig,
			},
			{
				Name:        "create_time",
				Description: "The time at which the deployment configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeDeployDeploymentConfig,
			},
			{
				Name:        "deployment_config_id",
				Description: "The deployment configuration ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployDeploymentConfig,
			},
			{
				Name:        "deployment_config_name",
				Description: "The deployment configuration name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "minimum_healthy_hosts",
				Description: "Information about the number or percentage of minimum healthy instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentConfig,
			},
			{
				Name:        "traffic_routing_config",
				Description: "The configuration that specifies how the deployment traffic is routed. Used for deployments with a Lambda or Amazon ECS compute platform only.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeDeployDeploymentConfig,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeploymentConfigName"),
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
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_config.listDeploymentConfigs", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}
	
	input := codedeploy.ListDeploymentConfigsInput{}

	paginator := codedeploy.NewListDeploymentConfigsPaginator(svc, &input, func(o *codedeploy.ListDeploymentConfigsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codedeploy_deployment_config.listDeploymentConfigs", "api_error", err)
			return nil, err
		}

		for _, deploymentconfig := range output.DeploymentConfigsList {
			d.StreamListItem(ctx, &types.DeploymentConfigInfo{
				DeploymentConfigName: aws.String(deploymentconfig),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeDeployDeploymentConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(*types.DeploymentConfigInfo).DeploymentConfigName
	} else {
		name = d.EqualsQuals["deployment_config_name"].GetStringValue()
	}

	// check if clusterName or nodegroupName is empty
	if name == "" {
		return nil, nil
	}

	// create service
	svc, err := CodeDeployClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codedeploy_deployment_config.getCodeDeployDeploymentConfig", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &codedeploy.GetDeploymentConfigInput{
		DeploymentConfigName: &name,
	}

	op, err := svc.GetDeploymentConfig(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codedeploy_deployment_config.getCodeDeployDeploymentConfig", "api_error", err)
		return nil, err
	}

	return op.DeploymentConfigInfo, nil
}

func getCodeDeployDeploymentConfigArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(*types.DeploymentConfigInfo).DeploymentConfigName
	} else {
		name = d.EqualsQuals["deployment_config_name"].GetStringValue()
	}
	region := d.EqualsQualString(matrixKeyRegion)
	logger := plugin.Logger(ctx)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		logger.Error("aws_codedeploy_deployment_config.CodeDeployDeploymentConfigArn", "caching_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	//arn:aws:codedeploy:region:account-id:deploymentconfig:deployment-configuration-name
	tableArn := "arn:" + commonColumnData.Partition + ":codedeploy:" + region + ":" + commonColumnData.AccountId + ":deploymentconfig:" + name
	return tableArn, nil
}
