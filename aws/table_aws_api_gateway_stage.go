package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	apigatewayEndpointId "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayStage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_stage",
		Description: "AWS API Gateway Stage",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"rest_api_id", "name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getAPIGatewayStage,
			Tags:    map[string]string{"service": "apigateway", "action": "GetStage"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listRestAPI,
			Hydrate:       listAPIGatewayStage,
			Tags:          map[string]string{"service": "apigateway", "action": "GetStages"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(apigatewayEndpointId.APIGATEWAYServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the stage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.StageName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the  stage.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAPIGatewayStageARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "rest_api_id",
				Description: "The id of the rest api which contains this stage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RestAPIId"),
			},
			{
				Name:        "deployment_id",
				Description: "The identifier of the Deployment that the stage points to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.DeploymentId"),
			},
			{
				Name:        "created_date",
				Description: "The timestamp when the stage was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Stage.CreatedDate"),
			},
			{
				Name:        "cache_cluster_enabled",
				Description: "Specifies whether a cache cluster is enabled for the stage.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.CacheClusterEnabled"),
			},
			{
				Name:        "tracing_enabled",
				Description: "Specifies whether active tracing with X-ray is enabled for the Stage.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Stage.TracingEnabled"),
			},
			{
				Name:        "access_log_settings",
				Description: "Settings for logging access in this stage.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.AccessLogSettings"),
			},
			{
				Name:        "cache_cluster_size",
				Description: "The size of the cache cluster for the stage, if enabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.CacheClusterSize"),
			},
			{
				Name:        "cache_cluster_status",
				Description: "The status of the cache cluster for the stage, if enabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.CacheClusterStatus"),
			},
			{
				Name:        "client_certificate_id",
				Description: "The identifier of a client certificate for an API stage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.ClientCertificateId"),
			},
			{
				Name:        "description",
				Description: "The stage's description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.Description"),
			},
			{
				Name:        "documentation_version",
				Description: "The version of the associated API documentation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.DocumentationVersion"),
			},
			{
				Name:        "last_updated_date",
				Description: "The timestamp when the stage last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Stage.LastUpdatedDate"),
			},
			{
				Name:        "canary_settings",
				Description: "A map of settings for the canary deployment in this stage.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.CanarySettings"),
			},
			{
				Name:        "method_settings",
				Description: "A map that defines the method settings for a Stage resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.MethodSettings"),
			},
			{
				Name:        "variables",
				Description: "A map that defines the stage variables for a Stage resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.Variables"),
			},
			{
				Name:        "web_acl_arn",
				Description: "The ARN of the WebAcl associated with the Stage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.WebAclArn"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Stage.StageName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Stage.Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAPIGatewayStageARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type stageRowData = struct {
	Stage     types.Stage
	RestAPIId *string
}

//// LIST FUNCTION

func listAPIGatewayStage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get Rest API details
	restAPI := h.Item.(types.RestApi)

	// Create Session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_stage.listAPIGatewayStage", "service_client_error", err)
		return nil, err
	}

	params := &apigateway.GetStagesInput{
		RestApiId: restAPI.Id,
	}

	op, err := svc.GetStages(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_stage.listAPIGatewayStage", "api_error", err)
		return nil, err
	}

	for _, stage := range op.Item {
		d.StreamLeafListItem(ctx, &stageRowData{stage, restAPI.Id})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAPIGatewayStage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_stage.getAPIGatewayStage", "service_client_error", err)
		return nil, err
	}

	stageName := d.EqualsQuals["name"].GetStringValue()
	restAPIID := d.EqualsQuals["rest_api_id"].GetStringValue()

	params := &apigateway.GetStageInput{
		RestApiId: aws.String(restAPIID),
		StageName: aws.String(stageName),
	}

	stageData, err := svc.GetStage(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_stage.getAPIGatewayStage", "api_error", err)
		return nil, err
	}

	return &stageRowData{types.Stage{
		AccessLogSettings:    stageData.AccessLogSettings,
		CacheClusterEnabled:  stageData.CacheClusterEnabled,
		CacheClusterSize:     stageData.CacheClusterSize,
		CacheClusterStatus:   stageData.CacheClusterStatus,
		CanarySettings:       stageData.CanarySettings,
		ClientCertificateId:  stageData.ClientCertificateId,
		CreatedDate:          stageData.CreatedDate,
		DeploymentId:         stageData.DeploymentId,
		Description:          stageData.Description,
		DocumentationVersion: stageData.DocumentationVersion,
		LastUpdatedDate:      stageData.LastUpdatedDate,
		MethodSettings:       stageData.MethodSettings,
		StageName:            stageData.StageName,
		Tags:                 stageData.Tags,
		TracingEnabled:       stageData.TracingEnabled,
		Variables:            stageData.Variables,
		WebAclArn:            stageData.WebAclArn,
	}, aws.String(restAPIID)}, nil
}

func getAPIGatewayStageARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	apiStage := h.Item.(*stageRowData)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/restapis/" + *apiStage.RestAPIId + "/stages/" + *apiStage.Stage.StageName
	return arn, nil
}
