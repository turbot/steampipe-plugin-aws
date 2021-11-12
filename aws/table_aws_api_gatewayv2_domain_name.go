package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayV2DomainName(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gatewayv2_domain_name",
		Description: "AWS API Gateway Version 2 Domain Name",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"domain_name"}),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			Hydrate:           getDomainName,
		},
		List: &plugin.ListConfig{
			Hydrate: listDomainNames,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "The name of the DomainName resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_name_configurations",
				Description: "The domain name configurations",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "mutual_tls_authentication",
				Description: "The mutual TLS authentication configuration for a custom domain name",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MutualTlsAuthentication"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getapiGatewayV2DomainNameAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listDomainNames(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := APIGatewayV2Service(ctx, d)
	if err != nil {
		return nil, err
	}
	params := &apigatewayv2.GetDomainNamesInput{}
	pagesLeft := true

	for pagesLeft {
		result, err := svc.GetDomainNames(params)
		if err != nil {
			return nil, err
		}

		for _, domainName := range result.Items {
			d.StreamListItem(ctx, domainName)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				pagesLeft = false
			}
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDomainName(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDomainName")

	// Create Session
	svc, err := APIGatewayV2Service(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Debug("getDomainName__", "ERROR", err)
		return nil, err
	}

	domainName := d.KeyColumnQuals["domain_name"].GetStringValue()
	input := &apigatewayv2.GetDomainNameInput{
		DomainName: aws.String(domainName),
	}

	op, err := svc.GetDomainName(input)
	if err != nil {
		return nil, err
	}

	if op != nil {
		domainName := &apigatewayv2.DomainName{
			DomainName:                    op.DomainName,
			Tags:                          op.Tags,
			ApiMappingSelectionExpression: op.ApiMappingSelectionExpression,
			DomainNameConfigurations:      op.DomainNameConfigurations,
		}
		return domainName, nil
	}

	return nil, nil
}

func getapiGatewayV2DomainNameAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	v2ApiDomain := h.Item.(*apigatewayv2.DomainName)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/domainnames/" + *v2ApiDomain.DomainName}

	return akas, nil
}
