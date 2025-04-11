package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsEc2TransitGatewayVpcAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_transit_gateway_vpc_attachment",
		Description: "AWS EC2 Transit Gateway VPC Attachment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("transit_gateway_attachment_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidTransitGatewayAttachmentID.NotFound", "InvalidTransitGatewayAttachmentID.Unavailable", "InvalidTransitGatewayAttachmentID.Malformed", "InvalidAction"}),
			},
			Hydrate: getEc2TransitGatewayVpcAttachment,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeTransitGatewayAttachments"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TransitGatewayVpcAttachment,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeTransitGatewayAttachments"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "association_state", Require: plugin.Optional},
				{Name: "association_transit_gateway_route_table_id", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "resource_owner_id", Require: plugin.Optional},
				{Name: "resource_type", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
				{Name: "transit_gateway_id", Require: plugin.Optional},
				{Name: "transit_gateway_owner_id", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAction"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "transit_gateway_attachment_id",
				Description: "The ID of the transit gateway attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_gateway_id",
				Description: "The ID of the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transit_gateway_owner_id",
				Description: "The ID of the AWS account that owns the transit gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The attachment state of the transit gateway attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the transit gateway attachment.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "resource_id",
				Description: "The ID of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The resource type of the transit gateway attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_owner_id",
				Description: "The ID of the AWS account that owns the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_state",
				Description: "The state of the association.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Association.State"),
			},
			{
				Name:        "association_transit_gateway_route_table_id",
				Description: "The ID of the route table for the transit gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Association.TransitGatewayRouteTableId"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			/// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(transitGatewayAttachmentRawTagsToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEc2TransitGatewayAttachmentTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2TransitGatewayVpcAttachmentAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2TransitGatewayVpcAttachment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_vpc_attachment.listEc2TransitGatewayVpcAttachment", "connection_error", err)
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

	input := &ec2.DescribeTransitGatewayAttachmentsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filters := buildEc2TransitGatewayVpcAttachmentFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeTransitGatewayAttachmentsPaginator(svc, input, func(o *ec2.DescribeTransitGatewayAttachmentsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_transit_gateway_vpc_attachment.listEc2TransitGatewayVpcAttachment", "api_error", err)
			return nil, err
		}

		for _, items := range output.TransitGatewayAttachments {
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

func getEc2TransitGatewayVpcAttachment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	transitGatewayAttachmentID := d.EqualsQuals["transit_gateway_attachment_id"].GetStringValue()

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_vpc_attachment.getEc2TransitGatewayVpcAttachment", "connection_error", err)
		return nil, err
	}

	// Build params
	params := &ec2.DescribeTransitGatewayAttachmentsInput{
		TransitGatewayAttachmentIds: []string{transitGatewayAttachmentID},
	}

	op, err := svc.DescribeTransitGatewayAttachments(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_transit_gateway_vpc_attachment.getEc2TransitGatewayVpcAttachment", "api_error", err)
		return nil, err
	}

	if len(op.TransitGatewayAttachments) > 0 {
		return op.TransitGatewayAttachments[0], nil
	}
	return nil, nil
}

func getAwsEc2TransitGatewayVpcAttachmentAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	transitGatewayAttachment := h.Item.(types.TransitGatewayAttachment)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get the resource akas
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":transit-gateway-attachment/" + *transitGatewayAttachment.TransitGatewayAttachmentId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func transitGatewayAttachmentRawTagsToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.TransitGatewayAttachment)
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

func getEc2TransitGatewayAttachmentTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.TransitGatewayAttachment)
	title := data.TransitGatewayAttachmentId
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
// Build ec2 transit gateway VPC attachment list call input filter
func buildEc2TransitGatewayVpcAttachmentFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"association_state":                          "association.state",
		"association_transit_gateway_route_table_id": "association.transit-gateway-route-table-id",
		"resource_id":                                "resource-id",
		"resource_owner_id":                          "resource-owner-id",
		"resource_type":                              "resource-type",
		"state":                                      "state",
		"transit_gateway_id":                         "transit-gateway-id",
		"transit_gateway_owner_id":                   "transit-gateway-owner-id",
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
