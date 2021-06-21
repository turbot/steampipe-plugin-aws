package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/macie2"
)

//// TABLE DEFINITION

func tableAwsMacie2ClassificationJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_macie2_classification_job",
		Description: "AWS Macie2 Classification Job",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("job_id"),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationException"}),
			Hydrate:           getMacie2ClassificationJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listMacie2ClassificationJobs,
		},
		GetMatrixItem: BuildRegionList,
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listMacie2ClassificationJobs", "AWS_REGION", region)

	// Create Session
	svc, err := Macie2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListClassificationJobsPages(
		&macie2.ListClassificationJobsInput{},
		func(page *macie2.ListClassificationJobsOutput, isLast bool) bool {
			for _, job := range page.Items {
				d.StreamListItem(ctx, job)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMacie2ClassificationJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMacie2ClassificationJob")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var id string
	if h.Item != nil {
		id = *h.Item.(*macie2.JobSummary).JobId
	} else {
		id = d.KeyColumnQuals["job_id"].GetStringValue()
	}

	// Create service
	svc, err := Macie2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &macie2.DescribeClassificationJobInput{
		JobId: &id,
	}

	// Get call
	op, err := svc.DescribeClassificationJob(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}
