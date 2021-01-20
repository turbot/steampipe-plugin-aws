package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/lambda"
)

func tableAwsLambdaAlias(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_alias",
		Description: "AWS Lambda Alias",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.AllColumns([]string{"name", "function_name"}),
			ItemFromKey: aliasFromKey,
			Hydrate:     getLambdaAlias,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsLambdaFunctions,
			Hydrate:       listLambdaAliases,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the alias",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.Name"),
			},
			{
				Name:        "function_name",
				Description: "The name of the function",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alias_arn",
				Description: "The Amazon Resource Name (ARN) of the alias.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.AliasArn"),
			},
			{
				Name:        "function_version",
				Description: "The function version that the alias invokes",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.FunctionVersion"),
			},
			{
				Name:        "revision_id",
				Description: "A unique identifier that changes when you update the alias",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.RevisionId"),
			},
			{
				Name:        "description",
				Description: "A description of the alias",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.Description"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Alias.AliasArn").Transform(arnToAkas),
			},
		}),
	}
}

type aliasRowData = struct {
	Alias        *lambda.AliasConfiguration
	FunctionName *string
}

//// ITEM FROM KEY

func aliasFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	functionName := quals["function_name"].GetStringValue()
	item := &aliasRowData{
		FunctionName: &functionName,
		Alias: &lambda.AliasConfiguration{
			Name: &name,
		},
	}
	return item, nil
}

//// LIST FUNCTION

func listLambdaAliases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listLambdaAliases", "AWS_REGION", defaultRegion)

	svc, err := LambdaService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	function := h.Item.(*lambda.FunctionConfiguration)

	err = svc.ListAliasesPages(
		&lambda.ListAliasesInput{FunctionName: function.FunctionName},
		func(page *lambda.ListAliasesOutput, lastPage bool) bool {
			for _, alias := range page.Aliases {
				d.StreamLeafListItem(ctx, &aliasRowData{alias, function.FunctionName})
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getLambdaAlias(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getLambdaAlias")
	defaultRegion := GetDefaultRegion()
	alias := h.Item.(*aliasRowData)

	// Create Session
	svc, err := LambdaService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &lambda.GetAliasInput{
		FunctionName: alias.FunctionName,
		Name:         alias.Alias.Name,
	}

	rowData, err := svc.GetAlias(params)
	if err != nil {
		logger.Debug("getLambdaAlias__", "ERROR", err)
		return nil, err
	}

	return &aliasRowData{rowData, alias.FunctionName}, nil
}
