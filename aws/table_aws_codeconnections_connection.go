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

func tableAwsCodeConnectionsConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeconnections_connection",
		Description: "AWS CodeConnections Connection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getCodeConnectionsConnection,
			Tags:    map[string]string{"service": "codeconnections", "action": "GetConnection"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeConnectionsConnections,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "provider_type", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "host_arn", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Tags: map[string]string{"service": "codeconnections", "action": "ListConnections"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listCodeConnectionsConnectionTags,
				Tags: map[string]string{"service": "codeconnections", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CODESTAR_CONNECTIONS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{Name: "arn", Description: "The Amazon Resource Name (ARN) of the connection.", Type: proto.ColumnType_STRING, Transform: transform.FromField("ConnectionArn")},
			{Name: "name", Description: "The name of the connection.", Type: proto.ColumnType_STRING, Transform: transform.FromField("ConnectionName")},
			{Name: "provider_type", Description: "The external provider where the third-party code repository is configured.", Type: proto.ColumnType_STRING},
			{Name: "owner_account_id", Description: "The identifier of the external provider account that owns the repository.", Type: proto.ColumnType_STRING},
			{Name: "connection_status", Description: "The current status of the connection.", Type: proto.ColumnType_STRING},
			{Name: "host_arn", Description: "The ARN of the host associated with the connection.", Type: proto.ColumnType_STRING},
			{Name: "tags_src", Description: "A list of tags associated with the connection.", Type: proto.ColumnType_JSON, Hydrate: listCodeConnectionsConnectionTags, Transform: transform.FromField("Tags")},
			{Name: "title", Description: resourceInterfaceDescription("title"), Type: proto.ColumnType_STRING, Transform: transform.FromField("ConnectionName")},
			{Name: "tags", Description: resourceInterfaceDescription("tags"), Type: proto.ColumnType_JSON, Hydrate: listCodeConnectionsConnectionTags, Transform: transform.From(codeConnectionsTurbotTags)},
			{Name: "akas", Description: resourceInterfaceDescription("akas"), Type: proto.ColumnType_JSON, Transform: transform.FromField("ConnectionArn").Transform(transform.EnsureStringArray)},
		}),
	}
}

func listCodeConnectionsConnections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_connection.listCodeConnectionsConnections", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	maxLimit := int32(100)
	if d.QueryContext.Limit != nil && int32(*d.QueryContext.Limit) < maxLimit {
		maxLimit = int32(*d.QueryContext.Limit)
	}
	input := &codeconnections.ListConnectionsInput{MaxResults: maxLimit}
	if providerType := d.EqualsQualString("provider_type"); providerType != "" {
		input.ProviderTypeFilter = types.ProviderType(providerType)
	}
	if hostArn := d.EqualsQualString("host_arn"); hostArn != "" {
		input.HostArnFilter = aws.String(hostArn)
	}

	paginator := codeconnections.NewListConnectionsPaginator(svc, input, func(o *codeconnections.ListConnectionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codeconnections_connection.listCodeConnectionsConnections", "api_error", err)
			return nil, err
		}
		for _, connection := range output.Connections {
			d.StreamListItem(ctx, connection)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

func getCodeConnectionsConnection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQualString("arn")
	if arn == "" {
		return nil, nil
	}
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_connection.getCodeConnectionsConnection", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	output, err := svc.GetConnection(ctx, &codeconnections.GetConnectionInput{ConnectionArn: aws.String(arn)})
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_connection.getCodeConnectionsConnection", "api_error", err)
		return nil, err
	}
	return output.Connection, nil
}

func listCodeConnectionsConnectionTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn *string
	switch item := h.Item.(type) {
	case types.Connection:
		arn = item.ConnectionArn
	case *types.Connection:
		arn = item.ConnectionArn
	}
	return listCodeConnectionsTags(ctx, d, arn)
}
