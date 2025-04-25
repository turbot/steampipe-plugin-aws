package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/macie2"
	"github.com/aws/aws-sdk-go-v2/service/macie2/types"

	macie2v1 "github.com/aws/aws-sdk-go/service/macie2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMacie2ClassificationJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_macie2_classification_job",
		Description: "AWS Macie2 Classification Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("job_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "InvalidParameter"}),
			},
			Hydrate: getMacie2ClassificationJob,
			Tags:    map[string]string{"service": "macie2", "action": "DescribeClassificationJob"},
		},
		List: &plugin.ListConfig{
			Hydrate: listMacie2ClassificationJobs,
			Tags:    map[string]string{"service": "macie2", "action": "ListClassificationJobs"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "job_status", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "job_type", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getMacie2ClassificationJob,
				Tags: map[string]string{"service": "macie2", "action": "DescribeClassificationJob"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(macie2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The custom name of the job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "job_id",
				Description: "The unique identifier for the job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2ClassificationJob,
				Transform:   transform.FromField("JobArn"),
			},
			{
				Name:        "job_status",
				Description: "The status of a classification job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "job_type",
				Description: "The schedule for running a classification job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The custom description of the job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "managed_data_identifier_selector",
				Description: "The selection type that determines which managed data identifiers the job uses when it analyzes data.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "initial_run",
				Description: "For a recurring job, specifies whether you configured the job to analyze all existing, eligible objects immediately after the job was created (true). ",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "client_token",
				Description: "The token that was provided to ensure the idempotency of the request to create the job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "created_at",
				Description: "The date and time, in UTC and extended ISO 8601 format, when the job was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_run_time",
				Description: "This value indicates when the most recent run started.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "sampling_percentage",
				Description: "The sampling depth, as a percentage, that determines the percentage of eligible objects that the job analyzes.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "bucket_definitions",
				Description: "The namespace of the AWS service that provides the resource, or a custom-resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "custom_data_identifier_ids",
				Description: "The custom data identifiers that the job uses to analyze data.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "last_run_error_status",
				Description: "Specifies whether any account- or bucket-level access errors occurred when a classification job ran.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "s3_job_definition",
				Description: "Specifies which S3 buckets contain the objects that a classification job analyzes, and the scope of that analysis.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "allow_list_ids",
				Description: "An array of unique identifiers, one for each allow list that the job uses when it analyzes data.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "managed_data_identifier_ids",
				Description: "An array of unique identifiers, one for each managed data identifier that the job is explicitly configured to include (use) or exclude (not use) when it analyzes data.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "schedule_frequency",
				Description: "Specifies the recurrence pattern for running a classification job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "statistics",
				Description: "Provides processing statistics for a classification job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "user_paused_details",
				Description: "Provides information about when a classification job was paused.",
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
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2ClassificationJob,
				Transform:   transform.FromField("JobArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listMacie2ClassificationJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := Macie2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("table_aws_macie2_classification_job.listMacie2ClassificationJobs", "client_error", err)
		return nil, err
	}
	// Service is not supported in the region
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(200)
	input := &macie2.ListClassificationJobsInput{
		MaxResults: aws.Int32(maxItems),
	}

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

	filterCriteria := buildMacie2ClassificationJobsFilterCriteria(d.Quals)

	if len(filterCriteria.Excludes) > 0 || len(filterCriteria.Includes) > 0 {
		input.FilterCriteria = filterCriteria
	}

	// List call
	paginator := macie2.NewListClassificationJobsPaginator(svc, input, func(o *macie2.ListClassificationJobsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Throws "AccessDeniedException: Macie is not enabled." when AWS Macie is not enabled in a region
			// also the API throws AccessDeniedException if the request does not have proper permission
			// with the below check we will only handle "Macie is not enabled"
			if strings.Contains(err.Error(), "Macie is not enabled.") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("table_aws_macie2_classification_job.listMacie2ClassificationJobs", "api_error", err)
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

//// HYDRATE FUNCTIONS

func getMacie2ClassificationJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var id string
	if h.Item != nil {
		id = *h.Item.(types.JobSummary).JobId
	} else {
		id = d.EqualsQuals["job_id"].GetStringValue()
	}

	// empty check for job id
	if id == "" {
		return nil, nil
	}

	// Create session
	svc, err := Macie2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("table_aws_macie2_classification_job.getMacie2ClassificationJob", "client_error", err)
		return nil, err
	}
	// Service is not supported in the region
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &macie2.DescribeClassificationJobInput{
		JobId: &id,
	}

	// Get call
	op, err := svc.DescribeClassificationJob(ctx, params)
	if err != nil {
		// Throws "AccessDeniedException: Macie is not enabled." when AWS Macie is not enabled in a region
		// also the API throws AccessDeniedException if the request does not have proper permission
		// with the below check we will only handle "Macie is not enabled"
		if strings.Contains(err.Error(), "Macie is not enabled.") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("table_aws_macie2_classification_job.listMacie2ClassificationJobs", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// UTILITY FUNCTION
//// Build macie2 list job classification job filter

func buildMacie2ClassificationJobsFilterCriteria(quals plugin.KeyColumnQualMap) *types.ListJobsFilterCriteria {
	filterCriteria := &types.ListJobsFilterCriteria{}

	filterQuals := map[string]string{
		"name":       "name",
		"job_type":   "jobType",
		"job_status": "jobStatus",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			for _, q := range quals[columnName].Quals {
				value := getQualsValueByColumn(quals, columnName, "string")

				filter := types.ListJobsFilterTerm{
					Comparator: types.JobComparatorEq,
				}

				val, ok := value.(string)
				if ok {
					filter.Values = []string{val}
				} else {
					filter.Values = value.([]string)
				}

				if filterName == "name" {
					filter.Key = types.ListJobsFilterKeyName
					switch q.Operator {
					case "<>":
						filterCriteria.Excludes = append(filterCriteria.Excludes, filter)
					case "=":
						filterCriteria.Includes = append(filterCriteria.Includes, filter)
					}
				}
				if filterName == "jobType" {
					filter.Key = types.ListJobsFilterKeyJobType
					switch q.Operator {
					case "<>":
						filterCriteria.Excludes = append(filterCriteria.Excludes, filter)
					case "=":
						filterCriteria.Includes = append(filterCriteria.Includes, filter)
					}
				}
				if filterName == "jobStatus" {
					filter.Key = types.ListJobsFilterKeyJobStatus
					switch q.Operator {
					case "<>":
						filterCriteria.Excludes = append(filterCriteria.Excludes, filter)
					case "=":
						filterCriteria.Includes = append(filterCriteria.Includes, filter)
					}
				}
			}
		}
	}

	return filterCriteria
}
