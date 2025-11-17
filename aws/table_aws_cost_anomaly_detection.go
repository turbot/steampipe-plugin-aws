package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCostAnomalyDetection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_anomaly_detection",
		Description: "AWS Cost Anomaly Detection",
		List: &plugin.ListConfig{
			Hydrate: listCostAnomalyDetections,
			Tags:    map[string]string{"service": "ce", "action": "ListCostAnomalyDetectors"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"monitor_arn"}),
			Hydrate:    getCostAnomalyDetection,
			Tags:       map[string]string{"service": "ce", "action": "DescribeCostAnomalyDetectors"},
		},
		HydrateConfig: []plugin.HydrateConfig{
		{
			Func: getCostAnomalyDetectionDetails,
			Tags: map[string]string{"service": "ce", "action": "DescribeCostAnomalyDetectors"},
		},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "monitor_arn",
				Description: "The ARN of the Cost Anomaly Detection monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorArn"),
			},
			{
				Name:        "name",
				Description: "The name of the Cost Anomaly Detection monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorName"),
			},
			{
				Name:        "status",
				Description: "The status of the Cost Anomaly Detection monitor (ACTIVE or INACTIVE).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorStatus"),
			},
			{
				Name:        "frequency",
				Description: "The frequency with which anomalies are analyzed (DAILY).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorFrequency"),
			},
			{
				Name:        "monitor_specification",
				Description: "The dimensions and tags used to detect cost anomalies.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCostAnomalyDetectionDetails,
				Transform:   transform.FromField("MonitorSpecification"),
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the anomaly monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "last_modified_date",
				Description: "The last modification date of the anomaly monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastModifiedDate"),
			},
			{
				Name:        "last_evaluation_date",
				Description: "The last evaluation date of the anomaly monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastEvaluationDate"),
			},
			{
				Name:        "next_evaluation_date",
				Description: "The next evaluation date of the anomaly monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NextEvaluationDate"),
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MonitorArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCostAnomalyDetections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_anomaly_detection.listCostAnomalyDetections", "connection_error", err)
		return nil, err
	}

	input := &costexplorer.ListCostAnomalyDetectorsInput{}

	paginator := costexplorer.NewListCostAnomalyDetectorsPaginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ce_cost_anomaly_detection.listCECostAnomalyDetections", "api_error", err)
			return nil, err
		}

		for _, detector := range output.CostAnomalyDetectors {
			d.StreamListItem(ctx, detector)

			// Context may get cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getCostAnomalyDetection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	monitorArn := d.EqualsQualString("monitor_arn")
	if monitorArn == "" {
		return nil, nil
	}

	client, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_anomaly_detection.getCostAnomalyDetection", "connection_error", err)
		return nil, err
	}

	input := &costexplorer.DescribeCostAnomalyDetectorsInput{
		MonitorArn: aws.String(monitorArn),
	}

	output, err := client.DescribeCostAnomalyDetectors(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ce_cost_anomaly_detection.getCECostAnomalyDetection", "api_error", err)
		return nil, err
	}

	if len(output.CostAnomalyDetectors) > 0 {
		return output.CostAnomalyDetectors[0], nil
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCostAnomalyDetectionDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	detector := h.Item.(types.CostAnomalyDetector)

	client, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cost_anomaly_detection.getCostAnomalyDetectionDetails", "connection_error", err)
		return nil, err
	}

	input := &costexplorer.DescribeCostAnomalyDetectorsInput{
		MonitorArn: detector.MonitorArn,
	}

	output, err := client.DescribeCostAnomalyDetectors(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ce_cost_anomaly_detection.getCECostAnomalyDetectionDetails", "api_error", err)
		return nil, err
	}

	if len(output.CostAnomalyDetectors) > 0 {
		return output.CostAnomalyDetectors[0], nil
	}

	return nil, nil
}

