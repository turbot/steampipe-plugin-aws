package aws

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayV2DomainName(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gatewayv2_domain_name",
		Description: "AWS API Gateway Version 2 Domain Name",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"domain_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"NotFoundException"}),
			},
			Hydrate: getDomainName,
		},
		List: &plugin.ListConfig{
			Hydrate: listDomainNames,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := APIGatewayV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_domain_name.listDomainNames", "service_client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	params := &apigatewayv2.GetDomainNamesInput{
		MaxResults: aws.String(fmt.Sprint(maxLimit)),
	}
	pagesLeft := true

	for pagesLeft {
		result, err := svc.GetDomainNames(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gatewayv2_domain_name.listDomainNames", "api_error", err)
			return nil, err
		}

		for _, domainName := range result.Items {
			d.StreamListItem(ctx, domainName)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
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

	// Create Session
	svc, err := APIGatewayV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_domain_name.getDomainName", "service_client_error", err)
		return nil, err
	}

	domainName := d.KeyColumnQuals["domain_name"].GetStringValue()
	input := &apigatewayv2.GetDomainNameInput{
		DomainName: aws.String(domainName),
	}

	op, err := svc.GetDomainName(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gatewayv2_domain_name.getDomainName", "api_error", err)
		return nil, err
	}

	if op != nil {
		domainName := &types.DomainName{
			DomainName:                    op.DomainName,
			Tags:                          op.Tags,
			ApiMappingSelectionExpression: op.ApiMappingSelectionExpression,
			MutualTlsAuthentication:       op.MutualTlsAuthentication,
			DomainNameConfigurations:      op.DomainNameConfigurations,
		}
		return domainName, nil
	}

	return nil, nil
}

func getapiGatewayV2DomainNameAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	domainName := ""

	switch h.Item.(type) {
	case *types.DomainName:
		domainName = *h.Item.(*types.DomainName).DomainName
	case types.DomainName:
		domainName = *h.Item.(types.DomainName).DomainName
	}

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/domainnames/" + domainName}

	return akas, nil
}
