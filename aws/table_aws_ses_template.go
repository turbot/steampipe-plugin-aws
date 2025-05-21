package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSESTemplate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ses_template",
		Description: "AWS SES Template",
		List: &plugin.ListConfig{
			Hydrate: listSESTemplates,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
			Tags: map[string]string{"service": "ses", "action": "ListTemplates"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSESTemplateDetails,
				Tags: map[string]string{"service": "ses", "action": "GetTemplate"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EMAIL_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the template.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "TemplateName"),
			},
			{
				Name:        "subject_part",
				Description: "The subject line of the email.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSESTemplateDetails,
			},
			{
				Name:        "text_part",
				Description: "The email body that will be visible to recipients whose email clients do not display HTML.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSESTemplateDetails,
			},
			{
				Name:        "html_part",
				Description: "The HTML body of the email.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSESTemplateDetails,
			},
			{
				Name:        "created_timestamp",
				Description: "The time and date the template was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "TemplateName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSESTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_template.listSESTemplates", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(1000)
	// Limiting the results
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	input := &ses.ListTemplatesInput{
		MaxItems: &maxItems,
	}

	// List call
	output, err := svc.ListTemplates(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_template.listSESTemplates", "api_error", err)
		return nil, err
	}

	for _, template := range output.TemplatesMetadata {
		d.StreamListItem(ctx, template)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSESTemplateDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	templateName := h.Item.(types.TemplateMetadata).Name

	if templateName == nil {
		return nil, nil
	}
	if d.EqualsQualString("name") != "" && *templateName != d.EqualsQualString("name") {
		return nil, nil
	}

	svc, err := SESClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_template.getSESTemplateDetails", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	params := &ses.GetTemplateInput{
		TemplateName: templateName,
	}

	op, err := svc.GetTemplate(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ses_template.getSESTemplateDetails", "api_error", err)
		return nil, err
	}

	return op.Template, nil
}
