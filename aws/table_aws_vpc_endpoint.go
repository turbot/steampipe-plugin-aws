package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_endpoint",
		Description: "AWS VPC Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("vpc_endpoint_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVpcEndpointId.NotFound", "InvalidVpcEndpointId.Malformed"}),
			ItemFromKey:       vpcEndpointFromKey,
			Hydrate:           getVpcEndpoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEndpoints,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vpc_endpoint_id",
				Description: "The ID of the VPC endpoint",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The name of the service to which the endpoint is associated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the VPC endpoint",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC to which the endpoint is associated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the VPC endpoint",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_dns_enabled",
				Description: "Indicates whether the VPC is associated with a private hosted zone",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "requester_managed",
				Description: "Indicates whether the VPC endpoint is being managed by its service",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "policy",
				Description: "The policy document associated with the endpoint, if applicable",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyDocument").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyDocument").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "subnet_ids",
				Description: "One or more subnets in which the endpoint is located",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "route_table_ids",
				Description: "One or more route tables associated with the endpoint",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "groups",
				Description: "Information about the security groups that are associated with the network interface",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interface_ids",
				Description: "One or more network interfaces for the endpoint",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dns_entries",
				Description: "The DNS entries for the endpoint",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "creation_timestamp",
				Description: "The date and time that the VPC endpoint was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the VPC endpoint",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			//standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getVpcEndpointTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getVpcEndpointTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcEndpointAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// ITEM FROM KEY

func vpcEndpointFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	vpcEndpointID := quals["vpc_endpoint_id"].GetStringValue()
	item := &ec2.VpcEndpoint{
		VpcEndpointId: &vpcEndpointID,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcEndpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listVpcEndpoints", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.DescribeVpcEndpointsPages(
		&ec2.DescribeVpcEndpointsInput{},
		func(page *ec2.DescribeVpcEndpointsOutput, lastPage bool) bool {
			for _, item := range page.VpcEndpoints {
				d.StreamListItem(ctx, item)
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcEndpoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcEndpoint")
	vpcEndpoint := h.Item.(*ec2.VpcEndpoint)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []*string{vpcEndpoint.VpcEndpointId},
	}

	//get call
	item, err := svc.DescribeVpcEndpoints(params)
	if err != nil {
		logger.Debug("getVpcEndpoint__", "Error", err)
		return nil, err
	}

	if item.VpcEndpoints != nil && len(item.VpcEndpoints) > 0 {
		return item.VpcEndpoints[0], nil
	}

	return nil, nil
}

func getVpcEndpointAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcEndpointAkas")
	vpcEndpoint := h.Item.(*ec2.VpcEndpoint)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":vpc-endpoint/" + *vpcEndpoint.VpcEndpointId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcEndpointTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpcEndpoint := d.HydrateItem.(*ec2.VpcEndpoint)
	param := d.Param.(string)

	// Get resource title
	title := vpcEndpoint.VpcEndpointId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if vpcEndpoint.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range vpcEndpoint.Tags {
			turbotTagsMap[*i.Key] = *i.Value
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	if param == "Tags" {
		return turbotTagsMap, nil
	}

	return title, nil
}
