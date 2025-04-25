package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/drs"
	"github.com/aws/aws-sdk-go-v2/service/drs/types"

	drsv1 "github.com/aws/aws-sdk-go/service/drs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsDRSJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_drs_job",
		Description: "AWS DRS Job",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "job_id", Require: plugin.Optional},
				{Name: "creation_date_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
				{Name: "end_date_time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UninitializedAccountException", "BadRequestException"}),
			},
			Hydrate: listAwsDRSJobs,
			Tags:    map[string]string{"service": "drs", "action": "DescribeJobs"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(drsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "job_id",
				Description: "The ID of the job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobID"),
			},
			{
				Name:        "arn",
				Description: "The ARN of a Job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "initiated_by",
				Description: "A string representing who initiated the Job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the Job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the Job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date_time",
				Description: "The date and time of when the Job was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_date_time",
				Description: "The date and time of when the Job ended.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "participating_servers",
				Description: "A list of servers that the Job is acting upon.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "participating_resources",
				Description: "A list of resources that the Job is acting upon.",
				Type:        proto.ColumnType_JSON,
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobID"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsDRSJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create service
	svc, err := DRSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_drs_job.listAwsDRSJobs", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(1000)
	input := drs.DescribeJobsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)

	filter := &types.DescribeJobsRequestFilters{}

	jobID := d.EqualsQualString("job_id")
	if jobID != "" {
		filter.JobIDs = []string{jobID}
	}

	quals := d.Quals
	if quals["creation_date_time"] != nil {
		for _, q := range quals["creation_date_time"].Quals {
			creationDateTime := q.Value.GetTimestampValue().AsTime().Format(time.RFC3339)
			switch q.Operator {
			case ">=", ">":
				filter.FromDate = aws.String(creationDateTime)
			case "<=", "<":
				filter.ToDate = aws.String(creationDateTime)
			case "=":
				filter.FromDate = aws.String(creationDateTime)
				filter.ToDate = aws.String(creationDateTime)
			}
		}
	}

	if quals["end_date_time"] != nil {
		for _, q := range quals["end_date_time"].Quals {
			endDateTime := q.Value.GetTimestampValue().AsTime().Format(time.RFC3339)
			switch q.Operator {
			case ">=", ">":
				filter.FromDate = aws.String(endDateTime)
			case "<=", "<":
				filter.ToDate = aws.String(endDateTime)
			case "=":
				filter.FromDate = aws.String(endDateTime)
				filter.ToDate = aws.String(endDateTime)
			}
		}
	}

	input.Filters = filter

	paginator := drs.NewDescribeJobsPaginator(svc, &input, func(o *drs.DescribeJobsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_drs_job.listAwsDRSJobs", "api_error", err)
			return nil, err
		}

		for _, job := range output.Items {
			d.StreamListItem(ctx, job)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
