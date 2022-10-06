package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_job",
		Description: "AWS Glue Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueJobs,
		},
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the GlueJob.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the GlueJob.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueJobArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "allocated_capacity",
				Description: "[DEPRECATED] This column has been deprecated and will be removed in a future release, use max_capacity instead. The number of Glue data processing units (DPUs) that can be allocated when this job runs.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "created_on",
				Description: "The time and date that this job definition was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "glue_version",
				Description: "Glue version determines the versions of Apache Spark and Python that Glue supports.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_on",
				Description: "The last point in time when this job definition was modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "log_uri",
				Description: "This field is reserved for future use.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_capacity",
				Description: "The number of Glue data processing units (DPUs) that can be allocated when this job runs.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "max_retries",
				Description: "The maximum number of times to retry this job after a JobRun fails.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "number_of_workers",
				Description: "The number of workers of a defined workerType that are allocated when a job runs.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "role",
				Description: "The name or Amazon Resource Name (ARN) of the IAM role associated with this job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_configuration",
				Description: "The name of the SecurityConfiguration structure to be used with this job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timeout",
				Description: "The job timeout in minutes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "worker_type",
				Description: "The type of predefined worker that is allocated when a job runs. Accepts a value of Standard, G.1X, or G.2X.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "command",
				Description: "The JobCommand that runs this job.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "connections",
				Description: "The connections used for this job.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "default_arguments",
				Description: "The default arguments for this job, specified as name-value pairs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_property",
				Description: "An ExecutionProperty specifying the maximum number of concurrent runs allowed for this job.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "job_bookmark",
				Description: "Defines a point that a job can resume processing.",
				Hydrate:     getGlueJobBookmark,
				Transform:   transform.FromValue(),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "non_overridable_arguments",
				Description: "Non-overridable arguments for this job, specified as name-value pairs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "notification_property",
				Description: "Specifies configuration properties of a job notification.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueJobArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_job.listGlueJobs", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit=1
			} else {
				maxLimit = int32(*limit) 
			}
		}
	}
	input := &glue.GetJobsInput{
		MaxResults: aws.Int32(maxLimit),
	}
	// List call
	paginator:=glue.NewGetJobsPaginator(svc ,input,func(o *glue.GetJobsPaginatorOptions) {
		o.Limit=maxLimit
		o.StopOnDuplicateToken=true
	}) 
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_job.listGlueJobs", "api_error", err)
			return nil, err
		}
		for _, job  := range output.Jobs{
			d.StreamListItem(ctx, job)
			plugin.Logger(ctx).Error("aws_glue_job.listGlueJobs", "api_error", err)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil,nil
			}
		}
		
	}
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_job.listGlueJobs", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueJob(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_job.getGlueJob", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// Build the params
	params := &glue.GetJobInput{
		JobName: aws.String(name),
	}

	// Get call
	data, err := svc.GetJob(ctx,params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_job.getGlueJob", "api_error", err)
		return nil, err
	}
	return *data.Job, nil
}

func getGlueJobBookmark(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := h.Item.(types.Job).Name

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_job.getGlueJobBookmark", "service_creation_error", err)
		return nil, err
	}

	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// Build the params
	params := &glue.GetJobBookmarkInput{
		JobName: name,
	}

	// Get call
	data, err := svc.GetJobBookmark(ctx,params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_job.getGlueJobBookmark", "api_error", err)
		return nil, err
	}
	return data.JobBookmarkEntry, nil
}

func getGlueJobArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(types.Job)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_job.getGlueJobArn", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn format - https://docs.aws.amazon.com/glue/latest/dg/glue-specifying-resource-arns.html
	// arn:aws:glue:region:account-id:job/job-name
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":job/" + *data.Name

	return arn, nil
}
