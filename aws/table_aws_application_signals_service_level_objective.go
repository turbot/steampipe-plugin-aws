package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationsignals"
	"github.com/aws/aws-sdk-go-v2/service/applicationsignals/types"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableAwsApplicationSignalsServiceLevelObjective(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_application_signals_service_level_objective",
		Description: "AWS CloudWatch Application Signals Service Level Objective (SLO)",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getApplicationSignalsSlo,
			Tags:       map[string]string{"service": "application-signals", "action": "GetServiceLevelObjective"},
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "include_linked_accounts",
					Require: plugin.Optional,
				},
				{
					Name:    "operation_name",
					Require: plugin.Optional,
				},
				{
					Name:    "slo_owner_aws_account_id",
					Require: plugin.Optional,
				},
				{
					Name:    "metric_source_type",
					Require: plugin.Optional,
				},

				// Fields derived from DependencyConfig
				// Ref: https://docs.aws.amazon.com/applicationsignals/latest/APIReference/API_DependencyConfig.html
				{Name: "dependency_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dependency_resource_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dependency_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dependency_identifier", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dependency_environment", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "dependency_operation_name", Require: plugin.Optional, Operators: []string{"="}},

				// Fields derived from KeyAttributes
				// Ref: https://docs.aws.amazon.com/applicationsignals/latest/APIReference/API_ListServiceLevelObjectives.html#applicationsignals-ListServiceLevelObjectives-request-KeyAttributes
				{Name: "slo_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "slo_resource_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "slo_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "slo_identifier", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "slo_environment", Require: plugin.Optional, Operators: []string{"="}},

				// Fields derived from MetricSource
				// Ref: https://docs.aws.amazon.com/applicationsignals/latest/APIReference/API_MetricSource.html
				{Name: "metric_source_key_attributes", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "metric_source_attributes", Require: plugin.Optional, Operators: []string{"="}},
			},
			Hydrate: listApplicationSignalsSlo,
			Tags:    map[string]string{"service": "application-signals", "action": "ListServiceLevelObjectives"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_APPLICATION_SIGNALS_SERVICE_ID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getApplicationSignalsSloTags,
				Tags: map[string]string{"service": "application-signals", "action": "ListTagsForResource"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The ARN of the SLO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "burn_rate_configurations",
				Description: "Array of configurations used to calculate the burn rate metrics of the SLO.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsSlo,
			},
			{
				Name:        "created_time",
				Description: "The creation timestamp of the SLO.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the SLO.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getApplicationSignalsSlo,
			},
			{
				Name:        "evaluation_type",
				Description: "Whether the SLO is a period-based SLO or a request-based SLO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "goal",
				Description: "Structure of attributes that define the goal of the SLO.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsSlo,
			},
			{
				Name:        "id",
				Description: "The ARN of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn"),
			},
			{
				Name:        "last_updated_time",
				Description: "The last update timestamp of the SLO.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getApplicationSignalsSlo,
			},
			{
				Name:        "metric_source_type",
				Description: "The metric source type of the SLO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the SLO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operation_name",
				Description: "If the SLO is specific to a single operation, this provides name of that operation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "request_based_sli",
				Description: "If this is a request-based SLO, this contains information about the performance metric.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsSlo,
			},
			{
				Name:        "sli",
				Description: "If this is a period-based SLO, this contains information about the performance metric.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsSlo,
			},
			{
				Name:        "tags",
				Description: "The tags assigned to the group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsSloTags,
				Transform:   transform.FromValue(),
			},

			// Attributes related to DependencyConfig
			{
				Name:        "dependency_config",
				Description: "The dependency config of the SLO.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dependency_type",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Type of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DependencyConfig.DependencyKeyAttributes.Type"),
			},
			{
				Name:        "dependency_resource_type",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.ResourceType of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DependencyConfig.DependencyKeyAttributes.ResourceType"),
			},
			{
				Name:        "dependency_name",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Name of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DependencyConfig.DependencyKeyAttributes.Name"),
			},
			{
				Name:        "dependency_identifier",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Identifier of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DependencyConfig.DependencyKeyAttributes.Identifier"),
			},
			{
				Name:        "dependency_environment",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Environment of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DependencyConfig.DependencyKeyAttributes.Environment"),
			},
			{
				Name:        "dependency_operation_name",
				Description: "Corresponds to DependencyConfig.DependencyOperationName of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DependencyConfig.DependencyOperationName"),
			},

			// Attributes related to KeyAttributes
			{
				Name:        "key_attributes",
				Description: "A string-to-string map of key attributes of the SLO.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "slo_type",
				Description: "Corresponds to KeyAttributes.Type of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyAttributes.Type"),
			},
			{
				Name:        "slo_resource_type",
				Description: "Corresponds to KeyAttributes.ResourceType of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyAttributes.ResourceType"),
			},
			{
				Name:        "slo_name",
				Description: "Corresponds to KeyAttributes.Name of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyAttributes.Name"),
			},
			{
				Name:        "slo_identifier",
				Description: "Corresponds to KeyAttributes.Identifier of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyAttributes.Identifier"),
			},
			{
				Name:        "slo_environment",
				Description: "Corresponds to KeyAttributes.Environment of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyAttributes.Environment"),
			},

			// Attributes related to MetricSource
			{
				Name:        "metric_source",
				Description: "The metric source of the SLO on resources other than Application Signals services.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metric_source_key_attributes",
				Description: "Corresponds to MetricSource.MetricSourceKeyAttributes of the SLO.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MetricSource.MetricSourceKeyAttributes"),
			},
			{
				Name:        "metric_source_attributes",
				Description: "Corresponds to MetricSource.MetricSourceAttributes of the SLO.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MetricSource.MetricSourceAttributes"),
			},

			// Attributes derived from qualifiers
			{
				Name:        "include_linked_accounts",
				Description: "Flag indicating whether the result set includes linked accounts.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("include_linked_accounts"),
			},
			{
				Name:        "slo_owner_aws_account_id",
				Description: "AWS account ID that owns the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("slo_owner_aws_account_id"),
			},
		}),
	}
}

//// LIST FUNCTION

func listApplicationSignalsSlo(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_slo.listApplicationSignalsSlo", "client_error", err)
		return nil, err
	}

	// Compile inputs
	input := &applicationsignals.ListServiceLevelObjectivesInput{}

	if d.EqualsQuals["include_linked_accounts"] != nil {
		input.IncludeLinkedAccounts = d.EqualsQuals["include_linked_accounts"].GetBoolValue()
	}

	if d.EqualsQuals["operation_name"] != nil {
		if d.EqualsQuals["operation_name"].GetStringValue() != "" {
			input.OperationName = aws.String(d.EqualsQuals["operation_name"].GetStringValue())
		}
	}

	if d.EqualsQuals["slo_owner_aws_account_id"] != nil {
		if d.EqualsQuals["slo_owner_aws_account_id"].GetStringValue() != "" {
			input.SloOwnerAwsAccountId = aws.String(d.EqualsQuals["slo_owner_aws_account_id"].GetStringValue())
		}
	}

	dependencyConfig := buildDependencyConfigParam(d.Quals)
	if dependencyConfig != nil {
		input.DependencyConfig = dependencyConfig
	}

	keyAttributes := buildKeyAttributesParam(d.Quals)
	if keyAttributes != nil {
		input.KeyAttributes = *keyAttributes
	}

	paginator := applicationsignals.NewListServiceLevelObjectivesPaginator(
		svc,
		input,
		func(o *applicationsignals.ListServiceLevelObjectivesPaginatorOptions) {
			o.StopOnDuplicateToken = true
		},
	)

	for paginator.HasMorePages() {
		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_application_signals_slo.listApplicationSignalsSlo", "api_error", err)
			return nil, err
		}

		for _, sloSummary := range output.SloSummaries {
			d.StreamListItem(ctx, sloSummary)

			// Context may be cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getApplicationSignalsSlo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var sloId string
	// The SLO ID may be supplied as a parameter for a Get hydration call
	// or as may have already been hydrated from the List hydration call.
	if d.EqualsQuals["id"].GetStringValue() != "" {
		sloId = d.EqualsQuals["id"].GetStringValue()
	} else {
		sloId = *h.Item.(types.ServiceLevelObjectiveSummary).Arn
	}

	if sloId == "" {
		return nil, nil
	}

	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_slo.getApplicationSignalsSlo", "client_error", err)
		return nil, err
	}

	output, err := svc.GetServiceLevelObjective(ctx, &applicationsignals.GetServiceLevelObjectiveInput{
		Id: aws.String(sloId),
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_slo.getApplicationSignalsSlo", "api_error", err)
		return nil, err
	}

	return output.Slo, nil
}

func getApplicationSignalsSloTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn, err := getSloArn(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_slo.getApplicationSignalsSloTags", "parse_error", err)
		return nil, err
	}

	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_slo.getApplicationSignalsSloTags", "client_error", err)
		return nil, err
	}

	output, err := svc.ListTagsForResource(ctx, &applicationsignals.ListTagsForResourceInput{
		ResourceArn: arn,
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_slo.getApplicationSignalsSloTags", "api_error", err)
		return nil, err
	}

	tags := make(map[string]string)
	for _, tag := range output.Tags {
		tags[*tag.Key] = *tag.Value
	}

	return tags, nil
}

// // UTILITY FUNCTIONS
func buildDependencyConfigParam(quals plugin.KeyColumnQualMap) *types.DependencyConfig {
	dependencyConfig := &types.DependencyConfig{}
	updated := false

	// Build DependencyConfig.DependencyKeyAttributes parameter
	dependencyKeyAttributeColumnNames := map[string]string{
		"dependency_type":          "Type",
		"dependency_resource_type": "ResourceType",
		"dependency_name":          "Name",
		"dependency_identifier":    "Identifier",
		"dependency_environment":   "Environment",
	}
	for qualName, paramName := range dependencyKeyAttributeColumnNames {
		if quals[qualName] == nil {
			continue
		}
		for _, q := range quals[qualName].Quals {
			value := q.Value.GetStringValue()
			if value == "" || q.Operator != "=" {
				continue
			}

			dependencyConfig.DependencyKeyAttributes[paramName] = value
			updated = true
		}
	}

	// Build DependencyConfig.DependencyOperationName
	if quals["dependency_operation_name"] != nil {
		for _, q := range quals["dependency_operation_name"].Quals {
			value := q.Value.GetStringValue()
			if value == "" || q.Operator != "=" {
				continue
			}

			dependencyConfig.DependencyOperationName = aws.String(value)
			updated = true
		}
	}

	// Return the filters only if it was updated by any of the qualifiers.
	if updated {
		return dependencyConfig
	} else {
		return nil
	}
}

func buildKeyAttributesParam(quals plugin.KeyColumnQualMap) *map[string]string {
	keyAttributes := make(map[string]string)

	// Build KeyAttributes parameter
	keyAttributeColumnNames := map[string]string{
		"slo_type":          "Type",
		"slo_resource_type": "ResourceType",
		"slo_name":          "Name",
		"slo_identifier":    "Identifier",
		"slo_environment":   "Environment",
	}
	for qualName, paramName := range keyAttributeColumnNames {
		if quals[qualName] == nil {
			continue
		}
		for _, q := range quals[qualName].Quals {
			value := q.Value.GetStringValue()
			if value == "" || q.Operator != "=" {
				continue
			}

			keyAttributes[paramName] = value
		}
	}

	// Return only if any key attributes are provided.
	if len(keyAttributes) > 0 {
		return &keyAttributes
	} else {
		return nil
	}
}

func getSloArn(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (*string, error) {
	if h.Item != nil {
		switch item := h.Item.(type) {
		case *types.ServiceLevelObjective:
			return item.Arn, nil
		case types.ServiceLevelObjectiveSummary:
			return item.Arn, nil
		}
	}
	return nil, nil
}
