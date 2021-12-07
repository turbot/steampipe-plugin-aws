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

func tableAwsLambdaVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_version",
		Description: "AWS Lambda Version",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"version", "function_name"}),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameter", "ResourceNotFoundException"}),
			Hydrate:           getFunctionVersion,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsLambdaFunctions,
			Hydrate:       listLambdaVersions,
		},
		GetMatrixItem: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listLambdaVersions")

	// Create service
	svc, err := LambdaService(ctx, d)
	if err != nil {
		return nil, err
	}

	function := h.Item.(*lambda.FunctionConfiguration)

	err = svc.ListVersionsByFunctionPages(
		&lambda.ListVersionsByFunctionInput{FunctionName: function.FunctionName},
		func(page *lambda.ListVersionsByFunctionOutput, lastPage bool) bool {
			for _, version := range page.Versions {
				d.StreamLeafListItem(ctx, version)
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

// do not have a get call
// using list api call to create get function
func getFunctionVersion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunctionVersion")

	// Create service
	svc, err := LambdaService(ctx, d)
	if err != nil {
		return nil, err
	}

	version := d.KeyColumnQuals["version"].GetStringValue()
	functionName := d.KeyColumnQuals["function_name"].GetStringValue()
	var functionVersion *lambda.FunctionConfiguration

	err = svc.ListVersionsByFunctionPages(
		&lambda.ListVersionsByFunctionInput{FunctionName: aws.String(functionName)},
		func(page *lambda.ListVersionsByFunctionOutput, lastPage bool) bool {
			for _, i := range page.Versions {
				if *i.Version == version {
					functionVersion = i
					return false
				}
			}
			return !lastPage
		},
	)

	if err != nil {
		return nil, err
	}

	if functionVersion != nil {
		return functionVersion, nil
	}

	return nil, nil
}

func getFunctionVersionPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunctionVersionPolicy")

	alias := h.Item.(*lambda.FunctionConfiguration)

	// Create Session
	svc, err := LambdaService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getFunctionVersionPolicy", "error_LambdaService", err)
		return nil, err
	}

	input := &lambda.GetPolicyInput{
		FunctionName: aws.String(*alias.FunctionName),
		Qualifier: aws.String(*alias.Version),
	}

	op, err := svc.GetPolicy(input)
	if err != nil {
		plugin.Logger(ctx).Error("getFunctionVersionPolicy", "error_GetPolicy", err)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return lambda.GetPolicyOutput{}, nil
			}
		}
		return nil, err
	}

	return op, nil
}
