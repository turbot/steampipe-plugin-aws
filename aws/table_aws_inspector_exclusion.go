package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type ExclusionInfo = struct {
	inspector.Exclusion
	AssessmentRunArn string
}

//// TABLE DEFINITION

func tableAwsInspectorExclusion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_exclusion",
		Description: "AWS Inspector Exclusion",
		List: &plugin.ListConfig{
			ParentHydrate: listInspectorAssessmentRuns,
			Hydrate:       listInspectorExclusions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "assessment_run_arn", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The ARN that specifies the exclusion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assessment_run_arn",
				Description: "The ARN that specifies the assessment run, the exclusion belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attributes",
				Description: "The system-defined attributes for the exclusion.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "description",
				Description: "The description of the exclusion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recommendation",
				Description: "The recommendation for the exclusion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scopes",
				Description: "The AWS resources for which the exclusion pertains.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
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

func listInspectorExclusions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listInspectorExclusions")

	// Create Session
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Exclusion is a sub resource of an assessment run, we need the assessment run ARN to list these.
	runArn := *h.Item.(*inspector.AssessmentRun).Arn
	equalQuals := d.KeyColumnQuals

	// Minimize the API call with the given assessment run ARN
	if equalQuals["assessment_run_arn"] != nil {
		if equalQuals["assessment_run_arn"].GetStringValue() != "" {
			if equalQuals["assessment_run_arn"].GetStringValue() != "" && equalQuals["assessment_run_arn"].GetStringValue() != runArn {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["assessment_run_arn"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["assessment_run_arn"].GetListValue())), runArn) {
				return nil, nil
			}
		}
	}

	// List all available exclusions
	var exclusions []*string
	err = svc.ListExclusionsPages(
		&inspector.ListExclusionsInput{
			AssessmentRunArn: &runArn,
			MaxResults:       aws.Int64(500),
		},
		func(page *inspector.ListExclusionsOutput, isLast bool) bool {
			exclusions = append(exclusions, page.ExclusionArns...)
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	passedExclusions := 0
	exclusionsLeft := true
	for exclusionsLeft {
		// DescribeExclusions API can take maximum 100 number of exclusions ARNs at a time.
		var arns []*string
		if len(exclusions) > passedExclusions {
			if (len(exclusions) - passedExclusions) >= 100 {
				arns = exclusions[passedExclusions : passedExclusions+100]
				passedExclusions += 100
			} else {
				arns = exclusions[passedExclusions:]
				exclusionsLeft = false
			}
		}

		// Build params
		params := &inspector.DescribeExclusionsInput{
			ExclusionArns: arns,
		}

		// Get details for all available exclusions
		result, err := svc.DescribeExclusions(params)
		if err != nil {
			return nil, err
		}
		for _, exclusion := range result.Exclusions {
			d.StreamListItem(ctx, ExclusionInfo{*exclusion, runArn})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
