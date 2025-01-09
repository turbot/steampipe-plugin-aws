package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

	sagemakerEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSageMakerTrainingJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_training_job",
		Description: "AWS SageMaker Training Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
			},
			Hydrate: getAwsSageMakerTrainingJob,
			Tags:    map[string]string{"service": "sagemaker", "action": "DescribeTrainingJob"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerTrainingJobs,
			Tags:    map[string]string{"service": "sagemaker", "action": "ListTrainingJobs"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "creation_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
				{Name: "last_modified_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
				{Name: "training_job_status", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsSageMakerTrainingJob,
				Tags: map[string]string{"service": "sagemaker", "action": "DescribeTrainingJob"},
			},
			{
				Func: getAwsSageMakerTrainingJobTags,
				Tags: map[string]string{"service": "sagemaker", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(sagemakerEndpoint.API_SAGEMAKERServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the training job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrainingJobName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the training job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrainingJobArn"),
			},
			{
				Name:        "training_job_status",
				Description: "The status of the training job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_ml_job_arn",
				Description: "The Amazon Resource Name (ARN) of an AutoML job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "billable_time_in_seconds",
				Description: "The billable time in seconds. Billable time refers to the absolute wall-clock time.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "creation_time",
				Description: "A timestamp that shows when the training job was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enable_managed_spot_training",
				Description: "A Boolean indicating whether managed spot training is enabled or not.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "enable_network_isolation",
				Description: "Specifies enable network isolation for training jobs.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "enable_inter_container_traffic_encryption",
				Description: "To encrypt all communications between ML compute instances in distributed training, choose True.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "failure_reason",
				Description: "If the training job failed, the reason it failed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "labeling_job_arn",
				Description: "The Amazon Resource Name (ARN) of the Amazon SageMaker Ground Truth labeling job that created the transform or training job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "last_modified_time",
				Description: "Timestamp when the training job was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "profiling_status",
				Description: "Profiling status of a training job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "role_arn",
				Description: "The AWS Identity and Access Management (IAM) role configured for the training job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "secondary_status",
				Description: "Provides detailed information about the state of the training job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "training_end_time",
				Description: "A timestamp that shows when the training job ended.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "training_start_time",
				Description: "Indicates the time when the training job starts on training instances.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "training_time_in_seconds",
				Description: "The training time in seconds.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "tuning_job_arn",
				Description: "The Amazon Resource Name (ARN) of the associated hyperparameter tuning job if the training job was launched by a hyperparameter tuning job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "enable_infra_check",
				Description: "Enables an infrastructure health check.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsSageMakerTrainingJob,
				Transform:   transform.FromField("InfraCheckConfig.EnableInfraCheck"),
			},
			{
				Name:        "enable_remote_debug",
				Description: "If set to True, enables remote debugging.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsSageMakerTrainingJob,
				Transform:   transform.FromField("InfraCheckConfig.EnableRemoteDebug"),
			},
			{
				Name:        "maximum_retry_attempts",
				Description: "The number of times to retry the job.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAwsSageMakerTrainingJob,
				Transform:   transform.FromField("RetryStrategy.MaximumRetryAttempts"),
			},
			{
				Name:        "warm_pool_status",
				Description: "The status of the warm pool associated with the training job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "algorithm_specification",
				Description: "Information about the algorithm used for training, and algorithm metadata.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "checkpoint_config",
				Description: "Contains information about the output location for managed spot training checkpoint data.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "debug_hook_config",
				Description: "Configuration information for the Debugger hook parameters, metric and tensor collections, and storage paths.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "debug_rule_configurations",
				Description: "Configuration information for Debugger rules for debugging output tensors.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "debug_rule_evaluation_statuses",
				Description: "Evaluation status of Debugger rules for debugging on a training job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "environment",
				Description: "The environment variables to set in the Docker container.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "experiment_config",
				Description: "Associates a SageMaker job as a trial component with an experiment and trial.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "final_metric_data_list",
				Description: "A collection of MetricData objects that specify the names, values, and dates and times that the training algorithm emitted to Amazon CloudWatch.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "hyper_parameters",
				Description: "Algorithm-specific parameters.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "input_data_config",
				Description: "An array of Channel objects that describes each data input channel.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "model_artifacts",
				Description: "Information about the Amazon S3 location that is configured for storing model artifacts.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "output_data_config",
				Description: "The S3 path where model artifacts that you configured when creating the job are stored.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "profiler_config",
				Description: "Configuration information for Debugger system monitoring,framework profiling and storage paths.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "profiler_rule_configurations",
				Description: "Configuration information for Debugger rules for profiling system and framework metrics.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "profiler_rule_evaluation_statuses",
				Description: "Evaluation status of Debugger rules for profiling on a training job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "resource_config",
				Description: "Resources, including ML compute instances and ML storage volumes, that are configured for model training.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "secondary_status_transitions",
				Description: "A history of all of the secondary statuses that the training job has transitioned through.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "stopping_condition",
				Description: "Specifies a limit to how long a model training job can run.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "tensor_board_output_config",
				Description: "Configuration of storage locations for the Debugger TensorBoard output data.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "vpc_config",
				Description: "A VpcConfig object that specifies the VPC that this training job has access to.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJob,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the training job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJobTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrainingJobName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerTrainingJobTags,
				Transform:   transform.FromValue().Transform(sageMakerTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TrainingJobArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSageMakerTrainingJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Client
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_training_job.listAwsSageMakerTrainingJobs", "connection_error", err)
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

	input := &sagemaker.ListTrainingJobsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	equalQuals := d.EqualsQuals
	if equalQuals["training_job_status"] != nil {
		input.StatusEquals = types.TrainingJobStatus(equalQuals["training_job_status"].GetStringValue())
	}

	quals := d.Quals
	if quals["creation_time"] != nil {
		for _, q := range quals["creation_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				input.CreationTimeAfter = aws.Time(timestamp)
			case "<", "<=":
				input.CreationTimeBefore = aws.Time(timestamp)
			}
		}
	}

	if quals["last_modified_time"] != nil {
		for _, q := range quals["last_modified_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				input.LastModifiedTimeAfter = aws.Time(timestamp)
			case "<", "<=":
				input.LastModifiedTimeBefore = aws.Time(timestamp)
			}
		}
	}

	paginator := sagemaker.NewListTrainingJobsPaginator(svc, input, func(o *sagemaker.ListTrainingJobsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_training_job.listAwsSageMakerTrainingJobs", "api_error", err)
			return nil, err
		}

		for _, items := range output.TrainingJobSummaries {
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

func getAwsSageMakerTrainingJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(types.TrainingJobSummary).TrainingJobName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Create Client
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_training_job.getAwsSageMakerTrainingJob", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &sagemaker.DescribeTrainingJobInput{
		TrainingJobName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeTrainingJob(ctx, params)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getAwsSageMakerTrainingJobTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := trainingJobArn(h.Item)
	// Create Client
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_training_job.getAwsSageMakerTrainingJobTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(arn),
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		keyTags, err := svc.ListTags(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_training_job.getAwsSageMakerTrainingJobTags", "api_error", err)
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

func trainingJobArn(item interface{}) string {
	switch item := item.(type) {
	case types.TrainingJobSummary:
		return *item.TrainingJobArn
	case *sagemaker.DescribeTrainingJobOutput:
		return *item.TrainingJobArn
	}
	return ""
}
