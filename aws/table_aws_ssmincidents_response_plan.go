package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssmincidents"
	"github.com/aws/aws-sdk-go-v2/service/ssmincidents/types"

	ssmincidentsv1 "github.com/aws/aws-sdk-go/service/ssmincidents"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSSMIncidentsResponseaPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssmincidents_response_plan",
		Description: "AWS SSMIncidents Response Plan",
		Get: &plugin.GetConfig{
			// The resource's ARN is 'arn:aws:ssm-incidents::333333333333:response-plan/test53', which lacks the region specification.
			// Therefore, omitting the 'region' column as a requirement leads to a 'Key column is not globally unique' error, despite the Get config functioning correctly.
			KeyColumns: plugin.AllColumns([]string{"arn", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getSSMIncidentsResponsePlan,
			Tags:    map[string]string{"service": "ssm-incidents", "action": "GetResponsePlan"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSSMIncidentsResponsePlans,
			Tags:    map[string]string{"service": "ssm-incidents", "action": "ListResponsePlans"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSSMIncidentsResponsePlan,
				Tags: map[string]string{"service": "ssm-incidents", "action": "GetResponsePlan"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssmincidentsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the response plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the response plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The human readable name of the response plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "incident_template",
				Description: "Details used to create the incident when using this response plan.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSSMIncidentsResponsePlan,
			},
			{
				Name:        "actions",
				Description: "The actions that this response plan takes at the beginning of the incident.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSSMIncidentsResponsePlan,
			},
			{
				Name:        "chat_channel",
				Description: "The Chatbot chat channel used for collaboration during an incident.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSSMIncidentsResponsePlan,
			},
			{
				Name:        "engagements",
				Description: "The Amazon Resource Name (ARN) for the contacts and escalation plans that the response plan engages during an incident.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSSMIncidentsResponsePlan,
			},
			{
				Name:        "integrations",
				Description: "Information about third-party services integrated into the Incident Manager response plan.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSSMIncidentsResponsePlan,
			},

			// Steampipe standard columns
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
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSSMIncidentsResponsePlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SSMIncidentsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssmincidents_response_plan.listSSMIncidentsResponsePlans", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &ssmincidents.ListResponsePlansInput{
		MaxResults: &maxLimit,
	}

	paginator := ssmincidents.NewListResponsePlansPaginator(svc, input, func(o *ssmincidents.ListResponsePlansPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssmincidents_response_plan.listSSMIncidentsResponsePlans", "api_error", err)
			return nil, err
		}

		for _, responsePlan := range output.ResponsePlanSummaries {
			d.StreamListItem(ctx, responsePlan)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSSMIncidentsResponsePlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SSMIncidentsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssmincidents_response_plan.getSSMIncidentsResponsePlan", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var responsePlanArn string
	if h.Item != nil {
		responsePlanArn = *h.Item.(types.ResponsePlanSummary).Arn
	} else {
		responsePlanArn = d.EqualsQualString("arn")
	}

	// We don't need region check here.
	// The function 'GetMatrixItemFunc' eliminates the need for a region check by retrieving results specific to the region defined in the query parameter.

	if responsePlanArn == "" {
		return nil, nil
	}

	params := &ssmincidents.GetResponsePlanInput{
		Arn: aws.String(responsePlanArn),
	}

	op, err := svc.GetResponsePlan(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssmincidents_response_plan.getSSMIncidentsResponsePlan", "api_error", err)
	}

	if op != nil {
		return *op, nil
	}

	return nil, nil
}
