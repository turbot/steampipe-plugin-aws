package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCEAnomalyMonitor(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ce_anomaly_monitor",
		Description: "AWS Cost Explorer Anomaly Monitor",
		List: &plugin.ListConfig{
			Hydrate: listCEAnomalyMonitors,
			Tags:    map[string]string{"service": "ce", "action": "GetAnomalyMonitors"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"monitor_arn"}),
			Hydrate:    getCEAnomalyMonitor,
			Tags:       map[string]string{"service": "ce", "action": "GetAnomalyMonitors"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "monitor_arn",
				Description: "The ARN of the anomaly monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorArn"),
			},
			{
				Name:        "monitor_name",
				Description: "The name of the anomaly monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorName"),
			},
			{
				Name:        "monitor_type",
				Description: "The type of the monitor (DIMENSIONAL or CUSTOM).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorType"),
			},
			{
				Name:        "monitor_dimension",
				Description: "The dimension to monitor for anomalies.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MonitorDimension"),
			},
			{
				Name:        "monitor_specification",
				Description: "The monitor specification with cost categories, tags, or dimensions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the anomaly monitor.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationDate"),
			},
			{
				Name:        "last_updated_date",
				Description: "The last update date of the anomaly monitor.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdatedDate"),
			},
			{
				Name:        "last_evaluated_date",
				Description: "The last evaluation date of the anomaly monitor.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastEvaluatedDate"),
			},
			{
				Name:        "dimensional_value_count",
				Description: "The number of evaluated dimensions.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DimensionalValueCount"),
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

func listCEAnomalyMonitors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ce_anomaly_monitor.listCEAnomalyMonitors", "connection_error", err)
		return nil, err
	}

	input := &costexplorer.GetAnomalyMonitorsInput{
		MaxResults: aws.Int32(100),
	}

	output, err := client.GetAnomalyMonitors(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ce_anomaly_monitor.listCEAnomalyMonitors", "api_error", err)
		return nil, err
	}

	for _, monitor := range output.AnomalyMonitors {
		d.StreamListItem(ctx, monitor)

		// Context may get cancelled
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getCEAnomalyMonitor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	monitorArn := d.EqualsQualString("monitor_arn")
	if monitorArn == "" {
		return nil, nil
	}

	client, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ce_anomaly_monitor.getCEAnomalyMonitor", "connection_error", err)
		return nil, err
	}

	// Use MonitorArnList to filter by specific ARN
	input := &costexplorer.GetAnomalyMonitorsInput{
		MonitorArnList: []string{monitorArn},
	}

	output, err := client.GetAnomalyMonitors(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ce_anomaly_monitor.getCEAnomalyMonitor", "api_error", err)
		return nil, err
	}

	if len(output.AnomalyMonitors) > 0 {
		return output.AnomalyMonitors[0], nil
	}

	return nil, nil
}
