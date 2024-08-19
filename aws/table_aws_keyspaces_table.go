package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/keyspaces"
	"github.com/aws/aws-sdk-go-v2/service/keyspaces/types"
	keyspacesv1 "github.com/aws/aws-sdk-go/service/keyspaces"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKeyspacesTable(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_keyspaces_table",
		Description: "AWS Keyspace Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"keyspace_name", "table_name"}),
			Hydrate:    getKeyspacesTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "keyspaces", "action": "GetTable"},
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyspacesTables,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "keyspace_name", Require: plugin.Required},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags: map[string]string{"service": "keyspaces", "action": "ListTables"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(keyspacesv1.EndpointsID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getKeyspacesTable,
				Tags: map[string]string{"service": "keyspaces", "action": "GetTable"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "keyspace_name",
				Description: "The name of the keyspace that the table is stored in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "table_name",
				Description: "The name of the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The unique identifier of the table in the format of an Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
			{
				Name:        "status",
				Description: "The current status of the specified table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyspacesTable,
			},
			{
				Name:        "ttl_status",
				Description: "The custom Time to Live settings of the specified table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyspacesTable,
				Transform:   transform.FromField("Ttl.Status"),
			},
			{
				Name:        "client_side_timestamps_status",
				Description: "Shows how to enable client-side timestamps settings for the specified table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyspacesTable,
				Transform:   transform.FromField("ClientSideTimestamps.Status"),
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the specified table.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getKeyspacesTable,
			},
			{
				Name:        "default_time_to_live",
				Description: "The default Time to Live settings in seconds of the specified table.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKeyspacesTable,
			},
			{
				Name:        "comment_message",
				Description: "An optional description of the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyspacesTable,
				Transform:   transform.FromField("Comment.Message"),
			},
			{
				Name:        "encryption_specification_type",
				Description: "The encryption option specified for the table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyspacesTable,
				Transform:   transform.FromField("EncryptionSpecification.Type"),
			},
			{
				Name:        "kms_key_identifier",
				Description: "The Amazon Resource Name (ARN) of the customer managed KMS key,",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyspacesTable,
				Transform:   transform.FromField("EncryptionSpecification.KmsKeyIdentifier"),
			},
			{
				Name:        "point_in_time_recovery",
				Description: "The point-in-time recovery status of the specified table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyspacesTable,
			},
			{
				Name:        "capacity_specification",
				Description: "The read/write throughput capacity mode for a table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyspacesTable,
			},
			{
				Name:        "replica_specifications",
				Description: "Returns the Amazon Web Services Region specific settings of all Regions a multi-Region table is replicated in.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyspacesTable,
			},
			{
				Name:        "schema_definition",
				Description: "The schema definition of the specified table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyspacesTable,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableName"),
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

func listKeyspacesTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	keyspaceName := d.EqualsQualString("keyspace_name")

	// Create Session
	svc, err := KeyspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_keyspaces_table.listKeyspacesTables", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(1000)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {

			maxItems = limit
		}
	}

	input := &keyspaces.ListTablesInput{
		MaxResults:   &maxItems,
		KeyspaceName: &keyspaceName,
	}

	// pagination not supported for ListTables API
	for {
		d.WaitForListRateLimit(ctx)

		result, err := svc.ListTables(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_keyspaces_table.listKeyspacesTables", "api_error", err)
			return nil, err
		}

		for _, op := range result.Tables {
			d.StreamListItem(ctx, op)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			input.NextToken = result.NextToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getKeyspacesTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	keyspaceName := d.EqualsQualString("keyspace_name")
	tableName := d.EqualsQualString("table_name")

	if h.Item != nil {
		keyspaceTable := h.Item.(types.TableSummary)
		keyspaceName = *keyspaceTable.KeyspaceName
		tableName = *keyspaceTable.TableName
	}

	// Empty id check
	if keyspaceName == "" || tableName == "" {
		return nil, nil
	}

	// Create Session
	svc, err := KeyspacesClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_keyspaces_table.getKeyspacesTable", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &keyspaces.GetTableInput{
		KeyspaceName: &keyspaceName,
		TableName:    &tableName,
	}

	op, err := svc.GetTable(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_keyspaces_table.getKeyspacesTable", "api_error", err)
		return nil, err
	}

	return op, nil
}
