package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	ssmincidents "github.com/aws/aws-sdk-go-v2/service/ssmincidents"
	"github.com/aws/aws-sdk-go-v2/service/ssmincidents/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSSMIncidentsResponseaPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssmincidents_response_plan",
		Description: "AWS SSMIncidents Response Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"arn", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFound", "InvalidParameter"}),
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
		GetMatrixItemFunc: SupportedRegionMatrix(ssmincidents.ServiceID),
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

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("arn").Transform(transform.EnsureStringArray),
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

	pagesLeft := true
	params := &ssmincidents.ListResponsePlansInput{
		MaxResults: aws.Int32(maxLimit),
	}

	for pagesLeft {
		result, err := svc.ListResponsePlans(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssmincidents_response_plan.listSSMIncidentsResponsePlans", "api_error", err)
			return nil, err
		}

		for _, plan := range result.ResponsePlanSummaries {
			d.StreamListItem(ctx, plan)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSSMIncidentsResponsePlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := d.EqualsQualString(matrixKeyRegion)

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

	var reportPlanArn, region string
	if h.Item != nil {
		reportPlanArn = *h.Item.(types.ResponsePlanSummary).Arn
		region = matrixRegion
	} else {
		reportPlanArn = d.EqualsQualString("arn")
		region = d.EqualsQualString("region")
	}

	if (reportPlanArn == "" || region == "") || region != matrixRegion {
		return nil, nil
	}

	params := &ssmincidents.GetResponsePlanInput{
		Arn: aws.String(reportPlanArn),
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
