package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/opensearch/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsOpenSearchClusterMetricHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_opensearch_cluster_metric_hourly",
		Description: "AWS OpenSearch Cluster Cloudwatch Metrics - Hourly",
		List: &plugin.ListConfig{
			ParentHydrate: listOpenSearchDomains,
			Hydrate:       listOpenSearchClusterMetricHourly,
			Tags:          map[string]string{"service": "cloudwatch", "action": "GetMetricStatistics"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "metric_name", Require: plugin.Required},
			},
		},
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "domain_name",
					Description: "The name of the domain.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue1"),
				},
				{
					Name:        "client_id",
					Description: "The client ID of the domain.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue2"),
				},
			})),
	}
}

func listOpenSearchClusterMetricHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	domain := h.Item.(types.DomainInfo)
	clusterMetric := d.EqualsQualString("metric_name")
	// Empty check
	if clusterMetric == "" {
		return nil, nil
	}

	account, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonData := account.(*awsCommonColumnData)
	
	return listOpenSearchCWMetricStatistics(ctx, d, "HOURLY", "AWS/ES", clusterMetric, "DomainName", *domain.DomainName, "ClientId", commonData.AccountId)
}
