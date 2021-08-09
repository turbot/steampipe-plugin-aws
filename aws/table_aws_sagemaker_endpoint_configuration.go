package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerEndpointConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_endpoint_configuration",
		Description: "AWS Sagemaker Endpoint Configuration",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException", "NotFoundException"}),
			Hydrate:           getSagemakerEndpointConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listSagemakerEndpointConfigurations,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the endpoint configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointConfigName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the endpoint configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointConfigArn"),
			},
			{
				Name:        "creation_time",
				Description: "A timestamp that shows when the endpoint configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kms_key_id",
				Description: "AWS KMS key ID Amazon SageMaker uses to encrypt data when storing it on the ML storage volume attached to the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSagemakerEndpointConfiguration,
			},
			{
				Name:        "data_capture_config",
				Description: "Specifies the parameters to capture input/output of Sagemaker models endpoints.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "production_variants",
				Description: "An array of ProductionVariant objects, one for each model that you want to host at this endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSagemakerEndpointConfiguration,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the endpoint configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSageMakerEndpointConfigurationTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointConfigName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSageMakerEndpointConfigurationTags,
				Transform:   transform.FromField("Tags").Transform(sageMakerEndpointConfigurationTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointConfigArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSagemakerEndpointConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSagemakerEndpointConfigurations")

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List Call
	err = svc.ListEndpointConfigsPages(
		&sagemaker.ListEndpointConfigsInput{},
		func(page *sagemaker.ListEndpointConfigsOutput, isLast bool) bool {
			for _, config := range page.EndpointConfigs {
				d.StreamListItem(ctx, config)
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getSagemakerEndpointConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get config name
	var configName string
	if h.Item != nil {
		configName = *h.Item.(*sagemaker.EndpointConfigSummary).EndpointConfigName
	} else {
		configName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.DescribeEndpointConfigInput{
		EndpointConfigName: aws.String(configName),
	}

	// Get call
	data, err := svc.DescribeEndpointConfig(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getSagemakerEndpointConfiguration", "ERROR", err)
		return nil, err
	}
	return data, nil

}

func listSageMakerEndpointConfigurationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSageMakerEndpointConfigurationTags")

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}
	configArn := endpointConfigARN(h.Item)

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(configArn),
	}

	// Get call
	op, err := svc.ListTags(params)
	if err != nil {
		plugin.Logger(ctx).Debug("listSageMakerEndpointConfigurationTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func sageMakerEndpointConfigurationTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*sagemaker.ListTagsOutput)

	if data.Tags == nil {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range data.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}
	return turbotTagsMap, nil
}

func endpointConfigARN(item interface{}) string {
	switch item := item.(type) {
	case *sagemaker.EndpointConfigSummary:
		return *item.EndpointConfigArn
	case *sagemaker.DescribeEndpointConfigOutput:
		return *item.EndpointConfigArn
	}
	return ""
}
