package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueDataQualityRuleset(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_data_quality_ruleset",
		Description: "AWS Glue Data Quality Ruleset",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException", "ValidationException", "UnknownOperationException"}),
			},
			Hydrate: getGlueDataQualityRuleset,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueDataQualityRulesets,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UnknownOperationException", "ValidationException"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the data quality ruleset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_name",
				Description: "The name of the database where the Glue table exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TargetTable.DatabaseName"),
			},
			{
				Name:        "table_name",
				Description: "The name of the Glue table.",
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
				Description: "An object representing an Glue table.",
				Type:        proto.ColumnType_INT,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointName"),
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

	// List call

	paginator := glue.NewListDataQualityRulesetsPaginator(svc, input, func(o *glue.ListDataQualityRulesetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_data_quality_ruleset.listGlueDataQualityRulesets", "api_error", err)
			return nil, err
		}
		for _, ruleset := range output.Rulesets {
			d.StreamListItem(ctx, ruleset)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
		name = d.KeyColumnQuals["name"].GetStringValue()
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
