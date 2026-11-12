package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
	"github.com/aws/aws-sdk-go-v2/service/codeconnections/types"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableAwsCodeConnectionsHost(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeconnections_host",
		Description: "AWS CodeConnections Host",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getCodeConnectionsHost,
			Tags:    map[string]string{"service": "codeconnections", "action": "GetHost"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeConnectionsHosts,
			Tags:    map[string]string{"service": "codeconnections", "action": "ListHosts"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listCodeConnectionsHostTags,
				Tags: map[string]string{"service": "codeconnections", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CODESTAR_CONNECTIONS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{Name: "arn", Description: "The Amazon Resource Name (ARN) of the host.", Type: proto.ColumnType_STRING, Transform: transform.FromField("HostArn")},
			{Name: "name", Description: "The name of the host.", Type: proto.ColumnType_STRING},
			{Name: "provider_type", Description: "The installed provider associated with the host.", Type: proto.ColumnType_STRING},
			{Name: "provider_endpoint", Description: "The endpoint of the infrastructure where the provider is installed.", Type: proto.ColumnType_STRING},
			{Name: "status", Description: "The current status of the host.", Type: proto.ColumnType_STRING},
			{Name: "status_message", Description: "The status description for the host.", Type: proto.ColumnType_STRING},
			{Name: "vpc_configuration", Description: "The VPC configuration provisioned for the host.", Type: proto.ColumnType_JSON},
			{Name: "tags_src", Description: "A list of tags associated with the host.", Type: proto.ColumnType_JSON, Hydrate: listCodeConnectionsHostTags, Transform: transform.FromField("Tags")},
			{Name: "title", Description: resourceInterfaceDescription("title"), Type: proto.ColumnType_STRING, Transform: transform.FromField("Name")},
			{Name: "tags", Description: resourceInterfaceDescription("tags"), Type: proto.ColumnType_JSON, Hydrate: listCodeConnectionsHostTags, Transform: transform.From(codeConnectionsTurbotTags)},
			{Name: "akas", Description: resourceInterfaceDescription("akas"), Type: proto.ColumnType_JSON, Transform: transform.FromField("HostArn").Transform(transform.EnsureStringArray)},
		}),
	}
}

func listCodeConnectionsHosts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_host.listCodeConnectionsHosts", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil && int32(*d.QueryContext.Limit) < maxLimit {
		maxLimit = int32(*d.QueryContext.Limit)
	}
	input := &codeconnections.ListHostsInput{MaxResults: maxLimit}
	paginator := codeconnections.NewListHostsPaginator(svc, input, func(o *codeconnections.ListHostsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_codeconnections_host.listCodeConnectionsHosts", "api_error", err)
			return nil, err
		}
		for _, host := range output.Hosts {
			d.StreamListItem(ctx, host)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

func getCodeConnectionsHost(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQualString("arn")
	if arn == "" {
		return nil, nil
	}
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_host.getCodeConnectionsHost", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	output, err := svc.GetHost(ctx, &codeconnections.GetHostInput{HostArn: aws.String(arn)})
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_host.getCodeConnectionsHost", "api_error", err)
		return nil, err
	}
	return &types.Host{
		HostArn:          aws.String(arn),
		Name:             output.Name,
		ProviderEndpoint: output.ProviderEndpoint,
		ProviderType:     output.ProviderType,
		Status:           output.Status,
		VpcConfiguration: output.VpcConfiguration,
	}, nil
}

func listCodeConnectionsHostTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn *string
	switch item := h.Item.(type) {
	case types.Host:
		arn = item.HostArn
	case *types.Host:
		arn = item.HostArn
	}
	return listCodeConnectionsTags(ctx, d, arn)
}
