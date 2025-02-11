package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector"
	"github.com/aws/aws-sdk-go-v2/service/inspector/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspectorAssessmentTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector_assessment_target",
		Description: "AWS Inspector Assessment Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getInspectorAssessmentTarget,
			Tags:       map[string]string{"service": "inspector", "action": "ListAssessmentTargets"},
		},
		List: &plugin.ListConfig{
			Hydrate: listInspectorAssessmentTargets,
			Tags:    map[string]string{"service": "inspector", "action": "ListAssessmentTargets"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_INSPECTOR_SERVICE_ID),
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

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_target.listInspectorAssessmentTargets", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

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

	input := &inspector.ListAssessmentTargetsInput{
		MaxResults: aws.Int32(maxLimit),
	}
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		if equalQuals["name"].GetStringValue() != "" {
			input.Filter = &types.AssessmentTargetFilter{AssessmentTargetNamePattern: aws.String(equalQuals["name"].GetStringValue())}
		}
	}

	paginator := inspector.NewListAssessmentTargetsPaginator(svc, input, func(o *inspector.ListAssessmentTargetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector_assessment_target.listInspectorAssessmentTargets", "api_error", err)
			return nil, err
		}

		for _, item := range output.AssessmentTargetArns {
			d.StreamListItem(ctx, &types.AssessmentTarget{
				Arn: aws.String(item),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getInspectorAssessmentTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var assessmentTargetArn string
	if h.Item != nil {
		assessmentTargetArn = *h.Item.(*types.AssessmentTarget).Arn
	} else {
		quals := d.EqualsQuals
		assessmentTargetArn = quals["arn"].GetStringValue()
	}

	// Create Session
	svc, err := InspectorClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_target.getInspectorAssessmentTarget", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &inspector.DescribeAssessmentTargetsInput{
		AssessmentTargetArns: []string{assessmentTargetArn},
	}

	// Get call
	data, err := svc.DescribeAssessmentTargets(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector_assessment_target.getInspectorAssessmentTarget", "api_error", err)
		return nil, err
	}
	if data != nil && len(data.AssessmentTargets) > 0 {
		return data.AssessmentTargets[0], nil
	}
	return nil, nil
}
