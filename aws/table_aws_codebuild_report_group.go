package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/aws/aws-sdk-go-v2/service/codebuild/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCodeBuildReportGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codebuild_report_group",
		Description: "AWS CodeBuild Report Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
			Hydrate: getCodeBuildReportGroup,
			Tags:    map[string]string{"service": "codebuild", "action": "BatchGetReportGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeBuildReportGroups,
			Tags:    map[string]string{"service": "codebuild", "action": "ListReportGroups"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCodeBuildReportGroup,
				Tags: map[string]string{"service": "codebuild", "action": "BatchGetReportGroups"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CODEBUILD_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the report group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildReportGroup,
			},
			{
				Name:        "arn",
				Description: "The ARN of the report group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the report group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildReportGroup,
			},
			{
				Name:        "created",
				Description: "The date and time this report group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeBuildReportGroup,
			},
			{
				Name:        "last_modified",
				Description: "The date and time this report group was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCodeBuildReportGroup,
			},
			{
				Name:        "export_config",
				Description: "Information about the destination where the raw data of this ReportGroup is exported.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeBuildReportGroup,
			},
			{
				Name:        "tags_src",
				Description: "A list of tag key and value pairs associated with this report group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeBuildReportGroup,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodeBuildReportGroup,
				Transform:   transform.From(codeBuildReportGroupTurbotTags),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeBuildReportGroup,
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

func listCodeBuildReportGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_report_group.listCodeBuildReportGroups", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxResults := int32(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxResults {
			maxResults = int32(limit)
		}
	}

	input := &codebuild.ListReportGroupsInput{
		MaxResults: &maxResults,
	}

	paginator := codebuild.NewListReportGroupsPaginator(svc, input, func(o *codebuild.ListReportGroupsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codebuild_report_group.listCodeBuildReportGroups", "api_error", err)
			return nil, err
		}

		for _, reportGroupArn := range output.ReportGroups {
			item := &types.ReportGroup{
				Arn: aws.String(reportGroupArn),
			}
			d.StreamListItem(ctx, *item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeBuildReportGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.ReportGroup).Arn
	} else {
		quals := d.EqualsQuals
		arn = quals["arn"].GetStringValue()
	}

	// Empty input check
	if arn == "" {
		return nil, nil
	}

	// Create service
	svc, err := CodeBuildClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_report_group.getCodeBuildReportGroup", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &codebuild.BatchGetReportGroupsInput{
		ReportGroupArns: []string{arn},
	}

	// Get call
	op, err := svc.BatchGetReportGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codebuild_report_group.getCodeBuildReportGroup", "api_error", err)
		return nil, err
	}

	if len(op.ReportGroups) > 0 {
		return op.ReportGroups[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func codeBuildReportGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.ReportGroup)

	if data.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if data.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
