package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayV2Stage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gatewayv2_stage",
		Description: "AWS API Gateway Version 2 Stage",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"api_id", "stage_name"}),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       apiGatewayV2StageFromKey,
			Hydrate:           getAPIGatewayV2Stage,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAPIGatewayV2API,
			Hydrate:       listAPIGatewayV2Stages,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "stage_name",
				Description: "The name of the stage",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.StageName"),
			},
			{
				Name:        "api_id",
				Description: "The id of the api which contains this stage",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("APIId"),
			},
			{
				Name:        "api_gateway_managed",
				Description: "Specifies whether a stage is managed by API Gateway",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.ApiGatewayManaged"),
			},
			{
				Name:        "auto_deploy",
				Description: "Specifies whether updates to an API automatically trigger a new deployment",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.AutoDeploy"),
			},
			{
				Name:        "client_certificate_id",
				Description: "The identifier of a client certificate for a Stage. Supported only for WebSocket APIs",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.ClientCertificateId"),
			},
			{
				Name:        "created_date",
				Description: "The timestamp when the stage was created",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.CreatedDate"),
			},
			{
				Name:        "deployment_id",
				Description: "The identifier of the Deployment that the Stage is associated with",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.DeploymentId"),
			},
			{
				Name:        "default_route_data_trace_enabled",
				Description: "Specifies whether (true) or not (false) data trace logging is enabled for this route. This property affects the log entries pushed to Amazon CloudWatch Logs. Supported only for WebSocket APIs",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.DefaultRouteSettings.DataTraceEnabled"),
			},
			{
				Name:        "default_route_detailed_metrics_enabled",
				Description: "Specifies whether detailed metrics are enabled",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.DefaultRouteSettings.DetailedMetricsEnabled"),
			},
			{
				Name:        "default_route_logging_level",
				Description: "Specifies the logging level for this route: INFO, ERROR, or OFF. This property affects the log entries pushed to Amazon CloudWatch Logs. Supported only for WebSocket APIs",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.DefaultRouteSettings.LoggingLevel"),
			},
			{
				Name:        "default_route_throttling_burst_limit",
				Description: "Throttling burst limit for default route settings",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Stage.DefaultRouteSettings.ThrottlingBurstLimit"),
			},
			{
				Name:        "default_route_throttling_rate_limit",
				Description: "Throttling rate limit for default route settings",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Stage.DefaultRouteSettings.ThrottlingRateLimit"),
			},
			{
				Name:        "last_deployment_status_message",
				Description: "Describes the status of the last deployment of a stage. Supported only for stages with autoDeploy enabled",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.LastDeploymentStatusMessage"),
			},
			{
				Name:        "last_updated_date",
				Description: "The timestamp when the stage was last updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.LastUpdatedDate"),
			},
			{
				Name:        "description",
				Description: "The stage's description",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.Description"),
			},
			{
				Name:        "stage_variables",
				Description: "A map that defines the stage variables for a stage resource",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.StageVariables"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.StageName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     apiGatewayV2StageAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type v2StageRowData = struct {
	Stage *apigatewayv2.Stage
	APIId *string
}

//// BUILD HYDRATE INPUT

func apiGatewayV2StageFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	stageName := quals["stage_name"].GetStringValue()
	apiID := quals["api_id"].GetStringValue()
	item := &v2StageRowData{
		APIId: &apiID,
		Stage: &apigatewayv2.Stage{
			StageName: &stageName,
		},
	}

	return item, nil
}

//// LIST FUNCTION

func listAPIGatewayV2Stages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var stages []*apigatewayv2.Stage

	plugin.Logger(ctx).Trace("listAPIGatewayV2Stages", "AWS_REGION", region)
	apiGatewayv2API := h.Item.(*apigatewayv2.Api)

	svc, err := APIGatewayV2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	pagesLeft := true
	params := &apigatewayv2.GetStagesInput{
		ApiId: apiGatewayv2API.ApiId,
	}

	for pagesLeft {
		result, err := svc.GetStages(params)
		if err != nil {
			return nil, err
		}

		stages = append(stages, result.Items...)
		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	for _, stage := range stages {
		d.StreamLeafListItem(ctx, &v2StageRowData{stage, apiGatewayv2API.ApiId})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAPIGatewayV2Stage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAPIGatewayStage")
	v2ApiStageRowData := h.Item.(*v2StageRowData)
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := APIGatewayV2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &apigatewayv2.GetStageInput{
		ApiId:     v2ApiStageRowData.APIId,
		StageName: v2ApiStageRowData.Stage.StageName,
	}

	stageData, err := svc.GetStage(input)
	if err != nil {
		logger.Debug("getAPIGatewayStage__", "ERROR", err)
		return nil, err
	}
	if stageData != nil {
		stage := &apigatewayv2.Stage{
			StageName:                   stageData.StageName,
			AccessLogSettings:           stageData.AccessLogSettings,
			ApiGatewayManaged:           stageData.ApiGatewayManaged,
			AutoDeploy:                  stageData.AutoDeploy,
			ClientCertificateId:         stageData.ClientCertificateId,
			CreatedDate:                 stageData.CreatedDate,
			DefaultRouteSettings:        stageData.DefaultRouteSettings,
			DeploymentId:                stageData.DeploymentId,
			Description:                 stageData.Description,
			LastDeploymentStatusMessage: stageData.LastDeploymentStatusMessage,
			LastUpdatedDate:             stageData.LastUpdatedDate,
			RouteSettings:               stageData.RouteSettings,
			StageVariables:              stageData.StageVariables,
			Tags:                        stageData.Tags,
		}
		rowData := &v2StageRowData{stage, v2ApiStageRowData.APIId}

		return rowData, nil
	}

	return nil, nil
}

func apiGatewayV2StageAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*v2StageRowData)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + "::/apis/" + *data.APIId + "/stages/" + *data.Stage.StageName}

	return akas, nil
}
