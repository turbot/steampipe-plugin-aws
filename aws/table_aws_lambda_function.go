package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func tableAwsLambdaFunction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_function",
		Description: "AWS Lambda Function",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getAwsLambdaFunction,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsLambdaFunctions,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionName"),
			},
			{
				Name:        "code_sha_256",
				Description: "The SHA256 hash of the function's deployment package.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "code_size",
				Description: "The size of the function's deployment package, in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "The function's description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "handler",
				Description: "The function that Lambda calls to begin executing your function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_arn",
				Description: "The KMS key that's used to encrypt the function's environment variables. This key is only returned if you've configured a customer managed CMK.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified",
				Description: "The date and time that the function was last updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timeout",
				Description: "The amount of time in seconds that Lambda allows a function to run before stopping it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The version of the Lambda function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_arn",
				Description: "For Lambda@Edge functions, the ARN of the master function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "memory_size",
				Description: "The memory that's allocated to the function.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "revision_id",
				Description: "The latest updated revision of the function or alias.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The function's execution role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "runtime",
				Description: "The runtime environment for the Lambda function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_reason",
				Description: "The reason for the function's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_reason_code",
				Description: "The reason code for the function's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_update_status",
				Description: "The status of the last update that was performed on the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_update_status_reason",
				Description: "The reason for the last update that was performed on the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_update_status_reason_code",
				Description: "The reason code for the last update that was performed on the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The VPC ID that is attached to Lambda function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpcConfig.VpcId"),
			},
			{
				Name:        "vpc_security_group_ids",
				Description: "A list of VPC security groups IDs attached to Lambda function.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VpcConfig.SecurityGroupIds"),
			},
			{
				Name:        "vpc_subnet_ids",
				Description: "A list of VPC subnet IDs attached to Lambda function.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VpcConfig.SubnetIds"),
			},
			{
				Name:        "policy",
				Description: "The resource-based iam policy of Lambda function.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunctionPolicy,
				Transform:   transform.FromField("Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunctionPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "arn",
				Description: "The function's Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionArn"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunctionTagging,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FunctionArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsLambdaFunctions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsLambdaFunctions", "AWS_REGION", region)

	// Create service
	svc, err := LambdaService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	err = svc.ListFunctionsPages(
		&lambda.ListFunctionsInput{},
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			for _, function := range page.Functions {
				d.StreamListItem(ctx, function)
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsLambdaFunction(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsLambdaFunction")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := LambdaService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	}

	rowData, err := svc.GetFunction(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsLambdaFunction__", "ERROR", err)
		return nil, err
	}

	return rowData.Configuration, nil
}

func getFunctionPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunctionPolicy")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	function := h.Item.(*lambda.FunctionConfiguration)

	// Create Session
	svc, err := LambdaService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &lambda.GetPolicyInput{
		FunctionName: function.FunctionName,
	}

	op, err := svc.GetPolicy(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return lambda.GetPolicyOutput{}, nil
			}
		}
		return nil, err
	}
	return op, nil
}

func getFunctionTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunctionTagging")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	function := h.Item.(*lambda.FunctionConfiguration)

	// Create Session
	svc, err := LambdaService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &lambda.GetFunctionInput{
		FunctionName: function.FunctionName,
	}

	op, err := svc.GetFunction(input)
	if err != nil {
		return nil, err
	}
	return op, nil
}
