package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_account",
		Description: "AWS API Gateway Account",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getAPIGatewayAccount,
			Tags:    map[string]string{"service": "apigateway", "action": "GetAccount"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_APIGATEWAY_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cloudwatch_role_arn",
				Description: "The ARN of an Amazon CloudWatch role for the current Account",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "api_key_version",
				Description: "The version of the API keys used for the account",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "throttle_burst_limit",
				Description: "The API target request burst rate limit. This allows more requests through for a period of time than the target rate limit",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ThrottleSettings.BurstLimit"),
			},
			{
				Name:        "throttle_rate_limit",
				Description: "The API target request rate limit",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("ThrottleSettings.RateLimit"),
			},
			{
				Name:        "features",
				Description: "A list of features supported for the account. When usage plans are enabled, the features list will include an entry of 'UsagePlans'",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func getAPIGatewayAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_account.getAPIGatewayAccount", "service_client_error", err)
		return nil, err
	}

	params := &apigateway.GetAccountInput{}

	detail, err := svc.GetAccount(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_account.getAPIGatewayAccount", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, detail)

	return nil, nil
}
