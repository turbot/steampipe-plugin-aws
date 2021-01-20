package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayStage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_stage",
		Description: "AWS API Gateway Stage",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"rest_api_id", "stage_name"}),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       accessKeyFromKey,
			Hydrate:           getAPIGatewayStage,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listRestAPI,
			Hydrate:       listAPIGatewayStage,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the stage",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.StageName"),
			},
			{
				Name:        "rest_api_id",
				Description: "The id of the rest api which contains this stage",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RestAPIId"),
			},
			{
				Name:        "deployment_id",
				Description: "The identifier of the Deployment that the stage points to",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.DeploymentId"),
			},
			{
				Name:        "created_date",
				Description: "The timestamp when the stage was created",
				Type:        proto.ColumnType_DATETIME,
				Transform:   transform.FromField("Stage.CreatedDate"),
			},
			{
				Name:        "cache_cluster_enabled",
				Description: "Specifies whether a cache cluster is enabled for the stage",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.CacheClusterEnabled"),
			},
			{
				Name:        "tracing_enabled",
				Description: "Specifies whether active tracing with X-ray is enabled for the Stage",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.TracingEnabled"),
			},
			{
				Name:        "access_log_settings",
				Description: "Settings for logging access in this stage",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.AccessLogSettings"),
			},
			{
				Name:        "cache_cluster_size",
				Description: "The size of the cache cluster for the stage, if enabled",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.CacheClusterSize"),
			},
			{
				Name:        "cache_cluster_status",
				Description: "The status of the cache cluster for the stage, if enabled",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.CacheClusterStatus"),
			},
			{
				Name:        "client_certificate_id",
				Description: "The identifier of a client certificate for an API stage",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.ClientCertificateId"),
			},
			{
				Name:        "description",
				Description: "The stage's description",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.Description"),
			},
			{
				Name:        "documentation_version",
				Description: "The version of the associated API documentation",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.DocumentationVersion"),
			},
			{
				Name:        "last_updated_date",
				Description: "The timestamp when the stage last updated",
				Type:        proto.ColumnType_DATETIME,
				Transform:   transform.FromField("Stage.LastUpdatedDate"),
			},
			{
				Name:        "variables",
				Description: "A map that defines the stage variables for a Stage resource",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.Variables"),
			},
			{
				Name:        "web_acl_arn",
				Description: "The ARN of the WebAcl associated with the Stage",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.WebAclArn"),
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
				Hydrate:     getAPIGatewayStageAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type stageRowData = struct {
	Stage     *apigateway.Stage
	RestAPIId *string
}

//// BUILD HYDRATE INPUT

func apiStageFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	stageName := quals["stage_name"].GetStringValue()
	RestAPIID := quals["rest_api_id"].GetStringValue()
	item := &stageRowData{
		RestAPIId: &RestAPIID,
		Stage: &apigateway.Stage{
			StageName: &stageName,
		},
	}

	return item, nil
}

//// LIST FUNCTION

func listAPIGatewayStage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listAPIGatewayStage", "AWS_REGION", defaultRegion)
	restAPI := h.Item.(*apigateway.RestApi)

	// Create Session
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetStagesInput{
		RestApiId: restAPI.Id,
	}

	op, err := svc.GetStages(params)
	if err != nil {
		return nil, err
	}

	for _, stage := range op.Item {
		d.StreamLeafListItem(ctx, &stageRowData{stage, restAPI.Id})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAPIGatewayStage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAPIGatewayStage")
	apiStage := h.Item.(stageRowData)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetStageInput{
		RestApiId: apiStage.RestAPIId,
		StageName: apiStage.Stage.StageName,
	}

	stageData, err := svc.GetStage(params)
	if err != nil {
		logger.Debug("getAPIGatewayStage__", "ERROR", err)
		return nil, err
	}

	return &stageRowData{stageData, apiStage.RestAPIId}, nil
}

func getAPIGatewayStageAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAPIGAtewayStageAkas")
	apiStage := h.Item.(*stageRowData)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + "::/restapis/" + *apiStage.RestAPIId + "/stages/" + *apiStage.Stage.StageName}
	return akas, nil
}
