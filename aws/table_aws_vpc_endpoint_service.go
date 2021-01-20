package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsVpcEndpointService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_endpoint_service",
		Description: "AWS VPC Endpoint Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("service_name"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidServiceName"}),
			ItemFromKey:       endpointServiceFromKey,
			Hydrate:           getVpcEndpointService,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEndpointServices,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Description: "The Amazon Resource Name (ARN) of the service",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_id",
				Description: "The ID of the endpoint service",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner",
				Description: "The AWS account ID of the service owner",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "acceptance_required",
				Description: "Indicates whether VPC endpoint connection requests to the service must be accepted by the service owner",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "manages_vpc_endpoints",
				Description: "Indicates whether the service manages its VPC endpoints. Management of the service VPC endpoints using the VPC endpoint API is restricted",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "private_dns_name",
				Description: "The private DNS name for the service",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_dns_name_verification_state",
				Description: "The verification state of the VPC endpoint service. Consumers of the endpoint service cannot use the private name when the state is not verified",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_endpoint_policy_supported",
				Description: "Indicates whether the service supports endpoint policies",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "service_type",
				Description: "The type of service",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "base_endpoint_dns_names",
				Description: "The DNS names for the service",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "availability_zones",
				Description: "The Availability Zones in which the service is available",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_raw",
				Description: "A list of tags assigned to the service",
				Type:        proto.ColumnType_JSON,
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

//// ITEM FROM KEY

func endpointServiceFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	serviceName := quals["service_name"].GetStringValue()
	item := &ec2.ServiceDetail{
		ServiceName: &serviceName,
	}
	return item, nil
}

//// LIST FUNCTION

func listVpcEndpointServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listVpcEndpointServices", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	data, err := svc.DescribeVpcEndpointServices(&ec2.DescribeVpcEndpointServicesInput{})
	for _, endpointService := range data.ServiceDetails {
		d.StreamListItem(ctx, endpointService)
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcEndpointService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcEndpointService")
	endpointService := h.Item.(*ec2.ServiceDetail)
	defaultRegion := GetDefaultRegion()

	// get service
	svc, err := Ec2Service(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcEndpointServicesInput{
		ServiceNames: []*string{endpointService.ServiceName},
	}

	// Get call
	op, err := svc.DescribeVpcEndpointServices(params)
	if err != nil {
		logger.Debug("getVpcEndpointService__", "ERROR", err)
		return nil, err
	}

	if op.ServiceDetails != nil && len(op.ServiceDetails) > 0 {
		return op.ServiceDetails[0], nil
	}
	return nil, nil
}

func getVpcEndpointServiceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcEndpointServiceAkas")
	endpointService := h.Item.(*ec2.ServiceDetail)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	splitServicName := strings.Split(*endpointService.ServiceName, ".")
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":vpc-endpoint-service/" + splitServicName[len(splitServicName)-1]}

	// plugin.Logger(ctx).Trace("getVpcEndpointServiceAkas", "Akas", akas)

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getVpcEndpointServiceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	endpointService := d.HydrateItem.(*ec2.ServiceDetail)
	return ec2TagsToMap(endpointService.Tags)
}
