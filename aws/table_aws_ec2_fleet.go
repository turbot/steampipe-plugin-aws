package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2Fleet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_fleet",
		Description: "AWS EC2 Fleet",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("fleet_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidFleetId.NotFound", "InvalidFleetId.Malformed"}),
			},
			Hydrate: getEc2Fleet,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeFleets"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2Fleets,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeFleets"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "fleet_state", Require: plugin.Optional},
				{Name: "activity_status", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "excess_capacity_termination_policy", Require: plugin.Optional},
				{Name: "replace_unhealthy_instances", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EC2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "fleet_id",
				Description: "The ID of the EC2 Fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the EC2 Fleet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEc2FleetARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "fleet_state",
				Description: "The state of the EC2 Fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activity_status",
				Description: "The progress of the EC2 Fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation date and time of the EC2 Fleet.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "type",
				Description: "The type of request. Indicates whether the EC2 Fleet only requests the target capacity or also attempts to maintain it.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_token",
				Description: "Unique, case-sensitive identifier that you provide to ensure the idempotency of the request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "context",
				Description: "Reserved.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "excess_capacity_termination_policy",
				Description: "Indicates whether running instances should be terminated if the target capacity of the EC2 Fleet is decreased below the current size of the EC2 Fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fulfilled_capacity",
				Description: "The number of units fulfilled by this request compared to the set target capacity.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "fulfilled_on_demand_capacity",
				Description: "The number of units fulfilled by this request compared to the set target On-Demand capacity.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "replace_unhealthy_instances",
				Description: "Indicates whether EC2 Fleet should replace unhealthy Spot Instances.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "terminate_instances_with_expiration",
				Description: "Indicates whether running instances should be terminated when the EC2 Fleet expires.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "valid_from",
				Description: "The start date and time of the request.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "valid_until",
				Description: "The end date and time of the request.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "target_capacity_specification",
				Description: "The number of units to request. You can choose to set the target capacity in terms of instances or a performance characteristic that is important to your application workload.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "spot_options",
				Description: "The configuration of Spot Instances in an EC2 Fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "on_demand_options",
				Description: "The allocation strategy of On-Demand Instances in an EC2 Fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "launch_template_configs",
				Description: "The launch template and overrides.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instances",
				Description: "Information about the instances that were launched by the fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "errors",
				Description: "Information about the instances that could not be launched by the fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the EC2 Fleet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FleetId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(ec2FleetTagsToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEc2FleetARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2Fleets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_fleet.listEc2Fleets", "connection_error", err)
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

	input := &ec2.DescribeFleetsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// Build filters
	filters := buildEc2FleetFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeFleetsPaginator(svc, input, func(o *ec2.DescribeFleetsPaginatorOptions) {
		o.Limit = maxLimit
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_fleet.listEc2Fleets", "api_error", err)
			return nil, err
		}

		for _, fleet := range output.Fleets {
			d.StreamListItem(ctx, fleet)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEc2Fleet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fleetId := d.EqualsQualString("fleet_id")

	// Empty check
	if fleetId == "" {
		return nil, nil
	}

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_fleet.getEc2Fleet", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeFleetsInput{
		FleetIds: []string{fleetId},
	}

	op, err := svc.DescribeFleets(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_fleet.getEc2Fleet", "api_error", err)
		return nil, err
	}

	if len(op.Fleets) > 0 {
		return op.Fleets[0], nil
	}

	return nil, nil
}

func getEc2FleetARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	fleet := h.Item.(types.FleetData)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_fleet.getEc2FleetARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Build ARN
	// arn:${Partition}:ec2:${Region}:${Account}:fleet/${FleetId}
	// https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazonec2.html#amazonec2-fleet
	arn := "arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":fleet/" + *fleet.FleetId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func ec2FleetTagsToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	fleet := d.HydrateItem.(types.FleetData)

	if fleet.Tags == nil {
		return nil, nil
	}

	turbotTags := make(map[string]string)
	for _, tag := range fleet.Tags {
		if tag.Key != nil && tag.Value != nil {
			turbotTags[*tag.Key] = *tag.Value
		}
	}
	return turbotTags, nil
}

//// UTILITY FUNCTIONS

func buildEc2FleetFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"fleet_state":                        "fleet-state",
		"activity_status":                    "activity-status",
		"type":                               "type",
		"excess_capacity_termination_policy": "excess-capacity-termination-policy",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			for _, q := range quals[columnName].Quals {
				value := q.Value.GetStringValue()
				if value != "" {
					filter.Values = append(filter.Values, value)
				}
			}
			if len(filter.Values) > 0 {
				filters = append(filters, filter)
			}
		}
	}

	// Handle boolean filter for replace_unhealthy_instances
	if quals["replace_unhealthy_instances"] != nil {
		for _, q := range quals["replace_unhealthy_instances"].Quals {
			value := q.Value.GetBoolValue()
			filter := types.Filter{
				Name:   aws.String("replace-unhealthy-instances"),
				Values: []string{fmt.Sprintf("%t", value)},
			}
			filters = append(filters, filter)
		}
	}

	return filters
}
