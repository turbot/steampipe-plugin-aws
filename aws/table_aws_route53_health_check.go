package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsRoute53HealthCheck(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_health_check",
		Description: "AWS Route53 Health Check",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getHealthCheck,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchHealthCheck"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listHealthChecks,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The identifier that Amazon Route 53 assigned to the health check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "caller_reference",
				Description: "A unique string that you specified when you created the health check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_version",
				Description: "The version of the health check.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "linked_service_principal",
				Description: "If the health check was created by another service, the service that created the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LinkedService.ServicePrincipal"),
			},
			{
				Name:        "linked_service_description",
				Description: "If the health check was created by another service, an configurationtional description that can be provided by the other service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LinkedService.Description"),
			},
			{
				Name:        "cloud_watch_alarm_configuration",
				Description: "A complex type that contains information about the CloudWatch alarm that Amazon Route 53 is monitoring for this health check.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(route53CloudWatchAlarmConfigurationEmptyCheck),
			},
			{
				Name:        "health_check_config",
				Description: "A complex type that contains detailed information about one health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "health_check_status",
				Description: "A list that contains one HealthCheckObservation element for each Amazon Route 53 health checker that is reporting a status about the health check endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHealthCheckStatus,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHealthCheckTags,
				Transform:   transform.FromField("ResourceTagSet.Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getHealthCheckTags,
				Transform:   transform.FromField("ResourceTagSet.Tags").Transform(route53HealthCheckTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRoute53HealthCheckTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

// A complex type that contains information about the CloudWatch alarm that Amazon
// Route 53 is monitoring for this health check.
type CloudWatchAlarmConfiguration struct {

	// For the metric that the CloudWatch alarm is associated with, the arithmetic
	// configurationeration that is used for the comparison.
	//
	// This member is required.
	ComparisonOperator types.ComparisonOperator

	// For the metric that the CloudWatch alarm is associated with, the number of
	// periods that the metric is compared to the threshold.
	//
	// This member is required.
	EvaluationPeriods *int32

	// The name of the CloudWatch metric that the alarm is associated with.
	//
	// This member is required.
	MetricName *string

	// The namespace of the metric that the alarm is associated with. For more
	// information, see Amazon CloudWatch Namespaces, Dimensions, and Metrics Reference
	// (https://docs.aws.amazon.com/AmazonCloudWatch/latest/DevelconfigurationerGuide/CW_Support_For_AWS.html)
	// in the Amazon CloudWatch User Guide.
	//
	// This member is required.
	Namespace *string

	// For the metric that the CloudWatch alarm is associated with, the duration of one
	// evaluation period in seconds.
	//
	// This member is required.
	Period *int32

	// For the metric that the CloudWatch alarm is associated with, the statistic that
	// is applied to the metric.
	//
	// This member is required.
	Statistic types.Statistic

	// For the metric that the CloudWatch alarm is associated with, the value the
	// metric is compared with.
	//
	// This member is required.
	Threshold *float64

	// For the metric that the CloudWatch alarm is associated with, a complex type that
	// contains information about the dimensions for the metric. For information, see
	// Amazon CloudWatch Namespaces, Dimensions, and Metrics Reference
	// (https://docs.aws.amazon.com/AmazonCloudWatch/latest/DevelconfigurationerGuide/CW_Support_For_AWS.html)
	// in the Amazon CloudWatch User Guide.
	Dimensions []types.Dimension
}

//// LIST FUNCTION

func listHealthChecks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_health_check.listHealthChecks", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)
	input := route53.ListHealthChecksInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxItems = aws.Int32(maxItems)
	paginator := route53.NewListHealthChecksPaginator(svc, &input, func(o *route53.ListHealthChecksPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_health_check.listHealthChecks", "api_error", err)
			return nil, err
		}

		for _, healthCheck := range output.HealthChecks {
			d.StreamListItem(ctx, healthCheck)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getHealthCheck(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_health_check.getHealthCheck", "client_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// Validate user input
	if len(id) < 1 {
		return nil, nil
	}

	params := &route53.GetHealthCheckInput{
		HealthCheckId: aws.String(id),
	}

	// Get call
	item, err := svc.GetHealthCheck(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_health_check.getHealthCheck", "api_error", err)
		return nil, err
	}

	return item.HealthCheck, nil
}

func getHealthCheckStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	healthCheck := h.Item.(types.HealthCheck)

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_health_check.getHealthCheckStatus", "client_error", err)
		return nil, err
	}

	params := &route53.GetHealthCheckStatusInput{
		HealthCheckId: aws.String(*healthCheck.Id),
	}

	// execute get call
	item, err := svc.GetHealthCheckStatus(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_health_check.getHealthCheckStatus", "api_error", err)
		return nil, err
	}

	return item.HealthCheckObservations, nil
}

func getHealthCheckTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	healthCheck := h.Item.(types.HealthCheck)

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_health_check.getHealthCheckTags", "client_error", err)
		return nil, err
	}

	params := &route53.ListTagsForResourceInput{
		ResourceId:   healthCheck.Id,
		ResourceType: "healthcheck",
	}

	// execute list call
	resp, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_health_check.getHealthCheckTags", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getRoute53HealthCheckTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	healthCheck := h.Item.(types.HealthCheck)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined prconfigurationerties
	akas := []string{"arn:" + commonColumnData.Partition + ":route53:::" + "healthcheck/" + *healthCheck.Id}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func route53HealthCheckTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(tags) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func route53CloudWatchAlarmConfigurationEmptyCheck(_ context.Context, d *transform.TransformData) (interface{}, error) {
	configuration := d.HydrateItem.(types.HealthCheck).CloudWatchAlarmConfiguration
	if len(configuration.Dimensions) == 0 {
		return CloudWatchAlarmConfiguration{configuration.ComparisonOperator, configuration.EvaluationPeriods, configuration.MetricName, configuration.Namespace, configuration.Period, configuration.Statistic, configuration.Threshold, nil}, nil
	}
	return CloudWatchAlarmConfiguration{configuration.ComparisonOperator, configuration.EvaluationPeriods, configuration.MetricName, configuration.Namespace, configuration.Period, configuration.Statistic, configuration.Threshold, configuration.Dimensions}, nil
}
