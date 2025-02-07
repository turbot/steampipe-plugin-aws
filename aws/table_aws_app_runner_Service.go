package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/apprunner/types"

	apprunnerEndpointId "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAppRunnerService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_app_runner_service",
		Description: "AWS App Runner Service",
		Get: &plugin.GetConfig{
			Hydrate:    getAwsAppRunnerService,
			KeyColumns: plugin.SingleColumn("arn"),
			// We need to handle the InvalidParameterException, as the provided ARN triggers an InvalidRequestException in regions where the resource is unavailable.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidRequestException"}),
			},
			Tags: map[string]string{"service": "apprunner", "action": "DescribeService"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsAppRunnerServices,
			Tags:    map[string]string{"service": "apprunner", "action": "ListServices"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(apprunnerEndpointId.AWS_APPRUNNER_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Description: "The customer-provided service name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_id",
				Description: "An ID that App Runner generated for this service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of this service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceArn"),
			},
			{
				Name:        "created_at",
				Description: "The time when the App Runner service was created. It's in the Unix time stamp format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service_url",
				Description: "A subdomain URL that App Runner generated for this service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The time when the App Runner service was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "deleted_at",
				Description: "The time when the App Runner service was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsAppRunnerService,
				Transform:   transform.FromField("DeletedAt").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "kms_key",
				Description: "The ARN of the KMS key that's used for encryption.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsAppRunnerService,
				Transform:   transform.FromField("EncryptionConfiguration.KmsKey"),
			},
			{
				Name:        "policy_type",
				Description: "The policy type. Currently supported values are TargetTrackingScaling and StepScaling",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_scaling_configuration_summary",
				Description: "Summary information for the App Runner automatic scaling configuration resource that's associated with this service.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAppRunnerService,
			},
			{
				Name:        "instance_configuration",
				Description: "The runtime configuration of instances (scaling units) of this service.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAppRunnerService,
			},
			{
				Name:        "network_configuration",
				Description: "Configuration settings related to network traffic of the web application that this service runs.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAppRunnerService,
			},
			{
				Name:        "source_configuration",
				Description: "The source deployed to the App Runner service. It can be a code or an image repository.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAppRunnerService,
			},
			{
				Name:        "health_check_configuration",
				Description: "The settings for the health check that App Runner performs to monitor the health of this service.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAppRunnerService,
			},
			{
				Name:        "observability_configuration",
				Description: "The observability configuration of this service.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsAppRunnerService,
			},

			// Standard standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsAppRunnerServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := AppRunnerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_app_runner_service.listAwsAppRunnerServices", "client_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil // Unsupported region check
	}

	// Limit the result
	input := &apprunner.ListServicesInput{
		MaxResults: aws.Int32(20),
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxResults {
			if limit < 1 {
				input.MaxResults = aws.Int32(1)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	paginator := apprunner.NewListServicesPaginator(svc, input, func(o *apprunner.ListServicesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_app_runner_service.listAwsAppRunnerServices", "api_error", err)
			return nil, err
		}

		for _, service := range output.ServiceSummaryList {
			d.StreamListItem(ctx, service)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsAppRunnerService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := AppRunnerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_app_runner_service.getAwsAppRunnerService", "client_error", err)
		return nil, err
	}

	if svc == nil {
		return nil, nil // Unsupported region check
	}

	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.ServiceSummary).ServiceArn
	} else {
		arn = d.EqualsQuals["arn"].GetStringValue()
	}

	params := &apprunner.DescribeServiceInput{
		ServiceArn: aws.String(arn),
	}

	service, err := svc.DescribeService(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_app_runner_service.getAwsAppRunnerService", "api_error", err)
		return nil, err
	}

	return service.Service, nil
}
