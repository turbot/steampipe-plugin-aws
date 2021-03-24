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

func tableAwsInspectorAssessmentTemplate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_assessment_template",
		Description: "AWS Inspector Assessment Template",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{}),
			Hydrate:           describeAssessmentTemplate,
		},
		List: &plugin.ListConfig{
			Hydrate: listAssessmentTemplates,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the assessment template.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "arn",
				Description: "The ARN of the assessment template.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "assessment_run_count",
				Description: "The number of existing assessment runs associated with this assessment template.",
				Type:        pb.ColumnType_INT,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "assessment_target_arn",
				Description: "The ARN of the assessment target that corresponds to this assessment template.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "created_at",
				Description: "The time at which the assessment template is created.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "duration_in_seconds",
				Description: "The duration in seconds specified for this assessment template.",
				Type:        pb.ColumnType_INT,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "last_assessment_run_arn",
				Description: "The Amazon Resource Name (ARN) of the most recent assessment run associated with this assessment template.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "rules_package_arns",
				Description: "The rules packages that are specified for this assessment template.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "user_attributes_for_findings",
				Description: "The user-defined attributes that are assigned to every generated finding from the assessment run that uses this assessment template.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeAssessmentTemplate,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the Assessment Template.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     getAwsInspectorAssessmentTemplateTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeAssessmentTemplate,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        pb.ColumnType_JSON,
				Hydrate:     getAwsInspectorAssessmentTemplateTags,
				Transform:   transform.FromField("Tags").Transform(inspectorTagListToTurbotTags),
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

func listAssessmentTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAssessmentTemplates", "AWS_REGION", region)

	// Create session
	svc, err := InspectorService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListAssessmentTemplatesPages(
		&inspector.ListAssessmentTemplatesInput{},
		func(page *inspector.ListAssessmentTemplatesOutput, isLast bool) bool {
			for _, assessmentTemplate := range page.AssessmentTemplateArns {
				d.StreamListItem(ctx, &inspector.AssessmentTemplate{
					Arn: assessmentTemplate,
				})
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeAssessmentTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("describeAssessmentTemplate")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var assessmentTemplateArn string
	if h.Item != nil {
		assessmentTemplateArn = *h.Item.(*inspector.AssessmentTemplate).Arn
	} else {
		quals := d.KeyColumnQuals
		assessmentTemplateArn = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := InspectorService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &inspector.DescribeAssessmentTemplatesInput{
		AssessmentTemplateArns: []*string{aws.String(assessmentTemplateArn)},
	}

	// Get call
	data, err := svc.DescribeAssessmentTemplates(params)
	if err != nil {
		logger.Debug("describeAssessmentTemplate__", "ERROR", err)
		return nil, err
	}
	if len(data.AssessmentTemplates) > 0 {
		return data.AssessmentTemplates[0], nil
	}
	return nil, nil
}

// API call for fetching tag list
func getAwsInspectorAssessmentTemplateTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsInspectorAssessmentTemplateTags")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	assessmentTemplateArn := *h.Item.(*inspector.AssessmentTemplate).Arn

	// Create Session
	svc, err := InspectorService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &inspector.ListTagsForResourceInput{
		ResourceArn: &assessmentTemplateArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getAwsInspectorAssessmentTemplateTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func inspectorTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("inspectorTagListToTurbotTags")
	tagList := d.Value.([]*inspector.Tag)

	if tagList == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
