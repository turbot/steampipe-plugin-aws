package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2TransitGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_transit_gateway",
		Description: "AWS EC2 Transit Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("transit_gateway_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidTransitGatewayID.NotFound", "InvalidTransitGatewayID.Unavailable", "InvalidTransitGatewayID.Malformed", "InvalidAction"}),
			},
			Hydrate: getEc2TransitGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TransitGateways,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "propagation_default_route_table_id", Require: plugin.Optional},
				{Name: "amazon_side_asn", Require: plugin.Optional},
				{Name: "association_default_route_table_id", Require: plugin.Optional},
				{Name: "auto_accept_shared_attachments", Require: plugin.Optional},
				{Name: "default_route_table_association", Require: plugin.Optional},
				{Name: "default_route_table_propagation", Require: plugin.Optional},
				{Name: "dns_support", Require: plugin.Optional},
				{Name: "vpn_ecmp_support", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidAction"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "transit_gateway_id",
				Description: "The ID of the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_gateway_arn",
				Description: "The Amazon Resource Name (ARN) of the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account ID that owns the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time when transit gateway was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "amazon_side_asn",
				Description: "A private Autonomous System Number (ASN) for the Amazon side of a BGP session. The range is 64512 to 65534 for 16-bit ASNs and 4200000000 to 4294967294 for 32-bit ASNs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Options.AmazonSideAsn"),
			},
			{
				Name:        "association_default_route_table_id",
				Description: "The ID of the default association route table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.AssociationDefaultRouteTableId"),
			},
			{
				Name:        "auto_accept_shared_attachments",
				Description: "Indicates whether attachment requests are automatically accepted.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.AutoAcceptSharedAttachments"),
			},
			{
				Name:        "default_route_table_association",
				Description: "Indicates whether resource attachments are automatically associated with the default association route table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.DefaultRouteTableAssociation"),
			},
			{
				Name:        "default_route_table_propagation",
				Description: "Indicates whether resource attachments are automatically associated with the default association route table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.DefaultRouteTablePropagation"),
			},
			{
				Name:        "dns_support",
				Description: "Indicates whether DNS support is enabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.DnsSupport"),
			},
			{
				Name:        "multicast_support",
				Description: "Indicates whether multicast is enabled on the transit gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.MulticastSupport"),
			},
			{
				Name:        "propagation_default_route_table_id",
				Description: "The ID of the default propagation route table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.PropagationDefaultRouteTableId"),
			},
			{
				Name:        "vpn_ecmp_support",
				Description: "Indicates whether Equal Cost Multipath Protocol support is enabled.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Options.VpnEcmpSupport"),
			},
			{
				Name:        "cidr_blocks",
				Description: "A list of transit gateway CIDR blocks.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Options.TransitGatewayCidrBlocks"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are assigned to the transit gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2TransitGatewayTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2TransitGatewayTurbotTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TransitGatewayArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2TransitGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway.listEc2TransitGateways", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeTransitGatewaysInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filters := buildEc2TransitGatewayFilter(d.Quals)

	equalQuals := d.KeyColumnQuals
	if equalQuals["amazon_side_asn"] != nil {
		filters = append(filters, types.Filter{Name: aws.String("options.amazon-side-asn"), Values: []string{fmt.Sprint(equalQuals["amazon_side_asn"].GetInt64Value())}})
	}

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeTransitGatewaysPaginator(svc, input, func(o *ec2.DescribeTransitGatewaysPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_transit_gateway.listEc2TransitGateways", "api_error", err)
			return nil, err
		}

		for _, items := range output.TransitGateways {
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

func getEc2TransitGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	transitGatewayID := d.KeyColumnQuals["transit_gateway_id"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway.getEc2TransitGateway", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeTransitGatewaysInput{
		TransitGatewayIds: []string{transitGatewayID},
	}

	op, err := svc.DescribeTransitGateways(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway.getEc2TransitGateway", "api_error", err)
		return nil, err
	}

	if op.TransitGateways != nil && len(op.TransitGateways) > 0 {
		return op.TransitGateways[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getEc2TransitGatewayTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.TransitGateway)
	var turbotTagsMap map[string]string
	if data.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range data.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func getEc2TransitGatewayTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.TransitGateway)
	title := data.TransitGatewayId
	if data.Tags != nil {
		for _, i := range data.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}
	return title, nil
}

// // UTILITY FUNCTION
// Build ec2 transit gateway list call input filter
func buildEc2TransitGatewayFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"propagation_default_route_table_id": "options.propagation-default-route-table-id",
		"association_default_route_table_id": "options.association-default-route-table-id",
		"auto_accept_shared_attachments":     "options.auto-accept-shared-attachments",
		"default_route_table_association":    "options.default-route-table-association",
		"default_route_table_propagation":    "options.default-route-table-propagation",
		"dns_support":                        "options.dns-support",
		"vpn_ecmp_support":                   "options.vpn-ecmp-support",
		"owner_id":                           "owner-id",
		"state":                              "state",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
