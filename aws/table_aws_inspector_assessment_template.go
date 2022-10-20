package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector"
	"github.com/aws/aws-sdk-go-v2/service/inspector/types"
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
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{}),
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

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_template.listInspectorAssessmentTemplates", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &inspector.ListAssessmentTemplatesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.Filter = &types.AssessmentTemplateFilter{
			NamePattern: aws.String(equalQuals["name"].GetStringValue()),
		}
	}

	if equalQuals["assessment_target_arn"] != nil {
		if equalQuals["assessment_target_arn"].GetStringValue() != "" {
			input.AssessmentTargetArns = []string{equalQuals["assessment_target_arn"].GetStringValue()}
		}
	}

	paginator := inspector.NewListAssessmentTemplatesPaginator(svc, input, func(o *inspector.ListAssessmentTemplatesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector_assessment_template.listInspectorAssessmentTemplates", "api_error", err)
			return nil, err
		}

		for _, items := range output.AssessmentTemplateArns {
			d.StreamListItem(ctx, &types.AssessmentTemplate{
				Arn: &items,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getInspectorAssessmentTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var assessmentTemplateArn string
	if h.Item != nil {
		assessmentTemplateArn = *h.Item.(*types.AssessmentTemplate).Arn
	} else {
		quals := d.KeyColumnQuals
		assessmentTemplateArn = quals["arn"].GetStringValue()
	}

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_template.listInspectorAssessmentTemplates", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &inspector.DescribeAssessmentTemplatesInput{
		AssessmentTemplateArns: []string{assessmentTemplateArn},
	}

	// Get call
	data, err := svc.DescribeAssessmentTemplates(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_template.listInspectorAssessmentTemplates", "api_error", err)
		return nil, err
	}
	if data.AssessmentTemplates != nil && len(data.AssessmentTemplates) > 0 {
		return &data.AssessmentTemplates[0], nil
	}
	return nil, nil
}

// API call for fetching tag list
func getAwsInspectorAssessmentTemplateTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	assessmentTemplateArn := *h.Item.(*types.AssessmentTemplate).Arn

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_template.getAwsInspectorAssessmentTemplateTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &inspector.ListTagsForResourceInput{
		ResourceArn: &assessmentTemplateArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_template.getAwsInspectorAssessmentTemplateTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

// API call for fetching event subscriptions
func listAwsInspectorAssessmentEventSubscriptions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	assessmentTemplateArn := *h.Item.(*types.AssessmentTemplate).Arn

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_template.listAwsInspectorAssessmentEventSubscriptions", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &inspector.ListEventSubscriptionsInput{
		ResourceArn: &assessmentTemplateArn,
	}

	var associatedEventSubscriptions []types.Subscription

	paginator := inspector.NewListEventSubscriptionsPaginator(svc, params, func(o *inspector.ListEventSubscriptionsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector_assessment_template.listAwsInspectorAssessmentEventSubscriptions", "api_error", err)
			return nil, err
		}

		associatedEventSubscriptions = append(associatedEventSubscriptions, output.Subscriptions...)
	}

	return associatedEventSubscriptions, err
}

//// TRANSFORM FUNCTIONS

func inspectorTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

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
