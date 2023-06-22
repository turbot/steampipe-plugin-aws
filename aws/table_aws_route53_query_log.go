package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRoute53QueryLog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_query_log",
		Description: "AWS Route53 Query Logging Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getRoute53QueryLog,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchQueryLoggingConfig"}),
			},
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "hosted_zone_id", Require: plugin.Optional},
			},
			Hydrate: listRoute53QueryLogs,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchHostedZone"}),
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID for a configuration for DNS query logging.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hosted_zone_id",
				Description: "The ID of the hosted zone that CloudWatch Logs is logging queries for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cloud_watch_logs_log_group_arn",
				Description: "The Amazon Resource Name (ARN) of the CloudWatch Logs log group that Amazon Route 53 is publishing logs to.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns

			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53QueryLogAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listRoute53QueryLogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	hostedZoneID := d.EqualsQualString("hosted_zone_id")

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_query_log.listRoute53QueryLogs", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := &route53.ListQueryLoggingConfigsInput{
		MaxResults: aws.Int32(maxItems),
	}

	if hostedZoneID != "" {
		input.HostedZoneId = &hostedZoneID
	}

	paginator := route53.NewListQueryLoggingConfigsPaginator(svc, input, func(o *route53.ListQueryLoggingConfigsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_query_log.listRoute53QueryLogs", "api_error", err)
			return nil, err
		}

		for _, config := range output.QueryLoggingConfigs {
			d.StreamListItem(ctx, config)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRoute53QueryLog(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_query_log.getRoute53QueryLog", "client_error", err)
		return nil, err
	}

	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	input := &route53.GetQueryLoggingConfigInput{
		Id: &id,
	}

	op, err := svc.GetQueryLoggingConfig(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_query_log.getRoute53QueryLog", "api_error", err)
		return nil, err
	}
	return *op.QueryLoggingConfig, nil
}

//// TRANSFORM FUNCTION

func getRoute53QueryLogAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logData := h.Item.(types.QueryLoggingConfig)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_query_log.getRoute53QueryLogAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	// arn:aws:route53:::query-log/<hosted-zone-ID>/<query-log-ID>
	arn := fmt.Sprintf("arn:%s:route53:::query-log/%s/%s", commonColumnData.Partition, *logData.HostedZoneId, *logData.Id)

	// Get data for turbot defined properties
	akas := []string{arn}

	return akas, nil
}
