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

func tableAwsSageMakerNotebookInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_notebook_instance",
		Description: "AWS Sagemaker Notebook Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
			},
			Hydrate: getAwsSageMakerNotebookInstance,
			Tags:    map[string]string{"service": "sagemaker", "action": "DescribeNotebookInstance"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerNotebookInstances,
			Tags:    map[string]string{"service": "sagemaker", "action": "ListNotebookInstances"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "default_code_repository", Require: plugin.Optional},
				{Name: "notebook_instance_lifecycle_config_name", Require: plugin.Optional},
				{Name: "notebook_instance_status", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsSageMakerNotebookInstance,
				Tags: map[string]string{"service": "sagemaker", "action": "DescribeNotebookInstance"},
			},
			{
				Func: listAwsSageMakerNotebookInstanceTags,
				Tags: map[string]string{"service": "sagemaker", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(sagemakerv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the notebook instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NotebookInstanceName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the notebook instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NotebookInstanceArn"),
			},
			{
				Name:        "creation_time",
				Description: "A timestamp that shows when the notebook instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "default_code_repository",
				Description: "The Git repository associated with the notebook instance as its default code repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "direct_internet_access",
				Description: "Describes whether Amazon SageMaker provides internet access to the notebook instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "failure_reason",
				Description: "If status is Failed, the reason it failed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "instance_type",
				Description: "The type of ML compute instance that the notebook instance is running on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The AWS KMS key ID Amazon SageMaker uses to encrypt data when storing it on the ML storage volume attached to the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "last_modified_time",
				Description: "A timestamp that shows when the notebook instance was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "network_interface_id",
				Description: "The network interface IDs that Amazon SageMaker created at the time of creating the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "notebook_instance_lifecycle_config_name",
				Description: "The name of a notebook instance lifecycle configuration associated with this notebook instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "notebook_instance_status",
				Description: "The status of the notebook instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role associated with the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "root_access",
				Description: "Whether root access is enabled or disabled for users of the notebook instance.Lifecycle configurations need root access to be able to set up a notebook instance",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "subnet_id",
				Description: "The ID of the VPC subnet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "url",
				Description: "The URL that you use to connect to the Jupyter notebook that is running in your notebook instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "volume_size_in_gb",
				Description: "The size, in GB, of the ML storage volume attached to the notebook instance.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "accelerator_types",
				Description: "The list of the Elastic Inference (EI) instance types associated with this notebook instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "additional_code_repositories",
				Description: "An array of up to three Git repositories associated with the notebook instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_groups",
				Description: "The IDs of the VPC security groups.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerNotebookInstance,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the notebook instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerNotebookInstanceTags,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NotebookInstanceName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerNotebookInstanceTags,
				Transform:   transform.From(sageMakerTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NotebookInstanceArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSageMakerNotebookInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_notebook_instance.listAwsSageMakerNotebookInstances", "connection_error", err)
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

	input := &sagemaker.ListNotebookInstancesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	equalQuals := d.EqualsQuals
	if equalQuals["default_code_repository"] != nil {
		if equalQuals["default_code_repository"].GetStringValue() != "" {
			input.DefaultCodeRepositoryContains = aws.String(equalQuals["default_code_repository"].GetStringValue())
		}
	}
	if equalQuals["notebook_instance_lifecycle_config_name"] != nil {
		if equalQuals["notebook_instance_lifecycle_config_name"].GetStringValue() != "" {
			input.NotebookInstanceLifecycleConfigNameContains = aws.String(equalQuals["notebook_instance_lifecycle_config_name"].GetStringValue())
		}
	}
	if equalQuals["notebook_instance_status"] != nil {
		if equalQuals["notebook_instance_status"].GetStringValue() != "" {
			input.StatusEquals = types.NotebookInstanceStatus(equalQuals["notebook_instance_status"].GetStringValue())
		}

	}

	paginator := sagemaker.NewListNotebookInstancesPaginator(svc, input, func(o *sagemaker.ListNotebookInstancesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_notebook_instance.listAwsSageMakerNotebookInstances", "api_error", err)
			return nil, err
		}

		for _, items := range output.NotebookInstances {
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

func getAwsSageMakerNotebookInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(types.NotebookInstanceSummary).NotebookInstanceName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_notebook_instance.getAwsSageMakerNotebookInstance", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &sagemaker.DescribeNotebookInstanceInput{
		NotebookInstanceName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeNotebookInstance(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_notebook_instance.getAwsSageMakerNotebookInstance", "api_error", err)
		return nil, err
	}
	return data, nil
}

func listAwsSageMakerNotebookInstanceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	resourceArn := notebookInstanceARN(h.Item)

	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_notebook_instance.listAwsSageMakerNotebookInstanceTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}
	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(resourceArn),
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		keyTags, err := svc.ListTags(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_notebook_instance.listAwsSageMakerNotebookInstanceTags", "api_error", err)
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

//// TRANSFORM FUNCTION

func notebookInstanceARN(item interface{}) string {
	switch item := item.(type) {
	case types.NotebookInstanceSummary:
		return *item.NotebookInstanceArn
	case *sagemaker.DescribeNotebookInstanceOutput:
		return *item.NotebookInstanceArn
	}
	return ""
}
