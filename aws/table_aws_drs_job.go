package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/drs"
	"github.com/aws/aws-sdk-go-v2/service/drs/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsDRSJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_drs_job",
		Description: "AWS DRS Job",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "job_id", Require: plugin.Optional},
				{Name: "from_date", Require: plugin.Optional},
				{Name: "to_date", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UninitializedAccountException", "BadRequestException"}),
			},
			Hydrate: listAwsDRSJobs,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "job_id",
				Description: "The ID of the Job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobID"),
			},
			{
				Name:        "arn",
				Description: "The ARN of a Job.",
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
				Name:        "initiated_by",
				Description: "A string representing who initiated the Job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "participating_servers",
				Description: "A list of servers that the Job is acting upon.",
				Type:        proto.ColumnType_JSON,
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
				Name:        "from_date",
				Description: "The start date in a date range query.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("from_date"),
			},
			{
				Name:        "to_date",
				Description: "The end date in a date range query.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("to_date"),
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

	input.MaxResults = int32(maxItems)
	jobID := d.KeyColumnQualString("job_id")
	fromDate := d.KeyColumnQualString("from_date")
	toDate := d.KeyColumnQualString("to_date")

	filter := &types.DescribeJobsRequestFilters{}

	if jobID != "" {
		filter.JobIDs = []string{jobID}
	}

	if fromDate != "" {
		filter.FromDate = aws.String(fromDate)
	}

	if toDate != "" {
		filter.ToDate = aws.String(toDate)
	}

	input.Filters = filter

	paginator := drs.NewDescribeJobsPaginator(svc, &input, func(o *drs.DescribeJobsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_drs_job.listAwsDRSJobs", "api_error", err)
			return nil, err
		}

		for _, job := range output.Items {
			d.StreamListItem(ctx, job)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
