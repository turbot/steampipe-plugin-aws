package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodepipelinePipeline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codepipeline_pipeline",
		Description: "AWS Codepipeline Pipeline",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"PipelineNotFoundException"}),
			Hydrate:           getCodepipelinePipeline,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodepipelinePipelines,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the pipeline.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the pipeline.",
				Hydrate:     getCodepipelinePipeline,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Metadata.PipelineArn"),
			},
			{
				Name:        "created_at",
				Description: "The date and time the pipeline was created, in timestamp format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Created"),
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) for AWS CodePipeline to use to either perform actions with no actionRoleArn, or to use to assume roles for actions with an actionRoleArn.",
				Hydrate:     getCodepipelinePipeline,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Pipeline.RoleArn"),
			},
			{
				Name:        "stages",
				Description: "The stage in which to perform the action.",
				Hydrate:     getCodepipelinePipeline,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.Stages"),
			},
			{
				Name:        "updated_at",
				Description: "The date and time of the last update to the pipeline, in timestamp format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Updated"),
			},
			{
				Name:        "version",
				Description: "The version number of the pipeline.",
				Type:        proto.ColumnType_INT,
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
				Hydrate:     getCodepipelinePipeline,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Pipeline.ArtifactStores"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tag key and value pairs associated with this pipeline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPipelineTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listCodepipelinePipelines", "AWS_REGION", region)

	// Create Session
	svc, err := CodePipelineService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListPipelinesPages(
		&codepipeline.ListPipelinesInput{},
		func(page *codepipeline.ListPipelinesOutput, isLast bool) bool {
			for _, result := range page.Pipelines {
				d.StreamListItem(ctx, result)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCodepipelinePipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCodepipelinepipeline")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(*codepipeline.PipelineSummary).Name
	} else {
		quals := d.KeyColumnQuals
		name = quals["name"].GetStringValue()
	}

	// Get service connection
	svc, err := CodePipelineService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	pipelineArn := pipelineARN(ctx, d, h)

	// Get service connection
	svc, err := CodePipelineService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build param
	params := &codepipeline.ListTagsForResourceInput{
		ResourceArn: aws.String(pipelineArn),
	}

	tags, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func pipelineARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) string {
	plugin.Logger(ctx).Trace("pipelineARN")

	// Get region, partition, account id
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return ""
	}
	commonColumnData := c.(*awsCommonColumnData)

	switch item := h.Item.(type) {
	case *codepipeline.PipelineSummary:
		return "arn:" + commonColumnData.Partition + ":codepipeline:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":" + *item.Name
	case *codepipeline.GetPipelineOutput:
		return *item.Metadata.PipelineArn
	}

	return ""
}

func codepipelineTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*codepipeline.ListTagsForResourceOutput)

	if data.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	turbotTagsMap = map[string]string{}
	for _, i := range data.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
