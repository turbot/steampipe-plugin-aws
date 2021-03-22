package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/inspector"
	pb "github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspectorAssessmentTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_assessment_target",
		Description: "AWS Inspector Assessment Target",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{}),
			Hydrate:           describeAssessmentTarget,
		},
		List: &plugin.ListConfig{
			Hydrate: listAssessmentTargets,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Amazon Inspector assessment target.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeAssessmentTarget,
			},
			{
				Name:        "arn",
				Description: "The ARN that specifies the Amazon Inspector assessment target.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "resource_group_arn",
				Description: "The ARN that specifies the resource group that is associated with the assessment target.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeAssessmentTarget,
			},
			{
				Name:        "created_at",
				Description: "The time at which the assessment target is created.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeAssessmentTarget,
			},
			{
				Name:        "updated_at",
				Description: "The time at which UpdateAssessmentTarget is called.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeAssessmentTarget,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeAssessmentTarget,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        pb.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAssessmentTargets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAssessmentTargets", "AWS_REGION", region)

	// Create session
	svc, err := InspectorService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListAssessmentTargetsPages(
		&inspector.ListAssessmentTargetsInput{},
		func(page *inspector.ListAssessmentTargetsOutput, isLast bool) bool {
			for _, assessmentTarget := range page.AssessmentTargetArns {
				d.StreamListItem(ctx, &inspector.AssessmentTarget{
					Arn: assessmentTarget,
				})
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeAssessmentTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("describeAssessmentTarget")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var assessmentTargetArn string
	if h.Item != nil {
		assessmentTargetArn = *h.Item.(*inspector.AssessmentTarget).Arn
	} else {
		quals := d.KeyColumnQuals
		assessmentTargetArn = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := InspectorService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &inspector.DescribeAssessmentTargetsInput{
		AssessmentTargetArns: []*string{aws.String(assessmentTargetArn)},
	}

	// Get call
	data, err := svc.DescribeAssessmentTargets(params)
	if err != nil {
		logger.Debug("describeAssessmentTarget__", "ERROR", err)
		return nil, err
	}
	if len(data.AssessmentTargets) > 0 {
		return data.AssessmentTargets[0], nil
	}
	return nil, nil
}