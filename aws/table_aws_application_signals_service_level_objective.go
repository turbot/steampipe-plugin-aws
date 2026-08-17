package aws

import (
	"context"
	"encoding/json"

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
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Tags: map[string]string{"service": "application-signals", "action": "GetServiceLevelObjective"},
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
				{Name: "dependency_type", Require: plugin.Optional},
				{Name: "dependency_resource_type", Require: plugin.Optional},
				{Name: "dependency_name", Require: plugin.Optional},
				{Name: "dependency_identifier", Require: plugin.Optional},
				{Name: "dependency_environment", Require: plugin.Optional},
				{Name: "dependency_operation_name", Require: plugin.Optional},

				// Fields derived from KeyAttributes
				// Ref: https://docs.aws.amazon.com/applicationsignals/latest/APIReference/API_ListServiceLevelObjectives.html#applicationsignals-ListServiceLevelObjectives-request-KeyAttributes
				{Name: "slo_type", Require: plugin.Optional},
				{Name: "slo_resource_type", Require: plugin.Optional},
				{Name: "slo_name", Require: plugin.Optional},
				{Name: "slo_identifier", Require: plugin.Optional},
				{Name: "slo_environment", Require: plugin.Optional},

				// Fields derived from MetricSource
				// Ref: https://docs.aws.amazon.com/applicationsignals/latest/APIReference/API_MetricSource.html
				{Name: "metric_source_key_attributes", Require: plugin.Optional},
				{Name: "metric_source_attributes", Require: plugin.Optional},
			},
			Hydrate: listApplicationSignalsSlo,
			Tags:    map[string]string{"service": "application-signals", "action": "ListServiceLevelObjectives"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_APPLICATION_SIGNALS_SERVICE_ID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getApplicationSignalsSlo,
				IgnoreConfig: &plugin.IgnoreConfig{
					ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
				},
				Tags: map[string]string{"service": "application-signals", "action": "GetServiceLevelObjective"},
			},
			{
				Func: getApplicationSignalsSloTags,
				IgnoreConfig: &plugin.IgnoreConfig{
					ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
				},
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
				Transform: transform.FromField(
					"OperationName",
					"Sli.SliMetric.OperationName",
					"RequestBasedSli.RequestBasedSliMetric.OperationName",
				),
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

			// Attributes related to DependencyConfig
			{
				Name:        "dependency_config",
				Description: "The dependency config of the SLO.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField(
					"DependencyConfig",
					"Sli.SliMetric.DependencyConfig",
					"RequestBasedSli.RequestBasedSliMetric.DependencyConfig",
				),
			},
			{
				Name:        "dependency_type",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Type of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"DependencyConfig.DependencyKeyAttributes.Type",
					"Sli.SliMetric.DependencyConfig.DependencyKeyAttributes.Type",
					"RequestBasedSli.RequestBasedSliMetric.DependencyConfig.DependencyKeyAttributes.Type",
				),
			},
			{
				Name:        "dependency_resource_type",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.ResourceType of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"DependencyConfig.DependencyKeyAttributes.ResourceType",
					"Sli.SliMetric.DependencyConfig.DependencyKeyAttributes.ResourceType",
					"RequestBasedSli.RequestBasedSliMetric.DependencyConfig.DependencyKeyAttributes.ResourceType",
				),
			},
			{
				Name:        "dependency_name",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Name of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"DependencyConfig.DependencyKeyAttributes.Name",
					"Sli.SliMetric.DependencyConfig.DependencyKeyAttributes.Name",
					"RequestBasedSli.RequestBasedSliMetric.DependencyConfig.DependencyKeyAttributes.Name",
				),
			},
			{
				Name:        "dependency_identifier",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Identifier of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"DependencyConfig.DependencyKeyAttributes.Identifier",
					"Sli.SliMetric.DependencyConfig.DependencyKeyAttributes.Identifier",
					"RequestBasedSli.RequestBasedSliMetric.DependencyConfig.DependencyKeyAttributes.Identifier",
				),
			},
			{
				Name:        "dependency_environment",
				Description: "Corresponds to DependencyConfig.DependencyKeyAttributes.Environment of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"DependencyConfig.DependencyKeyAttributes.Environment",
					"Sli.SliMetric.DependencyConfig.DependencyKeyAttributes.Environment",
					"RequestBasedSli.RequestBasedSliMetric.DependencyConfig.DependencyKeyAttributes.Environment",
				),
			},
			{
				Name:        "dependency_operation_name",
				Description: "Corresponds to DependencyConfig.DependencyOperationName of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"DependencyConfig.DependencyOperationName",
					"Sli.SliMetric.DependencyConfig.DependencyOperationName",
					"RequestBasedSli.RequestBasedSliMetric.DependencyConfig.DependencyOperationName",
				),
			},

			// Attributes related to KeyAttributes
			{
				Name:        "key_attributes",
				Description: "A string-to-string map of key attributes of the SLO.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField(
					"KeyAttributes",
					"Sli.SliMetric.KeyAttributes",
					"RequestBasedSli.RequestBasedSliMetric.KeyAttributes",
				),
			},
			{
				Name:        "slo_type",
				Description: "Corresponds to KeyAttributes.Type of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"KeyAttributes.Type",
					"Sli.SliMetric.KeyAttributes.Type",
					"RequestBasedSli.RequestBasedSliMetric.KeyAttributes.Type",
				),
			},
			{
				Name:        "slo_resource_type",
				Description: "Corresponds to KeyAttributes.ResourceType of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"KeyAttributes.ResourceType",
					"Sli.SliMetric.KeyAttributes.ResourceType",
					"RequestBasedSli.RequestBasedSliMetric.KeyAttributes.ResourceType",
				),
			},
			{
				Name:        "slo_name",
				Description: "Corresponds to KeyAttributes.Name of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"KeyAttributes.Name",
					"Sli.SliMetric.KeyAttributes.Name",
					"RequestBasedSli.RequestBasedSliMetric.KeyAttributes.Name",
				),
			},
			{
				Name:        "slo_identifier",
				Description: "Corresponds to KeyAttributes.Identifier of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"KeyAttributes.Identifier",
					"Sli.SliMetric.KeyAttributes.Identifier",
					"RequestBasedSli.RequestBasedSliMetric.KeyAttributes.Identifier",
				),
			},
			{
				Name:        "slo_environment",
				Description: "Corresponds to KeyAttributes.Environment of the SLO.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField(
					"KeyAttributes.Environment",
					"Sli.SliMetric.KeyAttributes.Environment",
					"RequestBasedSli.RequestBasedSliMetric.KeyAttributes.Environment",
				),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsSloTags,
				Transform:   transform.FromValue(),
			},

			// Attributes related to MetricSource
			{
				Name:        "metric_source",
				Description: "The metric source of the SLO on resources other than Application Signals services.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField(
					"MetricSource",
					"Sli.SliMetric.MetricSource",
					"RequestBasedSli.RequestBasedSliMetric.MetricSource",
				),
			},
			{
				Name:        "metric_source_key_attributes",
				Description: "Corresponds to MetricSource.MetricSourceKeyAttributes of the SLO.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField(
					"MetricSource.MetricSourceKeyAttributes",
					"Sli.SliMetric.MetricSource.MetricSourceKeyAttributes",
					"RequestBasedSli.RequestBasedSliMetric.MetricSource.MetricSourceKeyAttributes",
				),
			},
			{
				Name:        "metric_source_attributes",
				Description: "Corresponds to MetricSource.MetricSourceAttributes of the SLO.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField(
					"MetricSource.MetricSourceAttributes",
					"Sli.SliMetric.MetricSource.MetricSourceAttributes",
					"RequestBasedSli.RequestBasedSliMetric.MetricSource.MetricSourceAttributes",
				),
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
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getApplicationSignalsSloTags,
				Transform:   transform.From(getApplicationSignalsSloTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listApplicationSignalsSlo(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.listApplicationSignalsSlo", "client_error", err)
		return nil, err
	}

	// Compile inputs
	input := &applicationsignals.ListServiceLevelObjectivesInput{}

	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}
	input.MaxResults = aws.Int32(maxLimit)

	if d.EqualsQuals["include_linked_accounts"] != nil {
		input.IncludeLinkedAccounts = d.EqualsQuals["include_linked_accounts"].GetBoolValue()
	}

	if d.EqualsQualString("operation_name") != "" {
		input.OperationName = aws.String(d.EqualsQualString("operation_name"))
	}

	if d.EqualsQualString("slo_owner_aws_account_id") != "" {
		input.SloOwnerAwsAccountId = aws.String(d.EqualsQualString("slo_owner_aws_account_id"))
	}

	if d.EqualsQualString("metric_source_type") != "" {
		input.MetricSourceTypes = []types.MetricSourceType{types.MetricSourceType(d.EqualsQualString("metric_source_type"))}
	}

	dependencyConfig := buildDependencyConfigParam(d.Quals)
	if dependencyConfig != nil {
		input.DependencyConfig = dependencyConfig
	}

	keyAttributes := buildKeyAttributesParam(d.Quals)
	if keyAttributes != nil {
		input.KeyAttributes = *keyAttributes
	}

	metricSourceAttributesParam, err := buildParamsFromJson(d.Quals, "metric_source_attributes")
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.listApplicationSignalsSlo", "parse_error", err)
		return nil, err
	}

	metricSourceKeyAttributesParam, err := buildParamsFromJson(d.Quals, "metric_source_key_attributes")
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.listApplicationSignalsSlo", "parse_error", err)
		return nil, err
	}

	if len(metricSourceKeyAttributesParam) > 0 || len(metricSourceAttributesParam) > 0 {
		input.MetricSource = &types.MetricSource{}
		if len(metricSourceKeyAttributesParam) > 0 {
			input.MetricSource.MetricSourceKeyAttributes = metricSourceKeyAttributesParam
		}
		if len(metricSourceAttributesParam) > 0 {
			input.MetricSource.MetricSourceAttributes = metricSourceAttributesParam
		}
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
			plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.listApplicationSignalsSlo", "api_error", err)
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
	// or may have already been hydrated from the List hydration call.
	if d.EqualsQualString("id") != "" {
		sloId = d.EqualsQualString("id")
	} else {
		sloId = *h.Item.(types.ServiceLevelObjectiveSummary).Arn
	}

	if sloId == "" {
		return nil, nil
	}

	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.getApplicationSignalsSlo", "client_error", err)
		return nil, err
	}

	output, err := svc.GetServiceLevelObjective(ctx, &applicationsignals.GetServiceLevelObjectiveInput{
		Id: aws.String(sloId),
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.getApplicationSignalsSlo", "api_error", err)
		return nil, err
	}

	return output.Slo, nil
}

func getApplicationSignalsSloTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn, err := getSloArn(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.getApplicationSignalsSloTags", "parse_error", err)
		return nil, err
	}

	// Get client
	svc, err := ApplicationSignalsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.getApplicationSignalsSloTags", "client_error", err)
		return nil, err
	}

	output, err := svc.ListTagsForResource(ctx, &applicationsignals.ListTagsForResourceInput{
		ResourceArn: arn,
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_application_signals_service_level_objective.getApplicationSignalsSloTags", "api_error", err)
		return nil, err
	}

	if len(output.Tags) > 0 {
		return output.Tags, nil
	}

	return nil, nil
}

func getApplicationSignalsSloTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)

	if tags != nil {
		turbotTags := map[string]string{}
		for _, tag := range tags {
			turbotTags[*tag.Key] = *tag.Value
		}
		return turbotTags, nil
	}

	return nil, nil
}

//// UTILITY FUNCTIONS

func buildDependencyConfigParam(quals plugin.KeyColumnQualMap) *types.DependencyConfig {
	dependencyConfig := &types.DependencyConfig{
		DependencyKeyAttributes: make(map[string]string),
	}
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

func buildParamsFromJson(quals plugin.KeyColumnQualMap, qualName string) (map[string]string, error) {
	var parsed map[string]string

	if quals[qualName] != nil {
		for _, q := range quals[qualName].Quals {
			value := q.Value.GetJsonbValue()
			if value == "" || q.Operator != "=" {
				continue
			}

			err := json.Unmarshal([]byte(value), &parsed)
			if err != nil {
				return nil, err
			}
		}
	}

	return parsed, nil
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
