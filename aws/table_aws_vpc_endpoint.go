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

func tableAwsVpcEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_endpoint",
		Description: "AWS VPC Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vpc_endpoint_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidVpcEndpointId.NotFound", "InvalidVpcEndpointId.Malformed"}),
			},
			Hydrate: getVpcEndpoint,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVpcEndpoints"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEndpoints,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVpcEndpoints"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_name", Require: plugin.Optional},
				{Name: "vpc_id", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vpc_endpoint_id",
				Description: "The ID of the VPC endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The name of the service to which the endpoint is associated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the AWS account that owns the VPC endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC to which the endpoint is associated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_endpoint_type",
				Description: "The type of endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the VPC endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_dns_enabled",
				Description: "Indicates whether the VPC is associated with a private hosted zone.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "requester_managed",
				Description: "Indicates whether the VPC endpoint is being managed by its service.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "policy",
				Description: "The policy document associated with the endpoint, if applicable.",
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
				Description: "One or more subnets in which the endpoint is located.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "route_table_ids",
				Description: "One or more route tables associated with the endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "groups",
				Description: "Information about the security groups that are associated with the network interface.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interface_ids",
				Description: "One or more network interfaces for the endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dns_entries",
				Description: "The DNS entries for the endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "creation_timestamp",
				Description: "The date and time that the VPC endpoint was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the VPC endpoint.",
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

//// LIST FUNCTION

func listVpcEndpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint.listVpcEndpoints", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = int32(1)
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeVpcEndpointsInput{
		MaxResults: aws.Int32(maxLimit),
	}
	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "service_name", FilterName: "service-name", ColumnType: "string"},
		{ColumnName: "vpc_id", FilterName: "vpc-id", ColumnType: "string"},
		{ColumnName: "state", FilterName: "vpc-endpoint-state", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeVpcEndpointsPaginator(svc, input, func(o *ec2.DescribeVpcEndpointsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_endpoint.listVpcEndpoints", "api_error", err)
			return nil, err
		}

		for _, items := range output.VpcEndpoints {
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

func getVpcEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	vpcEndpointID := d.EqualsQuals["vpc_endpoint_id"].GetStringValue()

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint.getVpcEndpoint", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []string{vpcEndpointID},
	}

	//get call
	item, err := svc.DescribeVpcEndpoints(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint.getVpcEndpoint", "api_error", err)
		return nil, err
	}

	if item != nil && len(item.VpcEndpoints) > 0 {
		return item.VpcEndpoints[0], nil
	}

	return nil, nil
}

func getVpcEndpointAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	vpcEndpoint := h.Item.(types.VpcEndpoint)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint.getVpcEndpointAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":vpc-endpoint/" + *vpcEndpoint.VpcEndpointId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcEndpointTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpcEndpoint := d.HydrateItem.(types.VpcEndpoint)
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
