package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elasticachev1 "github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElastiCacheUpdateAction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elasticache_update_action",
		Description: "AWS ElastiCache Update Action",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cache_cluster_id", Require: plugin.Optional},
				{Name: "replication_group_id", Require: plugin.Optional},
				{Name: "engine", Require: plugin.Optional},
				{Name: "service_update_status", Require: plugin.Optional},
				{Name: "update_action_status", Require: plugin.Optional},
				{Name: "service_update_name", Require: plugin.Optional},
				{Name: "service_update_release_date", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
			},
			Hydrate: listElastiCacheUpdateActions,
			Tags:    map[string]string{"service": "elasticache", "action": "DescribeUpdateActions"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elasticachev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cache_cluster_id",
				Description: "The ID of the cache cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_group_id",
				Description: "The ID of the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine",
				Description: "The ElastiCache engine (Redis or Memcached).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "estimated_update_time",
				Description: "The estimated length of time for the update to complete.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "nodes_updated",
				Description: "The progress of the service update on the replication group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_name",
				Description: "The unique ID of the service update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_recommended_apply_by_date",
				Description: "Recommended date to apply the update for compliance.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service_update_release_date",
				Description: "The date the update was first available.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service_update_severity",
				Description: "Severity level of the service update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_status",
				Description: "Current status of the service update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_update_type",
				Description: "Type of service update (security update, etc).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sla_met",
				Description: "Indicates if nodes were updated by recommended date.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_action_available_date",
				Description: "Date the update became available to the replication group.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_action_status",
				Description: "Current status of the update action.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_action_status_modified_date",
				Description: "Last modification date of the update action status.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractElastiCacheUpdateActionId),
			},
		}),
	}
}

func listElastiCacheUpdateActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := ElastiCacheClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elasticache_update_action.listElastiCacheUpdateActions", "client_error", err)
		return nil, err
	}
	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}
	input := &elasticache.DescribeUpdateActionsInput{
		MaxRecords: &maxLimit,
	}
	if v, ok := d.EqualsQuals["cache_cluster_id"]; ok {
		input.CacheClusterIds = []string{v.GetStringValue()}
	}
	if v, ok := d.EqualsQuals["replication_group_id"]; ok {
		input.ReplicationGroupIds = []string{v.GetStringValue()}
	}
	if v, ok := d.EqualsQuals["engine"]; ok {
		input.Engine = aws.String(v.GetStringValue())
	}
	if v, ok := d.EqualsQuals["service_update_status"]; ok {
		input.ServiceUpdateStatus = []types.ServiceUpdateStatus{types.ServiceUpdateStatus(v.GetStringValue())}
	}
	if v, ok := d.EqualsQuals["update_action_status"]; ok {
		input.UpdateActionStatus = []types.UpdateActionStatus{types.UpdateActionStatus(v.GetStringValue())}
	}
	if v, ok := d.EqualsQuals["service_update_name"]; ok {
		input.ServiceUpdateName = aws.String(v.GetStringValue())
	}
	if val := d.Quals["service_update_release_date"]; val != nil {
		input.ServiceUpdateTimeRange = &types.TimeRangeFilter{}
		for _, q := range val.Quals {
			queryTime := aws.Time(q.Value.GetTimestampValue().AsTime())
			switch q.Operator {
			case ">=", ">":
				input.ServiceUpdateTimeRange.StartTime = queryTime
			case "<=", "<":
				input.ServiceUpdateTimeRange.EndTime = queryTime
			case "=":
				input.ServiceUpdateTimeRange.StartTime = queryTime
				input.ServiceUpdateTimeRange.EndTime = queryTime
			}
		}
	}

	paginator := elasticache.NewDescribeUpdateActionsPaginator(client, input)
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elasticache_update_action.listElastiCacheUpdateActions", "api_error", err)
			return nil, err
		}

		for _, action := range page.UpdateActions {
			d.StreamListItem(ctx, action)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// TRANSFORM FUNCTION

func extractElastiCacheUpdateActionId(ctx context.Context, data *transform.TransformData) (interface{}, error) {
	rs := data.HydrateItem.(types.UpdateAction)
	if rs.CacheClusterId != nil {
		return *rs.CacheClusterId, nil
	}
	return *rs.ReplicationGroupId, nil
}
