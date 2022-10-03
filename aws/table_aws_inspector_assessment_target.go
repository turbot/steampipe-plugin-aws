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

func tableAwsInspectorAssessmentTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_assessment_target",
		Description: "AWS Inspector Assessment Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getInspectorAssessmentTarget,
		},
		List: &plugin.ListConfig{
			Hydrate: listInspectorAssessmentTargets,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Amazon Inspector assessment target.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentTarget,
			},
			{
				Name:        "arn",
				Description: "The ARN that specifies the Amazon Inspector assessment target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_group_arn",
				Description: "The ARN that specifies the resource group that is associated with the assessment target.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentTarget,
			},
			{
				Name:        "created_at",
				Description: "The time at which the assessment target is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getInspectorAssessmentTarget,
			},
			{
				Name:        "updated_at",
				Description: "The time at which UpdateAssessmentTarget is called.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getInspectorAssessmentTarget,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentTarget,
				Transform:   transform.FromField("Name"),
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

func listInspectorAssessmentTargets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &inspector.ListAssessmentTargetsInput{
		MaxResults: aws.Int64(500),
	}
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		if equalQuals["name"].GetStringValue() != "" {
			input.Filter = &inspector.AssessmentTargetFilter{AssessmentTargetNamePattern: aws.String(equalQuals["name"].GetStringValue())}
		}
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
	err = svc.ListAssessmentTargetsPages(
		input,
		func(page *inspector.ListAssessmentTargetsOutput, isLast bool) bool {
			for _, assessmentTarget := range page.AssessmentTargetArns {
				d.StreamListItem(ctx, &inspector.AssessmentTarget{
					Arn: assessmentTarget,
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

func getInspectorAssessmentTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var assessmentTargetArn string
	if h.Item != nil {
		assessmentTargetArn = *h.Item.(*inspector.AssessmentTarget).Arn
	} else {
		quals := d.KeyColumnQuals
		assessmentTargetArn = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &inspector.DescribeAssessmentTargetsInput{
		AssessmentTargetArns: []*string{aws.String(assessmentTargetArn)},
	}

	// Get call
	data, err := svc.DescribeAssessmentTargets(params)
	if err != nil {
		plugin.Logger(ctx).Debug("describeAssessmentTarget__", "ERROR", err)
		return nil, err
	}
	if data.AssessmentTargets != nil && len(data.AssessmentTargets) > 0 {
		return data.AssessmentTargets[0], nil
	}
	return nil, nil
}
