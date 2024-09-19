package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/applicationsignals"
	applicationsignalsv1 "github.com/aws/aws-sdk-go/service/applicationsignals"

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
			Tags:       map[string]string{"service": "application-signals", "action": "GetApplicationSignalsServiceLevelObjective"},
			KeyColumns: plugin.AllColumns([]string{"arn", "name"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listApplicationSignalsServiceLevelObjectives,
			Tags:    map[string]string{"service": "application-signals", "action": "ListApplicationSignalsServiceLevelObjectives"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getApplicationSignalsServiceLevelObjective,
				Tags: map[string]string{"service": "application-signals", "action": "GetApplicationSignalsServiceLevelObjective"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(applicationsignalsv1.EndpointsID),
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
			//// Steampipe Standard Columns
			//{
			//	Name:        "tags",
			//	Description: resourceInterfaceDescription("tags"),
			//	Type:        proto.ColumnType_JSON,
			//	Hydrate:     getLogGroupTagging,
			//},
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

	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

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
