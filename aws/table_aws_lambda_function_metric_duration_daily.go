package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLambdaFunctionMetricDurationDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_function_metric_duration_daily",
		Description: "AWS Lambda Function Cloudwatch Metrics - Duration (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listAwsLambdaFunctions,
			Hydrate:       listLambdaFunctionMetricDurationDaily,
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

func listLambdaFunctionMetricDurationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lambdaFunctionConfiguration := h.Item.(*lambda.FunctionConfiguration)
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/Lambda", "Duration", "FunctionName", *lambdaFunctionConfiguration.FunctionName)
}
