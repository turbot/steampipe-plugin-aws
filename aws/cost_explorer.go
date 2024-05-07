package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// AllCostMetrics is a constant returning all the cost metrics
func AllCostMetrics() []string {
	return []string{
		"BlendedCost",
		"UnblendedCost",
		"NetUnblendedCost",
		"AmortizedCost",
		"NetAmortizedCost",
		"UsageQuantity",
		"NormalizedUsageAmount",
	}
}

// getCostMetricByMetricName returns the select metrics
func getCostMetricByMetricName(metricName string) []string {
	metrics := strings.Split(metricName, ",")
	var selectedMetric []string
	for _, m := range metrics {
		switch strings.ToLower(m) {
		case "blendedcost":
			selectedMetric = append(selectedMetric, "BlendedCost")
		case "unblendedcost":
			selectedMetric = append(selectedMetric, "UnblendedCost")
		case "netunblendedcost":
			selectedMetric = append(selectedMetric, "NetUnblendedCost")
		case "amortizedcost":
			selectedMetric = append(selectedMetric, "AmortizedCost")
		case "netamortizedcost":
			selectedMetric = append(selectedMetric, "NetAmortizedCost")
		case "usagequantity":
			selectedMetric = append(selectedMetric, "UsageQuantity")
		case "normalizedusageamount":
			selectedMetric = append(selectedMetric, "NormalizedUsageAmount")
		}
	}

	return selectedMetric
}

func getMetricsByQueryContext(qc *plugin.QueryContext) []string {
	queryColumns := qc.Columns
	var metrics []string

	for _, c := range queryColumns {
		switch c {
		case "blended_cost_amount", "blended_cost_unit":
			metrics = append(metrics, "BlendedCost")
		case "unblended_cost_amount", "unblended_cost_unit":
			metrics = append(metrics, "UnblendedCost")
		case "net_unblended_cost_amount", "net_unblended_cost_unit":
			metrics = append(metrics, "NetUnblendedCost")
		case "amortized_cost_amount", "amortized_cost_unit":
			metrics = append(metrics, "AmortizedCost")
		case "net_amortized_cost_amount", "net_amortized_cost_unit":
			metrics = append(metrics, "NetAmortizedCost")
		case "usage_quantity_amount", "usage_quantity_unit":
			metrics = append(metrics, "UsageQuantity")
		case "normalized_usage_amount", "normalized_usage_unit":
			metrics = append(metrics, "NormalizedUsageAmount")
		}
	}

	return removeDuplicates(metrics)
}

var costExplorerColumnDefs = []*plugin.Column{

	{
		Name:        "period_start",
		Description: "Start timestamp for this cost metric.",
		Type:        proto.ColumnType_TIMESTAMP,
	},
	{
		Name:        "period_end",
		Description: "End timestamp for this cost metric.",
		Type:        proto.ColumnType_TIMESTAMP,
	},

	{
		Name:        "estimated",
		Description: "Whether the result is estimated.",
		Type:        proto.ColumnType_BOOL,
	},
	{
		Name:        "blended_cost_amount",
		Description: "This cost metric reflects the average cost of usage across the consolidated billing family. If you use the consolidated billing feature in AWS Organizations, you can view costs using blended rates.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "blended_cost_unit",
		Description: "Unit type for blended costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "unblended_cost_amount",
		Description: "Unblended costs represent your usage costs on the day they are charged to you. In finance terms, they represent your costs on a cash basis of accounting.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "unblended_cost_unit",
		Description: "Unit type for unblended costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "net_unblended_cost_amount",
		Description: "This cost metric reflects the unblended cost after discounts.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "net_unblended_cost_unit",
		Description: "Unit type for net unblended costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "amortized_cost_amount",
		Description: "This cost metric reflects the effective cost of the upfront and monthly reservation fees spread across the billing period. By default, Cost Explorer shows the fees for Reserved Instances as a spike on the day that you're charged, but if you choose to show costs as amortized costs, the costs are amortized over the billing period. This means that the costs are broken out into the effective daily rate. AWS estimates your amortized costs by combining your unblended costs with the amortized portion of your upfront and recurring reservation fees.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "amortized_cost_unit",
		Description: "Unit type for amortized costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "net_amortized_cost_amount",
		Description: "This cost metric amortizes the upfront and monthly reservation fees while including discounts such as RI volume discounts.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "net_amortized_cost_unit",
		Description: "Unit type for net amortized costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "usage_quantity_amount",
		Description: "The amount of usage that you incurred. NOTE: If you return the UsageQuantity metric, the service aggregates all usage numbers without taking into account the units. For example, if you aggregate usageQuantity across all of Amazon EC2, the results aren't meaningful because Amazon EC2 compute hours and data transfer are measured in different units (for example, hours vs. GB).",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "usage_quantity_unit",
		Description: "Unit type for usage quantity.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "normalized_usage_amount",
		Description: "The amount of usage that you incurred, in normalized units, for size-flexible RIs. The NormalizedUsageAmount is equal to UsageAmount multiplied by NormalizationFactor.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "normalized_usage_unit",
		Description: "Unit type for normalized usage.",
		Type:        proto.ColumnType_STRING,
	},
}

// append the common aws cost explorer columns onto the column list
func costExplorerColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, costExplorerColumnDefs...)
}

//// LIST FUNCTION

func streamCostAndUsage(ctx context.Context, d *plugin.QueryData, params *costexplorer.GetCostAndUsageInput) (interface{}, error) {

	// Create session
	svc, err := CostExplorerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("streamCostAndUsage", "client_error", err)
		return nil, err
	}
	// List call
	for {
		output, err := svc.GetCostAndUsage(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("streamCostAndUsage", "api_error", err)
			return nil, err
		}

		// stream the results...
		for _, row := range buildCEMetricRows(ctx, output, d.EqualsQuals) {
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// get more pages if there are any...
		if output.NextPageToken == nil {
			break
		}
		params.NextPageToken = output.NextPageToken
	}

	return nil, nil
}

func buildCEMetricRows(ctx context.Context, costUsageData *costexplorer.GetCostAndUsageOutput, _ map[string]*proto.QualValue) []CEMetricRow {
	var rows []CEMetricRow

	for _, result := range costUsageData.ResultsByTime {

		// If there are no groupings, create a row from the totals
		if len(result.Groups) == 0 {
			var row CEMetricRow

			row.Estimated = result.Estimated
			row.PeriodStart = result.TimePeriod.Start
			row.PeriodEnd = result.TimePeriod.End
			row.setRowMetrics(result.Total)
			rows = append(rows, row)
		}
		// make a row per group
		for _, group := range result.Groups {
			var row CEMetricRow

			row.Estimated = result.Estimated
			row.PeriodStart = result.TimePeriod.Start
			row.PeriodEnd = result.TimePeriod.End

			if len(group.Keys) > 0 {
				row.Dimension1 = aws.String(group.Keys[0])
				if len(group.Keys) > 1 {
					row.Dimension2 = aws.String(group.Keys[1])
				}
			}
			row.setRowMetrics(group.Metrics)
			rows = append(rows, row)
		}
	}
	return rows
}

// CEMetricRow is the flattened, aggregated value for a metric.
type CEMetricRow struct {
	Estimated bool

	// The time period that the result covers.
	PeriodStart *string
	PeriodEnd   *string

	Dimension1 *string
	Dimension2 *string
	//Tag *string

	BlendedCostAmount      *string
	UnblendedCostAmount    *string
	NetUnblendedCostAmount *string
	AmortizedCostAmount    *string
	NetAmortizedCostAmount *string
	UsageQuantityAmount    *string
	NormalizedUsageAmount  *string

	BlendedCostUnit      *string
	UnblendedCostUnit    *string
	NetUnblendedCostUnit *string
	AmortizedCostUnit    *string
	NetAmortizedCostUnit *string
	UsageQuantityUnit    *string
	NormalizedUsageUnit  *string
}

func (row *CEMetricRow) setRowMetrics(metrics map[string]types.MetricValue) {

	row.BlendedCostAmount = metrics["BlendedCost"].Amount
	row.BlendedCostUnit = metrics["BlendedCost"].Unit

	row.UnblendedCostAmount = metrics["UnblendedCost"].Amount
	row.UnblendedCostUnit = metrics["UnblendedCost"].Unit

	row.NetUnblendedCostAmount = metrics["NetUnblendedCost"].Amount
	row.NetUnblendedCostUnit = metrics["NetUnblendedCost"].Unit

	row.AmortizedCostAmount = metrics["AmortizedCost"].Amount
	row.AmortizedCostUnit = metrics["AmortizedCost"].Unit

	row.NetAmortizedCostAmount = metrics["NetAmortizedCost"].Amount
	row.NetAmortizedCostUnit = metrics["NetAmortizedCost"].Unit

	row.UsageQuantityAmount = metrics["UsageQuantity"].Amount
	row.UsageQuantityUnit = metrics["UsageQuantity"].Unit

	row.NormalizedUsageAmount = metrics["NormalizedUsageAmount"].Amount
	row.NormalizedUsageUnit = metrics["NormalizedUsageAmount"].Unit

}

func getCEStartDateForGranularity(granularity string) time.Time {
	switch granularity {
	case "DAILY", "MONTHLY":
		// 1 year
		return time.Now().AddDate(-1, 0, 0)
	case "HOURLY":
		// 13 days
		return time.Now().AddDate(0, 0, -13)
	}
	return time.Now().AddDate(0, 0, -13)
}

type CEQuals struct {
	// Quals stuff
	SearchStartTime *timestamp.Timestamp
	SearchEndTime   *timestamp.Timestamp
	Metrics         string
	Granularity     string
	DimensionType1  string
	DimensionType2  string
	TagKey1         string
	TagKey2         string
}

func hydrateCostAndUsageQuals(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("hydrateKeyQuals", "d.EqualsQuals", d.EqualsQuals)

	return &CEQuals{
		SearchStartTime: d.EqualsQuals["search_start_time"].GetTimestampValue(),
		SearchEndTime:   d.EqualsQuals["search_end_time"].GetTimestampValue(),
		Granularity:     d.EqualsQuals["granularity"].GetStringValue(),
		DimensionType1:  d.EqualsQuals["dimension_type_1"].GetStringValue(),
		DimensionType2:  d.EqualsQuals["dimension_type_2"].GetStringValue(),
		TagKey1:         d.EqualsQuals["tag_key_1"].GetStringValue(),
		TagKey2:         d.EqualsQuals["tag_key_2"].GetStringValue(),
	}, nil
}