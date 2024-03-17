package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"

	batchv1 "github.com/aws/aws-sdk-go/service/batch"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBatchJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_batch_job",
		Description: "AWS Batch Job",
		List: &plugin.ListConfig{
			Hydrate: listAwsBatchJobs,
			Tags:    map[string]string{"service": "batch", "action": "ListJobs"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "job_queue",
					Require: plugin.Required,
				},
				{
					Name:    "status",
					Require: plugin.Optional,
				},
				{
					Name:    "job_id",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: describeBatchJobs,
				Tags: map[string]string{"service": "batch", "action": "DescribeJobs"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(batchv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobArn"),
			},
			{
				Name:        "job_id",
				Description: "The job ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "job_name",
				Description: "The job name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobName"),
			},
			{
				Name:        "created_at",
				Description: "The Unix timestamp (in milliseconds) for when the job was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "job_definition",
				Description: "The Amazon Resource Name (ARN) of the job definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "started_at",
				Description: "The Unix timestamp for when the job was started. More specifically, it's when the job transitioned from the STARTING state to the RUNNING state.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The current status for the job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_reason",
				Description: "A short, human-readable string to provide more details for the current status of the job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stopped_at",
				Description: "The Unix timestamp for when the job was stopped. More specifically, it's when the job transitioned from the RUNNING state to a terminal state, such as SUCCEEDED or FAILED.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "job_queue",
				Description: "The Amazon Resource Name (ARN) of the job queue that the job is associated with.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "is_cancelled",
				Description: "Indicates whether the job is canceled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "is_terminated",
				Description: "Indicates whether the job is terminated.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "platform_capabilities",
				Description: "The platform capabilities required by the job definition. If no value is specified, it defaults to EC2.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "propagate_tags",
				Description: "Specifies whether to propagate the tags from the job or job definition to the corresponding Amazon ECS task. If no value is specified, the tags aren't propagated.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "scheduling_priority",
				Description: "The scheduling policy of the job definition. This only affects jobs in job queues with a fair share policy. Jobs with a higher scheduling priority are scheduled before jobs with a lower scheduling priority.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "share_identifier",
				Description: "The share identifier for the job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeBatchJobs,
			},

			//JSON columns
			{
				Name:        "array_properties",
				Description: "The array properties of the job, if it's an array job.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "container",
				Description: "An object that represents the details of the container that's associated with the job.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "single_node_properties",
				Description: "The node properties for a single node in a job summary list.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NodeProperties"),
			},
			{
				Name:        "attempts",
				Description: "A list of job attempts that are associated with this job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "depends_on",
				Description: "A list of job IDs that this job depends on.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "eks_attempts",
				Description: "A list of job attempts that are associated with this job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "eks_properties",
				Description: "An object with various properties that are specific to Amazon EKS based jobs.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "node_details",
				Description: "An object that represents the details of a node that's associated with a multi-node parallel job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "multi_node_properties",
				Description: "An object that represents the node properties of a multi-node parallel job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
				Transform:   transform.FromField("NodeProperties"),
			},
			{
				Name:        "parameters",
				Description: "Additional parameters that are passed to the job that replace parameter substitution placeholders or override any corresponding parameter defaults from the job definition.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "retry_strategy",
				Description: "The retry strategy to use for this job if an attempt fails.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},
			{
				Name:        "timeout",
				Description: "The timeout configuration for the job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeBatchJobs,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("JobArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsBatchJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_batch_job.listAwsBatchJobs", "connection_error", err)
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
			maxLimit = limit
		}
	}

	input := &batch.ListJobsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := batch.NewListJobsPaginator(svc, input, func(o *batch.ListJobsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_batch_job.listAwsBatchJobs", "api_error", err)
			return nil, err
		}

		for _, items := range output.JobSummaryList {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func describeBatchJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	jobId := d.EqualsQualString("id")
	if jobId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := BatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_batch_job.describeBatchJobs", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &batch.DescribeJobsInput{
		Jobs: []string{jobId},
	}

	// Get call
	data, err := svc.DescribeJobs(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_batch_job.describeBatchJobs", "api_error", err)
		return nil, err
	}

	return data.Jobs, nil
}
