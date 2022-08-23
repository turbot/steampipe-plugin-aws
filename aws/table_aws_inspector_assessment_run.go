package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspectorAssessmentRun(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_assessment_run",
		Description: "AWS Inspector Assessment Run",
		List: &plugin.ListConfig{
			Hydrate: listInspectorAssessmentRuns,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "assessment_template_arn", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The auto-generated name for the assessment run.",
				Type:        proto.ColumnType_STRING,
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
			},
			{
				Name:        "state",
				Description: "The state of the assessment run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "completed_at",
				Description: "The assessment run completion time that corresponds to the rules packages evaluation completion time or failure.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_at",
				Description: "The time when StartAssessmentRun was called.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "data_collected",
				Description: "Boolean value (true or false) that specifies whether the process of collecting data from the agents is completed.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "duration_in_seconds",
				Description: "The duration of the assessment run.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "started_at",
				Description: "The time when StartAssessmentRun was called.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state_changed_at",
				Description: "The last time when the assessment run's state changed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "finding_counts",
				Description: "Provides a total count of generated findings per severity.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "notifications",
				Description: "A list of notifications for the event subscriptions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "rules_package_arns",
				Description: "The rules packages selected for the assessment run.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_changes",
				Description: "A list of the assessment run state changes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_attributes_for_findings",
				Description: "The user-defined attributes that are assigned to every generated finding.",
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
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listInspectorAssessmentRuns(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var assessmentRunArns []*string

	input := &inspector.ListAssessmentRunsInput{
		MaxResults: aws.Int64(500),
	}

	filter := &inspector.AssessmentRunFilter{}

	if d.KeyColumnQuals["assessment_template_arn"].GetStringValue() != "" {
		input.AssessmentTemplateArns = aws.StringSlice([]string{d.KeyColumnQuals["assessment_template_arn"].GetStringValue()})
	}
	if d.KeyColumnQuals["name"].GetStringValue() != "" {
		filter.NamePattern = aws.String(d.KeyColumnQuals["name"].GetStringValue())
	}
	if d.KeyColumnQuals["state"].GetStringValue() != "" {
		filter.States = aws.StringSlice([]string{d.KeyColumnQuals["state"].GetStringValue()})
	}

	input.Filter = filter

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
			if len(page.AssessmentRunArns) != 0 {
				assessmentRunArns = append(assessmentRunArns, page.AssessmentRunArns...)
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listInspectorAssessmentRuns", "ListAssessmentRunsPages", err)
		return nil, err
	}

	// check if there is any assessmentRunArn
	if len(assessmentRunArns) == 0 {
		return nil, nil
	}

	passedArns := 0
	arnLeft := true
	for arnLeft {
		// DescribeAssessmentRuns API can take maximum 10 arns at a time.
		var arns []*string
		if len(assessmentRunArns) > passedArns {
			if (len(assessmentRunArns) - passedArns) >= 10 {
				arns = assessmentRunArns[passedArns : passedArns+10]
				passedArns += 10
			} else {
				arns = assessmentRunArns[passedArns:]
				arnLeft = false
			}
		}

		// Build params
		input := &inspector.DescribeAssessmentRunsInput{
			AssessmentRunArns: arns,
		}

		// Get details for all available assessment runs
		result, err := svc.DescribeAssessmentRuns(input)
		if err != nil {
			return nil, err
		}

		for _, assessmentRun := range result.AssessmentRuns {
			d.StreamListItem(ctx, assessmentRun)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
