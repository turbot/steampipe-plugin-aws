package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSReservedDBInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_reserved_db_instance",
		Description: "AWS RDS Reserved DB Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("reserved_db_instance_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ReservedDBInstanceNotFound"}),
			},
			Hydrate: getRDSReservedDBInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSReservedDBInstances,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "class", Require: plugin.Optional},
				{Name: "duration", Require: plugin.Optional},
				{Name: "lease_id", Require: plugin.Optional},
				{Name: "multi_az", Require: plugin.Optional},
				{Name: "offering_type", Require: plugin.Optional},
				{Name: "reserved_db_instances_offering_id", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "reserved_db_instance_id",
				Description: "The unique identifier for the reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservedDBInstanceId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the reserved DB Instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservedDBInstanceArn"),
			},
			{
				Name:        "reserved_db_instances_offering_id",
				Description: "The offering identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservedDBInstancesOfferingId"),
			},
			{
				Name:        "state",
				Description: "The state of the reserved DB instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "class",
				Description: "The DB instance class for the reserved DB instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBInstanceClass"),
			},
			{
				Name:        "currency_code",
				Description: "The currency code for the reserved DB instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_instance_count",
				Description: "The number of reserved DB instances.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DBInstanceCount"),
			},
			{
				Name:        "duration",
				Description: "The duration of the reservation in seconds.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "fixed_price",
				Description: "The fixed price charged for this reserved DB instance.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "lease_id",
				Description: "The unique identifier for the lease associated with the reserved DB instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "multi_az",
				Description: "Indicates if the reservation applies to Multi-AZ deployments.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("MultiAZ"),
			},
			{
				Name:        "offering_type",
				Description: "The offering type of this reserved DB instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_description",
				Description: "The description of the reserved DB instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The time the reservation started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "usage_price",
				Description: "The hourly price charged for this reserved DB instance.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "recurring_charges",
				Description: "The recurring price charged to run this reserved DB instance.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReservedDBInstanceId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ReservedDBInstanceArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSReservedDBInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_reserved_db_instance.listRDSReservedDBInstances", "connection_error", err)
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

	input := &rds.DescribeReservedDBInstancesInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	if d.KeyColumnQuals["class"] != nil {
		input.DBInstanceClass = aws.String(d.KeyColumnQuals["class"].GetStringValue())
	}

	if d.KeyColumnQuals["duration"] != nil {
		input.Duration = aws.String(fmt.Sprintf("%v", d.KeyColumnQuals["duration"].GetInt64Value()))
	}

	if d.KeyColumnQuals["lease_id"] != nil {
		input.LeaseId = aws.String(d.KeyColumnQuals["lease_id"].GetStringValue())
	}

	if d.KeyColumnQuals["multi_az"] != nil {
		input.MultiAZ = aws.Bool(d.KeyColumnQuals["multi_az"].GetBoolValue())
	}

	if d.KeyColumnQuals["offering_type"] != nil {
		offeringType := d.KeyColumnQuals["offering_type"].GetStringValue()
		if offeringType != "Partial Upfront" && offeringType != "All Upfront" && offeringType != "No Upfront" {
			return nil, nil
		}
		input.OfferingType = aws.String(offeringType)
	}

	if d.KeyColumnQuals["reserved_db_instances_offering_id"] != nil {
		input.ReservedDBInstancesOfferingId = aws.String(d.KeyColumnQuals["reserved_db_instances_offering_id"].GetStringValue())
	}

	paginator := rds.NewDescribeReservedDBInstancesPaginator(svc, input, func(o *rds.DescribeReservedDBInstancesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_reserved_db_instance.listRDSReservedDBInstances", "api_error", err)
			return nil, err
		}

		for _, items := range output.ReservedDBInstances {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRDSReservedDBInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	dbInstanceIdentifier := d.KeyColumnQuals["reserved_db_instance_id"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_reserved_db_instance.getRDSReservedDBInstance", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeReservedDBInstancesInput{
		ReservedDBInstanceId: aws.String(dbInstanceIdentifier),
	}

	op, err := svc.DescribeReservedDBInstances(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_reserved_db_instance.getRDSReservedDBInstance", "api_error", err)
		return nil, err
	}

	if len(op.ReservedDBInstances) > 0 {
		return op.ReservedDBInstances[0], nil
	}
	return nil, nil
}
