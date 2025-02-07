package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"

	codepipelineEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodepipelinePipeline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codepipeline_pipeline",
		Description: "AWS Codepipeline Pipeline",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"PipelineNotFoundException"}),
			},
			Hydrate: getCodepipelinePipeline,
			Tags:    map[string]string{"service": "codepipeline", "action": "GetPipeline"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodepipelinePipelines,
			Tags:    map[string]string{"service": "codepipeline", "action": "ListPipelines"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getPipelineTags,
				Tags: map[string]string{"service": "codepipeline", "action": "ListTagsForResource"},
			},
			{
				Func: getCodepipelinePipeline,
				Tags: map[string]string{"service": "codepipeline", "action": "GetPipeline"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(codepipelineEndpoint.AWS_CODEPIPELINE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the pipeline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "Pipeline.Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the pipeline.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodepipelinePipeline,
				Transform:   transform.FromField("Metadata.PipelineArn"),
			},
			{
				Name:        "created_at",
				Description: "The date and time the pipeline was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Created", "Metadata.Created"),
			},
			{
				Name:        "execution_mode",
				Description: "The method that the pipeline will use to handle multiple executions.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExecutionMode", "Pipeline.ExecutionMode"),
			},
			{
				Name:        "pipeline_type",
				Description: "The pipeline type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PipelineType", "Pipeline.PipelineType"),
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) for AWS CodePipeline to use to either perform actions with no actionRoleArn, or to use to assume roles for actions with an actionRoleArn.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodepipelinePipeline,
				Transform:   transform.FromField("Pipeline.RoleArn"),
			},
			{
				Name:        "updated_at",
				Description: "The date and time of the last update to the pipeline.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Updated", "Metadata.Updated"),
			},
			{
				Name:        "version",
				Description: "The version number of the pipeline.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Version", "Pipeline.Version"),
			},
			{
				Name:        "encryption_key",
				Description: "The encryption key used to encrypt the data in the artifact store, such as an AWS Key Management Service (AWS KMS) key. If this is undefined, the default key for Amazon S3 is used.",
				Hydrate:     getCodepipelinePipeline,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.ArtifactStore.EncryptionKey"),
			},
			{
				Name:        "triggers",
				Description: "The trigger configuration specifying a type of event, such as Git tags, that starts the pipeline.",
				Hydrate:     getCodepipelinePipeline,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Triggers"),
			},
			{
				Name:        "variables",
				Description: "A list that defines the pipeline variables for a pipeline resource.",
				Hydrate:     getCodepipelinePipeline,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Variables"),
			},
			{
				Name:        "artifact_stores",
				Description: "A mapping of artifactStore objects and their corresponding AWS Regions. There must be an artifact store for the pipeline Region and for each cross-region action in the pipeline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodepipelinePipeline,
				Transform:   transform.FromField("Pipeline.ArtifactStore"),
			},
			{
				Name:        "stages",
				Description: "The stage in which to perform the action.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodepipelinePipeline,
				Transform:   transform.FromField("Pipeline.Stages"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tag key and value pairs associated with this pipeline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPipelineTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name", "Pipeline.Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPipelineTags,
				Transform:   transform.From(codepipelineTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodepipelinePipeline,
				Transform:   transform.FromField("Metadata.PipelineArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCodepipelinePipelines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := CodePipelineClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codepipeline_pipeline.listCodepipelinePipelines", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(1000)
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

	input := &codepipeline.ListPipelinesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := codepipeline.NewListPipelinesPaginator(svc, input, func(o *codepipeline.ListPipelinesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codepipeline_pipeline.listCodepipelinePipelines", "api_error", err)
			return nil, err
		}
		for _, items := range output.Pipelines {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCodepipelinePipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string
	if h.Item != nil {
		name = *h.Item.(types.PipelineSummary).Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Create session
	svc, err := CodePipelineClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codepipeline_pipeline.getCodepipelinePipeline", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &codepipeline.GetPipelineInput{
		Name: aws.String(name),
	}

	op, err := svc.GetPipeline(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codepipeline_pipeline.getCodepipelinePipeline", "api_error", err)
		return nil, err
	}

	if op != nil {
		return op, nil
	}

	return nil, nil
}

func getPipelineTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	pipelineArn := pipelineARN(ctx, d, h)

	// Create session
	svc, err := CodePipelineClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codepipeline_pipeline.getPipelineTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxLimit := aws.Int32(100)
	// Build params
	params := &codepipeline.ListTagsForResourceInput{
		ResourceArn: aws.String(pipelineArn),
		MaxResults:  aws.Int32(*maxLimit),
	}

	var tags []types.Tag

	paginator := codepipeline.NewListTagsForResourcePaginator(svc, params, func(o *codepipeline.ListTagsForResourcePaginatorOptions) {
		o.Limit = *maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codepipeline_pipeline.getPipelineTags", "api_error", err)
			return nil, err
		}

		tags = append(tags, output.Tags...)
	}
	if tags == nil {
		return make([]types.Tag, 0), nil
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func pipelineARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) string {
	region := d.EqualsQualString(matrixKeyRegion)

	// Get region, partition, account id

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return ""
	}
	commonColumnData := c.(*awsCommonColumnData)

	switch item := h.Item.(type) {
	case types.PipelineSummary:
		return "arn:" + commonColumnData.Partition + ":codepipeline:" + region + ":" + commonColumnData.AccountId + ":" + *item.Name
	case *codepipeline.GetPipelineOutput:
		return *item.Metadata.PipelineArn
	}

	return ""
}

func codepipelineTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)

	if len(tags) <= 0 {
		return map[string]string{}, nil
	}

	// Mapping the resource tags inside turbotTags
	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
