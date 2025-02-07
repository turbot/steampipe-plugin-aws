package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	glueEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueDataQualityRuleset(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_data_quality_ruleset",
		Description: "AWS Glue Data Quality Ruleset",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException", "UnknownOperationException"}),
			},
			Hydrate: getGlueDataQualityRuleset,
			Tags:    map[string]string{"service": "glue", "action": "GetDataQualityRuleset"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueDataQualityRulesets,
			Tags:    map[string]string{"service": "glue", "action": "ListDataQualityRulesets"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UnknownOperationException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				// We need to pass both database name and table name all together.
				// If the database name and table name are passed as input parameters, the API gives a validation error, thus they are eliminated from the optional quals.
				// api error ValidationException: 1 validation error detected: Value null at 'filter.targetTable.name' failed to satisfy constraint: Member must not be null
				{Name: "created_on", Require: plugin.Optional, Operators: []string{"<=", "<", ">=", ">"}},
				{Name: "last_modified_on", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getGlueDataQualityRuleset,
				Tags: map[string]string{"service": "glue", "action": "GetDataQualityRuleset"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(glueEndpoint.AWS_GLUE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the data quality ruleset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_name",
				Description: "The name of the database where the glue table exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TargetTable.DatabaseName"),
			},
			{
				Name:        "table_name",
				Description: "The name of the glue table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TargetTable.TableName"),
			},
			{
				Name:        "created_on",
				Description: "The date and time the data quality ruleset was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the data quality ruleset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_on",
				Description: "The date and time the data quality ruleset was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "recommendation_run_id",
				Description: "When a ruleset was created from a recommendation run, this run ID is generated to link the two together.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rule_count",
				Description: "The number of rules in the ruleset.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "rule_set",
				Description: "A Data Quality Definition Language (DQDL) ruleset.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueDataQualityRuleset,
			},
			{
				Name:        "target_table",
				Description: "An object representing a glue table.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueDataQualityRulesets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_data_quality_ruleset.listGlueDataQualityRulesets", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limit the results
	maxLimit := int32(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			maxLimit = int32(*limit)
		}
	}
	input := &glue.ListDataQualityRulesetsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filter := buildGlueDataQualityRulesetFilter(d.Quals)
	if filter != nil {
		input.Filter = filter
	}

	// List call
	paginator := glue.NewListDataQualityRulesetsPaginator(svc, input, func(o *glue.ListDataQualityRulesetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_data_quality_ruleset.listGlueDataQualityRulesets", "api_error", err)
			return nil, err
		}
		for _, ruleset := range output.Rulesets {
			d.StreamListItem(ctx, ruleset)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueDataQualityRuleset(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(types.DataQualityRulesetListDetails).Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_data_quality_ruleset.getGlueDataQualityRuleset", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &glue.GetDataQualityRulesetInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetDataQualityRuleset(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_data_quality_ruleset.getGlueDataQualityRuleset", "api_error", err)
		return nil, err
	}
	return data, nil
}

//// UTILITY FUNCTION

// Build glue data quality ruleset list call input filter
func buildGlueDataQualityRulesetFilter(quals plugin.KeyColumnQualMap) *types.DataQualityRulesetFilterCriteria {
	filter := &types.DataQualityRulesetFilterCriteria{}
	columns := map[string]string{
		"created_on":       "timestamp",
		"last_modified_on": "timestamp",
	}

	for columnName, dataType := range columns {
		if quals[columnName] != nil {
			switch dataType {
			case "timestamp":
				for _, q := range quals[columnName].Quals {
					value := q.Value.GetTimestampValue().AsTime()
					switch columnName {
					case "created_on":
						if helpers.StringSliceContains([]string{"<=", "<"}, q.Operator) {
							filter.CreatedBefore = &value
						} else if helpers.StringSliceContains([]string{">=", ">"}, q.Operator) {
							filter.CreatedAfter = &value
						}
					case "last_modified_on":
						if helpers.StringSliceContains([]string{"<=", "<"}, q.Operator) {
							filter.LastModifiedBefore = &value
						} else if helpers.StringSliceContains([]string{">=", ">"}, q.Operator) {
							filter.LastModifiedBefore = &value
						}
					}
				}

			}
		}
	}
	return filter
}
