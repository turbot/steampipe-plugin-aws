package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsElasticBeanstalkApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elastic_beanstalk_application",
		Description: "AWS Elastic Beanstalk Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsElasticBeanstalkApplication,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsElasticBeanstalkApplications,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationArn").Transform(arnToAkas),
			},
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
				Transform:   transform.FromField("ResourceTags").Transform(getElasticBeanstalkApplicationTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsElasticBeanstalkApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ElasticBeanstalkService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	params := &elasticbeanstalk.DescribeApplicationsInput{}

	op, err := svc.DescribeApplications(params)
	if err != nil {
		return nil, err
	}

	for _, application := range op.Applications {
		d.StreamListItem(ctx, application)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsElasticBeanstalkApplication(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsElasticBeanstalkApplication")

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := ElasticBeanstalkService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &elasticbeanstalk.DescribeApplicationsInput{
		ApplicationNames: []*string{aws.String(name)},
	}

	// Get call
	data, err := svc.DescribeApplications(params)

	if data.Applications != nil && len(data.Applications) > 0 {
		return data.Applications[0], nil
	}

	return nil, err
}

func listAwsElasticBeanstalkApplicationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsElasticBeanstalkApplicationTags")

	arn := *h.Item.(*elasticbeanstalk.ApplicationDescription).ApplicationArn

	// Create Session
	svc, err := ElasticBeanstalkService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &elasticbeanstalk.ListTagsForResourceInput{
		ResourceArn: &arn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getElasticBeanstalkApplicationTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getElasticBeanstalkApplicationTurbotTags")
	tagList := d.HydrateItem.(*elasticbeanstalk.ListTagsForResourceOutput)

	if tagList.ResourceTags == nil {
		return nil, nil
	}
	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.ResourceTags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
