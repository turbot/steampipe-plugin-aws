package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsElasticBeanstalkApplicationVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_elastic_beanstalk_application_version",
		Description: "AWS Elastic Beanstalk Application Version",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"application_name", "version_label"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getElasticBeanstalkApplicationVersion,
			Tags:    map[string]string{"service": "elasticbeanstalk", "action": "DescribeApplicationVersions"},
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "application_name",
					Require: plugin.Optional,
				},
				{
					Name:    "version_label",
					Require: plugin.Optional,
				},
			},
			Hydrate: listElasticBeanstalkApplicationVersions,
			Tags:    map[string]string{"service": "elasticbeanstalk", "action": "DescribeApplicationVersions"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listAwsElasticBeanstalkApplicationVersionTags,
				Tags: map[string]string{"service": "elasticbeanstalk", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ELASTICBEANSTALK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "application_name",
				Description: "The name of the application to which the application version belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_version_arn",
				Description: "The Amazon Resource Name (ARN) of the application version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "build_arn",
				Description: "Reference to the artifact from the AWS CodeBuild build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "date_created",
				Description: "The creation date of the application version.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "date_updated",
				Description: "The last modified date of the application version.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the application version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The processing status of the application version. Reflects the state of the application version during its creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_label",
				Description: "A unique identifier for the application version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_build_information",
				Description: "Information about the source code for the application version if the source code was retrieved from AWS CodeCommit.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_bundle",
				Description: "The storage location of the application version's source bundle in Amazon S3.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsElasticBeanstalkApplicationVersionTags,
				Transform:   transform.FromField("ResourceTags"),
			},
			// Standard columns for all tables
			{
				Name:        "title",
				Description: "A title for the resource, typically the resource name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VersionLabel"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsElasticBeanstalkApplicationVersionTags,
				Transform:   transform.FromField("ResourceTags").Transform(handleElasticBeanstalkApplicationVersionTurbotTags),
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationVersionArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listElasticBeanstalkApplicationVersions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application_version.listElasticBeanstalkApplicationVersions", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	pagesLeft := true

	// List call
	params := &elasticbeanstalk.DescribeApplicationVersionsInput{
		MaxRecords: aws.Int32(1000),
	}

	if d.EqualsQuals["application_name"] != nil {
		params.ApplicationName = aws.String(d.EqualsQuals["application_name"].GetStringValue())
	}

	if d.EqualsQuals["version_label"] != nil {
		params.VersionLabels = []string{d.EqualsQuals["version_label"].GetStringValue()}
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *params.MaxRecords {
			params.MaxRecords = aws.Int32(limit)
		}
	}

	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		op, err := svc.DescribeApplicationVersions(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_elastic_beanstalk_application_version.listElasticBeanstalkApplicationVersions", "api_error", err)
			return nil, err
		}

		for _, application_version := range op.ApplicationVersions {
			d.StreamListItem(ctx, application_version)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if op.NextToken != nil {
			pagesLeft = true
			params.NextToken = op.NextToken
		} else {
			pagesLeft = false
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getElasticBeanstalkApplicationVersion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	appname := d.EqualsQuals["application_name"].GetStringValue()
	versionlabel := d.EqualsQuals["version_label"].GetStringValue()

	// Return nil, if no input provided
	if appname == "" || versionlabel == "" {
		return nil, nil
	}

	// Create Session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application_version.getElasticBeanstalkApplicationVersion", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	namePtr := &appname
	// Build the params
	params := &elasticbeanstalk.DescribeApplicationVersionsInput{
		ApplicationName: namePtr,
		VersionLabels:   []string{versionlabel},
	}

	// Get call
	data, err := svc.DescribeApplicationVersions(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application_version.getElasticBeanstalkApplicationVersion", "api_error", err)
		return nil, err
	}

	if data != nil && len(data.ApplicationVersions) > 0 {
		return data.ApplicationVersions[0], nil
	}

	return nil, nil
}

func listAwsElasticBeanstalkApplicationVersionTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := h.Item.(types.ApplicationVersionDescription).ApplicationVersionArn

	// Create Session
	svc, err := ElasticBeanstalkClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application_version.listAwsElasticBeanstalkApplicationVersionTags", "connection_error", err)
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
		plugin.Logger(ctx).Error("aws_elastic_beanstalk_application_version.listAwsElasticBeanstalkApplicationVersionTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

func handleElasticBeanstalkApplicationVersionTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
