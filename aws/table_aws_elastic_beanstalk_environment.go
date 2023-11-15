package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"

	elasticbeanstalkv1 "github.com/aws/aws-sdk-go/service/elasticbeanstalk"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElasticBeanstalkEnvironment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elastic_beanstalk_environment",
		Description: "AWS ElasticBeanstalk Environment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("environment_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getElasticBeanstalkEnvironment,
			Tags:    map[string]string{"service": "elasticbeanstalk", "action": "DescribeEnvironments"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElasticBeanstalkEnvironments,
			Tags:    map[string]string{"service": "elasticbeanstalk", "action": "DescribeEnvironments"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "environment_id", Require: plugin.Optional},
				{Name: "application_name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getElasticBeanstalkEnvironment,
				Tags: map[string]string{"service": "elasticache", "action": "DescribeEnvironments"},
			},
			{
				Func: listElasticBeanstalkEnvironmentTags,
				Tags: map[string]string{"service": "elasticache", "action": "ListTagsForResource"},
			},
			{
				Func: getAwsElasticBeanstalkEnvironmentManagedActions,
				Tags: map[string]string{"service": "elasticache", "action": "DescribeEnvironmentManagedActions"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elasticbeanstalkv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "environment_name",
				Description: "The name of this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "environment_id",
				Description: "The ID of this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The environment's Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Describes this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "date_created",
				Description: "The creation date for this environment.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "abortable_operation_in_progress",
				Description: "Indicates if there is an in-progress environment configuration update or application version deployment that you can cancel.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "application_name",
				Description: "The name of the application associated with this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cname",
				Description: "The URL to the CNAME for this environment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getElasticBeanstalkEnvironment,
			},
			{
				Name:        "date_updated",
				Description: "The last modified date for this environment.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "endpoint_url",
				Description: "The URL to the LoadBalancer.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getElasticBeanstalkEnvironment,
			},
			{
				Name:        "health",
				Description: "The health status of the environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_status",
				Description: "Returns the health status of the application running in your environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operations_role",
				Description: "The Amazon Resource Name (ARN) of the environment's operations role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_arn",
				Description: "The ARN of the platform version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "solution_stack_name",
				Description: "The name of the SolutionStack deployed with this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current operational status of the environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "template_name",
				Description: "The name of the configuration template used to originally launch this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_label",
				Description: "The application version deployed in this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "configuration_settings",
				Description: "Returns a description of the settings for the specified configuration set, that is, either a configuration template or the configuration set associated with a running environment.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticBeanstalkConfigurationSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "environment_links",
				Description: "A list of links to other environments in the same group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getElasticBeanstalkEnvironment,
				Transform:   transform.FromField("EnvironmentLinks"),
			},
			{
				Name:        "managed_actions",
				Description: "A list of upcoming and in-progress managed actions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticBeanstalkEnvironmentManagedActions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "resources",
				Description: "The description of the AWS resources used by this environment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tier",
				Description: "Describes the current tier of this environment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Repository",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listElasticBeanstalkEnvironmentTags,
				Transform:   transform.FromField("ResourceTags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EnvironmentName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("ResourceTags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listElasticBeanstalkEnvironmentTags,
				Transform:   transform.FromField("ResourceTags").Transform(elasticBeanstalkEnvironmentTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EnvironmentArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listElasticBeanstalkEnvironments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_environment.listElasticBeanstalkEnvironments", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	pagesLeft := true
	params := &elasticbeanstalk.DescribeEnvironmentsInput{
		MaxRecords: aws.Int32(1000),
	}

	equalQuals := d.EqualsQuals
	if equalQuals["application_name"] != nil {
		params.ApplicationName = aws.String(equalQuals["application_name"].GetStringValue())
	}
	if equalQuals["environment_id"] != nil {
		params.EnvironmentIds = []string{equalQuals["environment_id"].GetStringValue()}
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *params.MaxRecords {
			if limit < 1 {
				params.MaxRecords = aws.Int32(1)
			} else {
				params.MaxRecords = aws.Int32(limit)
			}
		}
	}

	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		result, err := svc.DescribeEnvironments(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elastic_beanstalk_environment.listElasticBeanstalkEnvironments", "api_error", err)
			return nil, err
		}

		for _, environment := range result.Environments {
			d.StreamListItem(ctx, environment)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getElasticBeanstalkEnvironment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_environment.getElasticBeanstalkEnvironment", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(types.EnvironmentDescription).EnvironmentName
	} else {
		name = d.EqualsQuals["environment_name"].GetStringValue()
	}

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Build the params

	params := &elasticbeanstalk.DescribeEnvironmentsInput{
		EnvironmentNames: []string{name},
	}

	environmentData, err := svc.DescribeEnvironments(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_environment.getElasticBeanstalkEnvironment", "api_error", err)
		return nil, err
	}

	if len(environmentData.Environments) > 0 {
		return environmentData.Environments[0], nil
	}

	return nil, nil
}

func getAwsElasticBeanstalkEnvironmentManagedActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("elastic_beanstalk_environment.getAwsElasticBeanstalkEnvironmentManagedActions", "connection_error", err)
		return nil, err
	}

	env := h.Item.(types.EnvironmentDescription)
	// Build params
	params := &elasticbeanstalk.DescribeEnvironmentManagedActionsInput{
		EnvironmentName: env.EnvironmentName,
	}

	managedActions, err := svc.DescribeEnvironmentManagedActions(ctx, params)
	if err != nil {
		// The API throws InvalidParameterValue exception in the case if resource is not available.
		// Error: operation error Elastic Beanstalk: DescribeEnvironmentManagedActions, https response error StatusCode: 400, RequestID: b7503072-3694-4370-8a79-7182e5b1170a, api error InvalidParameterValue: No Environment found for EnvironmentName = 'Test32-envtwe'.
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "InvalidParameterValue" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("elastic_beanstalk_environment.getAwsElasticBeanstalkEnvironmentManagedActions", "api_error", err)
		return nil, err
	}

	if managedActions != nil && len(managedActions.ManagedActions) > 0 {
		return managedActions.ManagedActions, nil
	}
	return nil, nil
}

func getAwsElasticBeanstalkConfigurationSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("elastic_beanstalk_environment.getAwsElasticBeanstalkConfigurationSettings", "connection_error", err)
		return nil, err
	}

	env := h.Item.(types.EnvironmentDescription)
	// Build params
	params := &elasticbeanstalk.DescribeConfigurationSettingsInput{
		ApplicationName: env.ApplicationName,
		EnvironmentName: env.EnvironmentName,
	}

	configurationSettings, err := svc.DescribeConfigurationSettings(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "InvalidParameterValue" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("elastic_beanstalk_environment.getAwsElasticBeanstalkConfigurationSettings", "api_error", err)
		return nil, err
	}

	if configurationSettings != nil && len(configurationSettings.ConfigurationSettings) > 0 {
		return configurationSettings.ConfigurationSettings, nil
	}
	return nil, nil
}

func listElasticBeanstalkEnvironmentTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	resourceArn := h.Item.(types.EnvironmentDescription).EnvironmentArn

	// Create session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_environment.listElasticBeanstalkEnvironmentTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build param
	params := &elasticbeanstalk.ListTagsForResourceInput{
		ResourceArn: resourceArn,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_environment.listElasticBeanstalkEnvironmentTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

// // TRANSFORM FUNCTIONS
func elasticBeanstalkEnvironmentTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*elasticbeanstalk.ListTagsForResourceOutput)

	var turbotTagsMap map[string]string
	// Mapping the resource tags inside turbotTags
	if len(tags.ResourceTags) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.ResourceTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
