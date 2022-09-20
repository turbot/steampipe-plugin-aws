package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

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
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.FunctionName", "FunctionName"),
			},
			{
				Name:        "arn",
				Description: "The function's Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.FunctionArn", "FunctionArn"),
			},
			{
				Name:        "code_sha_256",
				Description: "The SHA256 hash of the function's deployment package.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.CodeSha256", "CodeSha256"),
			},
			{
				Name:        "code_size",
				Description: "The size of the function's deployment package, in bytes.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Configuration.CodeSize", "CodeSize"),
			},
			{
				Name:        "dead_letter_config_target_arn",
				Description: "The Amazon Resource Name (ARN) of an Amazon SQS queue or Amazon SNS topic.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.DeadLetterConfig.TargetArn", "DeadLetterConfig.TargetArn"),
			},
			{
				Name:        "description",
				Description: "The function's description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.Description", "Description"),
			},
			{
				Name:        "handler",
				Description: "The function that Lambda calls to begin executing your function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.Handler", "Handler"),
			},
			{
				Name:        "kms_key_arn",
				Description: "The KMS key that's used to encrypt the function's environment variables. This key is only returned if you've configured a customer managed CMK.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.KMSKeyArn", "KMSKeyArn"),
			},
			{
				Name:        "last_modified",
				Description: "The date and time that the function was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Configuration.LastModified", "LastModified"),
			},
			{
				Name:        "timeout",
				Description: "The amount of time in seconds that Lambda allows a function to run before stopping it.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.Timeout", "Timeout"),
			},
			{
				Name:        "version",
				Description: "The version of the Lambda function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.Version", "Version"),
			},
			{
				Name:        "package_type",
				Description: "The type of deployment package.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.PackageType", "PackageType"),
			},
			{
				Name:        "master_arn",
				Description: "For Lambda@Edge functions, the ARN of the master function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.MasterArn", "MasterArn"),
			},
			{
				Name:        "memory_size",
				Description: "The memory that's allocated to the function.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Configuration.MemorySize", "MemorySize"),
			},
			{
				Name:        "revision_id",
				Description: "The latest updated revision of the function or alias.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.RevisionId", "RevisionId"),
			},
			{
				Name:        "role",
				Description: "The function's execution role.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.Role", "Role"),
			},
			{
				Name:        "runtime",
				Description: "The runtime environment for the Lambda function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.Runtime", "Runtime"),
			},
			{
				Name:        "state",
				Description: "The current state of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.State", "State"),
			},
			{
				Name:        "state_reason",
				Description: "The reason for the function's current state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.StateReason", "StateReason"),
			},
			{
				Name:        "state_reason_code",
				Description: "The reason code for the function's current state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.StateReasonCode", "StateReasonCode"),
			},
			{
				Name:        "last_update_status",
				Description: "The status of the last update that was performed on the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.LastUpdateStatus", "LastUpdateStatus"),
			},
			{
				Name:        "last_update_status_reason",
				Description: "The reason for the last update that was performed on the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.LastUpdateStatusReason", "LastUpdateStatusReason"),
			},
			{
				Name:        "last_update_status_reason_code",
				Description: "The reason code for the last update that was performed on the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.LastUpdateStatusReasonCode", "LastUpdateStatusReasonCode"),
			},
			{
				Name:        "reserved_concurrent_executions",
				Description: "The number of concurrent executions that are reserved for this function.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsLambdaFunction,
				Transform:   transform.FromField("Concurrency.ReservedConcurrentExecutions"),
			},
			{
				Name:        "vpc_id",
				Description: "The VPC ID that is attached to Lambda function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.VpcConfig.VpcId", "VpcConfig.VpcId"),
			},
			{
				Name:        "architectures",
				Description: "The instruction set architecture that the function supports. Architecture is a string array with one of the valid values.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "code",
				Description: "The deployment package of the function or version.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsLambdaFunction,
			},
			{
				Name:        "environment_variables",
				Description: "The environment variables that are accessible from function code during execution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Configuration.Environment.Variables", "Environment.Variables"),
			},
			{
				Name:        "file_system_configs",
				Description: "Connection settings for an Amazon EFS file system.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsLambdaFunction,
				Transform:   transform.FromField("Configuration.FileSystemConfigs"),
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
				Name:        "url_config",
				Description: "The function URL configuration details of the function.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getLambdaFunctionUrlConfig,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "vpc_security_group_ids",
				Description: "A list of VPC security groups IDs attached to Lambda function.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Configuration.VpcConfig.SecurityGroupIds", "VpcConfig.SecurityGroupIds"),
			},
			{
				Name:        "vpc_subnet_ids",
				Description: "A list of VPC subnet IDs attached to Lambda function.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Configuration.VpcConfig.SubnetIds", "VpcConfig.SubnetIds"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.FunctionName", "FunctionName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsLambdaFunction,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Configuration.FunctionArn", "FunctionArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsLambdaFunctions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsLambdaFunctions")

	// Create service
	svc, err := LambdaService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &lambda.ListFunctionsInput{
		MaxItems: aws.Int64(10000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = aws.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	err = svc.ListFunctionsPages(
		input,
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			for _, function := range page.Functions {
				d.StreamListItem(ctx, function)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsLambdaFunction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsLambdaFunction")

	var name string
	if h.Item != nil {
		name = *h.Item.(*lambda.FunctionConfiguration).FunctionName
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := LambdaService(ctx, d)
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

	return rowData, nil
}

func getFunctionPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunctionPolicy")

	functionName := functionName(h.Item)

	// Create Session
	svc, err := LambdaService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &lambda.GetPolicyInput{
		FunctionName: aws.String(functionName),
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

func getLambdaFunctionUrlConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLambdaFunctionUrlConfig")

	functionName := functionName(h.Item)

	commonColumnData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getLambdaFunctionUrlConfig", "get_common_columns_error", err)
		return nil, err
	}

	awsCommonData := commonColumnData.(*awsCommonColumnData)
	// GovCloud does not support function URLs
	// https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/govcloud-lambda.html#govcloud-lambda-diffs
	if awsCommonData.Partition == "aws-us-gov" {
		return nil, nil
	}

	// Create Session
	svc, err := LambdaService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &lambda.GetFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
	}

	urlConfigs, err := svc.GetFunctionUrlConfig(input)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("getLambdaFunctionUrlConfig", "GetFunctionUrlConfig_error", err)
		return nil, err
	}

	return urlConfigs, nil
}

func functionName(item interface{}) string {
	switch item := item.(type) {
	case *lambda.FunctionConfiguration:
		return *item.FunctionName
	case *lambda.GetFunctionOutput:
		return *item.Configuration.FunctionName
	}
	return ""
}
