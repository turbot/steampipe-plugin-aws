package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsLambdaFunction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_function",
		Description: "AWS Lambda Function",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getAwsLambdaFunction,
			Tags:       map[string]string{"service": "lambda", "action": "GetFunction"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsLambdaFunctions,
			Tags:    map[string]string{"service": "lambda", "action": "ListFunctions"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_LAMBDA_SERVICE_ID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsLambdaFunction,
				Tags: map[string]string{"service": "lambda", "action": "GetFunction"},
			},
			{
				Func: getFunctionPolicy,
				Tags: map[string]string{"service": "lambda", "action": "GetPolicy"},
			},
			{
				Func: getLambdaFunctionUrlConfig,
				Tags: map[string]string{"service": "lambda", "action": "GetFunctionUrlConfig"},
			},
		},
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
				Transform:   transform.FromField("Configuration.Architectures", "Architectures"),
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
				Transform:   transform.FromField("Configuration.FileSystemConfigs", "FileSystemConfigs"),
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
				Name:        "tracing_config",
				Description: "The function's X-Ray tracing configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Configuration.TracingConfig", "TracingConfig"),
			},
			{
				Name:        "snap_start",
				Description: "Set ApplyOn to PublishedVersions to create a snapshot of the initialized execution environment when you publish a function version.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Configuration.SnapStart", "SnapStart"),
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
			{
				Name:        "layers",
				Description: resourceInterfaceDescription("layers"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Configuration.Layers", "Layers"),
			},
		}),
	}
}

// // LIST FUNCTION

func listAwsLambdaFunctions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.listAwsLambdaFunctions", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(10000)
	input := lambda.ListFunctionsInput{}

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
	paginator := lambda.NewListFunctionsPaginator(svc, &input, func(o *lambda.ListFunctionsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lambda_function.listAwsLambdaFunctions", "api_error", err)
			return nil, err
		}

		for _, function := range output.Functions {
			d.StreamListItem(ctx, function)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

// // HYDRATE FUNCTIONS

func getAwsLambdaFunction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		name = *h.Item.(types.FunctionConfiguration).FunctionName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Empty input check
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getAwsLambdaFunction", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	}

	rowData, err := svc.GetFunction(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getAwsLambdaFunction", "api_error", err)
		return nil, err
	}

	return rowData, nil
}

func getFunctionPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	functionName := functionName(h.Item)

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getFunctionPolicy", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &lambda.GetPolicyInput{
		FunctionName: aws.String(functionName),
	}

	op, err := svc.GetPolicy(ctx, input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// If the function does not exist or does not have a policy, the operation returns a 404 (ResourceNotFoundException) error.
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_lambda_function.getFunctionPolicy", "api_error", err)
		return nil, err
	}
	return op, nil
}

func getLambdaFunctionUrlConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	functionName := functionName(h.Item)

	commonColumnData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaFunctionUrlConfig", "get_common_columns_error", err)
		return nil, err
	}

	awsCommonData := commonColumnData.(*awsCommonColumnData)
	// GovCloud does not support function URLs
	// https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/govcloud-lambda.html#govcloud-lambda-diffs
	if awsCommonData.Partition == "aws-us-gov" {
		return nil, nil
	}

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &lambda.GetFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
	}

	urlConfigs, err := svc.GetFunctionUrlConfig(ctx, input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// If the function does not exist or does not have url config, the operation returns a 404 (ResourceNotFoundException) error.
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_lambda_function.getLambdaFunctionUrlConfig", "api_error", err)
		return nil, err
	}

	return urlConfigs, nil
}

func functionName(item interface{}) string {
	switch item := item.(type) {
	case types.FunctionConfiguration:
		return *item.FunctionName
	case *lambda.GetFunctionOutput:
		return *item.Configuration.FunctionName
	}
	return ""
}
