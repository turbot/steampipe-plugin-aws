package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice/types"

	databasemigrationserviceEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDmsReplicationTask(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dms_replication_task",
		Description: "AWS DMS Replication Task",
		List: &plugin.ListConfig{
			Hydrate: listDmsReplicationTasks,
			// If the ARN provided as an input parameter refers to a resource that is unavailable in the specified region, the API throws an InvalidParameterValueException exception.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException", "ResourceNotFoundFault"}),
			},
			Tags: map[string]string{"service": "dms", "action": "DescribeReplicationTasks"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "replication_task_identifier",
					Require: plugin.Optional,
				},
				{
					Name:    "arn",
					Require: plugin.Optional,
				},
				{
					Name:    "replication_instance_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "target_endpoint_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "source_endpoint_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "migration_type",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getDmsReplicationTaskTags,
				Tags: map[string]string{"service": "dms", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(databasemigrationserviceEndpoint.DMSServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "replication_task_identifier",
				Description: "The user-assigned replication task identifier or name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the replication task.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReplicationTaskArn"),
			},
			{
				Name:        "cdc_start_position",
				Description: "Indicates when you want a change data capture (CDC) operation to start.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cdc_stop_position",
				Description: "Indicates when you want a change data capture (CDC) operation to stop.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_failure_message",
				Description: "The last error (failure) message generated for the replication task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "migration_type",
				Description: "The type of migration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recovery_checkpoint",
				Description: "Indicates the last checkpoint that occurred during a change data capture (CDC) operation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_instance_arn",
				Description: "The Amazon Resource Name (ARN) of the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_task_creation_date",
				Description: "The date the replication task was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "replication_task_start_date",
				Description: "The date the replication task is scheduled to start.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "source_endpoint_arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the replication task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stop_reason",
				Description: "The reason the replication task was stopped.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "table_mappings",
				Description: "Table mappings specified in the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_endpoint_arn",
				Description: "The ARN that uniquely identifies the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_replication_instance_arn",
				Description: "The ARN of the replication instance to which this task is moved in response to running the MoveReplicationTask operation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "task_data",
				Description: "Supplemental information that the task requires to migrate the data for certain source and target endpoints.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_task_settings",
				Description: "The settings for the replication task.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replication_task_stats",
				Description: "The statistics for the task, including elapsed time, tables loaded, and table errors.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the replication instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsReplicationTaskTags,
				Transform:   transform.FromField("TagList"),
			},

			// Steampipe Standard Columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReplicationTaskIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsReplicationTaskTags,
				Transform:   transform.From(dmsReplicationTaskTagListToTagsMap),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ReplicationTaskArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDmsReplicationTasks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_task.listDmsReplicationTasks", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	// Build the params
	input := &databasemigrationservice.DescribeReplicationTasksInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	var filter []types.Filter

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["replication_task_identifier"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("replication-task-id"),
			Values: []string{equalQuals["replication_task_identifier"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	if equalQuals["arn"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("replication-task-arn"),
			Values: []string{equalQuals["arn"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	if equalQuals["replication_instance_arn"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("replication-instance-arn"),
			Values: []string{equalQuals["replication_instance_arn"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	if equalQuals["migration_type"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("migration-type"),
			Values: []string{equalQuals["migration_type"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	if equalQuals["target_endpoint_arn"] != nil || equalQuals["source_endpoint_arn"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("endpoint-arn"),
		}
		values := []string{}

		if equalQuals["target_endpoint_arn"].GetStringValue() != "" {
			values = append(values, equalQuals["target_endpoint_arn"].GetStringValue())
		}
		if equalQuals["source_endpoint_arn"].GetStringValue() != "" {
			values = append(values, equalQuals["source_endpoint_arn"].GetStringValue())
		}
		paramFilter.Values = values
		filter = append(filter, paramFilter)
	}
	input.Filters = filter

	paginator := databasemigrationservice.NewDescribeReplicationTasksPaginator(svc, input, func(o *databasemigrationservice.DescribeReplicationTasksPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dms_replication_task.listDmsReplicationTasks", "api_error", err)
			return nil, err
		}

		for _, items := range output.ReplicationTasks {
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

func getDmsReplicationTaskTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	taskArn := h.Item.(types.ReplicationTask).ReplicationTaskArn

	// Create service
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_task.getDmsReplicationTaskTags", "connection_error", err)
		return nil, err
	}

	params := &databasemigrationservice.ListTagsForResourceInput{
		ResourceArn: taskArn,
	}

	replicationInstanceTags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_task.getDmsReplicationTaskTags", "api_error", err)
		return nil, err
	}

	return replicationInstanceTags, nil
}

//// TRANSFORM FUNCTIONS

func dmsReplicationTaskTagListToTagsMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*databasemigrationservice.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	if data.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
