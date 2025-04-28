package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mwaa"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMwaaEnvironment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_mwaa_environment",
		Description: "AWS Managed Workflow for Apache Airflow (MWAA) Environment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getEnvironment,
		},
		List: &plugin.ListConfig{
			Hydrate: listMwaaEnvironments,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_AIRFLOW_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the MWAA environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the MWAA environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the MWAA environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "airflow_version",
				Description: "The version of Apache Airflow used in the environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the environment was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_update",
				Description: "The time when the environment was last updated",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdate.LastUpdateTime"),
			},
			{
				Name:        "execution_role_arn",
				Description: "The ARN of the execution role used by the environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_role_arn",
				Description: "The ARN of the service role used by the environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key",
				Description: "The ARN of the KMS key used to encrypt the environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_bucket_arn",
				Description: "The ARN of the S3 bucket used for storing DAGs and requirements",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dag_s3_path",
				Description: "The path to the DAGs folder in the S3 bucket",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "requirements_s3_path",
				Description: "The path to the requirements.txt file in the S3 bucket",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "environment_class",
				Description: "The environment class of the MWAA environment",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_workers",
				Description: "The maximum number of workers that can run in the environment",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_workers",
				Description: "The minimum number of workers that can run in the environment",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "schedulers",
				Description: "The number of schedulers in the environment",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_webservers",
				Description: "The minimum number of web servers in the environment",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_webservers",
				Description: "The maximum number of web servers in the environment",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "network_configuration",
				Description: "The network configuration of the environment",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "logging_configuration",
				Description: "The logging configuration of the environment",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "webserver_url",
				Description: "The URL of the Airflow web server",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "webserver_access_mode",
				Description: "The access mode of the web server (PUBLIC_ONLY or PRIVATE_ONLY)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "weekly_maintenance_window_start",
				Description: "The start time of the weekly maintenance window",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_management",
				Description: "The endpoint management type (CUSTOMER or SERVICE)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags",
				Description: "A map of tags for the resource",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

func listMwaaEnvironments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := MWAAClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_mwaa_environment.listMwaaEnvironments", "connection_error", err)
		return nil, err
	}

	// List all environments
	paginator := mwaa.NewListEnvironmentsPaginator(svc, &mwaa.ListEnvironmentsInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_mwaa_environment.listMwaaEnvironments", "api_error", err)
			return nil, err
		}

		for _, environmentName := range output.Environments {
			// Get environment details
			params := &mwaa.GetEnvironmentInput{
				Name: &environmentName,
			}
			environment, err := svc.GetEnvironment(ctx, params)
			if err != nil {
				plugin.Logger(ctx).Error("aws_mwaa_environment.listMwaaEnvironments", "api_error", err)
				return nil, err
			}
			d.StreamListItem(ctx, environment.Environment)
		}
	}

	return nil, nil
}

func getEnvironment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := MWAAClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_mwaa_environment.getEnvironment", "client_error", err)
		return nil, err
	}

	// Get environment name from query data
	name := d.EqualsQuals["name"].GetStringValue()
	if name == "" {
		return nil, nil
	}

	// Get environment details
	params := &mwaa.GetEnvironmentInput{
		Name: &name,
	}

	// Get call
	op, err := svc.GetEnvironment(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_mwaa_environment.getEnvironment", "api_error", err)
		return nil, err
	}

	return op.Environment, nil
}
