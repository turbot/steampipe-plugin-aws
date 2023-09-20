package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"

	elasticbeanstalkv1 "github.com/aws/aws-sdk-go/service/elasticbeanstalk"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsElasticBeanstalkApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elastic_beanstalk_application",
		Description: "AWS Elastic Beanstalk Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getElasticBeanstalkApplication,
			Tags: map[string]string{"service": "elasticache", "action": "DescribeApplications"},
		},
		List: &plugin.ListConfig{
			Hydrate: listElasticBeanstalkApplications,
			Tags:		map[string]string{"service": "elasticache", "action": "DescribeApplications"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listAwsElasticBeanstalkApplicationTags,
				Tags: map[string]string{"service": "elasticache", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(elasticbeanstalkv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationArn"),
			},
			{
				Name:        "description",
				Description: "User-defined description of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "date_created",
				Description: "The date when the application was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "date_updated",
				Description: "The date when the application was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "configuration_templates",
				Description: "The names of the configuration templates associated with this application.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConfigurationTemplates"),
			},
			{
				Name:        "versions",
				Description: "The names of the versions for this application.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_lifecycle_config",
				Description: "The lifecycle settings for the application.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsElasticBeanstalkApplicationTags,
				Transform:   transform.FromField("ResourceTags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsElasticBeanstalkApplicationTags,
				Transform:   transform.FromField("ResourceTags").Transform(handleElasticBeanstalkApplicationTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listElasticBeanstalkApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application.listElasticBeanstalkApplications", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// List call
	params := &elasticbeanstalk.DescribeApplicationsInput{}

	// DescribeApplications doesn't support pagination
	op, err := svc.DescribeApplications(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application.listElasticBeanstalkApplications", "api_error", err)
		return nil, err
	}

	for _, application := range op.Applications {
		d.StreamListItem(ctx, application)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getElasticBeanstalkApplication(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application.getElasticBeanstalkApplication", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &elasticbeanstalk.DescribeApplicationsInput{
		ApplicationNames: []string{name},
	}

	// Get call
	data, err := svc.DescribeApplications(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application.getElasticBeanstalkApplication", "api_error", err)
		return nil, err
	}

	if len(data.Applications) > 0 {
		return data.Applications[0], nil
	}

	return nil, nil
}

func listAwsElasticBeanstalkApplicationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := h.Item.(types.ApplicationDescription).ApplicationArn

	// Create Session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application.listAwsElasticBeanstalkApplicationTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &elasticbeanstalk.ListTagsForResourceInput{
		ResourceArn: arn,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application.listAwsElasticBeanstalkApplicationTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

func handleElasticBeanstalkApplicationTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(*elasticbeanstalk.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil && len(tagList.ResourceTags) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.ResourceTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
