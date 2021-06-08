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
		Description: "AWS Sagemaker Endpoint Configuratiion",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException", "NotFoundException"}),
			Hydrate:           getAwsSagemakerEndpointConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSagemakerEndpointConfiguration,
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
				Name:        "kms_key_id",
				Description: "AWS KMS key ID Amazon SageMaker uses to encrypt data when storing it on the ML storage volume attached to the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSagemakerEndpointConfiguration,
			},
			{
				Name:        "creation_time",
				Description: "The Amazon Resource Name (ARN) of the endpoint configuration.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "data_capture_config",
				Description: "The Amazon Resource Name (ARN) of the endpoint configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "production_variants",
				Description: "An array of ProductionVariant objects, one for each model that you want to host at this endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSagemakerEndpointConfiguration,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the endpoint configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerEndpointConfigurationTags,
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
				Hydrate:     listAwsSageMakerEndpointConfigurationTags,
				Transform:   transform.FromField("Tags").Transform(sageMakerEndpointConfigurationTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointConfigArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSagemakerEndpointConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsSagemakerEndpointConfiguration", "AWS_REGION", region)

	// Create Session
	svc, err := SageMakerService(ctx, d, region)
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

func getAwsSagemakerEndpointConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Get config name
	var configName string
	if h.Item != nil {
		configName = getConfigName(h.Item)
	} else {
		configName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerService(ctx, d, region)
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
		plugin.Logger(ctx).Debug("getAwsSagemakerEndpointConfiguration", "ERROR", err)
		return nil, err
	}
	return data, nil

}

func listAwsSageMakerEndpointConfigurationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsSageMakerEndpointConfigurationTags")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var configArn string
	switch h.Item.(type) {
	case *sagemaker.EndpointConfigSummary:
		configArn = *h.Item.(*sagemaker.EndpointConfigSummary).EndpointConfigArn
	case *sagemaker.DescribeEndpointConfigOutput:
		configArn = *h.Item.(*sagemaker.DescribeEndpointConfigOutput).EndpointConfigArn
	}

	// Create Session
	svc, err := SageMakerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(configArn),
	}

	// Get call
	op, err := svc.ListTags(params)
	if err != nil {
		logger.Debug("listAwsSageMakerEndpointConfigurationTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func sageMakerEndpointConfigurationTurbotTags(ctx context.Context, d *transform.TransformData) (interface{},
	error) {
	data := d.HydrateItem.(*sagemaker.ListTagsOutput)

	if data.Tags == nil {
		return nil, nil
	}

	if data.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

func getConfigName(item interface{}) string {
	switch item.(type) {
	case *sagemaker.EndpointConfigSummary:
		return *item.(*sagemaker.EndpointConfigSummary).EndpointConfigName
	case *sagemaker.DescribeEndpointConfigOutput:
		return *item.(*sagemaker.DescribeEndpointConfigOutput).EndpointConfigName
	}
	return ""
}
