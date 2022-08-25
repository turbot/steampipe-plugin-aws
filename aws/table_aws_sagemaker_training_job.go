package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSageMakerTrainingJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_training_job",
		Description: "AWS SageMaker Training Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "NotFoundException", "RecordNotFound"}),
			},
			Hydrate: getAwsSageMakerTrainingJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerTrainingJobs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "creation_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
				{Name: "last_modified_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
				{Name: "training_job_status", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listAwsSageMakerTrainingJobs")

	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &sagemaker.ListTrainingJobsInput{
		MaxResults: aws.Int64(100),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["training_job_status"] != nil {
		input.StatusEquals = aws.String(equalQuals["training_job_status"].GetStringValue())
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
	err = svc.ListTrainingJobsPages(
		input,
		func(page *sagemaker.ListTrainingJobsOutput, isLast bool) bool {
			for _, job := range page.TrainingJobSummaries {
				d.StreamListItem(ctx, job)

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

func getAwsSageMakerTrainingJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		name = *h.Item.(*sagemaker.TrainingJobSummary).TrainingJobName
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create service
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.DescribeTrainingJobInput{
		TrainingJobName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeTrainingJob(params)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getAwsSageMakerTrainingJobTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSageMakerTrainingJobTags")

	arn := trainingJobArn(h.Item)
	// Create Session
	svc, err := SageMakerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(arn),
	}

	pagesLeft := true
	tags := []*sagemaker.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListTags(params)
		if err != nil {
			plugin.Logger(ctx).Error("getAwsSageMakerTrainingJobTags", "ListTags_error", err)
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
	case *sagemaker.TrainingJobSummary:
		return *item.TrainingJobArn
	case *sagemaker.DescribeTrainingJobOutput:
		return *item.TrainingJobArn
	}
	return ""
}
