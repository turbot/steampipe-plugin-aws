package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsVpcEndpointService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_endpoint_service",
		Description: "AWS VPC Endpoint Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("service_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidServiceName"}),
			},
			Hydrate: getVpcEndpointService,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVpcEndpointServices"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEndpointServices,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeVpcEndpointServices"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Description: "The Amazon Resource Name (ARN) of the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_id",
				Description: "The ID of the endpoint service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner",
				Description: "The AWS account ID of the service owner.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "acceptance_required",
				Description: "Indicates whether VPC endpoint connection requests to the service must be accepted by the service owner.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "manages_vpc_endpoints",
				Description: "Indicates whether the service manages its VPC endpoints. Management of the service VPC endpoints using the VPC endpoint API is restricted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "private_dns_name",
				Description: "The private DNS name for the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_dns_name_verification_state",
				Description: "The verification state of the VPC endpoint service. Consumers of the endpoint service cannot use the private name when the state is not verified.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_endpoint_policy_supported",
				Description: "Indicates whether the service supports endpoint policies.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "availability_zones",
				Description: "The Availability Zones in which the service is available.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "base_endpoint_dns_names",
				Description: "The DNS names for the service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_type",
				Description: "The type of service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_endpoint_connections",
				Description: "Information about one or more VPC endpoint connections.",
				Hydrate:     getVpcEndpointConnections,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_endpoint_service_permissions",
				Description: "Information about one or more allowed principals.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listVpcEndpointServicePermissions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getVpcEndpointServiceTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcEndpointServiceAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcEndpointServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint_service.listVpcEndpointServices", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = int32(5)
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ec2.DescribeVpcEndpointServicesInput{
		MaxResults: &maxLimit,
	}

	// API doesn't support aws-sdk-go-v2 paginator as of date
	pagesLeft := true
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)
		
		result, err := svc.DescribeVpcEndpointServices(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_vpc_endpoint_service.listVpcEndpointServices", "api_error", err)
			return nil, err
		}

		for _, endpointService := range result.ServiceDetails {
			d.StreamListItem(ctx, endpointService)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			input.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVpcEndpointService(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	serviceName := d.EqualsQuals["service_name"].GetStringValue()

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint_service.getVpcEndpointService", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcEndpointServicesInput{
		ServiceNames: []string{serviceName},
	}

	// Get call
	op, err := svc.DescribeVpcEndpointServices(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint_service.getVpcEndpointService", "api_error", err)
		return nil, err
	}

	if op.ServiceDetails != nil && len(op.ServiceDetails) > 0 {
		return op.ServiceDetails[0], nil
	}
	return nil, nil
}

func listVpcEndpointServicePermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	serviceId := h.Item.(types.ServiceDetail).ServiceId

	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint_service.listVpcEndpointServicePermissions", "connection_error", err)
		return nil, err
	}

	// Build the input
	input := &ec2.DescribeVpcEndpointServicePermissionsInput{
		ServiceId:  serviceId,
		MaxResults: aws.Int32(1000),
	}

	paginator := ec2.NewDescribeVpcEndpointServicePermissionsPaginator(svc, input, func(o *ec2.DescribeVpcEndpointServicePermissionsPaginatorOptions) {
		o.Limit = 1000
		o.StopOnDuplicateToken = true
	})

	var allowedPrincipals []types.AllowedPrincipal
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			if err != nil {
				if strings.Contains(err.Error(), "NotFound") {
					return nil, nil
				}
				plugin.Logger(ctx).Error("aws_vpc_endpoint_service.listVpcEndpointServicePermissions", "api_error", err)
				return nil, err
			}

			allowedPrincipals = append(allowedPrincipals, output.AllowedPrincipals...)
		}
	}
	return allowedPrincipals, nil
}

func getVpcEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	endpointService := h.Item.(types.ServiceDetail)

	// get service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint_service.getVpcEndpointConnections", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcEndpointConnectionsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("service-id"),
				Values: []string{*endpointService.ServiceId},
			},
		},
	}

	// Get call
	op, err := svc.DescribeVpcEndpointConnections(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint_service.getVpcEndpointConnections", "api_error", err)
		return nil, err
	}

	if op.VpcEndpointConnections != nil && len(op.VpcEndpointConnections) > 0 {
		return op, nil
	}

	return nil, nil
}

func getVpcEndpointServiceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	region := d.EqualsQualString(matrixKeyRegion)
	endpointService := h.Item.(types.ServiceDetail)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_vpc_endpoint_service.getVpcEndpointServiceAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	splitServicName := strings.Split(*endpointService.ServiceName, ".")
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":vpc-endpoint-service/" + splitServicName[len(splitServicName)-1]}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcEndpointServiceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	endpointService := d.HydrateItem.(types.ServiceDetail)
	if len(endpointService.Tags) > 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range endpointService.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}
