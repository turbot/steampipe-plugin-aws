package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/keyspaces"

	keyspacesEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKeyspacesKeyspace(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_keyspaces_keyspace",
		Description: "AWS Keyspaces Keyspace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("keyspace_name"), // Identify the keyspace by its name
			Hydrate:    getKeyspacesKeyspace,                // Get function
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "keyspaces", "action": "GetKeyspace"},
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyspacesKeyspaces, // Parent hydrate function
			Tags:    map[string]string{"service": "keyspaces", "action": "ListKeyspaces"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(keyspacesEndpoint.CASSANDRAServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "keyspace_name",
				Description: "The name of the keyspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The unique identifier of the keyspace in the format of an Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
			{
				Name:        "replication_strategy",
				Description: "Returns the replication strategy of the keyspace. The options are SINGLE_REGION or MULTI_REGION .",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_regions",
				Description: "If the replication strategy of the keyspace is MULTI_REGION, a list of replication regions is returned.",
				Type:        proto.ColumnType_JSON,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyspaceName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listKeyspacesKeyspaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := KeyspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_keyspaces_keyspace.listKeyspacesKeyspaces", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Define max results for the request
	maxItems := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	input := &keyspaces.ListKeyspacesInput{
		MaxResults: &maxItems,
	}

	paginator := keyspaces.NewListKeyspacesPaginator(svc, input, func(o *keyspaces.ListKeyspacesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_keyspaces_keyspace.listKeyspacesKeyspaces", "api_error", err)
			return nil, err
		}

		for _, keyspace := range output.Keyspaces {
			d.StreamListItem(ctx, keyspace)

			// Stop processing if context is canceled or limit is reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getKeyspacesKeyspace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	keyspaceName := d.EqualsQualString("keyspace_name")

	// Empty check for keyspace name
	if keyspaceName == "" {
		return nil, nil
	}

	// Create session
	svc, err := KeyspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_keyspaces_keyspace.getKeyspacesKeyspace", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region
		return nil, nil
	}

	// Build the input parameters
	input := &keyspaces.GetKeyspaceInput{
		KeyspaceName: &keyspaceName,
	}

	// Make the GetKeyspace API call
	result, err := svc.GetKeyspace(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_keyspaces_keyspace.getKeyspacesKeyspace", "api_error", err)
		return nil, err
	}

	return result, nil
}
