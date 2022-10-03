package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueDevEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_dev_endpoint",
		Description: "AWS Glue Dev Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("endpoint_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueDevEndpoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueDevEndpoints,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "endpoint_name",
				Description: "The name of the DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the DevEndpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueDevEndpointArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "status",
				Description: "The current status of this DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_zone",
				Description: "The AWS Availability Zone where this DevEndpoint is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_timestamp",
				Description: "The point in time at which this DevEndpoint was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "extra_jars_s3_path",
				Description: "The path to one or more Java .jar files in an S3 bucket that should be loaded in your DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "extra_python_libs_s3_path",
				Description: "The paths to one or more Python libraries in an Amazon S3 bucket that should be loaded in your DevEndpoint. Multiple values must be complete paths separated by a comma.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "failure_reason",
				Description: "The reason for a current failure in this DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "glue_version",
				Description: "Glue version determines the versions of Apache Spark and Python that Glue supports.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_timestamp",
				Description: "The point in time at which this DevEndpoint was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_update_status",
				Description: "The status of the last update.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_nodes",
				Description: "The number of Glue Data Processing Units (DPUs) allocated to this DevEndpoint.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "number_of_workers",
				Description: "The number of workers of a defined workerType that are allocated to the development endpoint.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "private_address",
				Description: "A private IP address to access the DevEndpoint within a VPC if the DevEndpoint is created within one.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_address",
				Description: "The public IP address used by this DevEndpoint. The PublicAddress field is present only when you create a non-virtual private cloud (VPC) DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_key",
				Description: "The public key to be used by this DevEndpoint for authentication.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role used in this DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_configuration",
				Description: "The name of the SecurityConfiguration structure to be used with this DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_id",
				Description: "The subnet ID for this DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the virtual private cloud (VPC) used by this DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "worker_type",
				Description: "The type of predefined worker that is allocated to the development endpoint. Accepts a value of Standard, G.1X, or G.2X.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "yarn_endpoint_address",
				Description: "The YARN endpoint address used by this DevEndpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zeppelin_remote_spark_interpreter_port",
				Description: "The Apache Zeppelin port for the remote Apache Spark interpreter.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "public_keys",
				Description: "A list of public keys to be used by the DevEndpoints for authentication.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_group_ids",
				Description: "A list of security group identifiers used in this DevEndpoint.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueDevEndpointArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueDevEndpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_dev_endpoint.listGlueDevEndpoints", "service_creation_error", err)
		return nil, err
	}

	input := &glue.GetDevEndpointsInput{
		MaxResults: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.GetDevEndpointsPages(
		input,
		func(page *glue.GetDevEndpointsOutput, isLast bool) bool {
			for _, endpoint := range page.DevEndpoints {
				d.StreamListItem(ctx, endpoint)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_dev_endpoint.listGlueDevEndpoints", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueDevEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["endpoint_name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_dev_endpoint.getGlueDevEndpoint", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &glue.GetDevEndpointInput{
		EndpointName: aws.String(name),
	}

	// Get call
	data, err := svc.GetDevEndpoint(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_dev_endpoint.getGlueDevEndpoint", "api_error", err)
		return nil, err
	}
	return data.DevEndpoint, nil
}

func getGlueDevEndpointArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*glue.DevEndpoint)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn format - https://docs.aws.amazon.com/glue/latest/dg/glue-specifying-resource-arns.html
	// arn:aws:glue:region:account-id:devEndpoint/development-endpoint-name
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":devEndpoint/" + *data.EndpointName

	return arn, nil
}
