package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
	"github.com/aws/aws-sdk-go-v2/service/codeconnections/types"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v6/query_cache"
)

func tableAwsCodeConnectionsRepositoryLink(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeconnections_repository_link",
		Description: "AWS CodeConnections Repository Link",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("repository_link_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getCodeConnectionsRepositoryLink,
			Tags:    map[string]string{"service": "codeconnections", "action": "GetRepositoryLink"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeConnectionsRepositoryLinks,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "repository_link_id", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Tags: map[string]string{"service": "codeconnections", "action": "ListRepositoryLinks"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listCodeConnectionsRepositoryLinkTags,
				Tags: map[string]string{"service": "codeconnections", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CODESTAR_CONNECTIONS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{Name: "repository_link_id", Description: "The ID of the repository link.", Type: proto.ColumnType_STRING},
			{Name: "arn", Description: "The Amazon Resource Name (ARN) of the repository link.", Type: proto.ColumnType_STRING, Transform: transform.FromField("RepositoryLinkArn")},
			{Name: "repository_name", Description: "The name of the repository associated with the repository link.", Type: proto.ColumnType_STRING},
			{Name: "connection_arn", Description: "The ARN of the connection associated with the repository link.", Type: proto.ColumnType_STRING},
			{Name: "encryption_key_arn", Description: "The ARN of the encryption key for the repository.", Type: proto.ColumnType_STRING},
			{Name: "owner_id", Description: "The owner ID for the repository.", Type: proto.ColumnType_STRING},
			{Name: "provider_type", Description: "The provider type associated with the repository link.", Type: proto.ColumnType_STRING},
			{Name: "tags_src", Description: "A list of tags associated with the repository link.", Type: proto.ColumnType_JSON, Hydrate: listCodeConnectionsRepositoryLinkTags, Transform: transform.FromField("Tags")},
			{Name: "title", Description: resourceInterfaceDescription("title"), Type: proto.ColumnType_STRING, Transform: transform.FromField("RepositoryName")},
			{Name: "tags", Description: resourceInterfaceDescription("tags"), Type: proto.ColumnType_JSON, Hydrate: listCodeConnectionsRepositoryLinkTags, Transform: transform.From(codeConnectionsTurbotTags)},
			{Name: "akas", Description: resourceInterfaceDescription("akas"), Type: proto.ColumnType_JSON, Transform: transform.FromField("RepositoryLinkArn").Transform(transform.EnsureStringArray)},
		}),
	}
}

func listCodeConnectionsRepositoryLinks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_repository_link.listCodeConnectionsRepositoryLinks", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil && int32(*d.QueryContext.Limit) < maxLimit {
		maxLimit = int32(*d.QueryContext.Limit)
	}
	input := &codeconnections.ListRepositoryLinksInput{MaxResults: maxLimit}
	paginator := codeconnections.NewListRepositoryLinksPaginator(svc, input, func(o *codeconnections.ListRepositoryLinksPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	wantedID := d.EqualsQualString("repository_link_id")
	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codeconnections_repository_link.listCodeConnectionsRepositoryLinks", "api_error", err)
			return nil, err
		}
		for _, link := range output.RepositoryLinks {
			if wantedID != "" && aws.ToString(link.RepositoryLinkId) != wantedID {
				continue
			}
			d.StreamListItem(ctx, link)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

func getCodeConnectionsRepositoryLink(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("repository_link_id")
	if id == "" {
		return nil, nil
	}
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_repository_link.getCodeConnectionsRepositoryLink", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	output, err := svc.GetRepositoryLink(ctx, &codeconnections.GetRepositoryLinkInput{RepositoryLinkId: aws.String(id)})
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_repository_link.getCodeConnectionsRepositoryLink", "api_error", err)
		return nil, err
	}
	return output.RepositoryLinkInfo, nil
}

func listCodeConnectionsRepositoryLinkTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn *string
	switch item := h.Item.(type) {
	case types.RepositoryLinkInfo:
		arn = item.RepositoryLinkArn
	case *types.RepositoryLinkInfo:
		arn = item.RepositoryLinkArn
	}
	return listCodeConnectionsTags(ctx, d, arn)
}

func listCodeConnectionsTags(ctx context.Context, d *plugin.QueryData, arn *string) (interface{}, error) {
	if arn == nil || *arn == "" {
		return nil, nil
	}
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections.listCodeConnectionsTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	output, err := svc.ListTagsForResource(ctx, &codeconnections.ListTagsForResourceInput{ResourceArn: arn})
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections.listCodeConnectionsTags", "api_error", err)
		return nil, err
	}
	return output, nil
}

func codeConnectionsTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	output := d.HydrateItem.(*codeconnections.ListTagsForResourceOutput)
	tags := make(map[string]string, len(output.Tags))
	for _, tag := range output.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return tags, nil
}
