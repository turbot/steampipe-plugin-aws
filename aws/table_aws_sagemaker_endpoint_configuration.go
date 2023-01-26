package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

	sagemakerv1 "github.com/aws/aws-sdk-go/service/sagemaker"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerEndpointConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_endpoint_configuration",
		Description: "AWS Sagemaker Endpoint Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "NotFoundException"}),
			},
			Hydrate: getSagemakerEndpointConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listSagemakerEndpointConfigurations,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "creation_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(sagemakerv1.EndpointsID),
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
				Transform:   transform.FromValue(),
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
				Transform:   transform.FromValue().Transform(sageMakerTurbotTags),
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
	// Create client
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_endpoint_configuration.listSagemakerEndpointConfigurations", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &sagemaker.ListEndpointConfigsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	quals := d.Quals
	if quals["timestamp"] != nil {
		for _, q := range quals["timestamp"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				input.CreationTimeAfter = aws.Time(timestamp)
			case "<", "<=":
				input.CreationTimeBefore = aws.Time(timestamp)
			}
		}
	}

	paginator := sagemaker.NewListEndpointConfigsPaginator(svc, input, func(o *sagemaker.ListEndpointConfigsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_endpoint_configuration.listSagemakerEndpointConfigurations", "api_error", err)
			return nil, err
		}

		for _, items := range output.EndpointConfigs {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSagemakerEndpointConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get config name
	var configName string
	if h.Item != nil {
		configName = *h.Item.(types.EndpointConfigSummary).EndpointConfigName
	} else {
		configName = d.EqualsQuals["name"].GetStringValue()
	}

	// Create client
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_endpoint_configuration.getSagemakerEndpointConfiguration", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &sagemaker.DescribeEndpointConfigInput{
		EndpointConfigName: aws.String(configName),
	}

	// Get call
	data, err := svc.DescribeEndpointConfig(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_endpoint_configuration.getSagemakerEndpointConfiguration", "api_error", err)
		return nil, err
	}
	return data, nil

}

func listSageMakerEndpointConfigurationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	configArn := endpointConfigARN(h.Item)

	// Create client
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_endpoint_configuration.listSageMakerEndpointConfigurationTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(configArn),
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListTags(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_endpoint_configuration.listSageMakerEndpointConfigurationTags", "api_error", err)
			return nil, err
		}
		tags = append(tags, keyTags.Tags...)

		if keyTags.NextToken != nil {
			params.NextToken = keyTags.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func endpointConfigARN(item interface{}) string {
	switch item := item.(type) {
	case types.EndpointConfigSummary:
		return *item.EndpointConfigArn
	case *sagemaker.DescribeEndpointConfigOutput:
		return *item.EndpointConfigArn
	}
	return ""
}
