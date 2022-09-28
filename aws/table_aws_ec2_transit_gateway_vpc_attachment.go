package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsEc2TransitGatewayVpcAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "aws_ec2_transit_gateway_vpc_attachment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("transit_gateway_attachment_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidTransitGatewayAttachmentID.NotFound", "InvalidTransitGatewayAttachmentID.Unavailable", "InvalidTransitGatewayAttachmentID.Malformed", "InvalidAction"}),
			},
			Hydrate: getEc2TransitGatewayVpcAttachment,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2TransitGatewayVpcAttachment,
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
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidAction"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listEc2TransitGatewayVpcAttachment", "AWS_REGION", region)

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeTransitGatewayAttachmentsInput{
		MaxResults: aws.Int64(1000),
	}

	filters := buildEc2TransitGatewayVpcAttachmentFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 5 {
				input.MaxResults = aws.Int64(5)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeTransitGatewayAttachmentsPages(
		input,
		func(page *ec2.DescribeTransitGatewayAttachmentsOutput, isLast bool) bool {
			for _, transitGatewayAttachment := range page.TransitGatewayAttachments {
				d.StreamListItem(ctx, transitGatewayAttachment)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2TransitGatewayVpcAttachment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2TransitGatewayVpcAttachment")

	region := d.KeyColumnQualString(matrixKeyRegion)
	transitGatewayAttachmentID := d.KeyColumnQuals["transit_gateway_attachment_id"].GetStringValue()

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &ec2.DescribeTransitGatewayAttachmentsInput{
		TransitGatewayAttachmentIds: []*string{aws.String(transitGatewayAttachmentID)},
	}

	op, err := svc.DescribeTransitGatewayAttachments(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getEc2TransitGatewayVpcAttachment__", "ERROR", err)
		return nil, err
	}

	if op.TransitGatewayAttachments != nil && len(op.TransitGatewayAttachments) > 0 {
		return op.TransitGatewayAttachments[0], nil
	}
	return nil, nil
}

func getAwsEc2TransitGatewayVpcAttachmentAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2TransitGatewayVpcAttachmentAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	transitGatewayAttachment := h.Item.(*ec2.TransitGatewayAttachment)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
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
	data := d.HydrateItem.(*ec2.TransitGatewayAttachment)
	return ec2TagsToMap(data.Tags)
}

func getEc2TransitGatewayAttachmentTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.TransitGatewayAttachment)
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

//// UTILITY FUNCTION
// Build ec2 transit gateway VPC attachment list call input filter
func buildEc2TransitGatewayVpcAttachmentFilter(quals plugin.KeyColumnQualMap) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)

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
			filter := ec2.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []*string{aws.String(val)}
			} else {
				v := value.([]*string)
				filter.Values = v
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
