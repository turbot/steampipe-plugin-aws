package aws

import (
	"context"
	"errors"
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

func tableAwsLambdaVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_version",
		Description: "AWS Lambda Version",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"version", "function_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameter", "ResourceNotFoundException"}),
			},
			Hydrate: getFunctionVersion,
			Tags:    map[string]string{"service": "lambda", "action": "ListVersionsByFunction"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsLambdaFunctions,
			Hydrate:       listLambdaVersions,
			Tags:          map[string]string{"service": "lambda", "action": "ListVersionsByFunction"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "function_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lambdav1.EndpointsID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getFunctionVersionPolicy,
				Tags: map[string]string{"service": "lambda", "action": "GetPolicy"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "version",
				Description: "The version of the Lambda function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "function_name",
				Description: "The name of the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The function's Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionArn"),
			},
			{
				Name:        "master_arn",
				Description: "For Lambda@Edge functions, the ARN of the master function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the function.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "last_modified",
				Description: "The date and time that the function was last updated, in ISO-8601 format.",
				Type:        proto.ColumnType_TIMESTAMP,
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
				Name:        "runtime",
				Description: "The runtime environment for the Lambda function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timeout",
				Description: "The amount of time in seconds that Lambda allows a function to run before stopping it.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpcConfig.VpcId"),
			},
			{
				Name:        "kms_key_arn",
				Description: "The KMS key that's used to encrypt the function's environment variables.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KMSKeyArn"),
			},
			{
				Name:        "role",
				Description: "The function's execution role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "signing_job_arn",
				Description: "The ARN of the signing job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "signing_profile_version_arn",
				Description: "The ARN of the signing profile version.",
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
				Name:        "ephemeral_storage_size",
				Description: "The size of the function's /tmp directory in MB.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("EphemeralStorage.Size"),
			},
			{
				Name:        "environment_variables",
				Description: "The environment variables that are accessible from function code during execution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Environment.Variables"),
			},
			{
				Name:        "policy",
				Description: "Contains the resource-based policy.",
				Hydrate:     getFunctionVersionPolicy,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy_std",
				Description: "Contains the contents of the resource-based policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunctionVersionPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
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
				Name:        "architectures",
				Description: "The instruction set architecture that the function supports.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dead_letter_config",
				Description: "The function's dead letter queue configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "file_system_configs",
				Description: "Connection settings for an Amazon EFS file system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image_config_response",
				Description: "The function's image configuration values.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "layers",
				Description: "The function's layers.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "logging_config",
				Description: "The function's Amazon CloudWatch Logs configuration settings.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "runtime_version_config",
				Description: "The ARN of the runtime and any errors that occurred.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "snap_start",
				Description: "Configuration for creating a snapshot of the initialized execution environment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tracing_config",
				Description: "The function's X-Ray tracing configuration.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Version"),
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

func listLambdaVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	function := h.Item.(types.FunctionConfiguration)

	// Create service
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_version.listLambdaVersions", "connection_err", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	equalQuals := d.EqualsQuals

	// Minimize the API call with the given function name
	if equalQuals["function_name"] != nil {
		if equalQuals["function_name"].GetStringValue() != *function.FunctionName {
			return nil, nil
		}
	}

	maxItems := int32(50)

	input := &lambda.ListVersionsByFunctionInput{
		FunctionName: function.FunctionName,
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
	paginator := lambda.NewListVersionsByFunctionPaginator(svc, input, func(o *lambda.ListVersionsByFunctionPaginatorOptions) {
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

		for _, version := range output.Versions {
			d.StreamListItem(ctx, version)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

// do not have a get call
// using list api call to create get function
func getFunctionVersion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var functionVersion types.FunctionConfiguration

	version := d.EqualsQuals["version"].GetStringValue()
	functionName := d.EqualsQuals["function_name"].GetStringValue()
	if strings.TrimSpace(version) == "" || strings.TrimSpace(functionName) == "" {
		return nil, nil
	}

	// Create service
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_version.getFunctionVersion", "connection_err", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(functionName),
	}

	paginator := lambda.NewListVersionsByFunctionPaginator(svc, input, func(o *lambda.ListVersionsByFunctionPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lambda_function.getFunctionVersion", "api_error", err)
			return nil, err
		}

		for _, i := range output.Versions {
			if *i.Version == version {
				functionVersion = i
				return functionVersion, nil
			}
		}
	}

	return nil, nil
}

func getFunctionVersionPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	alias := h.Item.(types.FunctionConfiguration)

	// Create Session
	svc, err := LambdaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lambda_function.getFunctionVersionPolicy", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &lambda.GetPolicyInput{
		FunctionName: aws.String(*alias.FunctionName),
		Qualifier:    aws.String(*alias.Version),
	}

	op, err := svc.GetPolicy(ctx, input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return lambda.GetPolicyOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("aws_lambda_function.getFunctionVersionPolicy", "connection_error", err)
		return nil, err
	}

	return op, nil
}
