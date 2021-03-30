package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsElasticBeanstalkEnvironment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elastic_beanstalk_environment",
		Description: "AWS ElasticBeanstalk Environment",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("environment_name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getAwsElasticBeanstalkEnvironment,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsElasticBeanstalkEnvironments,
		},
		GetMatrixItem: BuildRegionList,
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
				Name:        "date-created",
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
				Hydrate:     getAwsElasticBeanstalkEnvironment,
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
				Hydrate:     getAwsElasticBeanstalkEnvironment,
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
				Name:        "environment_links",
				Description: "A list of links to other environments in the same group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsElasticBeanstalkEnvironment,
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
				Transform:   transform.FromField("EnvironmentArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsElasticBeanstalkEnvironments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsElasticBeanstalkEnvironments", "AWS_REGION", region)

	// Create session
	svc, err := ElasticBeanstalkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	pagesLeft := true
	params := &elasticbeanstalk.DescribeEnvironmentsInput{}

	for pagesLeft {
		result, err := svc.DescribeEnvironments(params)
		if err != nil {
			return nil, err
		}

		for _, environments := range result.Environments {
			d.StreamListItem(ctx, environments)
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

func getAwsElasticBeanstalkEnvironment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsElasticBeanstalkEnvironment")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := ElasticBeanstalkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(*elasticbeanstalk.EnvironmentDescription).EnvironmentName
	} else {
		name = d.KeyColumnQuals["environment_name"].GetStringValue()
	}

	// Build the params

	params := &elasticbeanstalk.DescribeEnvironmentsInput{
		EnvironmentNames: []*string{aws.String(name)},
	}

	environmentData, err := svc.DescribeEnvironments(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsElasticBeanstalkEnvironment__", "ERROR", err)
		return nil, err
	}

	if environmentData != nil && len(environmentData.Environments) > 0 {
		return environmentData.Environments[0], nil
	}

	return nil, nil
}

func listElasticBeanstalkEnvironmentTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	plugin.Logger(ctx).Trace("listElasticBeanstalkEnvironmentTags")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	resourceArn := h.Item.(*elasticbeanstalk.EnvironmentDescription).EnvironmentArn

	// Create session
	svc, err := ElasticBeanstalkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build param
	params := &elasticbeanstalk.ListTagsForResourceInput{
		ResourceArn: resourceArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("listElasticBeanstalkEnvironmentTags", "ERROR", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS
func elasticBeanstalkEnvironmentTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("elasticBeanstalkEnvironmentTagListToTurbotTags")
	tags := d.HydrateItem.(*elasticbeanstalk.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	if tags.ResourceTags == nil {
		return nil, nil
	}
	var turbotTagsMap map[string]string
	if tags.ResourceTags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.ResourceTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}