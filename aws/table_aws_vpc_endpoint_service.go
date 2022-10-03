package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsVpcEndpointService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_endpoint_service",
		Description: "AWS VPC Endpoint Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("service_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidServiceName"}),
			},
			Hydrate: getVpcEndpointService,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcEndpointServices,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcEndpointServices", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	pagesLeft := true
	params := &ec2.DescribeVpcEndpointServicesInput{
		MaxResults: aws.Int64(1000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			if *limit < 1 {
				params.MaxResults = aws.Int64(1)
			} else {
				params.MaxResults = limit
			}
		}
	}

	// List call
	for pagesLeft {
		result, err := svc.DescribeVpcEndpointServices(params)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcEndpointServices", "DescribeVpcEndpointServices_error", err)
			return nil, err
		}

		for _, endpointService := range result.ServiceDetails {
			d.StreamListItem(ctx, endpointService)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcEndpointService(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcEndpointService")

	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceName := d.KeyColumnQuals["service_name"].GetStringValue()

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcEndpointServicesInput{
		ServiceNames: []*string{aws.String(serviceName)},
	}

	// Get call
	op, err := svc.DescribeVpcEndpointServices(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getVpcEndpointService__", "ERROR", err)
		return nil, err
	}

	if op.ServiceDetails != nil && len(op.ServiceDetails) > 0 {
		return op.ServiceDetails[0], nil
	}
	return nil, nil
}

func listVpcEndpointServicePermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listVpcEndpointServicePermissions")

	region := d.KeyColumnQualString(matrixKeyRegion)
	serviceId := h.Item.(*ec2.ServiceDetail).ServiceId

	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcEndpointServicePermissionsInput{
		ServiceId:  serviceId,
		MaxResults: aws.Int64(1000),
	}

	allowedPrincipals := []*ec2.AllowedPrincipal{}
	err = svc.DescribeVpcEndpointServicePermissionsPages(
		params,
		func(page *ec2.DescribeVpcEndpointServicePermissionsOutput, isLast bool) bool {
			allowedPrincipals = append(allowedPrincipals, page.AllowedPrincipals...)
			return !isLast
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listVpcEndpointServicePermissions", "ERROR", err)
		return nil, err
	}

	return allowedPrincipals, nil
}

func getVpcEndpointConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	endpointService := h.Item.(*ec2.ServiceDetail)

	// get service
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Debug("aws_vpc_endpoint_service.getVpcEndpointConnections", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &ec2.DescribeVpcEndpointConnectionsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("service-id"),
				Values: []*string{endpointService.ServiceId},
			},
		},
	}

	// Get call
	op, err := svc.DescribeVpcEndpointConnections(params)
	if err != nil {
		plugin.Logger(ctx).Debug("aws_vpc_endpoint_service.getVpcEndpointConnections", "api_error", err)
		return nil, err
	}

	if op.VpcEndpointConnections != nil && len(op.VpcEndpointConnections) > 0 {
		return op, nil
	}
	return nil, nil
}

func getVpcEndpointServiceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcEndpointServiceAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	endpointService := h.Item.(*ec2.ServiceDetail)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
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
	endpointService := d.HydrateItem.(*ec2.ServiceDetail)
	return ec2TagsToMap(endpointService.Tags)
}
