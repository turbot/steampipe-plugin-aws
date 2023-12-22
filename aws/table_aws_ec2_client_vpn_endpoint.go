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

//// TABLE DEFINITION

func tableAwsEC2ClientVPNEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_client_vpn_endpoint",
		Description: "AWS EC2 Client VPN Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("client_vpn_endpoint_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationError", "InvalidQueryParameter", "InvalidParameterValue", "InvalidClientVpnEndpointId.NotFound"}),
			},
			Hydrate: getEC2ClientVPNEndpoint,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeClientVpnEndpoints"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEC2ClientVPNEndpoints,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeClientVpnEndpoints"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "transport_protocol", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "client_vpn_endpoint_id",
				Description: "The ID of the client VPN endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "transport_protocol",
				Description: "The transport protocol.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_cidr_block",
				Description: "The IPv4 address range, in CIDR notation, from which client IP addresses are assigned.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "creation_time",
				Description: "The date and time when the Client VPN endpoint was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "deletion_time",
				Description: "The date and time when the Client VPN endpoint was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A brief description of the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_name",
				Description: "The DNS name to be used by clients when connecting to the Client VPN endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_service_portal_url",
				Description: "The URL of the self-service portal.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_certificate_arn",
				Description: "The ARN of the server certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "session_timeout_hours",
				Description: "The maximum VPN session duration time in hours. Valid values: 8, 10, 12, 24. Defaults to 24.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "split_tunnel",
				Description: "Indicates whether split-tunnel is enabled in the Client VPN endpoint.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpn_port",
				Description: "The port number for the Client VPN endpoint.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "authentication_options",
				Description: "Information about the authentication method used by the Client VPN endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "client_connect_options",
				Description: "The options for managing connection authorization for new client connections.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "client_log_options",
				Description: "Information about the client connection logging options for the Client VPN endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "client_login_banner_options",
				Description: "Options for enabling a customizable text banner that will be displayed on Amazon Web Services provided clients when a VPN session is established.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dns_servers",
				Description: "Information about the DNS servers to be used for DNS resolution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_group_ids",
				Description: "The IDs of the security groups for the target network.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status",
				Description: "The current state of the Client VPN endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "Any tags assigned to the Client VPN endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "vpn_protocol",
				Description: "The protocol used by the VPN session.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getEC2ClientVPNEndpointTurbotTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEC2ClientVPNEndpointTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listEC2ClientVPNEndpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_client_vpn_endpoint.listEC2ClientVPNEndpoints", "connection_error", err)
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

	input := &ec2.DescribeClientVpnEndpointsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filters := buildEC2ClientVPNEndpointFilter(d.EqualsQuals)

	if len(filters) != 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeClientVpnEndpointsPaginator(svc, input, func(o *ec2.DescribeClientVpnEndpointsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_client_vpn_endpoint.listEC2ClientVPNEndpoints", "api_error", err)
			return nil, err
		}

		for _, items := range output.ClientVpnEndpoints {

			d.StreamListItem(ctx, items)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}

		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEC2ClientVPNEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	clientVpnEndpointId := d.EqualsQualString("client_vpn_endpoint_id")

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_client_vpn_endpoint.getEC2ClientVPNEndpoint", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeClientVpnEndpointsInput{
		ClientVpnEndpointIds: []string{clientVpnEndpointId},
	}

	op, err := svc.DescribeClientVpnEndpoints(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_client_vpn_endpoint.getEC2ClientVPNEndpoint", "api_error", err)
		return nil, err
	}

	return op.ClientVpnEndpoints[0], nil
}

//// UTILITY FUNCTIONS

// Build EC2 client VPN endpoint list call input filter

func buildEC2ClientVPNEndpointFilter(equalQuals plugin.KeyColumnEqualsQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"transport_protocol": "transport-protocol",
	}

	for columnName, filterName := range filterQuals {
		if equalQuals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := equalQuals[columnName]
			if value.GetStringValue() != "" {
				filter.Values = []string{equalQuals[columnName].GetStringValue()}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}

//// TRANSFORM FUNCTIONS

func getEC2ClientVPNEndpointTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.ClientVpnEndpoint)
	title := data.ClientVpnEndpointId
	if data.Tags != nil {
		for _, i := range data.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}
	return title, nil
}

func getEC2ClientVPNEndpointTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.ClientVpnEndpoint)
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
