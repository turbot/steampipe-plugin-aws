package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsBatchQueue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_batch_queue",
		Description: "AWS Batch Queue",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("job_queue_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"JobQueueNotFoundException"}),
			},
			Hydrate: getBatchQueue,
			Tags:    map[string]string{"service": "batch", "action": "DescribeJobQueues"},
		},
		List: &plugin.ListConfig{
			Hydrate: listBatchQueues,
			Tags:    map[string]string{"service": "batch", "action": "DescribeJobQueues"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BATCH_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "job_queue_name",
				Description: "The name of the job queue",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the job queue",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobQueueArn"),
			},
			{
				Name:        "state",
				Description: "The state of the job queue (ENABLED or DISABLED)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the job queue (CREATING, UPDATING, DELETING, or DELETED)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "priority",
				Description: "The priority of the job queue",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "compute_environment_order",
				Description: "The compute environments that are attached to the job queue and the order in which job placement is preferred",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "scheduling_policy_arn",
				Description: "The ARN of the scheduling policy",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_reason",
				Description: "A short, human-readable string to provide additional details about the current status of the job queue",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: "The tags assigned to the job queue",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobQueueName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("JobQueueArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

func listBatchQueues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service client
	svc, err := BatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_batch_queue.listBatchQueues", "client_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
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

	input := &batch.DescribeJobQueuesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := batch.NewDescribeJobQueuesPaginator(svc, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_batch_queue.listBatchQueues", "api_error", err)
			return nil, err
		}

		for _, queue := range output.JobQueues {
			d.StreamListItem(ctx, queue)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getBatchQueue(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	queueName := d.EqualsQualString("job_queue_name")
	if queueName == "" {
		return nil, nil
	}

	// Create service client
	svc, err := BatchClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_batch_queue.getBatchQueue", "client_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	input := &batch.DescribeJobQueuesInput{
		JobQueues: []string{queueName},
	}

	output, err := svc.DescribeJobQueues(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_batch_queue.getBatchQueue", "api_error", err)
		return nil, err
	}

	if len(output.JobQueues) == 0 {
		return nil, nil
	}

	return output.JobQueues[0], nil
}
