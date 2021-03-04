package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type awsServiceRegion struct {
	Name       string
	RegionName string
	Endpoint   string
}

const paramServiceRegionPrefix string = "/aws/service/global-infrastructure/regions/"

//// TABLE DEFINITION

func tableAwsServiceRegion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_service_region",
		Description: "AWS Service Region",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "region_name"}),
			ShouldIgnoreError: isNotFoundError([]string{"ValidationError", "ParameterNotFound"}),
			Hydrate:           getAwsServiceRegion,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsRegions,
			Hydrate:       listAwsRegionServices,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the AWS Service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "region_name",
				Description: "The name of the Region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint",
				Description: "The URL of the service's regional endpoint.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsServiceRegionEndpoint,
				Transform:   transform.FromField("Value"),
			},
			{
				Name:        "protocols",
				Description: "The available protocol (e.g. HTTPS, HTTP) for a given regional endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsServiceRegionProtocols,
				Transform:   transform.FromField("Value").Transform(csvToArray),
			},
		},
	}
}

//// LIST FUNCTION

func listAwsRegionServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)
	region := h.Item.(*ec2.Region)
	serviceRegionParamPath := paramServiceRegionPrefix + *region.RegionName + "/services/"

	// Create Session
	svc, err := SsmService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ssm.GetParametersByPathInput{
		Path: aws.String(serviceRegionParamPath),
	}

	pagesLeft := true

	for pagesLeft {
		result, err := svc.GetParametersByPath(params)
		if err != nil {
			return nil, err
		}

		for _, service := range result.Parameters {
			d.StreamLeafListItem(ctx, awsServiceRegion{
				Name:       *service.Value,
				RegionName: *region.RegionName,
			})
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsServiceRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)

	var serviceRegion string
	var serviceName string
	if h.Item != nil {
		i := h.Item.(*awsServiceRegion)
		serviceRegion = i.RegionName
		serviceName = i.Name
	} else {
		serviceRegion = d.KeyColumnQuals["region_name"].GetStringValue()
		serviceName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := SsmService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	// e.g. "/aws/service/global-infrastructure/regions/us-east-1/services/sqs"
	serviceRegionParamPath := paramServiceRegionPrefix + serviceRegion + "/services/" + serviceName

	params := &ssm.GetParameterInput{
		Name: aws.String(serviceRegionParamPath),
	}

	result, err := svc.GetParameter(params)
	if err != nil {
		return nil, err
	}

	if result != nil && result.Parameter != nil {
		return awsServiceRegion{
			Name:       *result.Parameter.Value,
			RegionName: serviceRegion,
		}, nil
	}

	return nil, nil

}

// Column Hydrate Functions

func getAwsServiceRegionEndpoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getAwsServiceRegionParam(ctx, d, h, "/endpoint")
}

func getAwsServiceRegionProtocols(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getAwsServiceRegionParam(ctx, d, h, "/protocols")
}

func getAwsServiceRegionParam(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, p string) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)
	var serviceRegion string
	var serviceName string
	if h.Item != nil {
		i := h.Item.(awsServiceRegion)
		serviceRegion = i.RegionName
		serviceName = i.Name
	} else {
		serviceRegion = d.KeyColumnQuals["region_name"].GetStringValue()
		serviceName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// e.g. "/aws/service/global-infrastructure/regions/us-east-1/services/sqs/endpoint"
	serviceRegionParamPath := paramServiceRegionPrefix + serviceRegion + "/services/" + serviceName + p
	plugin.Logger(ctx).Error("getAwsServiceRegionParam", "serviceRegionParamPath", serviceRegionParamPath)

	// Create Session
	svc, err := SsmService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ssm.GetParameterInput{
		Name: aws.String(serviceRegionParamPath),
	}

	param, err := svc.GetParameter(params)
	if err != nil {
		if strings.Contains(err.Error(), "ParameterNotFound") {
			return "", nil
		}
		return nil, err
	}

	if param != nil && param.Parameter != nil {
		return param.Parameter, nil
	}

	return nil, nil

}
