package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerNotebookInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_notebook_instance",
		Description: "AWS Sagemaker Notebook Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
			},
			Hydrate: getAwsSageMakerNotebookInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerNotebookInstances,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "default_code_repository", Require: plugin.Optional},
				{Name: "notebook_instance_lifecycle_config_name", Require: plugin.Optional},
				{Name: "notebook_instance_status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listAwsSageMakerNotebookInstances")

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &sagemaker.ListNotebookInstancesInput{
		MaxResults: aws.Int64(100),
	}

	equalQuals := d.KeyColumnQuals
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
			input.StatusEquals = aws.String(equalQuals["notebook_instance_status"].GetStringValue())
		}

	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListNotebookInstancesPages(
		input,
		func(page *sagemaker.ListNotebookInstancesOutput, isLast bool) bool {
			for _, notebookInstance := range page.NotebookInstances {
				d.StreamListItem(ctx, notebookInstance)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSageMakerNotebookInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(*sagemaker.NotebookInstanceSummary).NotebookInstanceName
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.DescribeNotebookInstanceInput{
		NotebookInstanceName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeNotebookInstance(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsSageMakerNotebookInstance", "ERROR", err)
		return nil, err
	}
	return data, nil
}

func listAwsSageMakerNotebookInstanceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsSageMakerNotebookInstanceTags")

	resourceArn := notebookInstanceARN(h.Item)

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(resourceArn),
	}

	pagesLeft := true
	tags := []*sagemaker.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListTags(params)
		if err != nil {
			plugin.Logger(ctx).Error("listAwsSageMakerNotebookInstanceTags", "ListTags_error", err)
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
	case *sagemaker.NotebookInstanceSummary:
		return *item.NotebookInstanceArn
	case *sagemaker.DescribeNotebookInstanceOutput:
		return *item.NotebookInstanceArn
	}
	return ""
}
