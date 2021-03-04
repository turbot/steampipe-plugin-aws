package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

type awsService struct {
	Name         string
	LongName     string
	MarketingUrl string
}

const paramServicePrefix string = "/aws/service/global-infrastructure/services/"

func tableAwsService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_service",
		Description: "AWS Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getAwsService,
			ShouldIgnoreError: isNotFoundError([]string{"ValidationError", "ParameterNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsServices,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The short name of the AWS Service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "long_name",
				Description: "The long marketing name of the AWS Service",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsServiceLongName,
				Transform:   transform.FromField("Value"),
			},
			{
				Name:        "marketing_url",
				Description: "The long marketing name of the AWS Service",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsServiceMarketingUrl,
				Transform:   transform.FromField("Value"),
			},
		},
	}
}

//// LIST FUNCTION

func listAwsServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)

	// Create Session
	svc, err := SsmService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	pagesLeft := true

	params := &ssm.GetParametersByPathInput{
		Path: aws.String(paramServicePrefix),
	}

	for pagesLeft {
		result, err := svc.GetParametersByPath(params)
		if err != nil {
			return nil, err
		}

		for _, service := range result.Parameters {
			d.StreamListItem(ctx, awsService{
				Name: *service.Value,
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

func getAwsService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)
	var serviceName string
	if h.Item != nil {
		i := h.Item.(*awsService)
		serviceName = i.Name
	} else {
		serviceName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// Create Session
	svc, err := SsmService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ssm.GetParameterInput{
		Name: aws.String(paramServicePrefix + serviceName),
	}

	result, err := svc.GetParameter(params)
	if err != nil {
		return nil, err
	}

	if result != nil && result.Parameter != nil {
		return awsService{
			Name: *result.Parameter.Value,
		}, nil
	}

	return nil, nil

}

func getAwsServiceLongName(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getAwsServiceParam(ctx, d, h, "/longName")
}

func getAwsServiceMarketingUrl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getAwsServiceParam(ctx, d, h, "/marketingHomeURL")
}

func getAwsServiceParam(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, p string) (interface{}, error) {
	defaultRegion := GetDefaultAwsRegion(d)
	var serviceName string
	if h.Item != nil {
		i := h.Item.(awsService)
		serviceName = i.Name
	} else {
		serviceName = d.KeyColumnQuals["name"].GetStringValue()
	}

	serviceParamPath := paramServicePrefix + serviceName + p

	// Create Session
	svc, err := SsmService(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ssm.GetParameterInput{
		Name: aws.String(serviceParamPath),
	}

	param, err := svc.GetParameter(params)
	if err != nil {
		// plugin.Logger(ctx).Error("getAwsServiceParam", "serviceParamPath", serviceParamPath)
		// plugin.Logger(ctx).Error("getAwsServiceParam", "err", err)
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
