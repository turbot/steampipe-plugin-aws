package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLambdaFunctionMetricInvocationsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_function_metric_invocations_daily",
		Description: "AWS Lambda Function Cloudwatch Metrics - Invocations (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsLambdaFunctions,
			Hydrate:       listLambdaFunctionMetricInvocationsDaily,
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

func listLambdaFunctionMetricInvocationsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lambdaFunctionConfiguration := h.Item.(types.FunctionConfiguration)
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/Lambda", "Invocations", "FunctionName", *lambdaFunctionConfiguration.FunctionName)
}
