package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsConfigRuleComplianceDetail(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_rule_compliance_detail",
		Description: "AWS Config Rule Compliance Detail",
		List: &plugin.ListConfig{
			ParentHydrate: listConfigRules,
			Hydrate:       listConfigRuleComplianceDetails,
			Tags:          map[string]string{"service": "config", "action": "GetComplianceDetailsByConfigRule"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "config_rule_name",
					Require: plugin.Optional,
				},
				{
					Name:    "compliance_type",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CONFIG_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "config_rule_name",
				Description: "The name of the AWS Config rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of AWS resource that was evaluated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.EvaluationResultIdentifier.EvaluationResultQualifier.ResourceType"),
			},
			{
				Name:        "resource_id",
				Description: "The ID of the AWS resource that was evaluated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.EvaluationResultIdentifier.EvaluationResultQualifier.ResourceId"),
			},
			{
				Name:        "compliance_type",
				Description: "Indicates whether the AWS resource complies with the Config rule (COMPLIANT, NON_COMPLIANT, NOT_APPLICABLE, INSUFFICIENT_DATA).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.ComplianceType"),
			},
			{
				Name:        "annotation",
				Description: "Supplementary information about how the evaluation determined the compliance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.Annotation"),
			},
			{
				Name:        "config_rule_invoked_time",
				Description: "The time when AWS Config invoked the AWS Config rule to evaluate your AWS resources.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EvaluationResult.ConfigRuleInvokedTime"),
			},
			{
				Name:        "result_recorded_time",
				Description: "The time when AWS Config recorded the evaluation result.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EvaluationResult.ResultRecordedTime"),
			},
			{
				Name:        "evaluation_mode",
				Description: "The mode of an evaluation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.EvaluationResultIdentifier.EvaluationResultQualifier.EvaluationMode"),
			},
			{
				Name:        "ordering_timestamp",
				Description: "The time of the event that triggered the evaluation of your AWS resources.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EvaluationResult.EvaluationResultIdentifier.OrderingTimestamp"),
			},
			{
				Name:        "resource_evaluation_id",
				Description: "A unique ResourceEvaluationId that is associated with a single execution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.EvaluationResultIdentifier.ResourceEvaluationId"),
			},
			{
				Name:        "result_token",
				Description: "An encrypted token that associates an evaluation with an AWS Config rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.ResultToken"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EvaluationResult.EvaluationResultIdentifier.EvaluationResultQualifier.ResourceId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listConfigRuleComplianceDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get config rule details from parent hydrate
	configRule := h.Item.(types.ConfigRule)
	configRuleName := *configRule.ConfigRuleName

	// Check if there's an optional filter for config_rule_name
	if d.EqualsQualString("config_rule_name") != "" {
		if d.EqualsQualString("config_rule_name") != configRuleName {
			return nil, nil
		}
	}

	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_rule_compliance_detail.listConfigRuleComplianceDetails", "client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
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

	input := &configservice.GetComplianceDetailsByConfigRuleInput{
		ConfigRuleName: &configRuleName,
		Limit:          maxLimit,
	}

	// Optional filter by compliance type
	if d.EqualsQuals["compliance_type"] != nil {
		complianceType := d.EqualsQualString("compliance_type")
		if complianceType != "" {
			input.ComplianceTypes = []types.ComplianceType{types.ComplianceType(complianceType)}
		}
	}

	paginator := configservice.NewGetComplianceDetailsByConfigRulePaginator(svc, input, func(o *configservice.GetComplianceDetailsByConfigRulePaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_config_rule_compliance_detail.listConfigRuleComplianceDetails", "api_error", err)
			return nil, err
		}

		for _, evaluationResult := range output.EvaluationResults {
			d.StreamListItem(ctx, &ConfigRuleComplianceDetail{
				ConfigRuleName:   configRuleName,
				EvaluationResult: evaluationResult,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

// ConfigRuleComplianceDetail is a struct to hold config rule name and evaluation result
type ConfigRuleComplianceDetail struct {
	ConfigRuleName   string
	EvaluationResult types.EvaluationResult
}
