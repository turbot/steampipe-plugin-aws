package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	lambdav1 "github.com/aws/aws-sdk-go/service/lambda"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsLambdaAlias(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_alias",
		Description: "AWS Lambda Alias",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "function_name", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameter", "ResourceNotFoundException"}),
			},
			Hydrate: getLambdaAlias,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsLambdaFunctions,
			Hydrate:       listLambdaAliases,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "function_version", Require: plugin.Optional},
				{Name: "function_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lambdav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the alias.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.Name"),
			},
			{
				Name:        "function_name",
				Description: "The name of the function.",
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
				Description: "The function version that the alias invokes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.FunctionVersion"),
			},
			{
				Name:        "revision_id",
				Description: "A unique identifier that changes when you update the alias.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.RevisionId"),
			},
			{
				Name:        "description",
				Description: "A description of the alias.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Alias.Description"),
			},
			{
				Name:        "policy",
				Description: "Contains the resource-based policy.",
				Hydrate:     getLambdaAliasPolicy,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy_std",
				Description: "Contains the contents of the resource-based policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLambdaAliasPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "url_config",
				Description: "The function URL configuration details of the alias.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLambdaAliasUrlConfig,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
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
	Alias        any
	FunctionName *string
}

//// LIST FUNCTION

func listLambdaAliases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_alias.listLambdaAliases", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	function := h.Item.(types.FunctionConfiguration)

	equalQuals := d.EqualsQuals
	// Minimize the API call with the given function name
	if equalQuals["function_name"] != nil {
		if equalQuals["function_name"].GetStringValue() != "" {
			if equalQuals["function_name"].GetStringValue() != "" && equalQuals["function_name"].GetStringValue() != *function.FunctionName {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["function_name"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["function_name"].GetListValue())), *function.FunctionName) {
				return nil, nil
			}
		}
	}

	maxItems := int32(10000)
	input := lambda.ListAliasesInput{
		FunctionName: function.FunctionName,
	}

	if equalQuals["function_version"] != nil {
		input.FunctionVersion = aws.String(equalQuals["function_version"].GetStringValue())
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxItems = aws.Int32(maxItems)
	paginator := lambda.NewListAliasesPaginator(svc, &input, func(o *lambda.ListAliasesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lambda_function.listAwsLambdaFunctions", "api_error", err)
			return nil, err
		}

		for _, alias := range output.Aliases {
			d.StreamListItem(ctx, &aliasRowData{alias, function.FunctionName})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLambdaAlias(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	functionName := d.EqualsQuals["function_name"].GetStringValue()

	// Empty check
	if name == "" || functionName == "" {
		return nil, nil
	}

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaAlias", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build params
	params := &lambda.GetAliasInput{
		FunctionName: aws.String(functionName),
		Name:         aws.String(name),
	}

	rowData, err := svc.GetAlias(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaAlias", "api_error", err)
		return nil, err
	}

	return &aliasRowData{rowData, aws.String(functionName)}, nil
}

func getLambdaAliasPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	alias := h.Item.(*aliasRowData)

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaAliasPolicy", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	qualifier := getAliasQualifier(h.Item)

	input := &lambda.GetPolicyInput{
		FunctionName: aws.String(*alias.FunctionName),
		Qualifier:    aws.String(qualifier),
	}

	op, err := svc.GetPolicy(ctx, input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// If the function alias does not exist or does not have resource policy, the operation returns a 404 (ResourceNotFoundException) error.
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaAliasPolicy", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getLambdaAliasUrlConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	alias := h.Item.(*aliasRowData)

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaAliasUrlConfig", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	qualifier := getAliasQualifier(h.Item)
	input := &lambda.GetFunctionUrlConfigInput{
		FunctionName: alias.FunctionName,
		Qualifier:    aws.String(qualifier),
	}

	urlConfigs, err := svc.GetFunctionUrlConfig(ctx, input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// If the function alias does not exist or does not have resource policy, the operation returns a 404 (ResourceNotFoundException) error.
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaAliasUrlConfig", "api_error", err)
		return nil, err
	}

	return urlConfigs, nil
}

func getAliasQualifier(aliasData any) string {
	alias := aliasData.(*aliasRowData)
	var qualifier string
	switch item := (alias.Alias).(type) {
	case types.AliasConfiguration:
		qualifier = *item.Name
	case *lambda.GetAliasOutput:
		qualifier = *item.Name
	}
	return qualifier
}
