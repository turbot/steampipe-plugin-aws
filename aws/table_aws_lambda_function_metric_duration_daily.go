package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		GetMatrixItemFunc: CloudWatchRegionsMatrix,
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
	lambdaFunctionConfiguration := h.Item.(types.FunctionConfiguration)
	return listCWMetricStatistics(ctx, d, "DAILY", "AWS/Lambda", "Duration", "FunctionName", *lambdaFunctionConfiguration.FunctionName)
}
