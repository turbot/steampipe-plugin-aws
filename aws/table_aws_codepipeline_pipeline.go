package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodepipelinePipeline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codepipeline_pipeline",
		Description: "AWS Codepipeline Pipeline",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"PipelineNotFoundException"}),
			},
			Hydrate: getCodepipelinePipeline,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodepipelinePipelines,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Name:        "artifact_stores",
				Description: "A mapping of artifactStore objects and their corresponding AWS Regions. There must be an artifact store for the pipeline Region and for each cross-region action in the pipeline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCodepipelinePipeline,
				Transform:   transform.FromField("Pipeline.ArtifactStores"),
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
	svc, err := CodePipelineService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &codepipeline.ListPipelinesInput{
		MaxResults: aws.Int64(1000),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
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
	err = svc.ListPipelinesPages(
		input,
		func(page *codepipeline.ListPipelinesOutput, isLast bool) bool {
			for _, result := range page.Pipelines {
				d.StreamListItem(ctx, result)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("ListPipelinesPages", "list", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCodepipelinePipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCodepipelinePipeline")

	var name string
	if h.Item != nil {
		name = *h.Item.(*codepipeline.PipelineSummary).Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create session
	svc, err := CodePipelineService(ctx, d)
	if err != nil {
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

	op, err := svc.GetPipeline(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getCodepipelinePipeline__", "ERROR", err)
		return nil, err
	}

	if op != nil {
		return op, nil
	}

	return nil, nil
}

func getPipelineTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPipelineTags")

	pipelineArn := pipelineARN(ctx, d, h)

	// Create session
	svc, err := CodePipelineService(ctx, d)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &codepipeline.ListTagsForResourceInput{
		ResourceArn: aws.String(pipelineArn),
	}

	tags := []*codepipeline.Tag{}

	err = svc.ListTagsForResourcePages(
		params,
		func(page *codepipeline.ListTagsForResourceOutput, isLast bool) bool {
			tags = append(tags, page.Tags...)
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("getPipelineTags", "ListTagsForResourcePages_error", err)
		return nil, err
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func pipelineARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) string {
	plugin.Logger(ctx).Trace("pipelineARN")
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get region, partition, account id
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return ""
	}
	commonColumnData := c.(*awsCommonColumnData)

	switch item := h.Item.(type) {
	case *codepipeline.PipelineSummary:
		return "arn:" + commonColumnData.Partition + ":codepipeline:" + region + ":" + commonColumnData.AccountId + ":" + *item.Name
	case *codepipeline.GetPipelineOutput:
		return *item.Metadata.PipelineArn
	}

	return ""
}

func codepipelineTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]*codepipeline.Tag)

	if tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
