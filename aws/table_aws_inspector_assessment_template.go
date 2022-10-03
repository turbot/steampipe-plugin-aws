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

func tableAwsInspectorAssessmentTemplate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_assessment_template",
		Description: "AWS Inspector Assessment Template",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{}),
			},
			Hydrate: getInspectorAssessmentTemplate,
		},
		List: &plugin.ListConfig{
			Hydrate: listInspectorAssessmentTemplates,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "assessment_target_arn", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the assessment template.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "arn",
				Description: "The ARN of the assessment template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assessment_run_count",
				Description: "The number of existing assessment runs associated with this assessment template.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "assessment_target_arn",
				Description: "The ARN of the assessment target that corresponds to this assessment template.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "created_at",
				Description: "The time at which the assessment template is created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "duration_in_seconds",
				Description: "The duration in seconds specified for this assessment template.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "last_assessment_run_arn",
				Description: "The Amazon Resource Name (ARN) of the most recent assessment run associated with this assessment template.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "rules_package_arns",
				Description: "The rules packages that are specified for this assessment template.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "user_attributes_for_findings",
				Description: "The user-defined attributes that are assigned to every generated finding from the assessment run that uses this assessment template.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInspectorAssessmentTemplate,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the Assessment Template.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsInspectorAssessmentTemplateTags,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "event_subscriptions",
				Description: "A list of event subscriptions associated with the Assessment Template.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsInspectorAssessmentEventSubscriptions,
				Transform:   transform.FromValue(),
			},
			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getInspectorAssessmentTemplate,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsInspectorAssessmentTemplateTags,
				Transform:   transform.FromField("Tags").Transform(inspectorTagListToTurbotTags),
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

func listInspectorAssessmentTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &inspector.ListAssessmentTemplatesInput{
		MaxResults: aws.Int64(500),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.Filter = &inspector.AssessmentTemplateFilter{
			NamePattern: aws.String(equalQuals["name"].GetStringValue()),
		}
	}

	if equalQuals["assessment_target_arn"] != nil {
		if equalQuals["assessment_target_arn"].GetStringValue() != "" {
			input.AssessmentTargetArns = []*string{aws.String(equalQuals["assessment_target_arn"].GetStringValue())}
		} else {
			input.AssessmentTargetArns = getListValues(equalQuals["assessment_target_arn"].GetListValue())
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
	err = svc.ListAssessmentTemplatesPages(
		input,
		func(page *inspector.ListAssessmentTemplatesOutput, isLast bool) bool {
			for _, assessmentTemplate := range page.AssessmentTemplateArns {
				d.StreamListItem(ctx, &inspector.AssessmentTemplate{
					Arn: assessmentTemplate,
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

func getInspectorAssessmentTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var assessmentTemplateArn string
	if h.Item != nil {
		assessmentTemplateArn = *h.Item.(*inspector.AssessmentTemplate).Arn
	} else {
		quals := d.KeyColumnQuals
		assessmentTemplateArn = quals["arn"].GetStringValue()
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
	params := &inspector.DescribeAssessmentTemplatesInput{
		AssessmentTemplateArns: []*string{aws.String(assessmentTemplateArn)},
	}

	// Get call
	data, err := svc.DescribeAssessmentTemplates(params)
	if err != nil {
		logger.Debug("describeAssessmentTemplate__", "ERROR", err)
		return nil, err
	}
	if data.AssessmentTemplates != nil && len(data.AssessmentTemplates) > 0 {
		return data.AssessmentTemplates[0], nil
	}
	return nil, nil
}

// API call for fetching tag list
func getAwsInspectorAssessmentTemplateTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsInspectorAssessmentTemplateTags")

	assessmentTemplateArn := *h.Item.(*inspector.AssessmentTemplate).Arn

	// Create Session
	svc, err := InspectorService(ctx, d)
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

// API call for fetching event subscriptions
func listAwsInspectorAssessmentEventSubscriptions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("ListAwsInspectorAssessmentTemplateEventSubscriptions")

	assessmentTemplateArn := *h.Item.(*inspector.AssessmentTemplate).Arn

	// Create Session
	svc, err := InspectorService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &inspector.ListEventSubscriptionsInput{
		ResourceArn: &assessmentTemplateArn,
	}

	var associatedEventSubscriptions []*inspector.Subscription

	err = svc.ListEventSubscriptionsPages(
		params,
		func(page *inspector.ListEventSubscriptionsOutput, lastPage bool) bool {
			associatedEventSubscriptions = append(associatedEventSubscriptions, page.Subscriptions...)
			return !lastPage
		},
	)

	return associatedEventSubscriptions, err
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
