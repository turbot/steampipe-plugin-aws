package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice/types"

	databasemigrationservicev1 "github.com/aws/aws-sdk-go/service/databasemigrationservice"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDmsReplicationInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dms_replication_instance",
		Description: "AWS DMS Replication Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException", "ResourceNotFoundFault", "InvalidParameterCombinationException"}),
			},
			Hydrate: getDmsReplicationInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listDmsReplicationInstances,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "replication_instance_identifier",
					Require: plugin.Optional,
				},
				{
					Name:    "arn",
					Require: plugin.Optional,
				},
				{
					Name:    "replication_instance_class",
					Require: plugin.Optional,
				},
				{
					Name:    "engine_version",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(databasemigrationservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "replication_instance_identifier",
				Description: "The identifier of the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the replication instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReplicationInstanceArn"),
			},
			{
				Name:        "replication_instance_class",
				Description: "The compute and memory capacity of the replication instance as defined for the specified replication instance class.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "The engine version number of the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "publicly_accessible",
				Description: "Specifies the accessibility options for the replication instance.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "allocated_storage",
				Description: "The amount of storage (in gigabytes) that is allocated for the replication instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "auto_minor_version_upgrade",
				Description: "Boolean value indicating if minor version upgrades will be automatically applied to the instance.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "availability_zone",
				Description: "The Availability Zone for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_name_servers",
				Description: "The DNS name servers supported for the replication instance to access your on-premise source or target database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "free_until",
				Description: "The expiration date of the free replication instance that is part of the Free DMS program.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "instance_create_time",
				Description: "The time the replication instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kms_key_id",
				Description: "An AWS KMS key identifier that is used to encrypt the data on the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "multi_az",
				Description: "Specifies whether the replication instance is a Multi-AZ deployment.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("MultiAZ"),
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "The maintenance window times for the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_instance_private_ip_address",
				Description: "The private IP address of the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_instance_public_ip_address",
				Description: "The public IP address of the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_instance_status",
				Description: "The status of the replication instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secondary_availability_zone",
				Description: "The Availability Zone of the standby replication instance in a Multi-AZ deployment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pending_modified_values",
				Description: "The pending modification values.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replication_instance_private_ip_addresses",
				Description: "One or more private IP addresses for the replication instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replication_instance_public_ip_addresses",
				Description: "One or more public IP addresses for the replication instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replication_subnet_group",
				Description: "The subnet group for the replication instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_security_groups",
				Description: "The VPC security group for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the replication instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsReplicationInstanceTags,
				Transform:   transform.FromField("TagList"),
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReplicationInstanceIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsReplicationInstanceTags,
				Transform:   transform.From(dmsReplicationInstanceTagListToTagsMap),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ReplicationInstanceArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listDmsReplicationInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_instance.listDmsReplicationInstances", "connection_error", err)
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
	input := &databasemigrationservice.DescribeReplicationInstancesInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	var filter []types.Filter

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["replication_instance_identifier"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("replication-instance-id"),
			Values: []string{equalQuals["replication_instance_identifier"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	if equalQuals["arn"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("replication-instance-arn"),
			Values: []string{equalQuals["arn"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	if equalQuals["replication_instance_class"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("replication-instance-class"),
			Values: []string{equalQuals["replication_instance_class"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	if equalQuals["engine_version"] != nil {
		paramFilter := types.Filter{
			Name:   aws.String("engine-version"),
			Values: []string{equalQuals["engine_version"].GetStringValue()},
		}
		filter = append(filter, paramFilter)
	}
	input.Filters = filter

	paginator := databasemigrationservice.NewDescribeReplicationInstancesPaginator(svc, input, func(o *databasemigrationservice.DescribeReplicationInstancesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dms_replication_instance.listDmsReplicationInstances", "api_error", err)
			return nil, err
		}

		for _, items := range output.ReplicationInstances {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDmsReplicationInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_instance.getDmsReplicationInstance", "connection_error", err)
		return nil, err
	}

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	params := &databasemigrationservice.DescribeReplicationInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("replication-instance-arn"),
				Values: []string{arn},
			},
		},
	}

	op, err := svc.DescribeReplicationInstances(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_instance.getDmsReplicationInstance", "api_error", err)
		return nil, err
	}

	if op.ReplicationInstances != nil && len(op.ReplicationInstances) > 0 {
		return op.ReplicationInstances[0], nil
	}
	return nil, nil
}

func getDmsReplicationInstanceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	replicationInstanceArn := h.Item.(types.ReplicationInstance).ReplicationInstanceArn

	// Create service
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_instance.getDmsReplicationInstance", "connection_error", err)
		return nil, err
	}

	params := &databasemigrationservice.ListTagsForResourceInput{
		ResourceArn: replicationInstanceArn,
	}

	replicationInstanceTags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_replication_instance.getDmsReplicationInstance", "api_error", err)
		return nil, err
	}

	return replicationInstanceTags, nil
}

//// TRANSFORM FUNCTIONS

func dmsReplicationInstanceTagListToTagsMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
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
