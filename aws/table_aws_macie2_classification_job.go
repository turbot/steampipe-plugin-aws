package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/macie2"
)

//// TABLE DEFINITION

func tableAwsMacie2ClassificationJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_macie2_classification_job",
		Description: "AWS Macie2 Classification Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("job_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "InvalidParameter"}),
			},
			Hydrate: getMacie2ClassificationJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listMacie2ClassificationJobs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "job_status", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "job_type", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listMacie2ClassificationJobs")

	// Create Session
	svc, err := Macie2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Service is not supported in the region
	if svc == nil {
		return nil, nil
	}

	input := &macie2.ListClassificationJobsInput{
		MaxResults: aws.Int64(100),
	}

	filterCriteris := buildMacie2ClassificationJobsFilterCriteria(d.Quals)

	if len(filterCriteris.Excludes) > 0 || len(filterCriteris.Includes) > 0 {
		input.FilterCriteria = filterCriteris
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
	err = svc.ListClassificationJobsPages(
		input,
		func(page *macie2.ListClassificationJobsOutput, isLast bool) bool {
			for _, job := range page.Items {
				d.StreamListItem(ctx, job)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Throws "AccessDeniedException: Macie is not enabled." when AWS Macie is not enabled in a region
			// also the API throws AccessDeniedException if the request does not have proper permission
			// with the below check we will only handle "Macie is not enabled"
			if awsErr.Message() == "Macie is not enabled." {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("listMacie2ClassificationJobs", "ListClassificationJobsPages_error", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMacie2ClassificationJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMacie2ClassificationJob")

	var id string
	if h.Item != nil {
		id = *h.Item.(*macie2.JobSummary).JobId
	} else {
		id = d.KeyColumnQuals["job_id"].GetStringValue()
	}

	// empty check for job id
	if id == "" {
		return nil, nil
	}

	// Create service
	svc, err := Macie2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Service is not supported in the region
	if svc == nil {
		return nil, nil
	}

	// Build params
	params := &macie2.DescribeClassificationJobInput{
		JobId: &id,
	}

	// Get call
	op, err := svc.DescribeClassificationJob(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Throws "AccessDeniedException: Macie is not enabled." when AWS Macie is not enabled in a region
			// also the API throws AccessDeniedException if the request does not have proper permission
			// with the below check we will only handle "Macie is not enabled"
			if awsErr.Message() == "Macie is not enabled." {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("getMacie2ClassificationJob", "DescribeClassificationJob_error", err)
		return nil, err
	}

	return op, nil
}

//// UTILITY FUNCTION
//// Build macie2 list job classification job filter

func buildMacie2ClassificationJobsFilterCriteria(quals plugin.KeyColumnQualMap) *macie2.ListJobsFilterCriteria {
	filterCriteria := &macie2.ListJobsFilterCriteria{}

	filterQuals := map[string]string{
		"name":       "name",
		"job_type":   "jobType",
		"job_status": "jobStatus",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			for _, q := range quals[columnName].Quals {
				value := getQualsValueByColumn(quals, columnName, "string")

				filter := &macie2.ListJobsFilterTerm{
					Comparator: aws.String(macie2.JobComparatorEq),
				}

				val, ok := value.(string)
				if ok {
					filter.Values = []*string{aws.String(val)}
				} else {
					filter.Values = value.([]*string)
				}

				if filterName == "name" {
					filter.Key = aws.String(macie2.ListJobsFilterKeyName)
					switch q.Operator {
					case "<>":
						filterCriteria.Excludes = append(filterCriteria.Excludes, filter)
					case "=":
						filterCriteria.Includes = append(filterCriteria.Includes, filter)
					}
				}
				if filterName == "jobType" {
					filter.Key = aws.String(macie2.ListJobsFilterKeyJobType)
					switch q.Operator {
					case "<>":
						filterCriteria.Excludes = append(filterCriteria.Excludes, filter)
					case "=":
						filterCriteria.Includes = append(filterCriteria.Includes, filter)
					}
				}
				if filterName == "jobStatus" {
					filter.Key = aws.String(macie2.ListJobsFilterKeyJobStatus)
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
