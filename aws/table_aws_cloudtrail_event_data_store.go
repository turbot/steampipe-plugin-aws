package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	cloudtrailEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudtrailEventDataStore(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_event_data_store",
		Description: "AWS CloudTrail Event Data Store",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getCloudTrailEventDataStore,
			Tags:       map[string]string{"service": "cloudtrail", "action": "GetEventDataStore"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudTrailEventDataStores,
			Tags:    map[string]string{"service": "cloudtrail", "action": "ListEventDataStores"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudTrailEventDataStore,
				Tags: map[string]string{"service": "cloudtrail", "action": "GetEventDataStore"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudtrailEndpoint.CLOUDTRAILServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the event data store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the event data store.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EventDataStoreArn"),
			},
			{
				Name:        "status",
				Description: "The status of an event data store.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "billing_mode",
				Description: "The billing mode for the event data store.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "federation_role_arn",
				Description: "If Lake query federation is enabled, provides the ARN of the federation role used to access the resources for the federated event data store.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "federation_status",
				Description: "Indicates the Lake query federation status.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "kms_key_id",
				Description: "Specifies the KMS key ID that encrypts the events delivered by CloudTrail.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "created_timestamp",
				Description: "The timestamp of the event data store's creation.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "multi_region_enabled",
				Description: "Indicates whether the event data store includes events from all regions, or only from the region in which it was created.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "organization_enabled",
				Description: "Indicates that an event data store is collecting logged events for an organization.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "retention_period",
				Description: "The retention period, in days.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "termination_protection_enabled",
				Description: "Indicates whether the event data store is protected from termination.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "updated_timestamp",
				Description: "The timestamp showing when an event data store was updated, if applicable. UpdatedTimestamp is always either the same or newer than the time shown in CreatedTimestamp.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudTrailEventDataStore,
			},
			{
				Name:        "advanced_event_selectors",
				Description: "The advanced event selectors that were used to select events for the data store.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudTrailEventDataStore,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EventDataStoreArn").Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudTrailEventDataStores(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_event_data_store.listCloudTrailEventDataStores", "client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &cloudtrail.ListEventDataStoresInput{
		// Default to the maximum allowed
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := cloudtrail.NewListEventDataStoresPaginator(svc, input, func(o *cloudtrail.ListEventDataStoresPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudtrail_event_data_store.listCloudTrailEventDataStores", "api_error", err)
			return nil, err
		}

		for _, item := range output.EventDataStores {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudTrailEventDataStore(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dataStore := h.Item.(types.EventDataStore)

	equalQuals := d.EqualsQuals
	if equalQuals["arn"] != nil {
		if equalQuals["arn"].GetStringValue() != *dataStore.EventDataStoreArn {
			return nil, nil
		}
	}

	// Create session
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_event_data_store.getCloudTrailEventDataStore", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.GetEventDataStoreInput{
		EventDataStore: dataStore.EventDataStoreArn,
	}

	// execute list call
	op, err := svc.GetEventDataStore(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_event_data_store.getCloudTrailEventDataStore", "api_error", err)
		return nil, err
	}

	return op, nil
}
