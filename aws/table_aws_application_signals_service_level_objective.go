package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/applicationsignals"
	cloudwatchEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationsignals/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsApplicationSignalsServiceLevelObjective(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_application_signals_service_level_objective",
		Description: "AWS Application Signals Service Level Objective",
		Get: &plugin.GetConfig{
			Hydrate:    getApplicationSignalsServiceLevelObjective,
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "application-signals", "action": "GetApplicationSignalsServiceLevelObjective"},
		},
		List: &plugin.ListConfig{
			Hydrate:    listApplicationSignalsServiceLevelObjectives,
			KeyColumns: plugin.OptionalColumns([]string{"operation_name"}),
			Tags:       map[string]string{"service": "application-signals", "action": "ListApplicationSignalsServiceLevelObjectives"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getApplicationSignalsServiceLevelObjective,
				Tags: map[string]string{"service": "application-signals", "action": "GetApplicationSignalsServiceLevelObjective"},
			},
		},
		// AWS Doesn't treat it as a separate service it is under CloudWatch service but the API package is different.
		// https://aws.amazon.com/about-aws/whats-new/2024/06/amazon-cloudwatch-application-signals-application-monitoring/
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchEndpoint.LOGSServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the service level objective.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the service level objective.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operation_name",
				Description: " If this service level objective is specific to a single operation, this field displays the name of that operation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The date and time that this SLO was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "evaluation_type",
				Description: "Displays whether this is a period-based SLO or a request-based SLO.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getApplicationSignalsServiceLevelObjective,
			},
			{
				Name:        "attainment_goal",
				Description: "The attainment goal of the service level objective.",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getApplicationSignalsServiceLevelObjective,
				Transform:   transform.FromField("Goal.AttainmentGoal"),
			},
			{
				Name:        "goal",
				Description: "The goal of the service level objective.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsServiceLevelObjective,
			},
			{
				Name:        "sli",
				Description: "The sli of the service level objective.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsServiceLevelObjective,
			},
			{
				Name:        "request_based_sli",
				Description: "A structure containing information about the performance metric that this SLO monitors, if this is a period-based SLO.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsServiceLevelObjective,
			},

			//// Steampipe Standard Columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listApplicationSignalsServiceLevelObjectives(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := ApplicationSignalsClient(ctx, d)

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.listApplicationSignalsServiceLevelObjectives", "client_error", err)
		return nil, err
	}

	maxItems := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := &applicationsignals.ListServiceLevelObjectivesInput{
		MaxResults: &maxItems,
	}

	if d.EqualsQualString("operation_name") != "" {
		input.OperationName = aws.String(d.EqualsQualString("operation_name"))
	}

	paginator := applicationsignals.NewListServiceLevelObjectivesPaginator(svc, input, func(o *applicationsignals.ListServiceLevelObjectivesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.listApplicationSignalsServiceLevelObjectives", "api_error", err)
			return nil, err
		}

		for _, sloSummary := range output.SloSummaries {
			d.StreamListItem(ctx, sloSummary)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getApplicationSignalsServiceLevelObjective(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	arn := ""

	if h.Item != nil {
		data := h.Item.(types.ServiceLevelObjectiveSummary)
		arn = *data.Arn
	} else {
		arn = d.EqualsQualString("arn")
	}

	// check if name is empty
	if strings.TrimSpace(arn) == "" {
		return nil, nil
	}

	// Restrict API call for other regions
	if len(strings.Split(arn, ":")) > 3 && strings.Split(arn, ":")[3] != region {
		return nil, nil
	}

	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.getApplicationSignalsServiceLevelObjective", "client_error", err)
		return nil, err
	}

	params := &applicationsignals.GetServiceLevelObjectiveInput{
		Id: aws.String(arn),
	}

	item, err := svc.GetServiceLevelObjective(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.getApplicationSignalsServiceLevelObjective", "api_error", err)
		return nil, err
	}

	return item.Slo, nil
}