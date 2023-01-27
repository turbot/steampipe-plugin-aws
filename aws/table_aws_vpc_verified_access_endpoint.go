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

func tableAwsVpcVerifiedAccessEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_verified_access_endpoint",
		Description: "AWS VPC verified access Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("verified_access_endpoint_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue", "InvalidVerifiedAccessEndpointId.NotFound", "InvalidAction"}),
			},
			Hydrate: getVpcVerifiedAccessEndpoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcVerifiedAccessEndpoints,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			// DescribeVerifiedAccessEndpoints API accept endpoint id as input param.
			// We have to pass MaxResults value to DescribeVerifiedAccessEndpoints as input to perform pagination.
			// We can not pass both MaxResults and VerifiedAccessEndpointId at a time in the same input, we have to pass either one. So verified_access_endpoint_id can not be added as optional quals and added get config for filtering out the endpoint by their id.
			KeyColumns: []*plugin.KeyColumn{
				{Name: "verified_access_group_id", Require: plugin.Optional},
				{Name: "verified_access_instance_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "verified_access_endpoint_id",
				Description: "The ID of the AWS verified access endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "verified_access_group_id",
				Description: "The ID of the AWS verified access group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "verified_access_instance_id",
				Description: "The ID of the AWS verified access instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status_code",
				Description: "The endpoint status code. Possible values are pending, active, updating, deleting or deleted.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.Code"),
			},
			{
				Name:        "application_domain",
				Description: "The DNS name for users to reach your application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attachment_type",
				Description: "The type of attachment used to provide connectivity between the AWS verified access endpoint and the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_time",
				Description: "The deletion time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description for the AWS verified access endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "device_validation_domain",
				Description: "Returned if endpoint has a device trust provider attached.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_certificate_arn",
				Description: "The ARN of a public TLS/SSL certificate imported into or created with ACM.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_domain",
				Description: "A DNS name that is generated for the endpoint..",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_type",
				Description: "The type of AWS verified access endpoint. Incoming application requests will be sent to an IP address, load balancer or a network interface depending on the endpoint type specified. Possible values are load-balancer or network-interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "load_balancer_options",
				Description: "The load balancer details if creating the AWS verified access endpoint as load-balancertype.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interface_options",
				Description: "The options for network-interface type endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status",
				Description: "The endpoint status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(endpointTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(endpointTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcVerifiedAccessEndpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_endpoint.listVpcVerifiedAccessEndpoints", "connection_error", err)
		return nil, err
	}

	// AWS doc says the MaxResults value can be between 5-1000(https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeVerifiedAccessEndpoints.html) but API throws 'InvalidParameterValue' error for MaxResults value 1000.
	// Error: api error InvalidParameterValue: The parameter MaxResults must be between 5 and 200
	// Limiting the results
	maxLimit := int32(200)
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

	input := &ec2.DescribeVerifiedAccessEndpointsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.EqualsQualString("verified_access_group_id") != "" {
		input.VerifiedAccessGroupId = aws.String(d.EqualsQualString("verified_access_group_id"))
	}
	if d.EqualsQualString("verified_access_instance_id") != "" {
		input.VerifiedAccessInstanceId = aws.String(d.EqualsQualString("verified_access_instance_id"))
	}

	for {
		// List call
		resp, err := svc.DescribeVerifiedAccessEndpoints(ctx, input)

		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_verified_access_endpoint.listVpcVerifiedAccessEndpoints", "api_error", err)
			return nil, nil
		}

		for _, endpoint := range resp.VerifiedAccessEndpoints {
			d.StreamListItem(ctx, endpoint)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.NextToken == nil {
			break
		} else {
			input.NextToken = resp.NextToken
		}
	}

	return nil, err
}

//// HYDRATED FUNCTION

func getVpcVerifiedAccessEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	endpointId := d.EqualsQuals["verified_access_endpoint_id"].GetStringValue()

	// Empty check
	if endpointId == "" {
		return nil, nil
	}

	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_endpoint.getVpcVerifiedAccessEndpoint", "connection_error", err)
		return nil, err
	}

	// Build the params
	input := &ec2.DescribeVerifiedAccessEndpointsInput{
		VerifiedAccessEndpointIds: []string{endpointId},
	}

	// Get call
	op, err := svc.DescribeVerifiedAccessEndpoints(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_verified_access_endpoint.getVpcVerifiedAccessEndpoint", "api_error", err)
		return nil, err
	}

	if op.VerifiedAccessEndpoints != nil && len(op.VerifiedAccessEndpoints) > 0 {
		return op.VerifiedAccessEndpoints[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func endpointTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	accessPoint := d.HydrateItem.(types.VerifiedAccessEndpoint)

	// Get the resource tags
	var turbotTagsMap map[string]string
	if accessPoint.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range accessPoint.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func endpointTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	accessPoint := d.HydrateItem.(types.VerifiedAccessEndpoint)
	title := accessPoint.VerifiedAccessEndpointId

	if accessPoint.Tags != nil {
		for _, i := range accessPoint.Tags {
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	return title, nil
}
