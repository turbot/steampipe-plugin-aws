package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspectorAssessmentRun(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_assessment_run",
		Description: "AWS Inspector Assessment Run",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getInspectorAssessmentRun,
		},
		List: &plugin.ListConfig{
			Hydrate: listInspectorAssessmentRuns,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The auto-generated name for the assessment run.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "arn",
				Description: "The ARN of the assessment run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assessment_template_arn",
				Description: "The ARN of the assessment template that is associated with the assessment run.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "state",
				Description: "The state of the assessment run.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "completed_at",
				Description: "The assessment run completion time that corresponds to the rules packages evaluation completion time or failure.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "created_at",
				Description: "The time when StartAssessmentRun was called.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "data_collected",
				Description: " Boolean value (true or false) that specifies whether the process of collecting data from the agents is completed.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "duration_in_seconds",
				Description: "The duration of the assessment run.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "started_at",
				Description: "The time when StartAssessmentRun was called.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "state_changed_at",
				Description: "The last time when the assessment run's state changed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "finding_counts",
				Description: "Provides a total count of generated findings per severity.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "notifications",
				Description: "A list of notifications for the event subscriptions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "rules_package_arns",
				Description: "The rules packages selected for the assessment run.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "state_changes",
				Description: "A list of the assessment run state changes.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "user_attributes_for_findings",
				Description: "The user-defined attributes that are assigned to every generated finding.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInspectorAssessmentRun,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
				Hydrate:     getInspectorAssessmentRun,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
				Hydrate:     getInspectorAssessmentRun,
			},
		}),
	}
}

//// LIST FUNCTION

func listInspectorAssessmentRuns(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &inspector.ListAssessmentRunsInput{
		MaxResults: aws.Int64(500),
	}
	// equalQuals := d.KeyColumnQuals
	// if equalQuals["name"] != nil {
	// 	if equalQuals["name"].GetStringValue() != "" {
	// 		input.Filter = &inspector.AssessmentRunFilter{AssessmentRunNamePattern: aws.String(equalQuals["name"].GetStringValue())}
	// 	}
	// }

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
	err = svc.ListAssessmentRunsPages(
		input,
		func(page *inspector.ListAssessmentRunsOutput, isLast bool) bool {
			for _, assessmentRunArn := range page.AssessmentRunArns {
				d.StreamListItem(ctx, &inspector.AssessmentRun{
					Arn: assessmentRunArn,
				})

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

func getInspectorAssessmentRun(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getInspectorAssessmentRun")

	var assessmentRunArn string
	if h.Item != nil {
		assessmentRunArn = *h.Item.(*inspector.AssessmentRun).Arn
	} else {
		quals := d.KeyColumnQuals
		assessmentRunArn = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &inspector.DescribeAssessmentRunsInput{
		AssessmentRunArns: []*string{aws.String(assessmentRunArn)},
	}

	// Get call
	data, err := svc.DescribeAssessmentRuns(params)
	if err != nil {
		logger.Debug("describeAssessmentRun__", "ERROR", err)
		return nil, err
	}
	if data.AssessmentRuns != nil && len(data.AssessmentRuns) > 0 {
		return data.AssessmentRuns[0], nil
	}
	return nil, nil
}
