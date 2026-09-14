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

func tableAwsCodeConnectionsSyncConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codeconnections_sync_configuration",
		Description: "AWS CodeConnections Sync Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"resource_name", "sync_type"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getCodeConnectionsSyncConfiguration,
			Tags:    map[string]string{"service": "codeconnections", "action": "GetSyncConfiguration"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCodeConnectionsRepositoryLinks,
			Hydrate:       listCodeConnectionsSyncConfigurations,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "repository_link_id", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "sync_type", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Tags: map[string]string{"service": "codeconnections", "action": "ListSyncConfigurations"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CODESTAR_CONNECTIONS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{Name: "resource_name", Description: "The name of the AWS resource associated with the sync configuration.", Type: proto.ColumnType_STRING},
			{Name: "sync_type", Description: "The type of sync for the configuration.", Type: proto.ColumnType_STRING},
			{Name: "repository_link_id", Description: "The ID of the repository link associated with the sync configuration.", Type: proto.ColumnType_STRING},
			{Name: "repository_name", Description: "The name of the repository associated with the sync configuration.", Type: proto.ColumnType_STRING},
			{Name: "owner_id", Description: "The owner ID for the repository.", Type: proto.ColumnType_STRING},
			{Name: "provider_type", Description: "The connection provider type associated with the sync configuration.", Type: proto.ColumnType_STRING},
			{Name: "branch", Description: "The branch associated with the sync configuration.", Type: proto.ColumnType_STRING},
			{Name: "config_file", Description: "The path to the configuration file in the linked repository.", Type: proto.ColumnType_STRING},
			{Name: "role_arn", Description: "The ARN of the IAM role associated with the sync configuration.", Type: proto.ColumnType_STRING},
			{Name: "publish_deployment_status", Description: "Whether deployment status is published to the source provider.", Type: proto.ColumnType_STRING},
			{Name: "pull_request_comment", Description: "Whether pull request comments are enabled for the sync configuration.", Type: proto.ColumnType_STRING},
			{Name: "trigger_resource_update_on", Description: "The changes that trigger Git sync to update the resource.", Type: proto.ColumnType_STRING},
			{Name: "title", Description: resourceInterfaceDescription("title"), Type: proto.ColumnType_STRING, Transform: transform.FromField("ResourceName")},
		}),
	}
}

func listCodeConnectionsSyncConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	link := h.Item.(types.RepositoryLinkInfo)
	wantedID := d.EqualsQualString("repository_link_id")
	if wantedID != "" && aws.ToString(link.RepositoryLinkId) != wantedID {
		return nil, nil
	}

	syncTypes := types.SyncConfigurationType("").Values()
	if wantedType := d.EqualsQualString("sync_type"); wantedType != "" {
		syncTypes = nil
		for _, syncType := range types.SyncConfigurationType("").Values() {
			if string(syncType) == wantedType {
				syncTypes = append(syncTypes, syncType)
			}
		}
	}
	if len(syncTypes) == 0 {
		return nil, nil
	}

	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_sync_configuration.listCodeConnectionsSyncConfigurations", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil && int32(*d.QueryContext.Limit) < maxLimit {
		maxLimit = int32(*d.QueryContext.Limit)
	}

	for _, syncType := range syncTypes {
		input := &codeconnections.ListSyncConfigurationsInput{
			RepositoryLinkId: link.RepositoryLinkId,
			SyncType:         syncType,
			MaxResults:       maxLimit,
		}
		paginator := codeconnections.NewListSyncConfigurationsPaginator(svc, input, func(o *codeconnections.ListSyncConfigurationsPaginatorOptions) {
			o.Limit = maxLimit
			o.StopOnDuplicateToken = true
		})
		for paginator.HasMorePages() {
			d.WaitForListRateLimit(ctx)
			output, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("aws_codeconnections_sync_configuration.listCodeConnectionsSyncConfigurations", "api_error", err)
				return nil, err
			}
			for _, configuration := range output.SyncConfigurations {
				d.StreamLeafListItem(ctx, configuration)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}
	return nil, nil
}

func getCodeConnectionsSyncConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	resourceName := d.EqualsQualString("resource_name")
	syncType := d.EqualsQualString("sync_type")
	if resourceName == "" || syncType == "" {
		return nil, nil
	}
	svc, err := CodeConnectionsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_sync_configuration.getCodeConnectionsSyncConfiguration", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	output, err := svc.GetSyncConfiguration(ctx, &codeconnections.GetSyncConfigurationInput{
		ResourceName: aws.String(resourceName),
		SyncType:     types.SyncConfigurationType(syncType),
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_codeconnections_sync_configuration.getCodeConnectionsSyncConfiguration", "api_error", err)
		return nil, err
	}
	return output.SyncConfiguration, nil
}
