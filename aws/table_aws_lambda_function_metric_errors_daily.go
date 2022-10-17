package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLambdaFunctionMetricErrorsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_function_metric_errors_daily",
		Description: "AWS Lambda Function Cloudwatch Metrics - Errors (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsLambdaFunctions,
			Hydrate:       listLambdaFunctionMetricErrorsDaily,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "name",
					Description: "The name of the function.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listLambdaFunctionMetricErrorsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lambdaFunctionConfiguration := h.Item.(types.FunctionConfiguration)
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/Lambda", "Errors", "FunctionName", *lambdaFunctionConfiguration.FunctionName)
}
